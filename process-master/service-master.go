/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package process_master

import (
	"context"

	"github.com/eolinker/eosc/service"
)

var _ service.MasterServer = (*MasterServiceServer)(nil)

type MasterServiceServer struct {
	service.UnimplementedMasterServer
}

func NewMasterServiceServer() *MasterServiceServer {
	return &MasterServiceServer{}
}

func (m *MasterServiceServer) Hello(ctx context.Context, request *service.HelloRequest) (*service.HelloResponse, error) {
	return &service.HelloResponse{
		Name: request.GetName(),
	}, nil

}

func (m *MasterServiceServer) Error(ctx context.Context, request *service.ErrorRequest) (*service.ErrorResponse, error) {
	return nil, nil
}
