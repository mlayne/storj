# storj

[![Build Status](https://travis-ci.org/mlayne/storj.svg?branch=master)](https://travis-ci.org/mlayne/storj)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A Go client library for the [Storj Bridge API](https://github.com/storj/bridge).

## Setup
You'll need to have an already set up storj account in order to leverage this package. The authentication key for the client needs to have the associated public key linked to the storj account. 

## Example
The example below will list each bucket's name and content.
```go
package main

import (
	"fmt"
	"log"

	"github.com/mlayne/storj"
)

func main() {
	// authKey is the file which contains a private key associated
	// with one's storj account.
	authKey := "private.key"

	c := storj.NewClient()

	err := c.LoadAuthKey(authKey)
	if err != nil {
		log.Fatalf("Failed to load auth key: %s", err)
	}

	buckets, err := c.Buckets.List()
	if err != nil {
		log.Fatalf("Error with listing buckets: %s", err)
	}

	// Now we'll print the bucket names and the files contained.
	for _, b := range buckets {
		fmt.Println(b.Name)

		files, err := c.Files.List(b.ID)
		if err != nil {
			log.Fatalf("Failed to list files: %s", err)
		}

		// Next we'll print the files contained in the bucket.
		for _, f := range files {
			fmt.Println("\t-" + f.Name)
		}
	}
}
```
