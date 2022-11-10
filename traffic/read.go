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
	"fmt"
	"io"
	"net"
	"os"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/utils"
	"google.golang.org/protobuf/proto"
)

func readTraffic(r io.Reader) ([]*PbTraffic, error) {

	frame, err := utils.ReadFrame(r)
	if err != nil {
		return nil, err
	}

	pts := new(PbTraffics)
	err = proto.Unmarshal(frame, pts)
	if err != nil {
		return nil, err
	}

	return pts.Traffic, nil
}

func toListeners(tfConf []*PbTraffic) (map[string]*net.TCPListener, error) {

	tfs := make(map[string]*net.TCPListener)
	for _, pt := range tfConf {
		name := fmt.Sprintf("%s:/%s", pt.Network, pt.Addr)

		log.DebugF("read traffic:%s=%d", name, pt.GetFD())
		switch pt.Network {

		case "tcp", "tcp4", "tcp6":

			f := os.NewFile(uintptr(pt.GetFD()), name)
			l, err := net.FileListener(f)
			if err != nil {
				log.Warn("error to read port-reqiure:", err)
				return nil, err
			}

			f.Close()
			tfs[pt.Addr] = l.(*net.TCPListener)
		}
	}

	return tfs, nil
}
