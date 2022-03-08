package xmlparser

import (
	"testing"
)

var mockedXmlFile = []byte("<SHOP><SHOPITEM><ITEM_ID><![CDATA[1]]></ITEM_ID><PRODUCTNAME><![CDATA[Test Product 01]]></PRODUCTNAME><PRODUCT><![CDATA[Test Product 01]]></PRODUCT><DESCRIPTION><![CDATA[Test Description 01]]></DESCRIPTION><URL><![CDATA[https://www.testurl.com]]></URL><IMGURL><![CDATA[https://testimgurl.com]]></IMGURL><MANUFACTURER><![CDATA[Test Manufacturer]]></MANUFACTURER><CATEGORYTEXT><![CDATA[Test Products]]></CATEGORYTEXT><EAN><![CDATA[4716123314660]]></EAN><PRODUCTNO><![CDATA[NF-F12 PWM]]></PRODUCTNO><DELIVERY_DATE><![CDATA[2]]></DELIVERY_DATE><PRICE_VAT><![CDATA[499]]></PRICE_VAT><IMGURL_ALTERNATIVE><![CDATA[https://testimgurl.com/2]]></IMGURL_ALTERNATIVE><DELIVERY><DELIVERY_ID><![CDATA[DPD]]></DELIVERY_ID><DELIVERY_PRICE><![CDATA[99]]></DELIVERY_PRICE><DELIVERY_PRICE_COD><![CDATA[144]]></DELIVERY_PRICE_COD></DELIVERY><DELIVERY><DELIVERY_ID><![CDATA[WEDO]]></DELIVERY_ID><DELIVERY_PRICE><![CDATA[99]]></DELIVERY_PRICE><DELIVERY_PRICE_COD><![CDATA[144]]></DELIVERY_PRICE_COD></DELIVERY><PARAM><PARAM_NAME><![CDATA[param_1]]></PARAM_NAME><VAL><![CDATA[val_1]]></VAL></PARAM><PARAM><PARAM_NAME><![CDATA[param_2]]></PARAM_NAME><VAL><![CDATA[val_2]]></VAL></PARAM><HEUREKA_CPC><![CDATA[6,1248]]></HEUREKA_CPC></SHOPITEM></SHOP>")

func TestParseXmlFeed(t *testing.T) {
	parser := XmlFeedParser{}

	shop, err := parser.ParseFile(mockedXmlFile)

	if err != nil {
		t.Fatalf(`ParseFile(mockedXmlFile), err = %v, want 1`, err)
	}
	if len(shop.ShopItems) != 1 {
		t.Fatalf(`ParseFeedXml(mockedXmlFile), len(result) = %q, want 1`, len(shop.ShopItems))
	}

}
