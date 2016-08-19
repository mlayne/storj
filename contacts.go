package storj

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ContactService struct {
	client *Client
}

type Contact struct {
	Address  string    `json:"address"`
	Port     int       `json:"port"`
	NodeID   string    `json:"nodeID"`
	LastSeen time.Time `json:"lastSeen"`
	Protocol string    `json:"protocol"`
}

func (s *ContactService) Get(nodeID string) (*Contact, error) {
	rel, err := url.Parse(fmt.Sprintf("/contacts/%s", nodeID))
	if err != nil {
		return nil, err
	}
	url := s.client.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var contact Contact
	resp, err := s.client.Do(req, &contact)
	if err != nil {
		return nil, err
	}

	status := resp.StatusCode
	if status < 200 || status >= 400 {
		return nil, fmt.Errorf("got bad status code")
	}

	return &contact, nil
}

func (s *ContactService) List() ([]Contact, error) {
	rel, _ := url.Parse("/contacts")
	url := s.client.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var contacts []Contact
	resp, err := s.client.Do(req, &contacts)
	if err != nil {
		return nil, err
	}

	status := resp.StatusCode
	if status < 200 || status >= 400 {
		return nil, errors.New("got bad status code")
	}

	return contacts, nil
}
