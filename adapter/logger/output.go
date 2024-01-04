package logger

// OutputType
type OutputType struct {
	string
}

type output struct {
	CONSOLE OutputType
	JSON    OutputType
}

// Output
var Output = &output{
	CONSOLE: OutputType{"CONSOLE"},
	JSON:    OutputType{"JSON"},
}

func (output OutputType) ToEncoding() string {
	switch output {
	case Output.CONSOLE:
		return "console"
	case Output.JSON:
		return "json"
	}
	return ""
}
