package feed

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// Atom feed structures
type atomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	Title   string      `xml:"title"`
	Subtitle string     `xml:"subtitle"`
	Link    []atomLink  `xml:"link"`
	ID      string      `xml:"id"`
	Updated string      `xml:"updated"`
	Entry   []atomEntry `xml:"entry"`
}

type atomEntry struct {
	Title   string     `xml:"title"`
	Link    []atomLink `xml:"link"`
	ID      string     `xml:"id"`
	Updated string     `xml:"updated"`
	Published string   `xml:"published"`
	Summary string     `xml:"summary"`
	Content string     `xml:"content"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

func (r *Reader) Read(url string) (*Channel, error) {
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading feed: %w", err)
	}

	// Detect feed type by checking the root element
	bodyStr := string(body)
	if strings.Contains(bodyStr, "<feed") || strings.Contains(bodyStr, "<feed>") {
		return r.parseAtom(body, url)
	}

	return r.parseRSS(body, url)
}

func (r *Reader) parseRSS(body []byte, feedURL string) (*Channel, error) {
	var feed rssFeed
	err := xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("error parsing RSS feed: %w", err)
	}

	feed.Channel.FeedURL = feedURL
	return &feed.Channel, nil
}

func (r *Reader) parseAtom(body []byte, feedURL string) (*Channel, error) {
	var atom atomFeed
	err := xml.Unmarshal(body, &atom)
	if err != nil {
		return nil, fmt.Errorf("error parsing Atom feed: %w", err)
	}

	// Convert Atom feed to Channel format
	channel := &Channel{
		Title:       atom.Title,
		Description: atom.Subtitle,
		FeedURL:     feedURL,
	}

	// Extract link from Atom feed (prefer alternate link, fallback to first link)
	for _, link := range atom.Link {
		if link.Rel == "alternate" || link.Rel == "" {
			channel.Link = link.Href
			break
		}
	}
	if channel.Link == "" && len(atom.Link) > 0 {
		channel.Link = atom.Link[0].Href
	}

	// Convert Atom entries to RSS items
	channel.Item = make([]Item, len(atom.Entry))
	for i, entry := range atom.Entry {
		item := Item{
			Title: entry.Title,
		}

		// Extract link from entry (prefer alternate link, fallback to first link)
		for _, link := range entry.Link {
			if link.Rel == "alternate" || link.Rel == "" {
				item.Link = link.Href
				break
			}
		}
		if item.Link == "" && len(entry.Link) > 0 {
			item.Link = entry.Link[0].Href
		}
		// Fallback to entry ID if no link found
		if item.Link == "" {
			item.Link = entry.ID
		}

		// Use published date if available, otherwise use updated date
		if entry.Published != "" {
			item.PubDate = entry.Published
		} else {
			item.PubDate = entry.Updated
		}

		channel.Item[i] = item
	}

	return channel, nil
}
