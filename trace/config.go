package trace

// TraceName represents the tracing name.
const TraceName = "mcdull"

// A Config is an opentelemetry config.
type Config struct {
	Name     string  `yaml:"name" json:",optional"`
	Endpoint string  `yaml:"endpoint" json:",optional"`
	Sampler  float64 `yaml:"sampler" json:",default=1.0"`
	Batcher  string  `yaml:"batcher" json:",default=jaeger,options=jaeger|zipkin|otlpgrpc|otlphttp|file"`
	// OtlpHeaders represents the headers for OTLP gRPC or HTTP transport.
	// For example:
	//  uptrace-dsn: 'http://project2_secret_token@localhost:14317/2'
	OtlpHeaders map[string]string `yaml:"otlp_headers" json:",optional"`
	// OtlpHttpPath represents the path for OTLP HTTP transport.
	// For example
	// /v1/traces
	OtlpHttpPath string `yaml:"otlp_http_path" json:",optional"`
	// Disabled indicates whether StartAgent starts the agent.
	Disabled bool `yaml:"disabled" json:",optional"`
}
