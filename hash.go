package eosc

type ICustomerVar interface {
	Get(key string, field string) (string, bool)
	Exists(key string, field string) bool
}
