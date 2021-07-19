package config

type ConfigDriver interface {
	Name() string
	// 解析成驱动可用的配置
	Decode(v string) (interface{}, error)
	// 格式化成可用表达当前配置的用于可读性的配置结构
	Format(v string) (interface{}, error)
	// 格式化成文本用于保存
	Encode(v interface{}) (string, error)
}
