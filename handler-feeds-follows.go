package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/hale-pretty/internal/database"
)

type responsedFF struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
}

func responseFF(dbff database.FeedsFollow) responsedFF {

	return responsedFF{
		ID:        dbff.ID,
		CreatedAt: dbff.CreatedAt,
		UpdatedAt: dbff.UpdatedAt,
		FeedID:    dbff.FeedID,
		UserID:    dbff.UserID,
	}
}

func responseFFs(dbffs []database.FeedsFollow) []responsedFF {
	resFFs := []responsedFF{}
	for _, dbff := range dbffs {
		resFFs = append(resFFs, responsedFF{
			ID:        dbff.ID,
			CreatedAt: dbff.CreatedAt,
			UpdatedAt: dbff.UpdatedAt,
			FeedID:    dbff.FeedID,
			UserID:    dbff.UserID,
		})
	}
	return resFFs
}

func (apiConfig *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	// Decode request body
	type ffParameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	param := ffParameters{}
	err := decoder.Decode(&param)
	if err != nil {
		respondwithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}

	// Create a feed follow
	ffParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    param.FeedID,
		UserID:    user.ID,
	}
	ff, err := apiConfig.DB.CreateFeedFollow(r.Context(), ffParams)
	if err != nil {
		respondwithError(w, 400, fmt.Sprintf("Cannot create new feed follow: %v", err))
		return
	}

	respondWithJSON(w, 201, responseFF(ff))
}

func (apiConfig *apiConfig) handlerGetFeedsFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	ffs, err := apiConfig.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondwithError(w, 400, fmt.Sprintf("Cannot get feeds follow %v", err))
		return
	}
	respondWithJSON(w, 200, responseFFs(ffs))
}

func (apiConfig *apiConfig) handlerDelFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	ffIdStr := chi.URLParam(r, "feedFollowID")
	ffId, err := uuid.Parse(ffIdStr)
	if err != nil {
		respondwithError(w, 400, fmt.Sprintf("Cannot parse feed follow id: %v", err))
		return
	}

	err = apiConfig.DB.DelFeedFollows(r.Context(), database.DelFeedFollowsParams{
		ID:     ffId,
		UserID: user.ID,
	})
	if err != nil {
		respondwithError(w, 400, fmt.Sprintf("Cannot delete feed follow %v", err))
		return
	}
	respondWithJSON(w, 200, struct{}{})
}
