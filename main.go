package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
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

	scopes := []string{"repo"}
	conf := oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       scopes,
		Endpoint:     oauth2github.Endpoint,
	}

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

	r.GET("/login", func(c *gin.Context) {
		ctx := appengine.NewContext(c.Request)
		log.Printf("%s", ctx)
		// TODO: stateをdbに保存してcallbackで確認する
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		log.Println(url)
		c.Header("Location", url)
		c.SecureJSON(http.StatusTemporaryRedirect, "")
	})

	r.GET("/callback", func(c *gin.Context) {
		ctx := appengine.NewContext(c.Request)
		log.Printf("%s", ctx)
		githubToken, _ := conf.Exchange(ctx, c.Query("code"))
		//　TODO これでアクセストークンが得られたので、セッションに情報を入れたい
		log.Println(githubToken.AccessToken)
		c.Header("Location", "/")
		c.SecureJSON(http.StatusTemporaryRedirect, "")
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