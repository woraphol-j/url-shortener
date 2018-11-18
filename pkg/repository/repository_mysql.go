package repository

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// ShortURLModel is the model
type ShortURLModel struct {
	gorm.Model
	Code string `gorm:"type:varchar(100);"`
	URL  string `gorm:"type:varchar(100);"`
}

type mysqlRepository struct {
	db *gorm.DB
}

// NewMySQLRepository creates a data access object for managing shortend urls
func NewMySQLRepository(mysqlURL string) Repository {
	connStr := os.Getenv("MYSQL_CONNECTION_STRING")
	db, err := gorm.Open("mysql", connStr)
	db.AutoMigrate(&ShortURLModel{})

	if err != nil {
		log.Fatal(err)
	}

	return &mysqlRepository{
		db: db,
	}
}

// Save saves url data in MongoDB
func (mr *mysqlRepository) Save(shortURL *ShortURL) error {
	shortURLModel := &ShortURLModel{
		Code: shortURL.Code,
		URL:  shortURL.URL,
	}
	mr.db.Create(shortURLModel)
	return nil
}

// Get fetchs url by code
func (mr *mysqlRepository) Get(code string) (*ShortURL, error) {
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
	mr.db.Delete(&ShortURLModel{})
	return 5, nil
}
