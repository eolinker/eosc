### traffic
> 描述：io 通信控制模块
> * 管理所有的需要热重启的监听管理（端口监听）
> * 只允许master执行新增，序列化成描述信息+文件描述符列表
> * 在fork worker时传递给worker
> * worker只允许使用传入进来的端口