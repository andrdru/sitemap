package sitemap

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"
)

var IndexName = "/sitemap.xml"
var URLSetName = "/sitemap.%d.xml.gz"

type (
	marshaler interface {
		Marshal() (data []byte, err error)
	}
)

type Sitemap struct {
	URL string

	mu sync.Mutex

	Index    *Index
	Sitemaps map[string]*URLSet

	indexData    []byte
	sitemapsData map[string][]byte
}

func NewSitemap(url string) *Sitemap {
	return &Sitemap{
		mu:       sync.Mutex{},
		Index:    NewIndex(),
		Sitemaps: make(map[string]*URLSet),

		sitemapsData: make(map[string][]byte),
		URL:          url,
	}
}

func (s *Sitemap) AddURLSet(u *URLSet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	name := fmt.Sprintf(URLSetName, len(s.Index.Sitemap)+1)

	s.Sitemaps[name] = u
	s.Index.Sitemap = append(s.Index.Sitemap, IndexSitemap{
		Loc:     strings.TrimSuffix(s.URL, "/") + path.Join("/", name),
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
		s.sitemapsData[name], err = gzipXML(s.Sitemaps[name])
		if err != nil {
			return fmt.Errorf("marshal sitemap %s: %w", name, err)
		}
	}

	return nil
}

func gzipXML(m marshaler) (data []byte, err error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	var tmpData []byte
	tmpData, err = m.Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	_, err = zw.Write(tmpData)
	if err != nil {
		return nil, fmt.Errorf("gzip write: %w", err)
	}

	if err = zw.Close(); err != nil {
		return nil, fmt.Errorf("gzip close: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *Sitemap) Store() error {
	panic("not implemented")
}

// Handler handle sitemap inmemory
func (s *Sitemap) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filePath := strings.ToLower(r.URL.Path)

		if filePath == IndexName {
			_, _ = w.Write(s.indexData)
			return
		}

		data, ok := s.sitemapsData[filePath]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/gzip")
		_, _ = w.Write(data)
		return
	}
}
