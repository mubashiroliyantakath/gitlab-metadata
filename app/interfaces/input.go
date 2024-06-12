package interfaces

type Input interface {
	Generate() string
	IsEnabled() bool

	Enable()
	Disable()

	SetPrefix(string)
	SetSuffix(string)

	SetEvent(string)

	GitlabFilterPassed() bool
	String() string
}
