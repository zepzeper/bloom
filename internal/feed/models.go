package feed

type Channel struct {
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`        // Website URL from RSS feed
	Language    string   `xml:"language"`
	Item        []Item   `xml:"item"`
	Category    string   `xml:"category"`
	Tags        []string `xml:"tags"`
	FeedURL     string   `xml:"-"`            // The feed URL used to fetch this channel (not in XML)
}

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
	Read    bool   `xml:"-"`
}

type Article struct {
	Title   string
	Content string
	Author  string
	URL     string
}
