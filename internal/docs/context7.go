package docs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Snippet struct {
	Title string `json:"codeTitle"`
	Code  string `json:"code"`
}

type DocsResult struct {
	LibraryID    string    `json:"-"`
	CodeSnippets []Snippet `json:"codeSnippets"`
	InfoSnippets []struct {
		Content string `json:"content"`
	} `json:"infoSnippets"`
}

type Client struct {
	apiKey  string
	baseURL string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://context7.com",
	}
}

type SearchHit struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Score       float64 `json:"score"`
}

type searchResponse struct {
	Results []SearchHit `json:"results"`
}

func (c *Client) Search(libraryName, query string) ([]SearchHit, error) {
	if c.apiKey == "" {
		return nil, nil
	}
	u := fmt.Sprintf("%s/api/v2/libs/search?libraryName=%s&query=%s",
		c.baseURL, url.QueryEscape(libraryName), url.QueryEscape(query))

	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("context7 search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, nil
	}

	var data searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("context7 decode: %w", err)
	}
	return data.Results, nil
}

func (c *Client) GetDocs(libraryID, query string) (*DocsResult, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("context7 api key required")
	}
	u := fmt.Sprintf("%s/api/v2/context?libraryId=%s&query=%s&type=json",
		c.baseURL, url.QueryEscape(libraryID), url.QueryEscape(query))

	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("context7 docs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("context7 docs status: %d", resp.StatusCode)
	}

	var result DocsResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("context7 docs decode: %w", err)
	}
	result.LibraryID = libraryID
	return &result, nil
}
