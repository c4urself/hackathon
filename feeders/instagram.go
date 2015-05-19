package main

import (
	"fmt"
	"github.com/jeffail/gabs"
	"net/http"
	"io/ioutil"
	"log"
	"strconv"
)

type CreatorPhoto struct {
	url string
	likes int64
}

type CreatorFeed struct {
	photos []CreatorPhoto
	audience []string
}

// Fetches photos/auedience from Instagram by username
func getCreatorFeed(username string) CreatorFeed {
	resp, err := http.Get(fmt.Sprintf("https://instagram.com/%s/media/", username))
	defer resp.Body.Close()

	if err != nil {
		log.Printf("Unable to load instagram feed", err)
		return CreatorFeed{}
	}

	bresp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read instagram feed", err)
		return CreatorFeed{}
	}

	jresp, err := gabs.ParseJSON(bresp)
	if err != nil {
		log.Printf("Unable to parse instagram feed", err)
		return CreatorFeed{}
	}

	// Fetch profile photos
	items, _ := jresp.Search("items").Children()

	var creatorPhotos []CreatorPhoto
	var audiencePhotos map[string]bool = make(map[string]bool, len(items) * 50)

	for _, item := range items  {
		url := item.Path("images.standard_resolution.url").String()
		likes, _ := strconv.ParseInt(item.Path("likes.count").String(), 10, 64)
		creatorPhotos = append(creatorPhotos, CreatorPhoto{url: url, likes: likes})

		// Parse likers
		likers, _ := item.Path("likes.data.profile_picture").Children()
		for _, liker := range likers {
			audiencePhotos[liker.String()] = true
		}

		// Parse commenters
		commenters, _ := item.Path("comments.data.from.profile_picture").Children()
		for _, commenter := range commenters {
			audiencePhotos[commenter.String()] = true
		}	
	}

	// Audience includes only unique urls and doesn't include placeholder photo
	audience := make([]string, len(audiencePhotos))
	for k := range audiencePhotos {
		if k != "https://instagramimages-a.akamaihd.net/profiles/anonymousUser.jpg" {
			audience = append(audience, k)
		}
	}

	return CreatorFeed{photos: creatorPhotos, audience: audience}
}
