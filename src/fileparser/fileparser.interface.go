package fileparser

import (
	"io"

	"github.com/MichalMitros/feed-parser/models"
)

type FeedFileParserInterface interface {
	ParseFile(io.ReadCloser)
	GetOutputChannel() chan models.ShopItem
}
