# FlexStore API Test Script

This Python script thoroughly tests all endpoints of the FlexStore API Server.

## Features

- Tests all API endpoints
- Performs CRUD operations on collections and documents
- Tests bulk operations and file uploads
- Provides detailed output with color coding
- Reports test statistics

## Requirements

- Python 3.6+
- `requests` library
- `termcolor` library

## Installation

Install the required dependencies:

```bash
pip install -r requirements.txt
```

## Usage

Make sure the FlexStore API Server is running before executing the test script.

```bash
python test_flexstore_api.py
```

For verbose output that shows all responses:

```bash
python test_flexstore_api.py --verbose
```

## What the Script Tests

1. **Health Endpoint**
   - Tests that `/health` returns status, version, and uptime

2. **Collection Operations**
   - Creating a collection
   - Listing all collections
   - Getting a specific collection

3. **Document Operations**
   - Creating a document
   - Getting a document
   - Updating a document
   - Listing documents
   - Pagination of document lists

4. **Bulk Operations**
   - Creating multiple documents at once
   - Verifying all documents are created

5. **File Upload**
   - Creating a temporary JSON file
   - Uploading the file
   - Verifying documents are created from the file

6. **Document Deletion**
   - Deleting a document
   - Verifying it no longer exists

7. **Collection Deletion**
   - Deleting a collection with its documents
   - Verifying all are removed

## Output

The script provides detailed color-coded output for each test:

- Green ✓ PASS for successful tests
- Red ✗ FAIL for failed tests
- Detailed error information when tests fail
- HTTP status codes and JSON responses
- Summary statistics at the end

The exit code will be 0 if all tests pass, 1 otherwise.
