package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
)

//go:embed public/*
var public embed.FS

type MoneroDStatus string

const (
	MoneroDStatusOnline  MoneroDStatus = "online"
	MoneroDStatusOffline MoneroDStatus = "offline"
	MoneroDStatusSyncing MoneroDStatus = "Syning"
)

type HealthCheckResponse struct {
	MoneroD MoneroDStatus `json:"monerod"`
}

type monerodInfoResponse struct {
	Offline     bool `json:"offline"`
	BusySyncing bool `json:"busy_syncing"`
	Synced      bool `json:"synchronized"`
}

func getMoneroDStatus() HealthCheckResponse {
	resp, err := http.Get("http://monero-node-monerod:18089/get_info")
	if err != nil {
		return HealthCheckResponse{
			MoneroD: MoneroDStatusOffline,
		}
	}

	if resp.StatusCode != http.StatusOK {
		return HealthCheckResponse{
			MoneroD: MoneroDStatusOffline,
		}
	}

	var respBody monerodInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return HealthCheckResponse{
			MoneroD: MoneroDStatusOffline,
		}
	}

	if respBody.BusySyncing || !respBody.Synced {
		return HealthCheckResponse{
			MoneroD: MoneroDStatusSyncing,
		}
	}

	if respBody.Offline {
		return HealthCheckResponse{
			MoneroD: MoneroDStatusOffline,
		}
	}

	return HealthCheckResponse{
		MoneroD: MoneroDStatusOnline,
	}
}

func main() {
	publicDir, err := fs.Sub(public, "public")
	if err != nil {
		log.Fatalln("Error accessing embedded directory:", err)
		return
	}

	http.Handle("/", http.FileServer(http.FS(publicDir)))

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(getMoneroDStatus())
	})

	log.Println("starting server at http://localhost:3000")
	log.Fatalln(http.ListenAndServe(":3000", nil))
}
