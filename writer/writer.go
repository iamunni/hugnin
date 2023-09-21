package writer

type Writer interface {
	Init(string) error
	Write(value string, tag string) error
}
