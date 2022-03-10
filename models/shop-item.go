package models

import "encoding/xml"

type ShopItem struct {
	XMLName           xml.Name                   `xml:"SHOPITEM" json:"-"`
	Text              string                     `xml:",chardata" json:"-"`
	ItemID            string                     `xml:"ITEM_ID" json:"itemId"`
	ProductName       string                     `xml:"PRODUCTNAME" json:"productName"`
	Product           string                     `xml:"PRODUCT" json:"product"`
	Description       string                     `xml:"DESCRIPTION" json:"description"`
	Url               string                     `xml:"URL" json:"url"`
	ImgUrl            string                     `xml:"IMGURL" json:"imgUrl"`
	ImgUrlAlternative string                     `xml:"IMGURL_ALTERNATIVE" json:"imgUrlAlternative"`
	VideoUrl          string                     `xml:"VIDEO_URL" json:"videoUrl"`
	PriceVat          string                     `xml:"PRICE_VAT" json:"priceVAT"`
	HeurekaCPC        string                     `xml:"HEUREKA_CPC" json:"heurekaCPC"`
	CategoryText      string                     `xml:"CATEGORYTEXT" json:"categoryText"`
	EAN               string                     `xml:"EAN" json:"ean"`
	ProductNo         string                     `xml:"PRODUCTNO" json:"productNo"`
	Params            []ShopItemParam            `xml:"PARAM" json:"param"`
	DelivaryDate      string                     `xml:"DELIVERY_DATE" json:"deliveryDate"`
	Deliveries        []ShopItemDelivery         `xml:"DELIVERY" json:"deliveries"`
	ItemGroupId       string                     `xml:"ITEMGROUP_ID" json:"itemGroupId"`
	Accessory         string                     `xml:"ACCESSORY" json:"accessory"`
	Gift              string                     `xml:"GIFT" json:"gift"`
	ExtendedWarranty  []ShopItemExtendedWarranty `xml:"EXTENDED_WARRANTY" json:"extendedWarranty"`
	SpecialService    string                     `xml:"SPECIAL_SERVICE" json:"specialService"`
	SalesVoucher      []ShopItemSalesVoucher     `xml:"SALES_VOUCHER" json:"salesVoucher"`
}

type ShopItemParam struct {
	XMLName   xml.Name `xml:"PARAM" json:"-"`
	Text      string   `xml:",chardata" json:"-"`
	ParamName string   `xml:"PARAM_NAME" json:"paramName"`
	Val       string   `xml:"VAL" json:"val"`
}

type ShopItemDelivery struct {
	XMLName          xml.Name `xml:"DELIVERY" json:"-"`
	Text             string   `xml:",chardata" json:"-"`
	DeliveryID       string   `xml:"DELIVERY_ID" json:"deliveryId"`
	DeliveryPrice    string   `xml:"DELIVERY_PRICE" json:"deliveryPrice"`
	DeliveryPriceCOD string   `xml:"DELIVERY_PRICE_COD" json:"deliveryPriceCOD"`
}

type ShopItemExtendedWarranty struct {
	XMLName xml.Name `xml:"EXTENDED_WARRANTY" json:"-"`
	Text    string   `xml:",chardata" json:"-"`
	Val     string   `xml:"VAL" json:"val"`
	Desc    string   `xml:"DESC" json:"desc"`
}

type ShopItemSalesVoucher struct {
	XMLName xml.Name `xml:"SALES_VOUCHER" json:"-"`
	Text    string   `xml:",chardata" json:"-"`
	Code    string   `xml:"CODE" json:"code"`
	Desc    string   `xml:"DESC" json:"desc"`
}
