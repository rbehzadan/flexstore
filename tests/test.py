#!/usr/bin/env python3
"""
Comprehensive test script for the Schemaless API Server.
This script tests all endpoints, including health check, collection operations,
document operations, and bulk operations.
"""

import requests
import json
import time
import random
import string
import sys
from termcolor import colored
from datetime import datetime

# Configuration
BASE_URL = "http://localhost:8080"
COLLECTION_NAME = f"test_collection_{int(time.time())}"
DOCUMENT_IDS = []

# Statistics
tests_run = 0
tests_passed = 0
tests_failed = 0

def random_string(length=10):
    """Generate a random string of letters and digits."""
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

def print_header(text):
    """Print a header for the test section."""
    print("\n" + "=" * 80)
    print(colored(f" {text} ", "white", "on_blue"))
    print("=" * 80)

def print_test(description, passed, response=None, details=None):
    """Print the result of a test."""
    global tests_run, tests_passed, tests_failed
    
    tests_run += 1
    if passed:
        tests_passed += 1
        status = colored("✓ PASS", "green")
    else:
        tests_failed += 1
        status = colored("✗ FAIL", "red")
    
    print(f"{status} - {description}")
    
    if not passed and details:
        print(f"       {colored('Error:', 'red')} {details}")
    
    if response:
        try:
            status_code = response.status_code
            response_json = response.json()
            print(f"       Status: {status_code}")
            if not passed or '--verbose' in sys.argv:
                print(f"       Response: {json.dumps(response_json, indent=2)}")
        except:
            print(f"       Response (non-JSON): {response.text}")

def test_health_endpoint():
    """Test the health endpoint."""
    print_header("Testing Health Endpoint")
    
    response = requests.get(f"{BASE_URL}/health")
    is_success = response.status_code == 200
    
    try:
        data = response.json().get('data', {})
        has_status = data.get('status') == 'ok'
        has_version = 'version' in data
        has_uptime = 'uptime' in data
        is_valid = has_status and has_version and has_uptime
    except:
        is_valid = False
    
    print_test("Health endpoint returns status OK", is_success and is_valid, response)
    
    return is_success

def test_collections():
    """Test collection operations."""
    print_header("Testing Collection Operations")
    
    # Test creating a collection
    print_header("Creating Collection")
    create_payload = {"name": COLLECTION_NAME}
    response = requests.post(f"{BASE_URL}/api/collections", json=create_payload)
    create_success = response.status_code == 201
    print_test(f"Create collection '{COLLECTION_NAME}'", create_success, response)
    
    # Test listing collections
    print_header("Listing Collections")
    response = requests.get(f"{BASE_URL}/api/collections")
    list_success = response.status_code == 200
    
    # Check if our collection is in the list
    collection_in_list = False
    if list_success:
        try:
            collections = response.json().get('data', {}).get('collections', [])
            collection_in_list = any(c.get('name') == COLLECTION_NAME for c in collections)
        except:
            pass
    
    print_test("List collections", list_success, response)
    print_test(f"Collection '{COLLECTION_NAME}' is in the list", collection_in_list, 
               details="Collection not found in the list" if not collection_in_list else None)
    
    # Test getting a specific collection
    print_header("Getting Collection")
    response = requests.get(f"{BASE_URL}/api/collections/{COLLECTION_NAME}")
    get_success = response.status_code == 200
    print_test(f"Get collection '{COLLECTION_NAME}'", get_success, response)
    
    return create_success and list_success and get_success

