package mongo

//go:generate mockgen -destination=./dao_mock.go -package=mongo github.com/woraphol-j/url-shortener/pkg/mongo DAO

import (
	"context"
	"log"
	"strings"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// DAO is the data access object interface used to access
// underlying data source
type DAO interface {
	Save(shortURL *ShortURL) error
	Get(code string) (*ShortURL, error)
	Truncate() (int64, error)
}

// DAO is the data access object
type dao struct {
	urlCollection *mongo.Collection
}

// ShortURL is the model to insert and get url data from the database
type ShortURL struct {
	Code     string `bson:"code"`
	URL      string `bson:"url"`
	NotFound bool
}

// NewDAO creates a data access object for managing shortend urls
func NewDAO(mongoURL, database, collection string) DAO {
	client, err := mongo.NewClient(mongoURL)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	coll := client.Database(database).Collection(collection)
	return &dao{
		urlCollection: coll,
	}
}

// Save saves url data in MongoDB
func (dao *dao) Save(shortURL *ShortURL) error {
	_, err := dao.urlCollection.InsertOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("code", shortURL.Code),
			bson.EC.String("url", shortURL.URL),
		),
	)
	if err != nil {
		return nil
	}
	return nil
}

// Get fetchs url by code
func (dao *dao) Get(code string) (*ShortURL, error) {
	docResult := dao.urlCollection.FindOne(
		nil,
		bson.NewDocument(bson.EC.String("code", code)),
	)
	var shortURL ShortURL
	err := docResult.Decode(&shortURL)
	if err != nil {
		//TODO: Find a better way to check the error type
		if strings.Contains(err.Error(), mongo.ErrNoDocuments.Error()) {
			return &ShortURL{
				NotFound: true,
			}, nil
		}
		return nil, err
	}

	return &shortURL, nil
}

// Truncate deletes the entire data in URL collection
func (dao *dao) Truncate() (int64, error) {
	result, err := dao.urlCollection.DeleteMany(
		context.Background(),
		bson.NewDocument(),
	)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
