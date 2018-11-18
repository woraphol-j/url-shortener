package repository

// ShortURL is the model to insert and get url data from the database
type ShortURL struct {
	Code     string `bson:"code"`
	URL      string `bson:"url"`
	NotFound bool
}
