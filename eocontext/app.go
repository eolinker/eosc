package eocontext

// NodeStatus 节点状态类型
type NodeStatus int

const (
	//Running 节点运行中状态
	Running NodeStatus = 1
	//Down 节点不可用状态
	Down NodeStatus = 2
	//Leave 节点离开状态
	Leave NodeStatus = 3
)

// Attrs 属性集合
type Attrs map[string]string

// IAttributes 属性接口
type IAttributes interface {
	GetAttrs() Attrs
	GetAttrByName(name string) (string, bool)
}

type EoApp interface {
	Nodes() []INode
}

// INode 节点接口
type INode interface {
	IAttributes
	ID() string
	IP() string
	Port() int
	Addr() string
	Status() NodeStatus
	Up()
	Down()
	Leave()
}
