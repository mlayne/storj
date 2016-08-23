package storj

import (
	"fmt"
	"net/url"
)

type FileService struct {
	client *Client
}

type File struct {
	ID       string `json:"id"`
	Bucket   string `json:"bucket"`
	MimeType string `json:"mimetype"`
	Name     string `json:"filename"`
	Size     int64  `json:"size"`
	Frame    string `json:"frame"`
}

func (s *FileService) List(bucketID string) ([]File, error) {
	req, err := s.client.newSignedRequest("GET", fmt.Sprintf("/buckets/%s/files", bucketID))
	if err != nil {
		return nil, err
	}

	var files []File
	_, err = s.client.Do(req, &files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (s *FileService) Delete(bucketID, fileID string) error {
	path := fmt.Sprintf("/buckets/%s/files/%s", bucketID, fileID)
	req, err := s.client.newSignedRequest("DELETE", path)
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

// TODO Reuse Contact here. json.Unmarshal doesn't handle the node-style timestamp that
// /buckets/_/files/_ returns.

type Farmer struct {
	Address  string `json:"address"`
	LastSeen int64  `json:"lastSeen"`
	NodeID   string `json:"nodeID"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type FilePointer struct {
	Hash      string `json:"hash"`
	Token     string `json:"token"`
	Operation string `json:"operation"`
	Farmer    Farmer `json:"farmer"`
}

func (s *FileService) ListPointers(bucketID, fileID, token string) ([]FilePointer, error) {
	rel, _ := url.Parse(fmt.Sprintf("/buckets/%s/files/%s", bucketID, fileID))
	url := s.client.BaseURL.ResolveReference(rel)
	req, err := s.client.newRequest("GET", url.String())
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-token", token)

	var fps []FilePointer
	_, err = s.client.Do(req, &fps)
	if err != nil {
		return nil, err
	}

	return fps, nil
}
