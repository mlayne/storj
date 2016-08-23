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

func TestFilesDelete(t *testing.T) {
	setup()
	defer teardown()

	err := client.Files.Delete("abc", "xyz")
	if err == nil || err.Error() != "authentication required" {
		t.Errorf("Files.Delete should require authentication")
	}

	enableAuth()
	defer disableAuth()

	pubFile := hex.EncodeToString(privKey.PubKey().SerializeCompressed())

	mux.HandleFunc("/buckets/abc/files/xyz", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, "DELETE")
		assertHeader(t, r, "x-pubkey", pubFile)
		if r.Header.Get("x-signature") == "" {
			t.Errorf(`missing "x-signature" header`)
		}
		w.WriteHeader(204)
	})

	err = client.Files.Delete("abc", "xyz")
	if err != nil {
		t.Errorf("Files.Delete returned error: %v", err)
	}
}

func TestFilesListPointers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/buckets/abc/files/xyz", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, "GET")
		assertHeader(t, r, "x-token", "a_token")
		fmt.Fprintf(w, `[
  {
    "hash": "ba084d3f143f2896809d3f1d7dffed472b39d8de",
    "token": "99cf1af00b552113a856f8ef44f58d22269389e8009d292bafd10af7cc30dcfa",
    "operation": "PULL",
    "farmer": {
      "address": "api.storj.io",
      "port": 8443,
      "nodeID": "32033d2dc11b877df4b1caefbffba06495ae6b18",
      "lastSeen": 1471922911187,
      "protocol": "0.7.0"
    }
  }]`)
	})

	fps, err := client.Files.ListPointers("abc", "xyz", "a_token")
	if err != nil {
		t.Errorf("Files.ListPointers returned error: %v", err)
	}

	expected := []FilePointer{{
		Hash:      "ba084d3f143f2896809d3f1d7dffed472b39d8de",
		Token:     "99cf1af00b552113a856f8ef44f58d22269389e8009d292bafd10af7cc30dcfa",
		Operation: "PULL",
		Farmer: Farmer{
			Address:  "api.storj.io",
			Port:     8443,
			NodeID:   "32033d2dc11b877df4b1caefbffba06495ae6b18",
			LastSeen: 1471922911187,
			Protocol: "0.7.0"}}}
	if !reflect.DeepEqual(fps, expected) {
		t.Errorf("Files.ListPointers returned %+v, expected %+v", fps, expected)
	}
}
