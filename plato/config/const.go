package config

const (
	PLATO                    = "plato"
	ODYSSEIA_PATH            = "odysseia-greek"
	DefaultSidecarService    = "localhost:50051"
	DefaultKubeConfig        = "/.kube/config"
	DefaultNamespace         = "odysseia"
	DefaultPodname           = "somepod-08595-383"
	DefaultSearchWord        = "greek"
	DefaultRoleName          = "solon"
	DefaultJobName           = "demokritos"
	DefaultServiceAddress    = "http://odysseia-greek.internal"
	DefaultCaValidity        = "3650"
	DefaultCrdName           = "perikles-mapping"
	EnvHealthCheckOverwrite  = "HEALTH_CHECK_OVERWRITE"
	EnvPodName               = "POD_NAME"
	EnvNamespace             = "NAMESPACE"
	EnvIndex                 = "ELASTIC_ACCESS"
	EnvSecondaryIndex        = "ELASTIC_SECONDARY_ACCESS"
	EnvVaultService          = "VAULT_SERVICE"
	EnvSolonService          = "SOLON_SERVICE"
	EnvPtolemaiosService     = "PTOLEMAIOS_SERVICE"
	EnvHerodotosService      = "HERODOTOS_SERVICE"
	EnvAlexandrosService     = "ALEXANDROS_SERVICE"
	EnvSokratesService       = "SOKRATES_SERVICE"
	EnvDionysiosService      = "DIONYSIOS_SERVICE"
	EnvRunOnce               = "RUN_ONCE"
	EnvTlSKey                = "TLS_ENABLED"
	EnvKey                   = "ENV"
	EnvSearchWord            = "SEARCH_WORD"
	EnvRole                  = "ELASTIC_ROLE"
	EnvRoles                 = "ELASTIC_ROLES"
	EnvIndexes               = "ELASTIC_INDEXES"
	EnvRootToken             = "VAULT_ROOT_TOKEN"
	EnvAuthMethod            = "AUTH_METHOD"
	EnvTLSEnabled            = "VAULT_TLS"
	EnvVaultRole             = "VAULT_ROLE"
	EnvKubePath              = "KUBE_PATH"
	EnvSidecarOverwrite      = "SIDECAR_OVERWRITE"
	EnvJobName               = "JOB_NAME"
	EnvCAValidity            = "CA_VALIDITY"
	EnvCrdName               = "CRD_NAME"
	EnvTLSFiles              = "TLS_FILES"
	EnvRootTlSDir            = "CERT_ROOT"
	EnvWaitTime              = "WAIT_TIME"
	EnvMetricsGathering      = "GATHER_METRICS"
	EnvMaxAge                = "MAX_AGE"
	AuthMethodKube           = "kubernetes"
	AuthMethodToken          = "token"
	baseDir                  = "base"
	configFileName           = "config.yaml"
	DefaultRoleAnnotation    = "odysseia-greek/role"
	DefaultAccessAnnotation  = "odysseia-greek/access"
	DefaultTLSFileLocation   = "/etc/certs"
	serviceAccountTokenPath  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	EnvTestOverWrite         = "TEST_OVERWRITE"
	EnvChannel               = "EUPALINOS_CHANNEL"
	EnvEupalinosService      = "EUPALINOS_SERVICE"
	EnvAggregatorAddress     = "ARISTARCHOS_SERVICE"
	DefaultAggregatorAddress = "aristarchos:50053"
	DefaultEupalinosService  = "eupalinos:50051"
	DefaultParmenidesChannel = "parmenides"
	DefaultDutchChannel      = "mouseion"
	DefaultTracingName       = "agreus"
	DefaultMetricsName       = "eumetros"
	HeaderKey                = "aischylos"
	CreatorElasticRole       = "creator"
	SeederElasticRole        = "seeder"
	HybridElasticRole        = "hybrid"
	ApiElasticRole           = "api"
	AliasElasticRole         = "alias"
	TracingElasticIndex      = "tracing"
	MetricsElasticIndex      = "metrics"
)

var serviceMapping = map[string]string{
	"SolonService": EnvSolonService,
}
