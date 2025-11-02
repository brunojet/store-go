package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	internal "github.com/brunojet/store-go/office/internal"
	"github.com/brunojet/store-go/office/internal/dto"
)

func parsePageParams(r *http.Request) (int, int) {
	q := r.URL.Query()
	page := 1
	pageSize := 20
	if p := q.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := q.Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}
	return page, pageSize
}

// Start launches a minimal HTTP server with two endpoints:
// GET /categories and GET /category_types
func Start(addr string) error {
	catSvc, tcSvc, appsSvc, err := internal.Bootstrap()
	if err != nil {
		return err
	}

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		page, pageSize := parsePageParams(r)
		items, err := catSvc.List(context.Background(), page, pageSize)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(items)
	})

	http.HandleFunc("/category_types", func(w http.ResponseWriter, r *http.Request) {
		page, pageSize := parsePageParams(r)
		items, err := tcSvc.List(context.Background(), page, pageSize)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(items)
	})

	// /apps supports GET (list) and POST (create)
	http.HandleFunc("/apps", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			page, pageSize := parsePageParams(r)
			items, err := appsSvc.List(context.Background(), page, pageSize)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(items)
		case http.MethodPost:
			var dtoObj dto.ApplicationCreate
			if err := json.NewDecoder(r.Body).Decode(&dtoObj); err != nil {
				http.Error(w, "invalid body", http.StatusBadRequest)
				return
			}
			if err := appsSvc.Create(context.Background(), &dtoObj); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Handlers for /apps/{id}
	http.HandleFunc("/apps/", func(w http.ResponseWriter, r *http.Request) {
		// extract the path after /apps/ and ensure there are no extra segments
		rem := strings.TrimPrefix(r.URL.Path, "/apps/")
		if rem == "" {
			http.NotFound(w, r)
			return
		}
		// reject requests with additional slashes (e.g. /apps/1/foo)
		if strings.Contains(rem, "/") {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.ParseInt(rem, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			appDto, err := appsSvc.GetByID(context.Background(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if appDto == nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(appDto)
		case http.MethodPut:
			var upd dto.ApplicationUpdate
			if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
				http.Error(w, "invalid body", http.StatusBadRequest)
				return
			}
			if err := appsSvc.Update(context.Background(), id, &upd); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			if err := appsSvc.Delete(context.Background(), id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Printf("starting office API on %s", addr)
	return http.ListenAndServe(addr, nil)
}
