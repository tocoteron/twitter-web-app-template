package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type TwitterOAuth1Context struct {
	echo.Context
	AuthConfig *oauth1.Config
}

func TwitterOAuth1Middleware(twitterOAuth1Config *oauth1.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authContext := &TwitterOAuth1Context{
				Context:    c,
				AuthConfig: twitterOAuth1Config,
			}
			return next(authContext)
		}
	}
}

func TwitterOAuth1Handler(c echo.Context) error {
	authContext := c.(*TwitterOAuth1Context)
	authConfig := authContext.AuthConfig

	requestToken, requestSecret, err := authConfig.RequestToken()
	if err != nil {
		return err
	}
	fmt.Println("RequestToken", requestToken)
	fmt.Println("RequestSecret", requestSecret)

	authorizationURL, err := authConfig.AuthorizationURL(requestToken)
	if err != nil {
		return err
	}
	fmt.Println(authorizationURL)

	return c.Redirect(http.StatusFound, authorizationURL.String())
}

func TwitterOAuth1CallbackHandler(c echo.Context) error {
	authContext := c.(*TwitterOAuth1Context)
	authConfig := authContext.AuthConfig

	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(c.Request())
	if err != nil {
		return err
	}
	fmt.Println("RequestToken", requestToken)
	fmt.Println("Verifier", verifier)

	accessToken, accessSecret, err := authConfig.AccessToken(requestToken, "", verifier)
	if err != nil {
		return err
	}
	fmt.Println("AccessToken", accessToken)
	fmt.Println("AccessSecret", accessSecret)

	token := oauth1.NewToken(accessToken, accessSecret)
	fmt.Println("NewToken", token)

	return c.JSON(http.StatusOK, nil)
}

func main() {
	twitterConsumerKey := os.Getenv("TWITTER_API_KEY")
	twitterConsumerSecret := os.Getenv("TWITTER_API_KEY_SECRET")

	fmt.Println("TwitterConsumerKey", twitterConsumerKey)
	fmt.Println("TWitterConsumerSecret", twitterConsumerSecret)

	twitterOAuth1Config := &oauth1.Config{
		ConsumerKey:    twitterConsumerKey,
		ConsumerSecret: twitterConsumerSecret,
		CallbackURL:    "http://localhost:8080/oauth/twitter/callback",
		Endpoint:       twitter.AuthorizeEndpoint,
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	twitterOAuth1Middleware := TwitterOAuth1Middleware(twitterOAuth1Config)
	e.GET("/oauth/twitter", TwitterOAuth1Handler, twitterOAuth1Middleware)
	e.GET("/oauth/twitter/callback", TwitterOAuth1CallbackHandler, twitterOAuth1Middleware)

	e.Logger.Fatal(e.Start(":8080"))
}
