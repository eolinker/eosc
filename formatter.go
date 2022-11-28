package eosc

type IEntry interface {
	Read(pattern string) string
	Children(child string) []IEntry
}

type FormatterConfig map[string][]string

type IFormatterFactory interface {
	Create(cfg FormatterConfig) (IFormatter, error)
}

// IFormatter format config
type IFormatter interface {
	Format(entry IEntry) []byte
}
