package feeders

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jeffail/gabs"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Photo struct {
	Id    string
	Url   string
	Likes int64
}

type Photos []Photo

func (s Photos) Len() int {
	return len(s)
}

func (s Photos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Photos) Less(i, j int) bool {
	return s[i].Likes < s[j].Likes
}

type CreatorFeed struct {
	Photos   Photos
	Audience Photos
}

func NewPhoto(url string, likes int64) Photo {
	return Photo{Url: url, Id: Hash(url), Likes: likes}
}

func Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func (f CreatorFeed) GetTopPhotos(count int) Photos {
	sort.Sort(sort.Reverse(Photos(f.Photos)))
	var topPhotos Photos

	if count > len(f.Photos) {
		topPhotos = f.Photos[:len(f.Photos)]
	} else {
		topPhotos = f.Photos[:count]
	}
	return topPhotos
}

// Fetches photos/auedience from Instagram by username
func GetCreatorFeed(username string) CreatorFeed {
	resp, err := http.Get(fmt.Sprintf("https://instagram.com/%s/media/", username))
	defer resp.Body.Close()

	if err != nil {
		log.Printf("Unable to load instagram feed %v", err)
		return CreatorFeed{}
	}

	bresp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read instagram feed %v", err)
		return CreatorFeed{}
	}

	jresp, err := gabs.ParseJSON(bresp)
	if err != nil {
		log.Printf("Unable to parse instagram feed %v", err)
		return CreatorFeed{}
	}

	// Fetch profile photos
	items, _ := jresp.Search("items").Children()

	var creatorPhotos Photos
	var audiencePhotos map[string]bool = make(map[string]bool, len(items)*50)

	for _, item := range items {
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
	var audience Photos
	for k := range audiencePhotos {
		k = strings.Trim(k, "\"")
		if k != "https://instagramimages-a.akamaihd.net/profiles/anonymousUser.jpg" {
			audience = append(audience, NewPhoto(k, 0))
		}
	}

	return CreatorFeed{Photos: creatorPhotos, Audience: audience}
}

func LoadPhotos(photos Photos, baseDir string) {
	var tasks chan Task = make(chan Task)

	workers := 10
	if len(photos) < workers {
		workers = len(photos)
	}

	// Run 10 parallel tasks
	for i := 0; i < workers; i++ {
		go loader(tasks)
	}

	// Populate channel with tasks
	for _, photo := range photos {
		tasks <- Task{
			url:  photo.Url,
			path: filepath.Join(baseDir, fmt.Sprintf("%s.png", photo.Id))}
	}
	close(tasks)
}

type Task struct {
	url  string
	path string
}

func loader(tasks <-chan Task) {
	for task := range tasks {
		log.Printf("Loading url %s", task.url)

		file, err := os.Create(task.path)
		defer file.Close()
		if err != nil {
			log.Printf("Unable to creator destination file %s, %s", task.path, err)
			continue
		}

		resp, err := http.Get(task.url)
		if err != nil {
			log.Printf("Unable to fetch url %s, %s", task.url, err)
			continue
		}

		img, err := jpeg.Decode(resp.Body)
		if err != nil {
			log.Printf("Unable to convert jpg -> png %s, %s", task.url, err)
			continue
		}

		png.Encode(file, img)
	}
}
