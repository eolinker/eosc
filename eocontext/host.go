package eocontext

type PassHostMod int

const (
	PassHost PassHostMod = iota
	NodeHost
	ReWriteHost
)

type UpstreamHostHandler interface {
	PassHost() (PassHostMod, string)
}
