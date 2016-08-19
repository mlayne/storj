package storj

import (
	"encoding/hex"
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
	if s.client.AuthKey == nil {
		return nil, fmt.Errorf("authentication required")
	}

	nonce, err := s.client.generateNonce()
	if err != nil {
		return nil, err
	}

	msg := []byte(fmt.Sprintf("GET\n/buckets\n__nonce=%s", nonce))
	sig, err := s.client.Sign(msg)
	if err != nil {
		return nil, err
	}

	rel, _ := url.Parse(fmt.Sprintf("/buckets?__nonce=%s", nonce))
	url := s.client.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-signature", sig)
	req.Header.Add("x-pubkey", hex.EncodeToString(s.client.AuthKey.PubKey().SerializeCompressed()))

	var buckets []Bucket
	_, err = s.client.Do(req, &buckets)
	if err != nil {
		return nil, err
	}

	return buckets, nil
}
