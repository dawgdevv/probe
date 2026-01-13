package models

type TestSuite struct {
	Env   map[string]string `yaml:"env"`
	Tests []TestCase
}

type TestCase struct {
	Name    string  `yaml:"name"`
	Request Request `yaml:"request"`
	Expect  Expect  `yaml:"expect"`
}

type Request struct {
	Method  string                 `yaml:"method"`
	Path    string                 `yaml:"path"`
	Headers map[string]string      `yaml:"headers"`
	Body    map[string]interface{} `yaml:"body"`
}
type Expect struct {
	Status int                    `yaml:"status"`
	JSON   map[string]interface{} `yaml:"json"`
}
