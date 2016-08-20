package storj

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const (
	fileJson = `
  {
    "id": "507f1f77bcf86cd799439011",
    "bucket": "607f1f77bcf86cd799439011",
    "mimetype": "video/mpeg",
    "filename": "big_buck_bunny.mp4",
    "size": 5071076,
    "frame": "707f1f77bcf86cd799439011"
  }`
)

var exFile = File{
	ID:       "507f1f77bcf86cd799439011",
	Bucket:   "607f1f77bcf86cd799439011",
	MimeType: "video/mpeg",
	Name:     "big_buck_bunny.mp4",
	Size:     5071076,
	Frame:    "707f1f77bcf86cd799439011"}

func TestFilesList(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Files.List("xyz")
	if err == nil || err.Error() != "authentication required" {
		t.Errorf("Files.List should require authentication")
	}

	enableAuth()
	defer disableAuth()

	pubKey := hex.EncodeToString(privKey.PubKey().SerializeCompressed())

	mux.HandleFunc("/buckets/xyz/files", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, "GET")
		assertHeader(t, r, "x-pubkey", pubKey)
		if r.Header.Get("x-signature") == "" {
			t.Errorf(`missing "x-signature" header`)
		}
		fmt.Fprintf(w, "[%s]", fileJson)
	})

	files, err := client.Files.List("xyz")
	if err != nil {
		t.Errorf("Files.List returned error: %v", err)
	}

	expected := []File{exFile}
	if !reflect.DeepEqual(files, expected) {
		t.Errorf("Files.List returned %+v, expected %+v", files, expected)
	}
}
