package storj

import (
	"fmt"
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
	req, err := s.client.newRequest("GET", fmt.Sprintf("/contacts/%s", nodeID))
	if err != nil {
		return nil, err
	}

	var contact Contact
	_, err = s.client.Do(req, &contact)
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (s *ContactService) List() ([]Contact, error) {
	req, err := s.client.newRequest("GET", "/contacts")
	if err != nil {
		return nil, err
	}

	var contacts []Contact
	_, err = s.client.Do(req, &contacts)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}
