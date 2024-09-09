package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hale-pretty/internal/database"
)

type responsedFeed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

type responsedFeedFF struct {
	FeedObject responsedFeed `json:"feed_object"`
	FeedFollow responsedFF   `json:"feed_follow"`
}

func responseFeed(dbfeed database.Feed, dbff database.FeedsFollow) responsedFeedFF {

	return responsedFeedFF{
		FeedObject: responsedFeed{
			ID:        dbfeed.ID,
			CreatedAt: dbfeed.CreatedAt,
			UpdatedAt: dbfeed.UpdatedAt,
			Name:      dbfeed.Name,
			Url:       dbfeed.Url,
			UserID:    dbfeed.UserID,
		},
		FeedFollow: responsedFF{
			ID:        dbff.ID,
			CreatedAt: dbff.CreatedAt,
			UpdatedAt: dbff.UpdatedAt,
			FeedID:    dbff.FeedID,
			UserID:    dbff.UserID,
		},
	}
}

func responseFeeds(dbfeeds []database.Feed) []responsedFeed {
	resFeeds := []responsedFeed{}
	for _, dbfeed := range dbfeeds {
		resFeeds = append(resFeeds, responsedFeed{
			ID:        dbfeed.ID,
			CreatedAt: dbfeed.CreatedAt,
			UpdatedAt: dbfeed.UpdatedAt,
			Name:      dbfeed.Name,
			Url:       dbfeed.Url,
			UserID:    dbfeed.UserID,
		})
	}
	return resFeeds
}

func (apiConfig *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {

	// Decode request body
	type feedParameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	param := feedParameters{}
	err := decoder.Decode(&param)
	if err != nil {
		respondwithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}

	// Create a feed
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      param.Name,
		Url:       param.Url,
		UserID:    user.ID,
	}
	feed, err := apiConfig.DB.CreateFeed(r.Context(), feedParams)
	if err != nil {
		respondwithError(w, 400, fmt.Sprintf("Cannot create new feed: %v", err))
		return
	}

	// Create an auto feed follow
	ffParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}
	ff, err := apiConfig.DB.CreateFeedFollow(r.Context(), ffParams)
	if err != nil {
		respondwithError(w, 400, fmt.Sprintf("Cannot create new feed follow: %v", err))
		return
	}

	respondWithJSON(w, 201, responseFeed(feed, ff))
}

func (apiConfig *apiConfig) handlerGetAllFeed(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiConfig.DB.GetAllFeed(r.Context())
	if err != nil {
		respondwithError(w, 400, fmt.Sprintf("Cannot get feeds %v", err))
		return
	}
	respondWithJSON(w, 200, responseFeeds(feeds))
}
