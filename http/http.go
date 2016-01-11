package http

import (
	"encoding/json"
	"net/http"
	"runtime"
)

func init() {
	http.HandleFunc("/debug/memstats", GetMemStats)
}

func GetMemStats(w http.ResponseWriter, r *http.Request) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	enc := json.NewEncoder(w)
	enc.Encode(&stats)
}
