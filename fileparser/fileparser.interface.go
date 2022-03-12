package fileparser

import (
	"io"

	"github.com/MichalMitros/feed-parser/models"
)

// File pareser used for parsing feed files from some format to objects
type FeedFileParserInterface interface {
	ParseFile(feedFile io.ReadCloser, shopItemsOutput chan models.ShopItem, errorsOutput chan error)
}
