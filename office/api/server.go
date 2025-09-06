package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	internal "github.com/brunojet/store-go/office/internal"
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
// GET /categorias and GET /category_types
func Start(addr string) error {
	catSvc, tcSvc, err := internal.Bootstrap()
	if err != nil {
		return err
	}

	http.HandleFunc("/categorias", func(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("starting office API on %s", addr)
	return http.ListenAndServe(addr, nil)
}
