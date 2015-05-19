package web

import (
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

	r.Run(":8080")
}
