package storj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TokenService struct {
	client *Client
}

type Token struct {
	Token     string    `json:"token"`
	Bucket    string    `json:"bucket"`
	Expires   time.Time `json:"expires"`
	Operation string    `json:"operation"`
}

func (s *TokenService) New(operation, bucketID string) (*Token, error) {
	nonce, err := s.client.generateNonce()
	if err != nil {
		return nil, err
	}

	b := struct {
		Operation string `json:"operation"`
		Nonce     string `json:"__nonce"`
	}{
		operation,
		nonce,
	}

	j, err := json.Marshal(&b)
	if err != nil {
		return nil, err
	}

	rel, _ := url.Parse(fmt.Sprintf("/buckets/%s/tokens", bucketID))
	url := s.client.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	msg := fmt.Sprintf("POST\n/buckets/%s/tokens\n%s", bucketID, j)
	err = s.client.signRequest(req, msg)
	if err != nil {
		return nil, err
	}

	var token Token
	_, err = s.client.Do(req, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
