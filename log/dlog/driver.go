package dlog

import "github.com/eolinker/eosc/log/config"

type ConfigDriver interface {
	config.ConfigDriver
	Title() string
	ConfigFields(ignoreNames ...string) []Field
}
type FullFieldsDriver struct {
	fields []Field
}

func NewFullFieldsDriver(fields []Field) *FullFieldsDriver {
	return &FullFieldsDriver{fields: fields}
}

func (d *FullFieldsDriver) ConfigFields(ignoreNames ...string) []Field {

	if len(ignoreNames) == 0 {
		return d.fields
	}

	ignores := make(map[string]string)
	for _, n := range ignoreNames {
		ignores[n] = n
	}
	fields := make([]Field, 0, len(d.fields))
	for _, f := range d.fields {
		if _, has := ignores[f.Name]; !has {
			fields = append(fields, f)
		}
	}
	return fields
}
