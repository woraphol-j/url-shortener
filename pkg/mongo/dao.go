package mongo

//go:generate mockgen -destination=../../mocks/mock_dao.go -package=mocks github.com/woraphol-j/url-shortener/pkg/mongo DAO

import (
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// DAO is the data access object
type DAO interface {
	Save(shortURL *ShortURL) error
	Get(code string) (*ShortURL, error)
}

// DAO is the data access object
type dao struct {
	mongoClient   *mongo.Client
	urlCollection *mongo.Collection
}

// ShortURL is the model to insert and get data from the database
type ShortURL struct {
	Code string `bson:"code"`
	URL  string `bson:"url"`
}

// NewDAO creates a data access object for managing shortend urls
func NewDAO(mongoURL, database, collection string) DAO {
	// client, err := mongo.NewClient("mongodb://mongo:27017")
	client, err := mongo.NewClient(mongoURL)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// coll := client.Database("url-shortener").Collection("urls")
	coll := client.Database(database).Collection(collection)
	return &dao{
		mongoClient:   client,
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

// Get fetch url by code
func (dao *dao) Get(code string) (*ShortURL, error) {
	result := dao.urlCollection.FindOne(
		context.Background(),
		bson.NewDocument(bson.EC.String("code", code)),
	)
	var shortURL ShortURL
	result.Decode(&shortURL)

	return &shortURL, nil
}
