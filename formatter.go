package eosc

import "strconv"

type IEntry interface {
	Read(pattern string) interface{}
	ReadLabel(pattern string) string
	Children(child string) []IEntry
}

func ReadStringFromEntry(entry IEntry, key string) string {
	value := entry.Read(key)
	return String(value)
}
func String(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	}
	return ""
}

type IMetricEntry interface {
	Read(pattern string) string
	GetFloat(pattern string) (float64, bool)
	Children(child string) []IMetricEntry
}

type FormatterConfig map[string][]string

type IFormatterFactory interface {
	Create(cfg FormatterConfig, extendCfg ...interface{}) (IFormatter, error)
}

// IFormatter format config
type IFormatter interface {
	Format(entry IEntry) []byte
}
