#!/bin/bash

protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto

mv message.pb.go ../
mv traffic.pb.go ../traffic/
mv listen.pb.go ../config/
