# Distributed Document Store

A distributed document store implementation in Go that supports real-time replication across multiple nodes. The system implements a hierarchical structure of Projects > Collections > Documents with automatic peer-to-peer synchronization and MongoDB-style querying.

## Features

- Hierarchical data organization (Projects > Collections > Documents)
- REST API for CRUD operations
- MongoDB-style query operations
- Real-time peer-to-peer replication
- Automatic retry mechanism for failed replications
- Thread-safe operations using mutex locks
- JSON document support
- Timestamp tracking for document creation and updates

### Components

- **Store**: Core data structure implementation with thread-safe CRUD operations
- **API**: REST endpoints for document operations
- **Query**: MongoDB-style query system with support for complex queries
- **Replication**: Peer-to-peer synchronization system

## API Endpoints
```
POST /{project}/{collection}/document # Create a new document

GET /{project}/{collection}/document # Get all documents in a collection

PUT /{project}/{collection}/document/{id} # Update a document

DELETE /{project}/{collection}/document/{id} # Delete a document

POST /{project}/{collection}/query # Query documents

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

### Query Documents
``` bash
# Find documents where age > 25
curl -X POST http://localhost:8080/project1/collection1/query \
  -H "Content-Type: application/json" \
  -d '{
    "age": {"$gt": 25}
  }'

# Find active users who know golang
curl -X POST http://localhost:8080/project1/collection1/query \
  -H "Content-Type: application/json" \
  -d '{
    "$and": [
      {"status": "active"},
      {"tags": {"$all": ["golang"]}}
    ]
  }'
```


## Query Operators
### Comparison Operators
$eq: Matches values that are equal to a specified value

$ne: Matches values that are not equal to a specified value

$gt: Matches values that are greater than a specified value

$gte: Matches values that are greater than or equal to a specified value

$lt: Matches values that are less than a specified value

$lte: Matches values that are less than or equal to a specified value

$in: Matches any of the values specified in an array

$nin: Matches none of the values specified in an array

### Logical Operators
$and: Joins query clauses with a logical AND

$or: Joins query clauses with a logical OR

$not: Inverts the effect of a query expression

$nor: Joins query clauses with a logical NOR

### Element Operators
$exists: Matches documents that have the specified field

$type: Selects documents if a field is of the specified type

### Evaluation Operators
$regex: Selects documents where values match a specified regular expression

$mod: Performs modulo operation on the value of a field

###  Array Operators
$all: Matches arrays that contain all elements specified in the query

$size: Matches arrays with the specified size

$elemMatch: Matches documents that contain an array field with at least one element that matches the specified query criteria

### Running the Service with Docker
``` bash
docker compose up --build
```