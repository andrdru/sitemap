package sitemap

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var IndexName = "/sitemap.xml"

type Sitemap struct {
	mu sync.Mutex

	Index    *Index
	Sitemaps map[string]*URLSet

	indexData    []byte
	sitemapsData map[string][]byte
}

func NewSitemap() *Sitemap {
	return &Sitemap{
		mu:       sync.Mutex{},
		Index:    NewIndex(),
		Sitemaps: make(map[string]*URLSet),

		sitemapsData: make(map[string][]byte),
	}
}

func (s *Sitemap) AddURLSet(name string, u *URLSet) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Sitemaps[name] = u
	// s.Index.Sitemap // todo add to index
}

func (s *Sitemap) Flush() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.indexData, err = s.Index.Marshal()
	if err != nil {
		return fmt.Errorf("marshal index: %w", err)
	}

	for name := range s.Sitemaps {
		s.sitemapsData[name], err = s.Sitemaps[name].Marshal()
		if err != nil {
			return fmt.Errorf("marshal sitemap %s: %w", name, err)
		}
	}

	return nil
}

func (s *Sitemap) Store() error {
	panic("not implemented")
}

func (s *Sitemap) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.ToLower(r.URL.Path)

		if path == IndexName {
			_, _ = w.Write(s.indexData)
			return
		}

		data, ok := s.sitemapsData[path]
		if ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, _ = w.Write(data)
		return
	}
}