def test_documents():
    """Test document operations."""
    print_header("Testing Document Operations")
    
    # Test creating a document
    print_header("Creating Document")
    doc_data = {
        "name": f"Test User {random_string()}",
        "email": f"test_{random_string(5)}@example.com",
        "age": random.randint(18, 80),
        "active": random.choice([True, False]),
        "created": datetime.now().isoformat(),
        "tags": ["test", random_string(), random_string()]
    }
    
    # Updated to use /documents path
    response = requests.post(f"{BASE_URL}/api/collections/{COLLECTION_NAME}/documents", json=doc_data)
    create_success = response.status_code == 201
    
    doc_id = None
    if create_success:
        try:
            doc_id = response.json().get('data', {}).get('id')
            DOCUMENT_IDS.append(doc_id)
        except:
            pass
    
    print_test("Create document", create_success, response)
    print_test("Document has ID", doc_id is not None, 
               details="No document ID returned" if doc_id is None else None)
    
    if not doc_id:
        return False
    
    # Test getting a document
    print_header("Getting Document")
    response = requests.get(f"{BASE_URL}/api/collections/{COLLECTION_NAME}/documents/{doc_id}")
    get_success = response.status_code == 200
    
    # Check if document data matches
    data_matches = False
    if get_success:
        try:
            returned_data = response.json().get('data', {}).get('data', {})
            for key, value in doc_data.items():
                if key not in returned_data or returned_data[key] != value:
                    data_matches = False
                    break
            else:
                data_matches = True
        except:
            pass
    
    print_test(f"Get document '{doc_id}'", get_success, response)
    print_test("Document data matches", data_matches, 
               details="Returned data does not match original" if not data_matches else None)
    
    # Test updating a document
    print_header("Updating Document")
    updated_doc_data = doc_data.copy()
    updated_doc_data["name"] = f"Updated User {random_string()}"
    updated_doc_data["age"] = random.randint(18, 80)
    updated_doc_data["tags"].append("updated")
    
    response = requests.put(f"{BASE_URL}/api/collections/{COLLECTION_NAME}/documents/{doc_id}", json=updated_doc_data)
    update_success = response.status_code == 200
    
    # Check if document was updated correctly
    update_correct = False
    if update_success:
        try:
            returned_data = response.json().get('data', {}).get('data', {})
            for key, value in updated_doc_data.items():
                if key not in returned_data or returned_data[key] != value:
                    update_correct = False
                    break
            else:
                update_correct = True
        except:
            pass
    
    print_test(f"Update document '{doc_id}'", update_success, response)
    print_test("Document updated correctly", update_correct, 
               details="Update didn't apply correctly" if not update_correct else None)
    
    # Test listing documents
    print_header("Listing Documents")
    # Updated to use /documents path
    response = requests.get(f"{BASE_URL}/api/collections/{COLLECTION_NAME}/documents")
    list_success = response.status_code == 200
    
    # Check if our document is in the list
    doc_in_list = False
    if list_success:
        try:
            documents = response.json().get('data', {}).get('documents', [])
            doc_in_list = any(d.get('id') == doc_id for d in documents)
        except:
            pass
    
    print_test("List documents", list_success, response)
    print_test(f"Document '{doc_id}' is in the list", doc_in_list, 
               details="Document not found in the list" if not doc_in_list else None)
    
    # Test pagination for document listing
    print_header("Testing Document Pagination")
    # Updated to use /documents path
    response = requests.get(f"{BASE_URL}/api/collections/{COLLECTION_NAME}/documents?limit=5&offset=0")
    pagination_success = response.status_code == 200
    
    # Check pagination parameters in response
    pagination_params_correct = False
    if pagination_success:
        try:
            data = response.json().get('data', {})
            pagination_params_correct = (
                data.get('limit') == 5 and
                data.get('offset') == 0
            )
        except:
            pass
    
    print_test("List documents with pagination", pagination_success, response)
    print_test("Pagination parameters correct", pagination_params_correct, 
               details="Pagination parameters not correctly reflected" if not pagination_params_correct else None)
    
    return create_success and get_success and update_success and list_success

def test_bulk_operations():
    """Test bulk operations."""
    print_header("Testing Bulk Operations")
    
    # Test bulk document creation
    print_header("Bulk Creating Documents")
    bulk_docs = []
    for i in range(3):
        bulk_docs.append({
            "name": f"Bulk User {i} {random_string(3)}",
            "email": f"bulk_{i}_{random_string(3)}@example.com",
            "age": random.randint(18, 80),
            "bulk_created": True
        })
    
    response = requests.post(f"{BASE_URL}/api/collections/{COLLECTION_NAME}/bulk", json=bulk_docs)
    bulk_create_success = response.status_code == 201
    
    # Check if all documents were created
    bulk_count_correct = False
    new_doc_ids = []
    if bulk_create_success:
        try:
            result = response.json().get('data', {})
            count = result.get('count', 0)
            documents = result.get('documents', [])
            bulk_count_correct = count == len(bulk_docs) and len(documents) == len(bulk_docs)
            
            # Get the new document IDs
            for doc in documents:
                if doc.get('id'):
                    new_doc_ids.append(doc.get('id'))
                    DOCUMENT_IDS.append(doc.get('id'))
        except:
            pass
    
    print_test("Bulk create documents", bulk_create_success, response)
    print_test("Correct number of documents created", bulk_count_correct, 
               details=f"Expected {len(bulk_docs)} documents" if not bulk_count_correct else None)
    
    # Check if we can get one of the bulk-created documents
    bulk_doc_get_success = False
    if new_doc_ids:
        response = requests.get(f"{BASE_URL}/api/collections/{COLLECTION_NAME}/documents/{new_doc_ids[0]}")
        bulk_doc_get_success = response.status_code == 200
        print_test(f"Get bulk-created document '{new_doc_ids[0]}'", bulk_doc_get_success, response)
    
    return bulk_create_success and bulk_count_correct and bulk_doc_get_success

