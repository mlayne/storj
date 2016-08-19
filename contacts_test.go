package storj

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var exContact = Contact{
	Address:  "api.storj.io",
	Port:     8443,
	NodeID:   "32033d2dc11b877df4b1caefbffba06495ae6b18",
	LastSeen: time.Date(2016, 5, 24, 15, 16, 1, 139000000, time.UTC),
	Protocol: "0.7.0"}

func TestContactsGet(t *testing.T) {
	setup()
	defer teardown()

	path := "/contacts/32033d2dc11b877df4b1caefbffba06495ae6b18"
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, "GET")
		fmt.Fprintf(w, `{"address": "api.storj.io", "port": 8443, "nodeID": "32033d2dc11b877df4b1caefbffba06495ae6b18", "lastSeen": "2016-05-24T15:16:01.139Z", "protocol": "0.7.0"}`)
	})

	contact, err := client.Contacts.Get(exContact.NodeID)
	if err != nil {
		t.Errorf("Contacts.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(contact, &exContact) {
		t.Errorf("Contacts.Get returned %+v, expected %+v", contact, exContact)
	}
}

func TestContactsList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, "GET")
		fmt.Fprintf(w, `[{"address": "api.storj.io", "port": 8443, "nodeID": "32033d2dc11b877df4b1caefbffba06495ae6b18", "lastSeen": "2016-05-24T15:16:01.139Z", "protocol": "0.7.0"}]`)
	})

	contacts, err := client.Contacts.List()
	if err != nil {
		t.Errorf("Contacts.List returned error: %v", err)
	}

	expected := []Contact{exContact}
	if !reflect.DeepEqual(contacts, expected) {
		t.Errorf("Contacts.List returned %+v, expected %+v", contacts, expected)
	}
}
