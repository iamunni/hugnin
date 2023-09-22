package writer

type Writer interface {
	Init(string) error
	Write(value string, tags []string) error
}
