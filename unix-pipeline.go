package unixpipeline

type Command struct {
	Name     Operation
	Argument string
}

type Operation string
