package xmlparser

import (
	"encoding/xml"
	"io"

	"github.com/MichalMitros/feed-parser/models"
)

// Parser for parsing feed files from XML to objects
type XmlFeedParser struct{}

// Creates new XmlFeedParser instance
func NewXmlFeedParser() *XmlFeedParser {
	return &XmlFeedParser{}
}

// Parses xml file and send shop items to shopItemsOutput channel
// Closes the channel when finished
func (p *XmlFeedParser) ParseFile(feedXmlFile io.ReadCloser, shopItemsOutput chan models.ShopItem) {
	// Close channel when finished parsing
	defer close(shopItemsOutput)

	decoder := xml.NewDecoder(feedXmlFile)

	for {
		// Get next xml token
		t, _ := decoder.Token()
		// Break when file is fully processed
		if t == nil {
			break
		}
		// When token is a xml start element...
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "SHOPITEM" {
				// Parse single ShopItem and send results to output channel
				var item models.ShopItem
				decoder.DecodeElement(&item, &se)
				shopItemsOutput <- item
			}
		}
	}
}
