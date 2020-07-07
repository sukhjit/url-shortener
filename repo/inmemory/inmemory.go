package inmemory

import (
	"github.com/sukhjit/url-shortener/model"
	"github.com/sukhjit/url-shortener/repo"
)

type svc struct {
	urlList []*model.Shortener
}

func loadUrls() []*model.Shortener {
	list := [][]string{
		{"go", "https://www.google.com"},
		{"yah", "https://www.yahoo.com"},
		{"hot", "https://www.hotmail.com"},
	}

	urlList := []*model.Shortener{}

	for _, row := range list {
		urlList = append(urlList, &model.Shortener{
			Slug: row[0],
			URL:  row[1],
		})
	}

	return urlList
}

// NewShortener func
func NewShortener() repo.Shortener {
	return &svc{
		urlList: loadUrls(),
	}
}

func (s *svc) Add(sh *model.Shortener) error {
	s.urlList = append(s.urlList, sh)

	return nil
}

func (s *svc) Info(slug string) (*model.Shortener, error) {
	for _, row := range s.urlList {
		if row.Slug == slug {
			return row, nil
		}
	}

	return nil, nil
}

func (s *svc) Load(slug string) (string, error) {
	for _, row := range s.urlList {
		if row.Slug != slug {
			continue
		}

		row.Visits++

		return row.URL, nil
	}

	return "", nil
}
