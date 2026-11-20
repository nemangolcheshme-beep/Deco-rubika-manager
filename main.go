
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"service": "DEKO Rubika Manager",
			"status":  "running",
		})
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true}`))
	})

	mux.HandleFunc("/rubika/status", func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(os.Getenv("RUBIKA_BOT_TOKEN"))
		if token == "" {
			http.Error(w, "Bot token is not configured", 503)
			return
		}

		body := bytes.NewBufferString("{}")
		req, err := http.NewRequest("POST",
			"https://botapi.rubika.ir/v3/"+token+"/getMe", body)
		if err != nil {
			http.Error(w, "Request creation failed", 500)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 20 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Rubika connection failed", 502)
			return
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(data)
	})

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
