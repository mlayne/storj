package storj

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

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
		fmt.Fprintf(w, `[
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
  }
]`)
	})

	buckets, err := client.Buckets.List()
	if err != nil {
		t.Errorf("Buckets.List returned error: %v", err)
	}

	expected := []Bucket{{
		ID:       "507f1f77bcf86cd799439011",
		Name:     "New Bucket",
		User:     "gordon@storj.io",
		PubKeys:  []string{"031a259ee122414f57a63bbd6887ee17960e9106b0adcf89a298cdad2108adf4d9"},
		Status:   "Active",
		Created:  time.Date(2016, 3, 4, 17, 1, 2, 629000000, time.UTC),
		Storage:  10,
		Transfer: 30}}
	if !reflect.DeepEqual(buckets, expected) {
		t.Errorf("Buckets.List returned %+v, expected %+v", buckets, expected)
	}
}
