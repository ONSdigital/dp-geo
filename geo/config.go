package geo

const (
	defaultConcurrency     = 10 // limit number of go routines to not put too much on heap
	defaultMaximumSegments = 180
)

// Config object to define geo configurations
type Config struct {
	defaultConcurrencyLimit int
	defaultMaxSegments      int
}

// Default geo configuration for methods on config receiver
var DefaultConfig = &Config{
	defaultConcurrencyLimit: defaultConcurrency,
	defaultMaxSegments:      defaultMaximumSegments,
}

// New instantiates a geo object with defined values to be used by geo methods
// e.g. concurrency limit and maximum number of segments when generating a circle from a point
func New(concurrencyLimit, maximumSegments int) (geo *Config) {
	geo = DefaultConfig

	geo.defaultConcurrencyLimit = concurrencyLimit
	geo.defaultMaxSegments = maximumSegments

	return
}
