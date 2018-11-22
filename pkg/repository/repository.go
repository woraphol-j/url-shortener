package repository

//go:generate mockgen -destination=./repository_mock.go -package=repository github.com/woraphol-j/url-shortener/pkg/repository Repository

// Repository is the data access object interface used to access
// underlying data source
type Repository interface {
	Save(shortURL *ShortURL) error
	Get(code string) (*ShortURL, error)
	Truncate() (int64, error)
}
