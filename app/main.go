package main

import (
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	ginoauth2 "github.com/zalando/gin-oauth2"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	// Set up hacky proxy for Dex, to have locally proper issuer/localhost easily
	go func() {
		dexUrl, _ := url.Parse("http://dex:5556")
		proxy := httputil.NewSingleHostReverseProxy(dexUrl)
		http.ListenAndServe(":5556", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			fmt.Println("Called", request.URL.Path, "in Dex proxy")
			proxy.ServeHTTP(writer, request)
		}))
	}()

	router := gin.New()
	router.Use(ginoauth2.RequestLogger([]string{"uid"}, "data"))
	router.Use(gin.Recovery())

	router.GET("/", func(ctx *gin.Context) {
		provider, err := oidc.NewProvider(ctx, "http://localhost:5556/dex")
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, err)
			return
		}

		// Configure the OAuth2 config with the client values.
		oauth2Config := oauth2.Config{
			// client_id and client_secret of the client.
			ClientID:     "example-app",
			ClientSecret: "example-app-secret",

			// The redirectURL.
			RedirectURL: "http://localhost:9999", // Redirect URL
			Endpoint:    provider.Endpoint(),
			Scopes:      []string{oidc.ScopeOpenID, "profile", "email"},
		}

		//idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: "example-app"})

		state := fmt.Sprintf("http://%s/finish", ctx.Request.Host)

		ctx.Redirect(http.StatusFound, strings.Replace(oauth2Config.AuthCodeURL(state), "http://dex", "http://localhost", 1))
	})

	router.GET("/finish", func(ctx *gin.Context) {
		ctx.JSON(200, "Hello!")
	})

	router.Run(":80")
}
