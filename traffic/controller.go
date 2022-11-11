/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package traffic

import (
	"github.com/eolinker/eosc/log"
	"io"
	"os"
)

func Export(data *TrafficData, startIndex int) ([]*PbTraffic, []*os.File) {
	log.Debug("traffic controller: Export: begin ", startIndex)
	ms := data.All()
	pts := make([]*PbTraffic, 0, len(ms))
	files := make([]*os.File, 0, len(ms))
	i := 0
	for addr, ln := range ms {

		file, err := ln.File()
		if err != nil {
			continue
		}
		pt := &PbTraffic{
			FD:      uint64(i + startIndex),
			Addr:    addr,
			Network: ln.Addr().Network(),
		}
		pts = append(pts, pt)
		files = append(files, file)
		i++

	}
	log.Debug("traffic controller: Export: size ", len(files))

	return pts, files
}

func ReadTraffic(r io.Reader, addrs ...string) (*TrafficData, error) {
	var tf *TrafficData
	if r != nil {
		traffics, err := readTraffic(r)
		if err != nil {
			return nil, err
		}
		listeners := toListeners(traffics)
		log.Debug("read listeners: ", len(listeners))
		tf = NewTrafficData(listeners)

	} else {
		tf = NewTrafficData(nil)
	}
	return tf.replace(addrs)

}
