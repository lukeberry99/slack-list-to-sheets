package handlers

import (
	"log"
	"net/http"

	"github.com/lukeberry99/slack-list-to-sheets/download"
	"github.com/lukeberry99/slack-list-to-sheets/extractor"
	"github.com/lukeberry99/slack-list-to-sheets/slack"
)

type CSVHandler struct {
	slackToken string
}

func NewCSVHandler(slackToken string) *CSVHandler {
	return &CSVHandler{
		slackToken: slackToken,
	}
}

func (h *CSVHandler) HandleGetCSV(w http.ResponseWriter, r *http.Request) {
	if h.slackToken == "" {
		http.Error(w, "SLACK_TOKEN is not set", http.StatusInternalServerError)
		return
	}

	fileId := r.URL.Query().Get("fileId")
	if fileId == "" {
		http.Error(w, "fileId parameter is required", http.StatusBadRequest)
		return
	}

	log.Printf("Getting file with ID: %s", fileId)

	slackClient := slack.NewClient(h.slackToken)
	fileInfo, err := slackClient.GetFileInfo(fileId)
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		http.Error(w, "Error getting file info", http.StatusInternalServerError)
		return
	}

	output, err := download.DownloadFile(fileInfo.URLPrivate, h.slackToken)
	if err != nil {
		log.Printf("Error downloading file: %v", err)
		http.Error(w, "Error downloading file", http.StatusInternalServerError)
		return
	}

	csvData, err := extractor.ConvertJSONToCSV(output)
	if err != nil {
		log.Printf("Error converting JSON to CSV: %v", err)
		http.Error(w, "Error converting JSON to CSV", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Write([]byte(csvData))

	log.Println("File downloaded successfully")
}
