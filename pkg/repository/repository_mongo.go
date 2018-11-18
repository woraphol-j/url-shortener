package repository

import (
	"context"
	"log"
	"strings"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type mongoRepository struct {
	urlCollection *mongo.Collection
}

// NewMongoRepository creates a data access object for managing shortend urls
func NewMongoRepository(mongoURL, database, collection string) Repository {
	client, err := mongo.NewClient(mongoURL)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	coll := client.Database(database).Collection(collection)
	return &mongoRepository{
		urlCollection: coll,
	}
}

// Save saves url data in MongoDB
func (dao *mongoRepository) Save(shortURL *ShortURL) error {
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
func (dao *mongoRepository) Get(code string) (*ShortURL, error) {
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
func (dao *mongoRepository) Truncate() (int64, error) {
	result, err := dao.urlCollection.DeleteMany(
		context.Background(),
		bson.NewDocument(),
	)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
