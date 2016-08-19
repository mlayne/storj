package storj

import (
	"encoding/hex"
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
	if s.client.AuthKey == nil {
		return nil, fmt.Errorf("authentication required")
	}

	nonce, err := s.client.generateNonce()
	if err != nil {
		return nil, err
	}

	msg := []byte(fmt.Sprintf("GET\n/keys\n__nonce=%s", nonce))
	sig, err := s.client.Sign(msg)
	if err != nil {
		return nil, err
	}

	rel, _ := url.Parse(fmt.Sprintf("/keys?__nonce=%s", nonce))
	url := s.client.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-signature", sig)
	req.Header.Add("x-pubkey", hex.EncodeToString(s.client.AuthKey.PubKey().SerializeCompressed()))

	var keys []Key
	_, err = s.client.Do(req, &keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
