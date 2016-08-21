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
	tokenJson = `
  {
    "token": "a_token",
    "bucket": "bucket_id",
    "expires": "2016-03-04T17:01:02.629Z",
    "operation": "PULL"
  }`
)

var exToken = Token{
	Token:     "a_token",
	Bucket:    "bucket_id",
	Expires:   time.Date(2016, 3, 4, 17, 1, 2, 629000000, time.UTC),
	Operation: "PULL"}

func TestTokensNew(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Tokens.New("PULL", "buket_id")
	if err == nil || err.Error() != "authentication required" {
		t.Errorf("Tokens.New should require authentication")
	}

	enableAuth()
	defer disableAuth()

	pubKey := hex.EncodeToString(privKey.PubKey().SerializeCompressed())

	mux.HandleFunc("/buckets/bucket_id/tokens", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, "POST")
		assertHeader(t, r, "x-pubkey", pubKey)
		assertHeader(t, r, "Content-Type", "application/json")
		if r.Header.Get("x-signature") == "" {
			t.Errorf(`missing "x-signature" header`)
		}
		var sent map[string]string
		if err := json.NewDecoder(r.Body).Decode(&sent); err != nil {
			t.Errorf("received bad JSON")
		}
		if sent["operation"] != "PULL" {
			t.Errorf("request missing operation")
		}
		if _, ok := sent["__nonce"]; !ok {
			t.Errorf("request did not contain a nonce")
		}
		fmt.Fprintf(w, tokenJson)
	})

	token, err := client.Tokens.New("PULL", "bucket_id")
	if err != nil {
		t.Errorf("Tokens.New returned error: %v", err)
	}

	if !reflect.DeepEqual(token, &exToken) {
		t.Errorf("Tokens.New returned %+v, expected %+v", token, exToken)
	}
}
