package storj

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/btcsuite/btcd/btcec"
)

type Client struct {
	client  *http.Client
	BaseURL *url.URL
	AuthKey *btcec.PrivateKey

	Keys     KeyService
	Files    FileService
	Buckets  BucketService
	Contacts ContactService
}

func NewClient() *Client {
	baseURL, _ := url.Parse("https://api.storj.io")

	c := &Client{client: http.DefaultClient, BaseURL: baseURL}

	c.Keys = KeyService{client: c}
	c.Files = FileService{client: c}
	c.Buckets = BucketService{client: c}
	c.Contacts = ContactService{client: c}

	return c
}

func (c *Client) LoadAuthKey(fileName string) error {
	keyHex, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	keyBytes, err := hex.DecodeString(string(keyHex))
	if err != nil {
		return err
	}

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), keyBytes)
	c.AuthKey = privKey

	return nil
}

func (c *Client) Do(req *http.Request, into interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	if status < 200 || status > 299 {
		return nil, fmt.Errorf("got status code %d", status)
	}

	if err := json.NewDecoder(resp.Body).Decode(into); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) generateNonce() (string, error) {
	b := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (c *Client) Sign(msg []byte) (string, error) {
	if c.AuthKey == nil {
		return "", fmt.Errorf("authentication required")
	}

	sha := sha256.Sum256(msg)
	sig, err := c.AuthKey.Sign(sha[:])
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig.Serialize()), nil
}

func (c *Client) signRequest(r *http.Request, msg string) error {
	sig, err := c.Sign([]byte(msg))
	if err != nil {
		return err
	}

	key := hex.EncodeToString(c.AuthKey.PubKey().SerializeCompressed())

	r.Header.Add("x-pubkey", key)
	r.Header.Add("x-signature", sig)
	return nil
}

func (c *Client) newRequest(method, path string) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	url := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) newSignedRequest(method, path string) (*http.Request, error) {
	if method != "GET" && method != "DELETE" && method != "OPTIONS" {
		return nil, fmt.Errorf("bad method")
	}

	nonce, err := c.generateNonce()
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(method, fmt.Sprintf("%s?__nonce=%s", path, nonce))
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("%s\n%s\n__nonce=%s", method, path, nonce)
	err = c.signRequest(req, msg)
	if err != nil {
		return nil, err
	}

	return req, nil
}
