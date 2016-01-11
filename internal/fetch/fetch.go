package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

func FetchMemStats(target string) (*runtime.MemStats, error) {
	resp, err := http.DefaultClient.Get(fmt.Sprintf("%s/debug/memstats", target))
	if err != nil {
		return nil, err
	}

	var stats runtime.MemStats
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&stats)

	return &stats, nil
}
