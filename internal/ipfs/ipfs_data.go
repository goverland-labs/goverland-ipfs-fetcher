package ipfs

import (
	"errors"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Data struct {
	IpfsID    string `gorm:"primary_key"`
	CreatedAt time.Time
	Type      string
	Data      datatypes.JSON
}

func (d Data) TableName() string {
	return "ipfs_data"
}

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(data Data) error {
	return r.db.Create(&data).Error
}

func (r *Repo) GetByID(id string) (*Data, error) {
	data := Data{IpfsID: id}
	err := r.db.Take(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &data, nil
}
