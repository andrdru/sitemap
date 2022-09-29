package sitemap

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var IndexName = "/sitemap.xml"
var URLSetName = "/sitemap.%d.xml"

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

func (s *Sitemap) AddURLSet(u *URLSet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	name := fmt.Sprintf(URLSetName, len(s.Index.Sitemap)+1)

	s.Sitemaps[name] = u
	s.Index.Sitemap = append(s.Index.Sitemap, IndexSitemap{
		Loc:     name,
		Lastmod: time.Now().Format("2006-01-02"),
	})

	return s.flush()
}

// Flush marshal sitemap data to files
func (s *Sitemap) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.flush()
}

func (s *Sitemap) flush() (err error) {
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

// Handler handle sitemap inmemory
func (s *Sitemap) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.ToLower(r.URL.Path)

		if path == IndexName {
			_, _ = w.Write(s.indexData)
			return
		}

		data, ok := s.sitemapsData[path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, _ = w.Write(data)
		return
	}
}
