package repo

import (
	"github.com/sukhjit/url-shortener/model"
)

// Shortener interface for shortener model
type Shortener interface {
	Load(slug string) (string, error)
	Add(*model.Shortener) error
	Info(slug string) (*model.Shortener, error)
}
