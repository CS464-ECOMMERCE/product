package main

import (
	"product/configs"
	_ "product/docs"
	"product/grpc"
	"product/storage"
	// "product/routes"
	// "backend/mq"
)

func main() {
	configs.InitEnv() // init env
	// mq.Init()           // init rabbitmq connection
	storage.GetStorageInstance() // init db
	grpc.Init()
}
