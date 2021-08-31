/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package traffic

import (
	"errors"
)

var (
	ErrorInvalidFiles = errors.New("invalid errors")
)

type FD struct {
	FD uintptr `json:"fd"`
	Address string `json:"address"`
}


type IController interface {
	Listener(network string,addr string)error
	export()[]*Traffic
}

type Controller struct {

}

func (c *Controller) Listener(network string, addr string) error {
	panic("implement me")
}

func (c *Controller) export() []*Traffic {
	panic("implement me")
}

func NewController() *Controller {
	return &Controller{}
}