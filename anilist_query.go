package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	anilistURL = "https://graphql.anilist.co/"
)

func sendPostQueryToAnilist() {
	fmt.Println("Url:", anilistURL)

	// Creating the request body.
	reqBody, err := json.Marshal(map[string]string{
		"query": `query ($page: Int = 1, $perPage: Int = 1, $id: Int, $type: MediaType = MANGA) {
			Page(page: $page, perPage: $perPage) {
			  pageInfo {
				total
				perPage
				currentPage
			  }
			  media(id: $id, type: $type) {
				id
				idMal
				coverImage {
				  large
				  medium
				}
			  }
			}
		  }`,
	})

	req, err := http.NewRequest("POST", anilistURL, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
		return
	}
	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("resp body:", string(body))
}
