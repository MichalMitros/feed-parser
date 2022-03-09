package xmlparser

import (
	"encoding/xml"
	"io"

	"github.com/MichalMitros/feed-parser/models"
)

type XmlFeedParser struct {
	output chan models.ShopItem
}

func NewXmlFeedParser(
	output chan models.ShopItem,
) *XmlFeedParser {
	return &XmlFeedParser{
		output: output,
	}
}

func (p *XmlFeedParser) GetOutputChannel() chan models.ShopItem {
	return p.output
}

func (p *XmlFeedParser) ParseFile(feedXmlFile io.ReadCloser) {
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
				p.output <- item
			}
		}
	}
}
