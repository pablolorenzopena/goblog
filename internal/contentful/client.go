package contentful

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const apiURL = "https://cdn.contentful.com"

type Client struct {
	SpaceID     string
	AccessToken string
	HTTPClient  *http.Client
}

func NewClient(spaceID, accessToken string) *Client {
	return &Client{
		SpaceID:     spaceID,
		AccessToken: accessToken,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type ContentfulPost struct {
	Sys struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
	} `json:"sys"`
	Fields struct {
		Title         string    `json:"title"`
		Slug          string    `json:"slug"`
		PublishedDate time.Time `json:"publishedDate"`
		Body          string    `json:"body"`
		FeatureImage  struct {
			Fields struct {
				File struct {
					URL string `json:"url"`
				} `json:"file"`
			} `json:"fields"`
		} `json:"featureImage"`
		Tags []string `json:"tags"`
	} `json:"fields"`
}

type getEntriesResponse struct {
	Items []ContentfulPost `json:"items"`
}

func (c *Client) GetAllPosts() ([]ContentfulPost, error) {
	endpoint := fmt.Sprintf("%s/spaces/%s/environments/master/entries", apiURL, c.SpaceID)
	params := url.Values{}
	params.Set("access_token", c.AccessToken)
	params.Set("content_type", "post")
	params.Set("order", "-fields.publishedDate")

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	resp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error haciendo request a Contentful: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	var data getEntriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %w", err)
	}

	return data.Items, nil
} 