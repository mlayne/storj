package storj

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestKeysList(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Keys.List()
	if err == nil || err.Error() != "authentication required" {
		t.Errorf("Keys.List should require authentication")
	}

	enableAuth()
	defer disableAuth()

	pubKey := hex.EncodeToString(privKey.PubKey().SerializeCompressed())

	mux.HandleFunc("/keys", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, "GET")
		assertHeader(t, r, "x-pubkey", pubKey)
		if r.Header.Get("x-signature") == "" {
			t.Errorf(`missing "x-signature" header`)
		}
		fmt.Fprintf(w, `[{"key": "031a259ee122414f57a63bbd6887ee17960e9106b0adcf89a298cdad2108adf4d9", "user": "gordon@storj.io"}]`)
	})

	keys, err := client.Keys.List()
	if err != nil {
		t.Errorf("Keys.List returned error: %v", err)
	}

	expected := []Key{{
		Key:  "031a259ee122414f57a63bbd6887ee17960e9106b0adcf89a298cdad2108adf4d9",
		User: "gordon@storj.io"}}
	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Keys.List returned %+v, expected %+v", keys, expected)
	}
}
