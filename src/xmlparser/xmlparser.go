package xmlparser

import (
	"encoding/xml"

	"github.com/MichalMitros/feed-parser/models"
)

type XmlParserInterface interface {
	ParseXML(string)
}

type XmlParser struct {
}

func (p *XmlParser) ParseFeedXml(feedXmlFile []byte) models.Shop {
	var shop models.Shop

	xml.Unmarshal(feedXmlFile, &shop)

	return shop
}
