# worker 进程管理优化

## 资源

1. 监听端口的句柄
2. 扩展配置
3. profession配置
4. woker(业务)配置

## 端口

端口配置对于master为固定，所以，worker进程不需要因为端口配置变更而结束

## profession &worker

profession与worker的配置可以可以在内部动态变更完成，worker进程提供grpc方法用于变更

## 扩展

由于golang的扩展加载无法卸载，且相代码pkg无法加载不同版本，所以，需要在扩展配置发生如下变更时重载

1. 卸载扩展
2. 变更扩展版本

由于卸载扩展可以通过不是用来规避，并且通过接口回收，所以只有在扩展版本变更时处罚强制重启

## 入参

### 传参方法

通过 stdin进行传递protobuf message，message 定义如下

```protobuf
message WorkerArgs {
	map<string,bytes> args=1;
}
```

> 注：message内容被 utils.WriteFrame 写入到std中， 需要用 utils.ReadFrame 读出来

### 子参数

上面 WorkerArgs中的args中，每个key的内容为一组独立结构， 内容为[]bytes



