package storj

import "fmt"

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
