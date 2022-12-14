package sitemap

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

const (
	xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

type IndexSitemap struct {
	Loc     string `xml:"loc"`
	Lastmod string `xml:"lastmod,omitempty"`
}

type Index struct {
	XMLName xml.Name       `xml:"sitemapindex"`
	Xmlns   string         `xml:"xmlns,attr"`
	Sitemap []IndexSitemap `xml:"sitemap"`
}

func NewIndex() *Index {
	return &Index{
		Xmlns: xmlns,
	}
}

func (x *Index) Marshal() (data []byte, err error) {
	buf := bytes.NewBuffer(data)

	_, err = buf.Write([]byte(xml.Header))
	if err != nil {
		return nil, fmt.Errorf("buf write header: %w", err)
	}

	err = xml.NewEncoder(buf).Encode(x)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	return buf.Bytes(), nil
}
