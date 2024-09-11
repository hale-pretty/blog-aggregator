package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hale-pretty/internal/database"
)

func scraper(db *database.Queries, concurrency int, interval time.Duration) {

	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		wg := &sync.WaitGroup{}
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error getting next feeds to fetch: %v", err)
			return
		}
		for _, feed := range feeds {
			wg.Add(1)
			go scraping(db, wg, feed)
		}
		wg.Wait()
	}
}

func scraping(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}

	rssFeed, err := urltoFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed: %v", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		ti, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("cannot parse date %v with err %v", item.PubDate, err)
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Tittle,
			Url:         item.Link,
			Description: description,
			PublishedAt: ti,
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Printf("cannot create post: %v", err)
		}
	}

}
