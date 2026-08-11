# Diogenes

Diogenes is the shared HashiCorp Vault client for Odysseia. It provides Vault
authentication, token, KV secret, policy, initialization, unseal, and Raft
operations.

## Migrating from v0.1.x to v0.2.0

Version `v0.2.0` requires a `context.Context` for every operation that can make
a Vault request. Pass the context as the first argument:

```go
// v0.1.x
secret, err := client.GetSecret(name)

// v0.2.0
secret, err := client.GetSecret(ctx, name)
```

The client constructors that can authenticate or perform a health check also
require context:

```go
client, err := diogenes.CreateVaultClient(ctx, true)

kubernetesClient, err := diogenes.CreateVaultClientKubernetes(
	ctx,
	vaultAddress,
	vaultRole,
	serviceAccountToken,
	tlsConfig,
)
```

Common migrations:

| v0.1.x | v0.2.0 |
| --- | --- |
| `Health()` | `Health(ctx)` |
| `GetSecret(name)` | `GetSecret(ctx, name)` |
| `CreateNewSecret(name, payload)` | `CreateNewSecret(ctx, name, payload)` |
| `ListSecrets()` | `ListSecrets(ctx)` |
| `CreateOneTimeToken(policies)` | `CreateOneTimeToken(ctx, policies)` |
| `ReadPolicy(name)` | `ReadPolicy(ctx, name)` |
| `WritePolicy(name, policy)` | `WritePolicy(ctx, name, policy)` |
| `Status()` | `Status(ctx)` |
| `Unseal(keys)` | `Unseal(ctx, keys)` |
| `Initialize(shares, threshold)` | `Initialize(ctx, shares, threshold)` |
| `RaftJoin(address, cert, key, ca)` | `RaftJoin(ctx, address, cert, key, ca)` |

Use an incoming request or job context when available. At process boundaries,
derive a cancellable or timeout context instead of using an unbounded
`context.Background()` for Vault operations.

## Kubernetes authentication token

By default, Kubernetes authentication reads the ServiceAccount token from
`/var/run/secrets/kubernetes.io/serviceaccount/token`. Set
`VAULT_KUBERNETES_TOKEN_PATH` when Vault uses a dedicated audience and its
projected ServiceAccount token is mounted elsewhere:

```yaml
env:
  - name: VAULT_KUBERNETES_TOKEN_PATH
    value: /var/run/secrets/vault/token
```

## Testing

The embedded mock fixtures remain available without a live Vault server:

```go
client, err := diogenes.CreateMockVaultClient(
	[]string{"retrieveSecret"},
	http.StatusOK,
)
```

The returned mock implements the same context-required `Client` interface.
