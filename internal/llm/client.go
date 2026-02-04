package llm

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func GenerateResponse(prompt string) string {
	url := fmt.Sprintf("%s/generate", "http://localhost:11434")

	var body io.Reader

	req, err := http.NewRequest(http.MethodGet, url, body)
	if err != nil {
		log.Printf("NewRequest Error: %v", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Response Error: %v", err)
		return ""
	}
	defer resp.Body.Close()

	return prompt
}
