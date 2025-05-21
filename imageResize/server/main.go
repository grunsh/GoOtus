package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"imageResize/internal/cache"
	"imageResize/internal/processor"
	Storage "imageResize/internal/storage"
)

func main() {
	cacheCapacity := 100
	if envCap := os.Getenv("CACHE_CAPACITY"); envCap != "" {
		if cap, err := strconv.Atoi(envCap); err == nil {
			cacheCapacity = cap
		}
	}

	var storage Storage.Storage
	var err error
	if os.Getenv("STORAGE_TYPE") == "memory" {
		storage = Storage.NewMemoryStorage()
	} else {
		storage, err = Storage.NewFileStorage("./image_cache")
		if err != nil {
			fmt.Printf("Failed to initialize file storage: %v\n", err)
			os.Exit(1)
		}
	}

	cache := cache.NewLRUCache(cacheCapacity, storage)
	processor := processor.NewImageProcessor(cache)

	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/fill/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 5 {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		width, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "Invalid width", http.StatusBadRequest)
			return
		}

		height, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "Invalid height", http.StatusBadRequest)
			return
		}

		url := strings.Join(parts[4:], "/")
		if url == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		data, contentType, err := processor.ProcessImage(r.Context(), url, width, height)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(data); err != nil {
			fmt.Printf("Failed to write response: %v\n", err)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server listening on :%s (cache capacity: %d)\n", port, cacheCapacity)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
