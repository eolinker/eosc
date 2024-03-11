package eosc

import "strconv"

type IEntry interface {
	Read(pattern string) interface{}
	Children(child string) []IEntry
}

func ReadStringFromEntry(entry IEntry, key string) string {
	var data string
	value := entry.Read(key)
	switch v := value.(type) {
	case string:
		data = v
	case []byte:
		data = string(v)
	case int:
		data = strconv.Itoa(v)
	case int64:
		data = strconv.FormatInt(v, 10)
	case float32:
		data = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		data = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		data = strconv.FormatBool(v)
	}
	return data
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
