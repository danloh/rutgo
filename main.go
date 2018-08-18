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
	"fmt"
	"rutgo/serve"
)

func index(w serve.ResponseWriter, r *serve.Request, _ serve.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func sayHi(w serve.ResponseWriter, r *serve.Request, _ serve.Params) {
	w.WriteJson(map[string]string{"say": "hello world"})
}

func main() {
	api := serve.NewServe()
	// api.Use(rest.DefaultDevStack...)
	// api.Handle("GET", "/", index)
	api.GET("/", index)
	api.GET("/hi", sayHi)
	api.Run(":8080")
}
