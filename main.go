package main

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func main() {

// 	router := gin.Default()

// 	router.GET("/", func(c *gin.Context) {
// 		c.String(http.StatusOK, "Hello World")
// 	})

// 	router.GET("/user/:name", func(c *gin.Context) {
// 		name := c.Param("name")
// 		c.String(http.StatusOK, "Hello %s", name)
// 	})
// 	router.Run(":8000")
// }

import (
	"rutgo/serve"
)

func main() {
	api := serve.NewServe()
	// api.Use(rest.DefaultDevStack...)
	api.SetApp(serve.AppSimple(func(w serve.ResponseWriter, r *serve.Request) {
		w.WriteJson(map[string]string{"Body": "Hello World!"})
	}))
	api.Run(":5000")
}
