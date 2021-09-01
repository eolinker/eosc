/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 *
 */
/*

io 通信控制模块
管理所有的需要热重启的监听管理（端口监听）， 只允许master执行新增， 序列化成描述信息+文件描述符列表，在fork worker时传递给worker，
worker只允许使用传入进来的端口

 */
package traffic

import (
	"net"
	"os"
)

type Traffic struct {
	Addr net.Addr
	File *os.File
}

type Traffics []*Traffic
