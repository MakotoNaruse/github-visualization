package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
	"google.golang.org/appengine"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//デバッグ用キー
	//oauthKey := os.Getenv("OAUTH_KEY")

	r := gin.New()

	//r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")

	// セッションの設定
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("github-visualization", store))

	var scopes = []string{"repo:status", "read:user"}
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
		session := sessions.Default(c)
		token := session.Get("accessToken")
		if token == nil {
			log.Printf("redirect")
			c.Redirect(http.StatusMovedPermanently, "/top")
			c.Abort()
		} else {
			log.Println("logged in")
			accessToken, _ := token.(string)
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: accessToken},
			)
			tc := oauth2.NewClient(ctx, ts)

			client := github.NewClient(tc)
			repos, _, _ := client.Repositories.List(ctx, "MakotoNaruse", nil)
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"title":   "index",
				"message": "hello world!!",
				"repos":   repos,
			})
		}
	})

	r.GET("/top", func(c *gin.Context) {
		ctx := appengine.NewContext(c.Request)
		log.Printf("%s", ctx)
		c.HTML(http.StatusOK, "landing.tmpl", gin.H{
			"title": "top",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		ctx := appengine.NewContext(c.Request)
		log.Printf("%s", ctx)
		session := sessions.Default(c)
		state := createRand()
		session.Clear()
		session.Set("state", state)
		session.Save()
		url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
		log.Println(url)
		c.Header("Location", url)
		c.SecureJSON(http.StatusTemporaryRedirect, "")
	})

	r.GET("/callback", func(c *gin.Context) {
		ctx := appengine.NewContext(c.Request)
		log.Printf("%s", ctx)
		githubToken, _ := conf.Exchange(ctx, c.Query("code"))
		log.Println(githubToken.AccessToken)
		session := sessions.Default(c)
		state := session.Get("state")
		if state == nil {
			log.Printf("redirect")
			c.Redirect(http.StatusMovedPermanently, "/top")
		} else{
			st, _ := state.(string)
			if st == c.Query("state"){
				log.Println("authorized")
				session.Clear()
				session.Set("accessToken", githubToken.AccessToken)
				session.Save()
				c.Redirect(http.StatusMovedPermanently, "/")
			} else{
				log.Printf("redirect")
				c.Redirect(http.StatusMovedPermanently, "/top")
			}
		}
	})

	r.GET("/logout", func(c *gin.Context) {
		ctx := appengine.NewContext(c.Request)
		log.Printf("%s", ctx)
		//セッションからデータを破棄する
		session := sessions.Default(c)
		session.Clear()
		log.Println("セッション破棄")
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/")
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

const (
	letters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	indexBit  = 6
	indexMask = 1<<indexBit - 1
	indexMax  = 63 / indexBit
)

func createRand() (randVal string) {
	randSource := rand.NewSource(time.Now().UnixNano())
	n := 32
	b := make([]byte, n)
	cache, remain := randSource.Int63(), indexMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSource.Int63(), indexMax
		}
		index := int(cache & indexMask)
		if index < len(letters) {
			b[i] = letters[index]
			i--
		}
		cache >>= indexBit
		remain--
	}
	randVal = string(b)
	return
}