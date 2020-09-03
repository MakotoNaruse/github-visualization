package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
	"google.golang.org/appengine"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var scopes = []string{"repo", "read:repo_hook", "read:user"}
var conf = oauth2.Config{}

func add(a float64, b float64) float64 {
	switch a {
	case 1:
		a++
	case 2:
		a++
	case 3:
		a++
	case 4:
		a++
	case 5:
		a++
	case 6:
		a++
	case 7:
		a++
	case 8:
		a++
	case 9:
		if a*10 == 0 {
			a++
		}
	case 10:
		a++
	}
	return a + b
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	conf = oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       scopes,
		Endpoint:     oauth2github.Endpoint,
	}

	r := gin.New()
	r.SetFuncMap(template.FuncMap{
		"add": add,
	})
	r.Use()

	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")

	// セッションの設定
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("github-visualization", store))

	r.GET("/", index)
	r.GET("/top", top)
	r.GET("/login", login)
	r.GET("/callback", callback)
	r.GET("/logout", logout)

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

func login(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	log.Printf("%s", ctx)
	session := sessions.Default(c)
	state := createRand()
	session.Clear()
	session.Save()
	session.Set("state", state)
	session.Save()
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	log.Println(url)
	c.Header("Location", url)
	c.SecureJSON(http.StatusTemporaryRedirect, "")
}

func index(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	log.Printf("%s", ctx)
	session := sessions.Default(c)
	token := session.Get("accessToken")
	if token == nil {
		fmt.Printf("redirect")
		c.Redirect(http.StatusMovedPermanently, "/top")
		c.Abort()
		return
	}
	accessToken, _ := token.(string)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	query, err := ioutil.ReadFile("query.txt")
	if err != nil {
		log.Fatal(err)
	}
	req, _ := client.NewRequest("POST", "/graphql", gin.H{"query": string(query)})
	resp, _ := tc.Do(req)
	result, _ := ioutil.ReadAll(resp.Body)
	var githubWrap GithubWrap
	if err := json.Unmarshal(result, &githubWrap); err != nil {
		log.Fatal(err)
	}
	if githubWrap.GithubData == nil {
		c.Redirect(http.StatusMovedPermanently, "/logout")
		c.Abort()
	} else {
		visualizer := Visualize(*githubWrap.GithubData)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":   "index",
			"message": "hello world!!",
			"data":    visualizer,
		})
	}
}

func top(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	log.Printf("%s", ctx)
	c.HTML(http.StatusOK, "landing.tmpl", gin.H{
		"title": "top",
	})
}

func callback(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	log.Printf("%s", ctx)
	githubToken, error := conf.Exchange(ctx, c.Query("code"))
	if error != nil {
		c.Redirect(http.StatusMovedPermanently, "/logout")
		c.Abort()
		return
	}
	session := sessions.Default(c)
	state := session.Get("state")
	if state == nil {
		c.Redirect(http.StatusMovedPermanently, "/logout")
		c.Abort()
		return
	}
	st, _ := state.(string)
	if st == c.Query("state") {
		session.Clear()
		session.Save()
		session.Set("accessToken", githubToken.AccessToken)
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/")
		c.Abort()
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/logout")
	c.Abort()
}

func logout(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	log.Printf("%s", ctx)
	//セッションからデータを破棄する
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusMovedPermanently, "/")
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
