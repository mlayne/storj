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

	Contacts ContactService
	Keys     KeyService
}

func NewClient() *Client {
	baseURL, _ := url.Parse("https://api.storj.io")

	c := &Client{client: http.DefaultClient, BaseURL: baseURL}
	c.Contacts = ContactService{client: c}
	c.Keys = KeyService{client: c}

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

	if err := json.NewDecoder(resp.Body).Decode(into); err != nil {
		fmt.Printf("%v\n", resp)
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
