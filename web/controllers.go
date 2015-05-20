package web

import (
	"fmt"
	"github.com/c4urself/hackathon/feeders"
	"github.com/c4urself/hackathon/mosaic"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"html/template"
	"strings"
)

const BASE_TMPL = "templates/base.tmpl"
const REDIS_ADDRESS = "127.0.0.1:6379"

func StartApp() {
	r := gin.Default()
	r.Static("/static", "./static")

	// Redis connection
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", REDIS_ADDRESS)

		if err != nil {
			return nil, err
		}

		return c, err
	}, 10)
	defer redisPool.Close()

	// Homepage
	r.GET("/", func(c *gin.Context) {
		obj := gin.H{"title": "Main website"}
		tmpl := template.Must(template.ParseFiles(BASE_TMPL, "templates/home.tmpl"))
		r.SetHTMLTemplate(tmpl)
		c.HTML(200, "base", obj)
	})

	// Search page
	r.GET("/search", func(c *gin.Context) {
		c.Request.ParseForm()
		username := c.Request.Form.Get("username")

		var topPhotos feeders.Photos
		feed := feeders.GetCreatorFeed(username)
		topPhotos = feed.GetTopPhotos(5)

		obj := gin.H{"top_photos": topPhotos}
		tmpl := template.Must(template.ParseFiles(BASE_TMPL, "templates/result.tmpl"))
		r.SetHTMLTemplate(tmpl)
		c.HTML(200, "base", obj)
	})

	// Result page
	r.GET("/mosaic/:username", func(c *gin.Context) {
		c.Request.ParseForm()
		resetCache := c.Request.Form.Get("reset")
		username := c.Params.ByName("username")

		// Check for cache
		client := redisPool.Get()
		cachecKey := fmt.Sprintf("hackathon:%s", username)
		if resetCache == "true" {
			client.Do("DEL", cachecKey)
		}
		cache, err := redis.String(client.Do("GET", cachecKey))

		var mosaics []mosaic.Mosaic
		if err == nil {
			mosaics = deserialize(cache)
		} else {
			mosaics = mosaic.MakeInstagramMosaic(
				username,
				fmt.Sprintf("/tmp/hack/%s/photos/", username),
				fmt.Sprintf("/tmp/hack/%s/audience/", username),
				fmt.Sprintf("./static/mosaic/%s/", username))
			client.Do("SET", cachecKey, serialize(mosaics))
		}

		obj := gin.H{
			"username": c.Params.ByName("username"),
			"baseUrl":  fmt.Sprintf("/static/mosaic/%s/", username),
			"mosaics":  mosaics}

		tmpl := template.Must(template.ParseFiles(BASE_TMPL, "templates/result.tmpl"))
		r.SetHTMLTemplate(tmpl)
		c.HTML(200, "base", obj)
	})

	r.Run(":8080")
}

func serialize(mosaics []mosaic.Mosaic) string {
	b, _ := json.Marshal(mosaics)
	return string(b)
}

func deserialize(data string) []mosaic.Mosaic {
	var mosaics []mosaic.Mosaic
	dec := json.NewDecoder(strings.NewReader(data))
	dec.Decode(&mosaics)
	return mosaics
}
