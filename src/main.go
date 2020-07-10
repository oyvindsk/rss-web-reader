package main

import (
	"bufio"
	"context"
	"crypto/subtle"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	logFormat    = `${time_rfc3339} ${level} prefix:${prefix},file:${short_file}:${line}: `
	logFormatReq = `${time_rfc3339} ${status} ${method} ${uri} ${host} ${remote_ip} id:${id},error:"${error}",latency_human:${latency_human},bytes_in:${bytes_in},bytes_out:${bytes_out}` + "\n"
)

func main() {

	// Create an echo instance and set log options
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG) // gommon
	e.Debug = true

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader(logFormat)
	}

	// Allow overwriting of the feeds.txt
	feedsFilepath := os.Getenv("RSS_FEED_FEEDSFILE")
	if feedsFilepath == "" {
		feedsFilepath = "feeds.txt"
	}

	// Load feeds to check from a file
	// if this fails we just give up
	s, err := serverFromFeeds(feedsFilepath)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Printf("Loaded %d feeds:", len(s.feeds))
	for _, f := range s.feeds {
		e.Logger.Printf("\t\t%s", f)
	}

	// Look for username and password for the server in ENV vars
	// if this fails we just give up
	s.username = os.Getenv("RSS_FEED_USERNAME")
	s.password = os.Getenv("RSS_FEED_PASSWORD")
	if s.username == "" || s.password == "" {
		e.Logger.Fatal("I need two environment variables: RSS_FEED_USERNAME and RSS_FEED_PASSWORD")
	}

	// Init a Datastore connections
	s.ds, err = dsInit(s.ctx, "rss-test-281216")
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Parse html templates from disk
	t := &tmpl{}
	t.tmpls = template.Must(template.ParseGlob("./templates/*.html"))
	e.Renderer = t

	//
	//
	// Echo routes and middleware

	// Middleware for all

	// log requests
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: logFormatReq,
	}))

	//
	// Routes with Basic Auth
	ba := e.Group("")
	ba.Use(middleware.BasicAuth(func(user, pass string, e echo.Context) (bool, error) {
		e.Echo().Logger.Debugf("BasicAuth checking: %q %q", user, pass)
		if subtle.ConstantTimeCompare([]byte(user), []byte(s.username)) == 1 && subtle.ConstantTimeCompare([]byte(pass), []byte(s.password)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	ba.GET("/", s.showUnseen)
	ba.GET("/show-one", s.showOne)
	ba.GET("/status", s.showStatus)
	ba.POST("/seen", s.setSeen) // reads ids from a cookie
	ba.POST("/refresh", s.refresh)

	//
	// Other, special, routes
	// (none ATM)

	//
	// Echo serve
	e.Logger.Fatal(e.Start(":8080"))
}

type feed struct {
	name, url string
}

func (f feed) String() string {
	return fmt.Sprintf("Name: %q,  URL: %q", f.name, f.url)
}

type server struct {

	// these are set at startup and never changed again, so no mutex

	username string
	password string
	ctx      context.Context
	ds       ds

	// The feeds we want to check
	feeds []feed

	// Status will cange, so it's protected
	// status struct {
	// 	sync.RWMutex
	// }
}

// serverFromFeeds loads feeds from a file. The format is one per line:
// name url
//
// Empty lines line starting with // are ignored
// there can't be any spaces in the name or url
func serverFromFeeds(filepath string) (server, error) {

	//
	// Init server
	// handle the non-ffeds fields first
	s := server{ctx: context.Background()}

	//
	// Feeds - Open file
	file, err := os.Open(filepath)
	if err != nil {
		return s, err
	}

	// loop the lines
	lineS := bufio.NewScanner(file)
	for lineS.Scan() {

		// Skip empty lines and lines starting with //
		if strings.HasPrefix(lineS.Text(), "//") || len(lineS.Text()) < 3 {
			continue
		}

		// result, one feed (name + url)
		var f feed

		// Scan two "words"
		wordS := bufio.NewScanner(strings.NewReader(lineS.Text()))
		wordS.Split(bufio.ScanWords)

		// feed name
		success := wordS.Scan()
		if !success {
			if err := wordS.Err(); err != nil {
				return s, fmt.Errorf("serverFromFeeds: can't parse line, err: %w", err)
			}
			return s, fmt.Errorf("serverFromFeeds: can't parse line")
		}
		f.name = wordS.Text()

		// feed url
		success = wordS.Scan()
		if !success {
			if err := wordS.Err(); err != nil {
				return s, fmt.Errorf("serverFromFeeds: can't parse line, err: %w", err)
			}
			return s, fmt.Errorf("serverFromFeeds: can't parse line")
		}
		f.url = wordS.Text()

		// Succcess!
		s.feeds = append(s.feeds, f)
	}

	// Any errors when scanning lines?
	if err := lineS.Err(); err != nil {
		return s, fmt.Errorf("serverFromFeeds: can't parse feeds file, err: %w", err)
	}

	// Succcess!
	return s, nil
}

// refreshFeedsSecret returns a secret for protecting the protect POST /refresh
// It's derived from the username+password, mostly just so we don't throw the password around to many places
// It's, by definition, not "more secure" than the username+password, and if you have those you can get this with GET /status
// Edit: found a way to use Basic Auth
// func (s *server) refreshFeedsSecret() (string, error) {
// 	if s.username == "" || s.password == "" {
// 		return "", fmt.Errorf("refreshFeedsSecret: empty username and/or password")
// 	}

// 	// use a hmac to add something "unknown" so the hash is not enough to test if you have the correct username+password
// 	// of course the source is public, so.. but even this very much overkill for this application
// 	// guess an even better way would be to generate a random value at first start, or something..
// 	mac := hmac.New(sha256.New, []byte("b0M{)m[-L&RqJag*)|1Uv:"))

// 	_, err := mac.Write([]byte(s.username + s.password))
// 	if err != nil {
// 		return "", fmt.Errorf("refreshFeedsSecret: mac: %w", err)
// 	}

// 	hmac := mac.Sum(nil)

// 	return base64.StdEncoding.EncodeToString(hmac), nil
// }
