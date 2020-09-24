package main

import (
	"context"
	"fmt"
	"html/template"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/labstack/gommon/log"
	"github.com/mmcdole/gofeed"
)

// Item - Stored in datastore
// Based on github.com/mmcdole/gofeed .Item
type Item struct {
	Title           string        `datastore:",noindex"`
	Description     template.HTML `datastore:",noindex"` // Unsafe / unencoded. Input must be safe, a it is here since it comes from ascidoc(tor)
	Content         template.HTML `datastore:",noindex"` // Unsafe / unencoded. Input must be safe, a it is here since it comes from ascidoc(tor)
	Link            string        `datastore:",noindex"`
	UpdatedParsed   *time.Time
	PublishedParsed *time.Time
	// Author          *Person
	// Image           *Image
	Categories []string
	// Enclosures      []*Enclosure

	FeedDesc    string
	FeedURL     string
	ShownToUser bool
	FirstSeen   time.Time

	K *datastore.Key `datastore:"__key__"`

	// Deprecated fields
	// remove when they are no longer in ds, or they all have ShownToUser == true
	// or handle the error, see https://godoc.org/cloud.google.com/go/datastore#ErrFieldMismatch
	// Feed     string    `datastore:",noindex,omitempty"` // old feed url
	// GUID     string    `datastore:",noindex,omitempty"` // we have K if we need it
}

type ds struct {
	ctx    context.Context
	client *datastore.Client
}

func dsInit(ctx context.Context, projectid string) (ds, error) {
	log.Print("dsInit: Project: ", projectid)
	if ctx == nil {
		return ds{}, fmt.Errorf("dsInit: ctx can't be nil")
	}

	d := ds{ctx: ctx}

	// Create a datastore client, a single client which is reused for every datastore operation.
	var err error
	d.client, err = datastore.NewClient(ctx, projectid)
	if err != nil {
		return ds{}, err
	}

	return d, nil
}

func (d *ds) storeItem(feedInfo feed, item *gofeed.Item) error {

	// We assume the rss lib alays figure out a Global Unique ID,
	// but check that it's not empty anyway
	// it it's duplicated somehow it will result in lost item/article
	if item.GUID == "" {
		return fmt.Errorf("storeItem: Item.GUID can't be empty")
	}

	k := datastore.NameKey("items", item.GUID, nil)
	e := &Item{

		// GUID:            item.GUID,
		Title:           item.Title,
		Description:     template.HTML(item.Description),
		Content:         template.HTML(item.Content),
		Link:            item.Link,
		UpdatedParsed:   item.UpdatedParsed,
		PublishedParsed: item.PublishedParsed,
		Categories:      item.Categories,

		FeedDesc:    feedInfo.name,
		FeedURL:     feedInfo.url,
		FirstSeen:   time.Now(),
		ShownToUser: false,

		// Deprecated
		// ..
	}

	// TODO many at a time
	_, err := d.client.RunInTransaction(d.ctx, func(tx *datastore.Transaction) error {

		// We first check that there is no entity stored with the given key.
		var empty Item
		if err := tx.Get(k, &empty); err != datastore.ErrNoSuchEntity {
			return err // also err == nil !
		}

		// If there was no matching entity, store it now.
		_, err := tx.Put(k, e)
		return err
	})

	if err != nil {
		return err
	}

	// log.Printf("Updated value for feed: %s , item GUID: %q,  FirstSeen: %s", feedInfo, item.GUID, e.FirstSeen)
	return nil
}

func (d *ds) setSeenMany(guids []string) error {

	var mutations []*datastore.Mutation

	for _, guid := range guids {
		e := &Item{ShownToUser: true}
		mutations = append(mutations, datastore.NewUpdate(datastore.NameKey("items", guid, nil), e))
	}

	_, err := d.client.RunInTransaction(d.ctx, func(tx *datastore.Transaction) error {
		_, err := tx.Mutate(mutations...)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (d *ds) getByGUID(guid string) (Item, bool, error) {
	k := datastore.NameKey("items", guid, nil)
	var i Item
	err := d.client.Get(d.ctx, k, &i)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return i, false, nil
		}
		return i, false, err
	}
	return i, true, nil
}

func (d *ds) getAll() ([]Item, error) {
	q := datastore.NewQuery("items").Order("-FirstSeen")

	var items []Item
	_, err := d.client.GetAll(d.ctx, q, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (d *ds) getUnseen(cnt int) ([]Item, error) {

	q := datastore.NewQuery("items").Filter("ShownToUser =", false).Limit(cnt).Order("FirstSeen")

	var items []Item
	_, err := d.client.GetAll(d.ctx, q, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (d *ds) cntUnseen(limit int) (int, error) {

	// limit to avoid the count taking too long in extreme cases
	q := datastore.NewQuery("items").Filter("ShownToUser =", false).Limit(limit)

	cnt, err := d.client.Count(d.ctx, q)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
