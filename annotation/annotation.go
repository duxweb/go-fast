package annotation

var Annotations = make([]*File, 0)

type File struct {
	Name        string
	Annotations []*Annotation
}

type Annotation struct {
	Name   string
	Params map[string]any
	Func   any
}

type Import struct {
	Name string
	As   string
}