def test_file_upload():
    """Test file upload functionality."""
    print_header("Testing File Upload")
    
    # Create a temporary JSON file
    upload_docs = []
    for i in range(2):
        upload_docs.append({
            "name": f"Uploaded User {i} {random_string(3)}",
            "email": f"upload_{i}_{random_string(3)}@example.com",
            "source": "file_upload",
            "timestamp": datetime.now().isoformat()
        })
    
    # Create a temporary file
    file_name = f"test_upload_{random_string()}.json"
    with open(file_name, 'w') as f:
        json.dump(upload_docs, f)
    
    # Upload the file
    files = {'file': open(file_name, 'rb')}
    response = requests.post(f"{BASE_URL}/api/upload/{COLLECTION_NAME}", files=files)
    upload_success = response.status_code == 201
    
    # Check if all documents were created
    upload_count_correct = False
    if upload_success:
        try:
            result = response.json().get('data', {})
            count = result.get('count', 0)
            documents = result.get('documents', [])
            upload_count_correct = count == len(upload_docs) and len(documents) == len(upload_docs)
            
            # Get the document IDs
            for doc in documents:
                if doc.get('id'):
                    DOCUMENT_IDS.append(doc.get('id'))
        except:
            pass
    
    print_test("Upload JSON file", upload_success, response)
    print_test("Correct number of documents created from file", upload_count_correct, 
               details=f"Expected {len(upload_docs)} documents" if not upload_count_correct else None)
    
    # Clean up the file
    import os
    os.remove(file_name)
    
    return upload_success and upload_count_correct

def test_document_deletion():
    """Test document deletion."""
    print_header("Testing Document Deletion")
    
    # We'll delete one document and keep the rest for collection deletion test
    if not DOCUMENT_IDS:
        print_test("Skip document deletion test", False, details="No document IDs available")
        return False
    
    doc_id = DOCUMENT_IDS.pop(0)
    response = requests.delete(f"{BASE_URL}/api/collections/{COLLECTION_NAME}/documents/{doc_id}")
    delete_success = response.status_code == 200
    
    # Verify document no longer exists
    verification_response = requests.get(f"{BASE_URL}/api/collections/{COLLECTION_NAME}/documents/{doc_id}")
    verification_success = verification_response.status_code == 404
    
    print_test(f"Delete document '{doc_id}'", delete_success, response)
    print_test("Document no longer exists", verification_success, verification_response)
    
    return delete_success and verification_success

def test_collection_deletion():
    """Test collection deletion."""
    print_header("Testing Collection Deletion")
    
    response = requests.delete(f"{BASE_URL}/api/collections/{COLLECTION_NAME}")
    delete_success = response.status_code == 200
    
    # Verify collection no longer exists
    verification_response = requests.get(f"{BASE_URL}/api/collections/{COLLECTION_NAME}")
    verification_success = verification_response.status_code == 404
    
    # Verify all documents in the collection are deleted
    if DOCUMENT_IDS and delete_success:
        for doc_id in DOCUMENT_IDS:
            doc_response = requests.get(f"{BASE_URL}/api/collections/{COLLECTION_NAME}/documents/{doc_id}")
            if doc_response.status_code != 404:
                verification_success = False
                break
    
    print_test(f"Delete collection '{COLLECTION_NAME}'", delete_success, response)
    print_test("Collection no longer exists", verification_success, verification_response)
    
    return delete_success and verification_success

def run_all_tests():
    """Run all tests and report results."""
    print(colored(f"\nTESTING SCHEMALESS API SERVER AT {BASE_URL}", "white", "on_blue"))
    print(f"Test Collection: {COLLECTION_NAME}")
    print("=" * 80)
    
    # Run tests
    tests = [
        ("Health Endpoint", test_health_endpoint),
        ("Collection Operations", test_collections),
        ("Document Operations", test_documents),
        ("Bulk Operations", test_bulk_operations),
        ("File Upload", test_file_upload),
        ("Document Deletion", test_document_deletion),
        ("Collection Deletion", test_collection_deletion),
    ]
    
    results = {}
    for name, test_func in tests:
        try:
            result = test_func()
            results[name] = result
        except Exception as e:
            print(colored(f"\nERROR in {name} test: {str(e)}", "red"))
            results[name] = False
    
    # Print summary
    print("\n" + "=" * 80)
    print(colored("TEST SUMMARY", "white", "on_blue"))
    print("=" * 80)
    
    for name, result in results.items():
        status = colored("PASS", "green") if result else colored("FAIL", "red")
        print(f"{status} - {name}")
    
    print(f"\nTotal Tests Run: {tests_run}")
    print(f"Tests Passed: {colored(tests_passed, 'green')}")
    print(f"Tests Failed: {colored(tests_failed, 'red')}")
    
    # Return success if all tests passed
    return all(results.values())

if __name__ == "__main__":
    try:
        # Check if the server is available
        response = requests.get(f"{BASE_URL}/health", timeout=5)
        if response.status_code != 200:
            print(colored(f"ERROR: Server at {BASE_URL} does not appear to be running", "red"))
            sys.exit(1)
    except requests.exceptions.RequestException:
        print(colored(f"ERROR: Cannot connect to server at {BASE_URL}", "red"))
        print("Make sure the server is running before executing this script.")
        sys.exit(1)
    
    success = run_all_tests()
    sys.exit(0 if success else 1)
