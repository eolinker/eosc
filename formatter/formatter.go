package formatter

type IEntry interface {
	Read(pattern string) (string, bool)
	Children(name string) []IEntry
}

type Config map[string][]string

type IFormatterFactory interface {
	Create(cfg Config) IFormatter
}

//IFormatter format config
type IFormatter interface {
	Format(entry IEntry) []byte
}
