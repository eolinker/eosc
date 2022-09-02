### traffic
> 描述：io 通信控制模块
> * 管理所有的需要热重启的监听管理（端口监听）
> * 只允许master执行新增，序列化成描述信息+文件描述符列表
> * 在fork worker时传递给worker
> * worker只允许使用传入进来的端口

## 用法

```golang
// master
// 首次启动
tf, err := traffic.ReadController(os.stdint, &net.TCPAddr{
IP:   "0.0.0.0",
Port: 8080,
})
// 导出
traffics, files := tf.Export(3)
trafficsData, err := proto.Marshal(&traffic.PbTraffics{Traffic: traffics})
// 新进程
cmd := exec.Command(path)
cmd.Stdin = bytes.NewReader(data)
cmd.SysProcAttr = &syscall.SysProcAttr{
Setsid: true,
}
cmd.ExtraFiles = files


// 热重启，新进程，会在复用的前提下监听新的tcpAddr列表，并关闭不再使用的listener
traffic, err := traffic.ReadController(os.stdint, &net.TCPAddr{
IP:   "0.0.0.0",
Port: 8080,
},
&net.TCPAddr{ // 新增端口
IP:   "0.0.0.0",
Port: 9090,
})

//worker 中
trafficConfigs = read from os.Stdin
tf:=traffic.NewTraffic(tfConf)

// 启用服务

l := tf.ListenTcp(port, traffic.Http1)
if l == nil {
    panic(fmt.SprintF("port:%s is not listener",port))
}

err:=http.Serve(l,handler)


```