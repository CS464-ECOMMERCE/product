package storage

import (
	"fmt"
	"product/configs"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	read    *gorm.DB
	write   *gorm.DB
	Product ProductInterface
}

func (s *Storage) InitDB() {
	var err error
	s.write, err = gorm.Open(postgres.Open(configs.POSTGRESQL_CONN_STRING_MASTER), &gorm.Config{})
	if err != nil {
		fmt.Println("status: ", err)
	}
	writeDB, err := s.write.DB()
	if err != nil {
		fmt.Println("status: ", err)
	}
	writeDB.SetMaxOpenConns(configs.POSTGRESQL_MAX_OPEN_CONNS)
	writeDB.SetMaxIdleConns(configs.POSTGRESQL_MAX_IDLE_CONNS)

	s.read, err = gorm.Open(postgres.Open(configs.POSTGRESQL_CONN_STRING_SLAVE), &gorm.Config{})
	if err != nil {
		fmt.Println("status: ", err)
	}
	readDB, err := s.read.DB()
	if err != nil {
		fmt.Println("status: ", err)
	}
	readDB.SetMaxOpenConns(configs.POSTGRESQL_MAX_OPEN_CONNS)
	readDB.SetMaxIdleConns(configs.POSTGRESQL_MAX_IDLE_CONNS)
}

func (s *Storage) GetWriteDB() *gorm.DB {
	return s.write
}

func (s *Storage) GetReadDB() *gorm.DB {
	return s.read
}

func (s *Storage) AutoMigrate(model interface{}) {
	s.write.AutoMigrate(model)
	s.read.AutoMigrate(model)
}

var StorageInstance *Storage
var once sync.Once

func GetStorageInstance() *Storage {
	once.Do(func() {
		StorageInstance = &Storage{}
		StorageInstance.InitDB()
		StorageInstance.Product = NewProductTable(StorageInstance.read, StorageInstance.write)
	})
	return StorageInstance
}
