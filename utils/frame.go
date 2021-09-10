/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package utils

import (
	"encoding/binary"
	"errors"
	"io"
)

const frameCode uint32 = 0x656f7363 // frameCode = "eosc"
var (
	ErrorInvalidFrame = errors.New("invalid frame")
)
func ReadFrame(r io.Reader)([]byte,error)  {

	heater := make([]byte,4)
	_, err := io.ReadFull(r, heater)
	if err != nil {
		return nil, err
	}

	code:=binary.BigEndian.Uint32(heater)
	if code != frameCode{
		return nil,ErrorInvalidFrame
	}

	_, err = io.ReadFull(r, heater)
	if err != nil {
		return nil, err
	}
	size:=binary.BigEndian.Uint32(heater)

	buf :=make([]byte,size)
	_,err =io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	return buf,nil
}
func WriteFrame(w io.Writer,data []byte) error {
	size := len(data)

	err:=binary.Write(w,binary.BigEndian,frameCode)
	if err != nil{
		return err
	}
	err = binary.Write(w,binary.BigEndian,uint32(size))
	if err != nil{
		return err
	}
	_,err =w.Write(data)
	return err
}
func EncodeFrame(data []byte)[]byte  {
	size := len(data)
	buf:= make([]byte,size+8)
	binary.BigEndian.PutUint32(buf[0:4],frameCode)
	binary.BigEndian.PutUint32(buf[4:8],uint32(size))
  	copy(buf[8:],data)
	return buf
}