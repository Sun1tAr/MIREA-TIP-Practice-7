package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/pz7-redis/internal/cache"
)

func main() {
    redisPassword := os.Getenv("REDIS_PASS")
    redisIp := os.Getenv("REDIS_IP")

    if redisPassword == "" || redisIp == "" {
      log.Fatal("REDIS_PASS and REDIS_IP environment variable required")
      return
    } else {
      log.Printf("%s", "REDIS_PASS=" + redisPassword + " redisIp=" + redisIp)
    }

    c := cache.New(redisIp + ":6379", redisPassword)

    mux := http.NewServeMux()

    mux.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
        key := r.URL.Query().Get("key")
        value := r.URL.Query().Get("value")
        if key == "" || value == "" {
            http.Error(w, "key and value required", http.StatusBadRequest)
            return
        }
        err := c.Set(key, value, 10*time.Second) // TTL = 10 сек
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(w, "OK: %s=%s (TTL 10s)", key, value)
    })

    mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
        key := r.URL.Query().Get("key")
        val, err := c.Get(key)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        fmt.Fprintf(w, "VALUE: %s=%s", key, val)
    })

    mux.HandleFunc("/ttl", func(w http.ResponseWriter, r *http.Request) {
        key := r.URL.Query().Get("key")
        ttl, err := c.TTL(key)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(w, "TTL for %s: %v", key, ttl)
    })

    log.Println("Listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
