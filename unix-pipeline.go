package unix_pipeline

type Command struct {
	Name     Operation
	Argument string
}

type Operation string
