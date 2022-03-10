package models

import "encoding/xml"

type Shop struct {
	XMLName   xml.Name   `xml:"SHOP" json:"-"`
	Text      string     `xml:",chardata" json:"-"`
	ShopItems []ShopItem `xml:"SHOPITEM" json:"shopItems"`
}
