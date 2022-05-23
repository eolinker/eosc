package eosc

const (
	//ProcessMaster master进程，守护进程
	ProcessMaster = "master"
	//ProcessWorker worker进程，负责网关主流程的执行
	ProcessWorker = "worker"
	//ProcessHelper helper进程，临时进程，用于检测插件下载操作
	ProcessHelper = "helper"
	//ProcessAdmin admin进程，缓存配置信息，常驻进程
	ProcessAdmin = "admin"
)
