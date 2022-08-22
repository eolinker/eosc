# cmux match

基于 github.com/soheilhy/cmux 实现的动态添加协议的 match管理器

## 用法

```go
package main

import cmuxMatch "github.com/eolinker/eosc/traffic/cmux-match"
	match := cmuxMatch.NewMatch(listener)
	h1l:=match.Match(cmuxMatch.Http1)
	wslmatch.Match(cmuxMatch.Websocket)
	wslmatch.Match(cmuxMatch.GRPC)
```