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
	"log"
	"rutgo/serve"
)

func index(w serve.ResponseWriter, r *serve.Request, _ serve.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func sayHi(w serve.ResponseWriter, r *serve.Request, _ serve.Params) {
	w.WriteJson(map[string]string{"say": "hello world"})
}
func hello(w serve.ResponseWriter, r *serve.Request, ps serve.Params) {
	w.WriteJson(map[string]string{"Hello": ps.ByName("name")})
}

// MidOne a middleware
type MidOne struct{}

// MiddlewareFunc is MiddleWare
func (one *MidOne) MiddlewareFunc(handler serve.HandlerFunc) serve.HandlerFunc {
	return serve.HandlerFunc(func(w serve.ResponseWriter, r *serve.Request, ps serve.Params) {
		log.Printf("first mid")
		handler(w, r, ps)
	})
}

// MidTwo a middleware
type MidTwo struct{}

// MiddlewareFunc is MiddleWare
func (two *MidTwo) MiddlewareFunc(handler serve.HandlerFunc) serve.HandlerFunc {
	return serve.HandlerFunc(func(w serve.ResponseWriter, r *serve.Request, ps serve.Params) {
		log.Printf("second mid")
		handler(w, r, ps)
	})
}

var chainM = []serve.Middleware{
	&MidOne{},
	&MidTwo{},
}

func main() {
	api := serve.NewServe()
	api.Use(&MidOne{}) // galoble mid
	// api.Handle("GET", "/", index)
	api.GET("/", index)
	midHi := serve.WrapMws(chainM, sayHi)
	api.GET("/hi", midHi)
	api.GET("/hi/:name", hello)
	api.Run(":8080")
}
