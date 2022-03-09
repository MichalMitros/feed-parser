package xmlparser

import (
	"encoding/xml"
	"io"

	"github.com/MichalMitros/feed-parser/models"
)

type XmlFeedParser struct{}

func NewXmlFeedParser() *XmlFeedParser {
	return &XmlFeedParser{}
}

func (p *XmlFeedParser) ParseFile(feedXmlFile io.ReadCloser, shopItemsOutput chan models.ShopItem) {
	decoder := xml.NewDecoder(feedXmlFile)

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:

			if se.Name.Local == "SHOPITEM" {
				var item models.ShopItem
				decoder.DecodeElement(&item, &se)
				shopItemsOutput <- item
			}
		}
	}

	close(shopItemsOutput)
}
