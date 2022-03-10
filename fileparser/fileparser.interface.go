package fileparser

import (
	"io"

	"github.com/MichalMitros/feed-parser/models"
)

type FeedFileParserInterface interface {
	ParseFile(feedFile io.ReadCloser, shopItemsOutput chan models.ShopItem)
}
