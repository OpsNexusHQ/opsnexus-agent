package collector

// SystemCollector collects system metrics such as CPU, memory, and load average.
type SystemCollector struct{}

// NewSystemCollector creates a new instance of SystemCollector.
func NewSystemCollector() *SystemCollector {
	return &SystemCollector{}
}

// Name returns the identifier of this collector.
func (s *SystemCollector) Name() string {
	return "system"
}
