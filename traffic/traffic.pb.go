package traffic

type PbTraffic struct {
	FD      uint64 `json:"FD,omitempty"`
	Addr    string `json:"Addr,omitempty"`
	Network string `json:"Network,omitempty"`
}

type PbTraffics struct {
	Traffic []*PbTraffic `json:"traffic,omitempty"`
}
