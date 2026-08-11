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
	"context"

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
ctx := context.Background()
queryBuilder := client.Builder()
query := queryBuilder.MatchQuery("word", "λόγος")

response, err := client.Query().Match(ctx, "dictionary", query)
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
res, err := client.Index().CreateDocument(ctx, "dictionary", body)
```

## Migrating from v0.3.x to v0.4.0

Version `v0.4.0` requires a `context.Context` for every Elasticsearch operation.
The non-context methods and the `XWithContext` method names were removed. Pass
the context as the first argument to the original method name:

| v0.3.x | v0.4.0 |
| --- | --- |
| `MatchWithContext(ctx, index, query)` | `Match(ctx, index, query)` |
| `CreateWithContext(ctx, index, body)` | `Create(ctx, index, body)` |
| `CreateWithIdWithContext(ctx, index, id, body)` | `CreateWithId(ctx, index, id, body)` |
| `UpdateWithContext(ctx, index, id, body)` | `Update(ctx, index, id, body)` |
| `DeleteWithContext(ctx, index)` | `Delete(ctx, index)` |
| `CreateUserWithContext(ctx, name, request)` | `CreateUser(ctx, name, request)` |

Code using the old background-context convenience methods must now supply a
context explicitly:

```go
// v0.3.x
response, err := client.Query().Match(index, query)

// v0.4.0
response, err := client.Query().Match(ctx, index, query)
```

Use the request context in handlers and jobs whenever one is available. Only
use `context.Background()` at an application boundary that genuinely has no
parent context.

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
