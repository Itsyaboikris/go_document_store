# Distributed Document Store

A distributed document store implementation in Go that supports real-time replication across multiple nodes. The system implements a hierarchical structure of Projects > Collections > Documents with automatic peer-to-peer synchronization.

## Features

- Hierarchical data organization (Projects > Collections > Documents)
- REST API for CRUD operations
- Real-time peer-to-peer replication
- Automatic retry mechanism for failed replications
- Thread-safe operations using mutex locks
- JSON document support
- Timestamp tracking for document creation and updates

### Components

- **Store**: Core data structure implementation with thread-safe CRUD operations
- **API**: REST endpoints for document operations
- **Replication**: Peer-to-peer synchronization system

## API Endpoints
```
POST /{project}/{collection}/document # Create a new document

GET /{project}/{collection}/document # Get all documents in a collection

PUT /{project}/{collection}/document/{id} # Update a document

DELETE /{project}/{collection}/document/{id} # Delete a document

POST /replicate # Internal endpoint for replication
```

## Usage Examples

### Create a Document
```bash
curl -X POST http://localhost:8080/project1/collection1/document \
  -H "Content-Type: application/json" \
  -d '{"name": "test document", "value": 123}'
```

### Update a Document
``` bash
curl -X PUT http://localhost:8080/project1/collection1/document/123 \
  -H "Content-Type: application/json" \
  -d '{"name": "updated document", "value": 456}'

```

### Delete a Document
``` bash
curl -X DELETE http://localhost:8080/project1/collection1/document/123

```

### Get all Documents
``` bash
curl http://localhost:8080/project1/collection1/document

```

### Running the Service with Docker
``` bash
docker compose up --build
```