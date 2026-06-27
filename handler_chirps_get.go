package main

import (
	"net/http"
	"server/internal/database"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChripGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpId")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChrip, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirp", err)
		return
	}

	chirp := Chirp{
		ID:        dbChrip.ID,
		CreatedAt: dbChrip.CreatedAt,
		UpdatedAt: dbChrip.UpdatedAt,
		UserID:    dbChrip.UserID,
		Body:      dbChrip.Body,
	}
	respondWithJSON(w, http.StatusOK, chirp)
}

func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString == "" {
		return uuid.Nil, nil
	}

	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return uuid.Nil, err
	}

	return authorID, nil
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorID, err := authorIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		return
	}

	var dbChirps []database.Chirp

	if authorID != uuid.Nil {
		dbChirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorID)
	} else {
		dbChirps, err = cfg.db.GetChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	sortDirection := "asc"
	sortDirectionParam := r.URL.Query().Get("sort")
	if sortDirectionParam == "desc" {
		sortDirection = "desc"
	}

	chirps := []Chirp{}
	for _, dbChrip := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChrip.ID,
			CreatedAt: dbChrip.CreatedAt,
			UpdatedAt: dbChrip.UpdatedAt,
			UserID:    dbChrip.UserID,
			Body:      dbChrip.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortDirection == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, chirps)

}
