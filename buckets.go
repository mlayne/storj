package storj

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type BucketService struct {
	client *Client
}

type Bucket struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	User     string    `json:"user"`
	PubKeys  []string  `json:"pubkeys"`
	Status   string    `json:"status"`
	Created  time.Time `json:"created"`
	Storage  int       `json:"storage"`
	Transfer int       `json:"transfer"`
}

func (s *BucketService) List() ([]Bucket, error) {
	nonce, err := s.client.generateNonce()
	if err != nil {
		return nil, err
	}

	rel, _ := url.Parse(fmt.Sprintf("/buckets?__nonce=%s", nonce))
	url := s.client.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("GET\n/buckets\n__nonce=%s", nonce)
	err = s.client.signRequest(req, msg)
	if err != nil {
		return nil, err
	}

	var buckets []Bucket
	_, err = s.client.Do(req, &buckets)
	if err != nil {
		return nil, err
	}

	return buckets, nil
}
