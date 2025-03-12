package models

type Product struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Price       uint   `json:"price" gorm:"not null"`
	Picture     string `json:"picture"`
	Stock       uint   `json:"stock" gorm:"not null"`
	CreatedAt   uint   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   uint   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *uint  `json:"deleted_at" gorm:"autoDeleteTime"`
}
