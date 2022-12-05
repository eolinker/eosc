# 架构

![image](https://user-images.githubusercontent.com/25589530/129505228-121923df-e4d0-4fa6-b216-4e815a5b8dbb.png)

# 抽象概念

## profession

profession:职业,定义抽象分类

1. profession定义名是唯一的,在框架中等于常量,不区分大小写, 只能用 字母、数字、下划线
    * 例如: upstream,service,router,plugin
1. profession 定义列表字段,以及列表默认值
1. 所有的配置项目都是profession实例,实例必须具有如下属性
    * id: uuid,全局唯一
    * name: name, 同profession内唯一
    * driver: 实现该实例的驱动名
1. profession 实例需要实现销毁方法,框架会检查依赖关系并中断销毁

目前已知可能有的 profession 定义有

* upstream
* service
* router
* service discovery
* auth

### driver

driver:驱动,定义一个profession并实现能力

1. driver需要定义一个render 给 admin ui 处理界面
1. driver 实现检查 profession 的属性
1. driver 实现通过 profession 的属性实例话运行的实例并完成运行状态的
1. profession 实例的能力和属性由driver来定义
1. driver 声明一个能力清单, 在实例依赖另一个实例时,通过查询对方能力来决定是否可用,并在执行时使用该能力
