package contracts

import "github.com/MichalMitros/feed-parser/models"

type ParseFeedResponse struct {
	Statuses []models.FeedParsingResult `json:"statuses"`
}
