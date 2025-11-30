package feed

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Reader struct {
	client *http.Client
}

func NewReader() *Reader {
	return &Reader{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type rssFeed struct {
	Channel Channel `xml:"channel"`
}

func (r *Reader) Read(url string) (*Channel, error) {
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching RSS feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading RSS feed: %w", err)
	}

	var feed rssFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("error parsing RSS feed: %w", err)
	}

	return &feed.Channel, nil
}
