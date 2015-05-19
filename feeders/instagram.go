package feeders

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jeffail/gabs"
	"net/http"
	"image/jpeg"
 	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Photo struct {
	Id string
	Url string
	Likes int64
}

type CreatorFeed struct {
	Photos []Photo
	Audience []Photo
}

func NewPhoto(url string, likes int64) Photo {
	return Photo{Url: url, Id: Hash(url), Likes: likes}
}

func Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// Fetches photos/auedience from Instagram by username
func GetCreatorFeed(username string) CreatorFeed {
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

	var creatorPhotos []Photo
	var audiencePhotos map[string]bool = make(map[string]bool, len(items) * 50)

	for _, item := range items  {
		url := item.Path("images.standard_resolution.url").String()
		url = strings.Trim(url, "\"")
		likes, _ := strconv.ParseInt(item.Path("likes.count").String(), 10, 64)
		creatorPhotos = append(creatorPhotos, NewPhoto(url, likes))

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
	var audience []Photo
	for k := range audiencePhotos {
		k = strings.Trim(k, "\"")
		if k != "https://instagramimages-a.akamaihd.net/profiles/anonymousUser.jpg" {
			audience = append(audience, NewPhoto(k, 0))
		}
	}

	return CreatorFeed{Photos: creatorPhotos, Audience: audience}
}

func LoadPhotos(photos []Photo, baseDir string) {
	for _, photo := range photos {
		log.Printf("Loading photo %s", photo.Url)

		fPath := filepath.Join(baseDir, fmt.Sprintf("%s.png", photo.Id))
		f, err := os.Create(fPath)
		defer f.Close()
		if err != nil {
			log.Printf("Unable to create photo for loading %s, %s", fPath, err)
			continue
		}

		resp, err := http.Get(photo.Url)
		if err != nil {
			log.Printf("Unable to fetch photo %s, %s", photo.Url, err)
			continue
		}

		img, err := jpeg.Decode(resp.Body)
		if err != nil {
			log.Printf("Unable to convert jpg -> png %s, %s", photo.Url, err)
		}

		png.Encode(f, img)
	}
}
