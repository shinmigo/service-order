package main

import "goshop/service-order/pkg/grpc/gclient"

func initService() {
	go gclient.DialGrpcService()
	//go user.Hello()
}
