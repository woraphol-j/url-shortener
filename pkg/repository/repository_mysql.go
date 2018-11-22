package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	logger "github.com/sirupsen/logrus"
)

// ShortURLModel is the model
type ShortURLModel struct {
	gorm.Model
	Code string `gorm:"type:varchar(100);"`
	URL  string `gorm:"type:varchar(100);"`
}

// TableName sets table name
func (ShortURLModel) TableName() string {
	return "short_urls"
}

type mysqlRepository struct {
	db *gorm.DB
}

// NewMySQLRepository creates a data access object for managing shortend urls
func NewMySQLRepository(mysqlURL string) Repository {
	db, err := gorm.Open("mysql", mysqlURL)
	if err != nil {
		logger.Fatal("error creating database", err)
	}
	db.AutoMigrate(&ShortURLModel{})

	return &mysqlRepository{
		db: db,
	}
}

// Save saves url data in MongoDB
func (mr *mysqlRepository) Save(shortURL *ShortURL) error {
	logger.WithFields(logger.Fields{
		"url": shortURL.URL,
	}).Info("Generate short url")
	shortURLModel := &ShortURLModel{
		Code: shortURL.Code,
		URL:  shortURL.URL,
	}
	mr.db.Create(shortURLModel)
	return nil
}

// Get fetchs url by code
func (mr *mysqlRepository) Get(code string) (*ShortURL, error) {
	logger.WithFields(logger.Fields{
		"code": code,
	}).Info("Retrieve short url")
	shortURLModel := ShortURLModel{}
	mr.db.Where("code = ?", code).First(&shortURLModel)
	return &ShortURL{
		Code:     shortURLModel.Code,
		URL:      shortURLModel.URL,
		NotFound: false,
	}, nil
}

// Truncate deletes the entire data in URL collection
func (mr *mysqlRepository) Truncate() (int64, error) {
	fmt.Println("*************************Trancate.....")
	mr.db.Exec("DELETE * FROM short_urls")
	return 5, nil
}
