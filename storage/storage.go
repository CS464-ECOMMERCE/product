package storage

import (
	"fmt"
	"log"
	"os"
	"product/configs"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	read    *gorm.DB
	write   *gorm.DB
	Product ProductInterface
	S3      S3Interface
	Order   OrderInterface
}

func (s *Storage) InitDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	var err error
	s.write, err = gorm.Open(postgres.Open(configs.POSTGRESQL_CONN_STRING_MASTER), &gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Println("status: ", err)
	}
	writeDB, err := s.write.DB()
	if err != nil {
		fmt.Println("status: ", err)
	}
	writeDB.SetMaxOpenConns(configs.POSTGRESQL_MAX_OPEN_CONNS)
	writeDB.SetMaxIdleConns(configs.POSTGRESQL_MAX_IDLE_CONNS)

	s.read, err = gorm.Open(postgres.Open(configs.POSTGRESQL_CONN_STRING_SLAVE), &gorm.Config{Logger: newLogger})
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

func (s *Storage) BeginTransaction() *gorm.DB {
	return s.write.Begin()
}

var StorageInstance *Storage
var once sync.Once

func GetStorageInstance() *Storage {
	once.Do(func() {
		StorageInstance = &Storage{}
		StorageInstance.InitDB()
		StorageInstance.Product = NewProductTable(StorageInstance.read, StorageInstance.write)
		if configs.ENVIRONMENT == "prod" {
			StorageInstance.S3 = NewS3()
		} else {
			StorageInstance.S3 = NewMinio()
		}
		StorageInstance.Order = NewOrderTable(StorageInstance.write)
	})
	return StorageInstance
}
