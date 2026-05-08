package siem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func Emit(event any, endpoint, token string) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshalling SIEM event: %w", err)
	}
	fmt.Println(string(data))

	if endpoint == "" {
		return nil
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "siem: building request: %v\n", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "siem: POST failed: %v\n", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "siem: unexpected status %d\n", resp.StatusCode)
	}
	return nil
}
