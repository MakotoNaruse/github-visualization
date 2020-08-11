package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"google.golang.org/appengine"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	oauthKey := os.Getenv("OAUTH_KEY")

	r := gin.New()

	//r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		ctx := appengine.NewContext(c.Request)
		log.Printf("%s", ctx)
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: oauthKey},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)
		repos, _, _ := client.Repositories.List(ctx, "MakotoNaruse", nil)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "test",
			"message":"hello world!!",
			"repos": repos,
		})
	})
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}