package main

import (
	grpcserver "goshop/service-order/pkg/grpc/server"
)

func initService() {
	go grpcserver.Run()
	//go user.Hello()
}
