package storj

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

const (
	bucketJson = `
  {
    "storage": 10,
    "transfer": 30,
    "status": "Active",
    "pubkeys": [
      "031a259ee122414f57a63bbd6887ee17960e9106b0adcf89a298cdad2108adf4d9"
    ],
    "user": "gordon@storj.io",
    "name": "New Bucket",
    "created": "2016-03-04T17:01:02.629Z",
    "id": "507f1f77bcf86cd799439011"
  }`
)

var exBucket = Bucket{
	ID:       "507f1f77bcf86cd799439011",
	Name:     "New Bucket",
	User:     "gordon@storj.io",
	PubKeys:  []string{"031a259ee122414f57a63bbd6887ee17960e9106b0adcf89a298cdad2108adf4d9"},
	Status:   "Active",
	Created:  time.Date(2016, 3, 4, 17, 1, 2, 629000000, time.UTC),
	Storage:  10,
	Transfer: 30}

func TestBucketsList(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Buckets.List()
	if err == nil || err.Error() != "authentication required" {
		t.Errorf("Buckets.List should require authentication")
	}

	enableAuth()
	defer disableAuth()

	pubKey := hex.EncodeToString(privKey.PubKey().SerializeCompressed())

	mux.HandleFunc("/buckets", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, "GET")
		assertHeader(t, r, "x-pubkey", pubKey)
		if r.Header.Get("x-signature") == "" {
			t.Errorf(`missing "x-signature" header`)
		}
		fmt.Fprintf(w, "[%s]", bucketJson)
	})

	buckets, err := client.Buckets.List()
	if err != nil {
		t.Errorf("Buckets.List returned error: %v", err)
	}

	expected := []Bucket{exBucket}
	if !reflect.DeepEqual(buckets, expected) {
		t.Errorf("Buckets.List returned %+v, expected %+v", buckets, expected)
	}
}

func TestBucketsNew(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Buckets.New("test bucket", 42, 43)
	if err == nil || err.Error() != "authentication required" {
		t.Errorf("Buckets.List should require authentication")
	}

	enableAuth()
	defer disableAuth()

	pubKey := hex.EncodeToString(privKey.PubKey().SerializeCompressed())

	mux.HandleFunc("/buckets", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, "POST")
		assertHeader(t, r, "x-pubkey", pubKey)
		assertHeader(t, r, "Content-Type", "application/json")
		if r.Header.Get("x-signature") == "" {
			t.Errorf(`missing "x-signature" header`)
		}
		var sent map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&sent); err != nil {
			t.Errorf("received bad JSON")
		}
		if sent["name"] != "test bucket" || sent["storage"] != float64(42) || sent["transfer"] != float64(43) {
			t.Errorf("invalid request parameters")
		}
		if _, ok := sent["__nonce"]; !ok {
			t.Errorf("request did not contain a nonce")
		}
		fmt.Fprintf(w, bucketJson)
	})

	bucket, err := client.Buckets.New("test bucket", 42, 43)
	if err != nil {
		t.Errorf("Buckets.New returned error: %v", err)
	}

	if !reflect.DeepEqual(bucket, &exBucket) {
		t.Errorf("Buckets.New returned %+v, expected %+v", bucket, exBucket)
	}
}
