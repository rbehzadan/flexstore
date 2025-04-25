from urllib.parse import urlparse, urljoin
import json

import requests


BASE_URL = "http://localhost:8080"


def is_absolute_url(url: str) -> bool:
    parsed = urlparse(url)
    return bool(parsed.scheme and parsed.netloc)


def get_full_url(url_or_path: str) -> str:
    return url_or_path if is_absolute_url(url_or_path) else urljoin(BASE_URL, url_or_path)

def handle_response(response: requests.Response, method: str, url: str) -> None:
    print(f"{method.upper()} {url}")
    print(f"Response: {response.status_code} {response.reason}")
    print("Headers:", response.headers)

    try:
        response.raise_for_status()
    except requests.HTTPError as e:
        print(f"HTTP error: {e}")
    except requests.RequestException as e:
        print(f"Request failed: {e}")

    content_type = response.headers.get("Content-Type", "").lower()
    content = response.content.strip()

    if not content:
        print("No content in response.")
        return

    # Attempt to decode JSON if content type is JSON-ish
    if "json" in content_type or content.startswith(b"{") or content.startswith(b"["):
        try:
            parsed_json = response.json()
            print(json.dumps(parsed_json, indent=2))
            return
        except json.JSONDecodeError:
            print("Invalid JSON in response.")

    # If not JSON, and content is decodable (i.e., not binary), show it as text
    try:
        text = content.decode(response.encoding or "utf-8")
        print("Raw response:")
        print(text)
    except UnicodeDecodeError:
        print("Binary content (not displayed).")

def get(path: str) -> None:
    url = get_full_url(path)
    try:
        response = requests.get(url)
        handle_response(response, "GET", url)
    except requests.RequestException as e:
        print(f"GET request error: {e}")


def post(path: str, data: dict) -> None:
    url = get_full_url(path)
    try:
        response = requests.post(url, json=data)
        handle_response(response, "POST", url)
    except requests.RequestException as e:
        print(f"POST request error: {e}")


def put(path: str, data: dict) -> None:
    url = get_full_url(path)
    try:
        response = requests.put(url, json=data)
        handle_response(response, "PUT", url)
    except requests.RequestException as e:
        print(f"PUT request error: {e}")


def delete(path: str) -> None:
    url = get_full_url(path)
    try:
        response = requests.delete(url)
        handle_response(response, "DELETE", url)
    except requests.RequestException as e:
        print(f"DELETE request error: {e}")

