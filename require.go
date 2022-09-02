package eosc

type IRequires interface {
	Set(id string, requires []string)
	Del(id string)
	RequireByCount(requireId string) int
	Requires(id string) []string
	RequireBy(requireId string) []string
}
