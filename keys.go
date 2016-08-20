package storj

import (
	"bytes"
	"encoding/json"
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
	req, err := s.client.newSignedRequest("GET", "/keys")
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

func (s *KeyService) Register(key string) error {
	nonce, err := s.client.generateNonce()
	if err != nil {
		return err
	}

	k := struct {
		Key   string `json:"key"`
		Nonce string `json:"__nonce"`
	}{
		key,
		nonce,
	}

	j, err := json.Marshal(&k)
	if err != nil {
		return err
	}

	rel, _ := url.Parse("/keys")
	url := s.client.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(j))
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("POST\n/keys\n%s", j)
	err = s.client.signRequest(req, msg)
	if err != nil {
		return err
	}

	var respKey Key
	_, err = s.client.Do(req, &respKey)
	if err != nil {
		return err
	}
	if respKey.Key != key {
		return fmt.Errorf("received non-matching key from server")
	}

	return nil
}

func (s *KeyService) Delete(key string) error {
	req, err := s.client.newSignedRequest("DELETE", fmt.Sprintf("/keys/%s", key))
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
