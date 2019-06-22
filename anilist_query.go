package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
		"query": query,
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
	mediaData, err := UnmarshalAnilist(body)
	if err != nil {
		log.Fatal(err)
	}
	// testing the parser
	fmt.Println(mediaData.Data.Page.Media[0].Relations.Edges[0].Node.Type)
	//fmt.Println("resp body:", string(body))
}

// Main query for anilist API.
// FIXME: @perPage and @page static or variable? What was the requirement for this?
var query = `
query ($page: Int = 1, $perPage: Int = 1, $id: Int, $type: MediaType = MANGA) {
	Page(page: $page, perPage: $perPage) {
	  pageInfo {
		total
		perPage
		currentPage
		lastPage
		hasNextPage
	  }
	  media(id: $id, type: $type) {
		id
		idMal
		coverImage {
		  large
		  medium
		}
		bannerImage
		title {
		  romaji
		  english
		  native
		}
		startDate {
		  year
		  month
		  day
		}
		endDate {
		  year
		  month
		  day
		}
		status
		chapters
		volumes
		genres
		tags {
		  name
		  rank
		  category
			  isGeneralSpoiler
		  isMediaSpoiler
		}
		popularity
		staff {
		  edges {
			id
			role
			node{
			  name {
				first
				last
				native
			  }
			  image {
				large
				medium
			  }
			}
		  }
		}
		characters {
		  edges {
			id
			role
			node{
			  image {
				large
				medium
			  }
			  name {
				first
				last
				native
			  }
			}
		  }
		}
		relations {
		  edges {
			id
			  relationType
			node{
			  bannerImage
			  title {
				romaji
				english
				native
			  }
			  type
			  status
			  idMal
			}
		  }
		}
	  }
	}
  }
  
`
