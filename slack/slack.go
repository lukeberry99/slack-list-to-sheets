package slack

import (
	"github.com/slack-go/slack"
)

type SlackClient struct {
	client *slack.Client
}

func NewClient(token string) *SlackClient {
	return &SlackClient{
		client: slack.New(token),
	}
}

func (s *SlackClient) GetFileInfo(fileId string) (*slack.File, error) {
	file, _, _, err := s.client.GetFileInfo(fileId, 0, 0)

	if err != nil {
		return nil, err
	}

	return file, nil
}
