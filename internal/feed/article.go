package feed

import (
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	readability "github.com/go-shiori/go-readability"
	"net/http"
	"net/url"
	"time"
)

type ArticleFetcher struct {
	client *http.Client
}

func NewArticleFetcher() *ArticleFetcher {
	return &ArticleFetcher{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *ArticleFetcher) Extract(u string) (Article, error) {
	resp, err := a.client.Get(u)
	if err != nil {
		return Article{}, err
	}

	defer resp.Body.Close()

	parsedURL, err := url.Parse(u)
	if err != nil {
		return Article{}, err
	}

	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return Article{}, err
	}

	markdown, err := htmltomarkdown.ConvertString(article.Content)
	if err != nil {
		return Article{}, err
	}

	return Article{
		Title:   article.Title,
		Content: markdown,
		Author:  article.Byline,
		URL:     u,
	}, nil

}
