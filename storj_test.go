package storj

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/btcsuite/btcd/btcec"
)

var (
	mux    *http.ServeMux
	server *httptest.Server

	client *Client

	privKey *btcec.PrivateKey
)

func init() {
	var err error
	privKey, err = btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		panic("failed to generate test key")
	}
}

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient()
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func enableAuth() {
	client.AuthKey = privKey
}

func disableAuth() {
	client.AuthKey = nil
}

func assertMethod(t *testing.T, r *http.Request, expected string) {
	if r.Method != expected {
		t.Errorf("got request method %v, expected %v", r.Method, expected)
	}
}

func assertHeader(t *testing.T, r *http.Request, key, value string) {
	if h := r.Header.Get(key); h != value {
		t.Errorf("expected header %q to be %q, but got %q", key, value, h)
	}
}
