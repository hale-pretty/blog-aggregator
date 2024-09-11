package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hale-pretty/internal/database"
)

type respPost struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"tittle"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func (apiConfig *apiConfig) handlerGetPostByUser(w http.ResponseWriter, r *http.Request, user database.User) {

	postParams := database.GetPostByUserParams{
		UserID: user.ID,
		Limit:  10,
	}
	posts, err := apiConfig.DB.GetPostByUser(r.Context(), postParams)
	if err != nil {
		respondwithError(w, 400, fmt.Sprintf("Cannot get post: %v", err))
		return
	}

	respPosts := make([]respPost, 0, len(posts))
	for _, post := range posts {
		respPosts = append(respPosts, respPost{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       post.Title,
			Url:         post.Url,
			Description: post.Description.String,
			PublishedAt: post.PublishedAt,
			FeedID:      post.FeedID,
		})
	}
	respondWithJSON(w, 200, respPosts)
}
