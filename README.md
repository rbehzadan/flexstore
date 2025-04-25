# FlexStore

A RESTful API server in Go with SQLite for schemaless data storage.

## Features

- Collection-based schemaless data storage
- Full CRUD operations
- Bulk insert via JSON files or POST requests
- SQLite database backend
- HTTP Basic Authentication for protected endpoints

## Getting Started

### Prerequisites

- Go 1.18 or newer
- SQLite3

### Installation

1. Clone this repository
2. Navigate to the project directory
3. Build the application

```bash
make build
```

### Usage

Run the server:

```bash
make run
```

Or with custom parameters:

```bash
./build/schemaless-api -addr :9000
```

Available command line options:

- `-help`: Show help information
- `-version`: Show version information
- `-addr`: Set HTTP service address (default: ":8080")
- `-auth`: Enable HTTP Basic Authentication for protected endpoints
- `-username`: Set username for authentication (default: "admin")
- `-password`: Set password for authentication (default: "password")

### API Endpoints

- `GET /health`: Health check endpoint that returns status, version, and uptime

More endpoints will be implemented in future commits.

## Development

### Project Structure

```
/
├── main.go                  # Main entry point with VERSION embedding
├── VERSION                  # Version file
├── cmd/
│   └── server/              # Server package with Run function
├── internal/
│   ├── api/
│   │   ├── types.go         # Common API types
│   │   ├── handlers/        # API handlers
│   │   └── middleware/      # HTTP middleware
│   ├── db/                  # Database operations
│   ├── models/              # Data models
│   └── service/             # Business logic
└── pkg/
    ├── config/              # Configuration
    └── httputils/           # HTTP utilities
```

### Building

```bash
make build
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage
# Coverage reports will be generated in the coverage/ directory
```

### API Endpoints

#### System Endpoints

- `GET /health`: Health check endpoint that returns status, version, and uptime

#### Collection Endpoints

- `GET /api/collections`: List all collections
- `POST /api/collections`: Create a new collection
- `GET /api/collections/{name}`: Get a specific collection
- `DELETE /api/collections/{name}`: Delete a collection

#### Document Endpoints

- `GET /api/collections/{name}/documents`: List documents in a collection
  - Query parameters:
    - `limit`: Maximum number of documents to return (default: 100)
    - `offset`: Number of documents to skip (default: 0)
- `POST /api/collections/{name}/documents`: Create a new document in a collection
- `GET /api/collections/{name}/documents/{id}`: Get a specific document
- `PUT /api/collections/{name}/documents/{id}`: Update a document
- `DELETE /api/collections/{name}/documents/{id}`: Delete a document

#### Bulk Operations

- `POST /api/collections/{name}/bulk`: Bulk insert documents from a JSON array
- `POST /api/upload/{name}`: Upload and process a JSON file for bulk insertion

### Data Format

All data is stored and returned as JSON. Documents are schemaless and can contain any valid JSON structure.

#### Collection Creation Example

```json
POST /api/collections
{
  "name": "users"
}
```

#### Document Creation Example

```json
POST /api/collections/users/documents
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30,
  "address": {
    "city": "New York",
    "country": "USA"
  }
}
```
