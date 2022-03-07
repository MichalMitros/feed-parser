package models

import "encoding/xml"

type ShopItem struct {
	XMLName           xml.Name                   `xml:"SHOPITEM"`
	Text              string                     `xml:",chardata"`
	ItemID            string                     `xml:"ITEM_ID"`
	ProductName       string                     `xml:"PRODUCTNAME"`
	Product           string                     `xml:"PRODUCT"`
	Description       string                     `xml:"DESCRIPTION"`
	Url               string                     `xml:"URL"`
	ImgUrl            string                     `xml:"IMGURL"`
	ImgUrlAlternative string                     `xml:"IMGURL_ALTERNATIVE"`
	VideoUrl          string                     `xml:"VIDEO_URL"`
	PriceVat          string                     `xml:"PRICE_VAT"`
	HeurekaCPC        string                     `xml:"HEUREKA_CPC"`
	CategoryText      string                     `xml:"CATEGORYTEXT"`
	EAN               string                     `xml:"EAN"`
	ProductNo         string                     `xml:"PRODUCTNO"`
	Params            []ShopItemParam            `xml:"PARAM"`
	DelivaryDate      string                     `xml:"DELIVERY_DATE"`
	Deliveries        []ShopItemDelivery         `xml:"DELIVERY"`
	ItemGroupId       string                     `xml:"ITEMGROUP_ID"`
	Accessory         string                     `xml:"ACCESSORY"`
	Gift              string                     `xml:"GIFT"`
	ExtendedWarranty  []ShopItemExtendedWarranty `xml:"EXTENDED_WARRANTY"`
	SpecialService    string                     `xml:"SPECIAL_SERVICE"`
	SalesVoucher      []ShopItemSalesVoucher     `xml:"SALES_VOUCHER"`
}

type ShopItemParam struct {
	XMLName   xml.Name `xml:"PARAM"`
	Text      string   `xml:",chardata"`
	ParamName string   `xml:"PARAM_NAME"`
	Val       string   `xml:"VAL"`
}

type ShopItemDelivery struct {
	XMLName          xml.Name `xml:"DELIVERY"`
	Text             string   `xml:",chardata"`
	DeliveryID       string   `xml:"DELIVERY_ID"`
	DeliveryPrice    string   `xml:"DELIVERY_PRICE"`
	DeliveryPriceCOD string   `xml:"DELIVERY_PRICE_COD"`
}

type ShopItemExtendedWarranty struct {
	XMLName xml.Name `xml:"EXTENDED_WARRANTY"`
	Text    string   `xml:",chardata"`
	Val     string   `xml:"VAL"`
	Desc    string   `xml:"DESC"`
}

type ShopItemSalesVoucher struct {
	XMLName xml.Name `xml:"SALES_VOUCHER"`
	Text    string   `xml:",chardata"`
	Code    string   `xml:"CODE"`
	Desc    string   `xml:"DESC"`
}
