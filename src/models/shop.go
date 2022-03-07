package models

import "encoding/xml"

type Shop struct {
	XMLName   xml.Name   `xml:"SHOP"`
	Text      string     `xml:",chardata"`
	ShopItems []ShopItem `xml:"SHOPITEM"`
}
