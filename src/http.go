package main

import (
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/mmcdole/gofeed"
)

//
//
// Custom Middleware

//
//
// Handlers

func (s server) refresh(c echo.Context) error {

	// Check THE secret - only used here, so no middleware
	// see s.refreshFeedsSecret()

	for _, f := range s.feeds {
		log.Printf("Fetching feed: %s", f)

		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(f.url)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(feed.Title)

		for _, v := range feed.Items {
			log.Printf("Title: %q UUID: %q\n", v.Title, v.GUID)
			log.Printf("%q %q\n", v.GUID, v.Title)

			err = s.ds.storeItem(f, v)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return c.String(http.StatusOK, "ok")
}

func (s server) setSeen(c echo.Context) error {

	// Read the cookie or die trying
	cookie, err := c.Cookie("rss-feed-seen-items")
	if err != nil {
		return err
	}

	// parse cookie
	var seenGUIDs []string
	for _, s := range strings.Split(cookie.Value, "||") {
		if len(s) == 0 {
			continue
		}
		seenGUIDs = append(seenGUIDs, s)
		// c.Logger().Debugf("%q", s)
	}

	// Update datastore
	err = s.ds.setSeenMany(seenGUIDs)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (s server) showUnseen(c echo.Context) error {

	// Count unseen
	cnt, err := s.ds.cntUnseen(999) // limit to avoid the count taking too long in extreme cases
	if err != nil {
		return err
	}

	// Get X of them
	items, err := s.ds.getUnseen(10)
	if err != nil {
		return err
	}

	// Set a cookie with all the feed item id's we're returning
	ids := strings.Builder{}
	for _, i := range items {
		_, err := ids.WriteString(i.K.Name + "||")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "writing cookie failed: ", err)
		}
	}

	cookie := &http.Cookie{}
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Name = "rss-feed-seen-items"
	cookie.Value = ids.String()
	c.Logger().Debugf("cookie val:\n%q\n", cookie.Value)
	c.SetCookie(cookie)

	c.Logger().Debugf("showing %d items", len(items))
	return c.Render(http.StatusOK, "show-all", struct {
		Items []Item
		Cnt   int
	}{items, cnt})
}

func (s server) showOne(c echo.Context) error {

	guid := c.QueryParam("guid")
	if guid == "" {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	c.Logger().Debugf("showOne guid: %q", guid)

	item, found, err := s.ds.getByGUID(guid)
	if err != nil {
		return err
	}
	if !found {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.Render(http.StatusOK, "show-one", item)
}

func (s server) showStatus(c echo.Context) error {

	// secret, err := s.refreshFeedsSecret()
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError)
	// }

	return c.String(
		http.StatusOK,
		"OK =)", // \n\nRefresh Secret for HTTP POST: %q", secret),
	)
}

//
//
// Templates
type tmpl struct {
	tmpls *template.Template
}

// Render: render a go template by name
// returning error is kind of useless?? It already sent 200?
// log all errors instead??
// https://github.com/labstack/echo/issues/225
func (t *tmpl) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.tmpls.ExecuteTemplate(w, name, data)
	if err != nil {
		c.Logger().Errorf("tmpl: Render: err: %s", err)
		return err
	}
	return nil
}
