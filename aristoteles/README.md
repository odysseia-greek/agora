# Aristoteles | Ἀριστοτέλης

`aristoteles` is the Elasticsearch interface layer for the **odysseia-greek** project. It provides a structured and idiomatic Go way to interact with Elasticsearch, abstracting away the complexities of raw HTTP requests and providing a clean set of interfaces for querying, indexing, and management.

## Installation

```bash
go get github.com/odysseia-greek/agora/aristoteles
```

## Features

- **Modular Interfaces**: Specialized interfaces for `Query`, `Index`, `Document`, `Access`, `Health`, and `Policy`.
- **Structured Models**: Type-safe request and response models for common Elasticsearch operations.
- **Easy Initialization**: Support for both standard and TLS-based connections.
- **Mocking Support**: Built-in mock client for unit testing downstream services without a live Elasticsearch instance.

## Usage

### Initialization

To use `aristoteles`, first define a configuration and create a new client:

```go
import (
    "github.com/odysseia-greek/agora/aristoteles"
    "github.com/odysseia-greek/agora/aristoteles/models"
)

func main() {
    config := models.Config{
        Service:     "http://localhost:9200",
        Username:    "elastic",
        Password:    "password",
        ElasticCERT: "", // Path to cert if using TLS
    }

    client, err := aristoteles.NewClient(config)
    if err != nil {
        // handle error
    }
}
```

### Performing a Query

The `Query()` sub-interface handles searches:

```go
queryBuilder := client.Builder()
query := queryBuilder.MatchQuery("word", "λόγος")

response, err := client.Query().Match("dictionary", query)
if err != nil {
    // handle error
}

for _, hit := range response.Hits.Hits {
    fmt.Printf("Found: %v\n", hit.Source)
}
```

### Managing Documents

Use the `Document()` or `Index()` sub-interfaces for CRUD operations:

```go
body := []byte(`{"word": "λόγος", "meaning": "word, reason"}`)
res, err := client.Index().CreateDocument("dictionary", body)
```

## Sub-Interfaces

- **`Query()`**: Search operations (Match, Count, Scroll, Aggregate).
- **`Index()`**: Index-level management (Create, Delete, Exists) and document creation.
- **`Document()`**: Document-level operations (Update, Get).
- **`Access()`**: Security management (Create User/Role, List Users).
- **`Health()`**: Cluster health checks and information.
- **`Builder()`**: Helper methods for constructing Elasticsearch query DSL maps.

## Testing

### Run standard tests

To run all tests in the package:

```bash
go test -v ./...
```

### Run tests with coverage

To see how much of the code is covered by tests:

```bash
go test -cover ./...
```

To generate a detailed HTML coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Using the Mock Client

Downstream services can use `NewMockClient` for testing:

```go
// Load your fixture JSON files
fixtures := []string{"./fixtures/search_response.json"}
mockClient, err := aristoteles.NewMockClient(fixtures, 200)

// Use mockClient as if it were a real aristoteles.Client
```