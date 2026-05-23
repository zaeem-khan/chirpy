package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
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

	respondWithJSON(w, http.StatusOK, chirps)

}
