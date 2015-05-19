package web

import (
	"fmt"
	"github.com/c4urself/hackathon/mosaic"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

const BASE_TMPL = "templates/base.tmpl"

func StartApp() {
	r := gin.Default()
	r.Static("/static", "./static")

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

		message := "Hello " + username
		c.String(http.StatusOK, message)
	})

	// Result page
	r.GET("/result", func(c *gin.Context) {
		c.String(http.StatusOK, "arst")
	})

	// Mosaic page
	r.GET("/mosaic/:username", func(c *gin.Context) {
		username := c.Params.ByName("username")
		mosaics := mosaic.MakeInstagramMosaic(
			username,
			fmt.Sprintf("/tmp/hack/%s/photos/", username),
			fmt.Sprintf("/tmp/hack/%s/audience/", username),
			fmt.Sprintf("./static/mosaic/%s/", username))

		c.JSON(200, gin.H{
			"username": username,
			"baseUrl": fmt.Sprintf("/static/mosaic/%s/", username),
			"mosaics": mosaics})
	})

	r.Run(":8080")
}
