package formatter

type IEntry interface {
	Read(pattern string) string
	Children(child string) []IEntry
}

type Config map[string][]string

type IFormatterFactory interface {
	Create(cfg Config) (IFormatter, error)
}

//IFormatter format config
type IFormatter interface {
	Format(entry IEntry) []byte
}