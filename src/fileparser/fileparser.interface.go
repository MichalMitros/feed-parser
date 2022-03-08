package fileparser

import "github.com/MichalMitros/feed-parser/models"

type FeedFileParserInterface interface {
	ParseFile([]byte) (*models.Shop, error)
}
