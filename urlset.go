package sitemap

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

const (
	FreqAlways  = "always" // The value "always" should be used to describe documents that change each time they are accessed.
	FreqHourly  = "hourly"
	FreqDaily   = "daily"
	FreqWeekly  = "weekly"
	FreqMonthly = "monthly"
	FreqYearly  = "yearly"
	FreqNewer   = "never" // The value "never" should be used to describe archived URLs.
)

type URL struct {
	Loc        string `xml:"loc"`
	Lastmod    string `xml:"lastmod,omitempty"`
	Changefreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URL     []URL    `xml:"url"`
}

func NewURLSet() *URLSet {
	return &URLSet{
		Xmlns: xmlns,
	}
}

func (x *URLSet) Marshal() (data []byte, err error) {
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
