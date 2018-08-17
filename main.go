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
	"net/http"
	"rutgo/serve"
)

func index(w http.ResponseWriter, r *http.Request, _ serve.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

//func jsn(w http.ResponseWriter, r *http.Request, _ serve.Params) {
//	w.WriteJson(map[string]string{'geret': 'hello world'})
//}

func main() {
	api := serve.NewServe()
	// api.Use(rest.DefaultDevStack...)
	api.GET("/", index)
	// api.GET("/geret", jsn)
	api.Run(":8080")
}
