package main

import (
	"product/configs"
	"product/grpc"
	"product/storage"
)

func main() {
	configs.InitEnv()            // init env
	storage.GetStorageInstance() // init db
	grpc.Init()
}
