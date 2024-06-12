package app

import (
	"encoding/json"
	"fmt"
	"gredis/internal/cache"
	"gredis/internal/config"
	"gredis/internal/db"
	"gredis/pkg/logging"
	"math/rand"
	"net/http"
	"strconv"
)

type DbService interface {
	GetArticleById(id string) db.Article
}

func getTrending(db *db.DB, rc *cache.RedisClient, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := rand.Intn(100)
		ids := strconv.Itoa(id)

		// Check if the result is in the Redis cache
		cachedResult, err := rc.Get(ids)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			logger.Error(err)
		}
		if cachedResult != "" {
			// If the result is in the cache, return it
			fmt.Fprint(w, cachedResult)
		} else {
			// If the result is not in the cache, query the database
			a := db.GetArticleById(string(ids))
			resp, err := json.Marshal(a)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				logger.Error(err)
			}

			// Store the result in the Redis cache
			err = rc.Set(ids, string(resp))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				logger.Error(err)
			}

			// Return the result
			fmt.Fprint(w, string(resp))
		}
	}
}

func getArticlesById(db DbService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		a := db.GetArticleById(id)
		resp, err := json.Marshal(a)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		fmt.Fprintf(w, string(resp))
	}
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome")
}

func StartApp(cfg config.Config, logger logging.Logger, db *db.DB, rc *cache.RedisClient) {
	// Setup server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /articles/{id}", getArticlesById(db))
	mux.HandleFunc("GET /trending", getTrending(db, rc, logger))
	mux.HandleFunc("GET /", welcome)
	logger.Infof("server is working on %s:%s", cfg.Server.Host, cfg.Server.Port)
	err := http.ListenAndServe(cfg.Server.Port, mux)
	if err != nil {
		logger.Fatal(err)
	}
}
