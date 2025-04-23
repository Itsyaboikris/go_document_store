package replication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func Replicate(peers []string, projectID string, collection string, id string, doc map[string]interface{}) {
	println("I was called")

	replicationData := map[string]interface{}{
		"project":    projectID,
		"collection": collection,
		"id":         id,
		"data":       doc["data"],
		"created_at": doc["created_at"],
		"updated_at": doc["updated_at"],
	}

	for _, peer := range peers {
		go func(url string) {
			maxRetries := 3

			for i := 0; i < maxRetries; i++ {
				if err := replicateToPeer(url, replicationData); err != nil {
					log.Printf("Failed to replicate to %s (attempt %d/%d): %v", url, i+1, maxRetries, err)
					time.Sleep(time.Second * time.Duration(i+1))
					continue
				}
				log.Printf("Successfully replicated to %s", url)
				return
			}

		}(peer)
	}
}

func replicateToPeer(peer string, data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	resp, err := http.Post("http://"+peer+"/replicate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK status: %d", resp.StatusCode)
	}

	return nil
}
