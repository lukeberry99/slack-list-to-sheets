package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lukeberry99/slack-list-to-sheets/download"
	"github.com/lukeberry99/slack-list-to-sheets/extractor"
	"github.com/lukeberry99/slack-list-to-sheets/slack"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Start the HTTP server
	http.HandleFunc("/get-file", handleGetCSV)
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGetCSV(w http.ResponseWriter, r *http.Request) {

	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		http.Error(w, "SLACK_TOKEN is not set", http.StatusInternalServerError)
		return
	}

	slackClient := slack.NewClient(token)
	fileId := r.URL.Query().Get("fileId")

	log.Printf("Getting file with ID: %s", fileId)

	if fileId == "" {
		http.Error(w, "fileId parameter is required", http.StatusBadRequest)
		return
	}

	fileInfo, err := slackClient.GetFileInfo(fileId)
	if err != nil {
		log.Fatalf("Error getting file info: %v", err)
		http.Error(w, "Error getting file info", http.StatusInternalServerError)
		return
	}

	output, err := download.DownloadFile(fileInfo.URLPrivate, token)
	if err != nil {
		log.Fatalf("Error downloading file: %v", err)
		http.Error(w, "Error downloading file", http.StatusInternalServerError)
		return
	}

	csvData, err := extractor.ConvertJSONToCSV(output)
	if err != nil {
		log.Fatalf("Error converting JSON to CSV: %v", err)
		http.Error(w, "Error converting JSON to CSV", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Write([]byte(csvData))

	log.Println("File downloaded successfully")
}
