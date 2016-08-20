package storj

import (
	"bytes"
	"encoding/json"
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
	req, err := s.client.newSignedRequest("GET", "/buckets")
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

func (s *BucketService) New(name string, storage, transfer int) (*Bucket, error) {
	nonce, err := s.client.generateNonce()
	if err != nil {
		return nil, err
	}

	b := struct {
		Name     string `json:"name"`
		Storage  int    `json:"storage"`
		Transfer int    `json:"transfer"`
		Nonce    string `json:"__nonce"`
	}{
		name,
		storage,
		transfer,
		nonce,
	}

	j, err := json.Marshal(&b)
	if err != nil {
		return nil, err
	}

	rel, _ := url.Parse("/buckets")
	url := s.client.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	msg := fmt.Sprintf("POST\n/buckets\n%s", j)
	err = s.client.signRequest(req, msg)
	if err != nil {
		return nil, err
	}

	var bucket Bucket
	_, err = s.client.Do(req, &bucket)
	if err != nil {
		return nil, err
	}

	return &bucket, nil
}

func (s *BucketService) Delete(bucketID string) error {
	req, err := s.client.newSignedRequest("DELETE", fmt.Sprintf("/buckets/%s", bucketID))
	if err != nil {
		return err
	}

	resp, err := s.client.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return fmt.Errorf("expected status 204, got %d", resp.StatusCode)
	}

	return nil
}
