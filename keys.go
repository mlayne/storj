package storj

import (
	"fmt"
	"net/http"
	"net/url"
)

type KeyService struct {
	client *Client
}

type Key struct {
	Key  string `json:"key"`
	User string `json:"user"`
}

func (s *KeyService) List() ([]Key, error) {
	nonce, err := s.client.generateNonce()
	if err != nil {
		return nil, err
	}

	rel, _ := url.Parse(fmt.Sprintf("/keys?__nonce=%s", nonce))
	url := s.client.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("GET\n/keys\n__nonce=%s", nonce)
	err = s.client.signRequest(req, msg)
	if err != nil {
		return nil, err
	}

	var keys []Key
	_, err = s.client.Do(req, &keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
