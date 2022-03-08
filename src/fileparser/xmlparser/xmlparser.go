package xmlparser

import (
	"encoding/xml"

	"github.com/MichalMitros/feed-parser/models"
)

type XmlFeedParser struct {
}

func (p XmlFeedParser) ParseFile(feedXmlFile []byte) (*models.Shop, error) {
	var shop models.Shop

	if err := xml.Unmarshal(feedXmlFile, &shop); err != nil {
		return nil, err
	}

	return &shop, nil
}
