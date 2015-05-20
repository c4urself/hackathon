package web

import (
	"fmt"
	"github.com/c4urself/hackathon/feeders"
	"github.com/c4urself/hackathon/mosaic"
	"github.com/gin-gonic/gin"
	"html/template"
)

const BASE_TMPL = "templates/base.tmpl"

func StartApp() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

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
		username := c.Params.ByName("username")
		mosaics := mosaic.MakeInstagramMosaic(
			username,
			fmt.Sprintf("/tmp/hack/%s/photos/", username),
			fmt.Sprintf("/tmp/hack/%s/audience/", username),
			fmt.Sprintf("./static/mosaic/%s/", username))

		c.HTML(200, "result.tmpl", gin.H{
			"username": c.Params.ByName("username"),
			"baseUrl":  fmt.Sprintf("/static/mosaic/%s/", username),
			"mosaics": mosaics})
	})

	r.Run(":8080")
}
