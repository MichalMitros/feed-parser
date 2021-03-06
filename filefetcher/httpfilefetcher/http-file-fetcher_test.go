package httpfilefetcher

import (
	"encoding/xml"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/MichalMitros/feed-parser/models"
)

func TestFetchFile(t *testing.T) {
	client := MockedHttpClient{}
	filesFetcher := NewHttpFileFetcher(
		client,
	)
	result, _, err := filesFetcher.FetchFile("some_test_url")

	// Check error
	if err != nil {
		t.Fatalf(`FetchFiles(string), err = %v, want nil`, err)
	}

	// Check if returned io.ReadCloser isn't changed
	if !reflect.DeepEqual(mockedReadCloser, *result) {
		t.Fatalf(
			"HttpFileFetcher.FetchFiles(string), should return unchanged io.ReadCloser",
		)
	}
}

func TestDefaultHttpFileFetcher(t *testing.T) {
	// Create default fetcher
	fetcher := DefaultHttpFileFetcher()

	// Check if returned io.ReadCloser isn't changed
	if !reflect.DeepEqual(fetcher.httpClient, http.DefaultClient) {
		t.Fatalf(
			"DefaultHttpFileFetcher(), should return fetcher with http.DefaultClient",
		)
	}
}

func TestNewHttpFileFetcher(t *testing.T) {
	// Create default fetcher
	fetcher := NewHttpFileFetcher(http.DefaultClient)

	// Check if returned io.ReadCloser isn't changed
	if !reflect.DeepEqual(fetcher.httpClient, http.DefaultClient) {
		t.Fatalf(
			"DefaultHttpFileFetcher(), should return fetcher with HttpClientInterface from the input",
		)
	}
}

// MOCKED DATA

// Mocked http.Client as struct implementing FileFetcher interface
type MockedHttpClient struct{}

func (c MockedHttpClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       mockedReadCloser,
	}, nil
}

// io.ReadCloser with mockedCorrectShop
var mockedXmlFileBytes, _ = xml.Marshal(mockedCorrectShop)
var mockedReadCloser = io.NopCloser(strings.NewReader(string(mockedXmlFileBytes)))

// Correct shop with 3 items
var mockedCorrectShop models.Shop = models.Shop{
	ShopItems: []models.ShopItem{
		{
			XMLName:           xml.Name{Local: "SHOPITEM"},
			ItemID:            "testId_1",
			ProductName:       "testProductName_1",
			Product:           "testProduct_1",
			Description:       "testDescription_1",
			Url:               "testUrl_1",
			ImgUrl:            "testImgUrl_1",
			ImgUrlAlternative: "testImgUrlAlt_1",
			VideoUrl:          "testVideoUrl_1",
			PriceVat:          "testPriceVat_1",
			HeurekaCPC:        "testHeurekaCPC_1",
			CategoryText:      "testCategoryText_1",
			EAN:               "testEAN_1",
			ProductNo:         "testProductNo_1",
			Params: []models.ShopItemParam{
				{XMLName: xml.Name{Local: "PARAM"}, ParamName: "testParam_1_1", Val: "testVal_1_1"},
				{XMLName: xml.Name{Local: "PARAM"}, ParamName: "testParam_1_2", Val: "testVal_1_2"},
			},
			DelivaryDate: "testDelDate_1",
			Deliveries: []models.ShopItemDelivery{
				{XMLName: xml.Name{Local: "DELIVERY"}, DeliveryID: "testDelId_1_1", DeliveryPrice: "testDelPrice_1_1", DeliveryPriceCOD: "testDelPriceCOD_1_1"},
				{XMLName: xml.Name{Local: "DELIVERY"}, DeliveryID: "testDelId_1_2", DeliveryPrice: "testDelPrice_1_2", DeliveryPriceCOD: "testDelPriceCOD_1_2"},
			},
			ItemGroupId: "testGroupId_1",
			Accessory:   "testAccesory_1",
			Gift:        "testGift_1",
			ExtendedWarranty: []models.ShopItemExtendedWarranty{
				{XMLName: xml.Name{Local: "EXTENDED_WARRANTY"}, Val: "testWarrantyVal_1_1", Desc: "testWarrantyDesc_1_1"},
				{XMLName: xml.Name{Local: "EXTENDED_WARRANTY"}, Val: "testWarrantyVal_1_2", Desc: "testWarrantyDesc_1_2"},
				{XMLName: xml.Name{Local: "EXTENDED_WARRANTY"}, Val: "testWarrantyVal_1_3", Desc: "testWarrantyDesc_1_3"},
			},
			SpecialService: "testSpecSvc_1",
			SalesVoucher: []models.ShopItemSalesVoucher{
				{XMLName: xml.Name{Local: "SALES_VOUCHER"}, Code: "testVoucherCode_1_1", Desc: "testVoucherDesc_1_1"},
			},
		},
		{
			XMLName:           xml.Name{Local: "SHOPITEM"},
			ItemID:            "testId_2",
			ProductName:       "testProductName_2",
			Product:           "testProduct_2",
			Description:       "testDescription_2",
			Url:               "testUrl_2",
			ImgUrl:            "testImgUrl_2",
			ImgUrlAlternative: "testImgUrlAlt_2",
			VideoUrl:          "testVideoUrl_2",
			PriceVat:          "testPriceVat_2",
			HeurekaCPC:        "testHeurekaCPC_2",
			CategoryText:      "testCategoryText_2",
			EAN:               "testEAN_2",
			ProductNo:         "testProductNo_2",
			Params: []models.ShopItemParam{
				{XMLName: xml.Name{Local: "PARAM"}, ParamName: "testParam_2_1", Val: "testVal_2_1"},
				{XMLName: xml.Name{Local: "PARAM"}, ParamName: "testParam_2_2", Val: "testVal_2_2"},
			},
			DelivaryDate: "testDelDate_2",
			Deliveries: []models.ShopItemDelivery{
				{XMLName: xml.Name{Local: "DELIVERY"}, DeliveryID: "testDelId_2_1", DeliveryPrice: "testDelPrice_2_1", DeliveryPriceCOD: "testDelPriceCOD_2_1"},
				{XMLName: xml.Name{Local: "DELIVERY"}, DeliveryID: "testDelId_2_2", DeliveryPrice: "testDelPrice_2_2", DeliveryPriceCOD: "testDelPriceCOD_2_2"},
			},
			ItemGroupId: "testGroupId_2",
			Accessory:   "testAccesory_2",
			Gift:        "testGift_2",
			ExtendedWarranty: []models.ShopItemExtendedWarranty{
				{XMLName: xml.Name{Local: "EXTENDED_WARRANTY"}, Val: "testWarrantyVal_2_1", Desc: "testWarrantyDesc_2_1"},
				{XMLName: xml.Name{Local: "EXTENDED_WARRANTY"}, Val: "testWarrantyVal_2_2", Desc: "testWarrantyDesc_2_2"},
				{XMLName: xml.Name{Local: "EXTENDED_WARRANTY"}, Val: "testWarrantyVal_2_3", Desc: "testWarrantyDesc_2_3"},
			},
			SpecialService: "testSpecSvc_2",
			SalesVoucher: []models.ShopItemSalesVoucher{
				{XMLName: xml.Name{Local: "SALES_VOUCHER"}, Code: "testVoucherCode_2_1", Desc: "testVoucherDesc_2_1"},
			},
		},
		{
			XMLName:           xml.Name{Local: "SHOPITEM"},
			ItemID:            "testId_3",
			ProductName:       "testProductName_3",
			Product:           "testProduct_3",
			Description:       "testDescription_3",
			Url:               "testUrl_3",
			ImgUrl:            "testImgUrl_3",
			ImgUrlAlternative: "testImgUrlAlt_3",
			VideoUrl:          "testVideoUrl_3",
			PriceVat:          "testPriceVat_3",
			HeurekaCPC:        "testHeurekaCPC_3",
			CategoryText:      "testCategoryText_3",
			EAN:               "testEAN_3",
			ProductNo:         "testProductNo_3",
			Params: []models.ShopItemParam{
				{XMLName: xml.Name{Local: "PARAM"}, ParamName: "testParam_3_1", Val: "testVal_3_1"},
				{XMLName: xml.Name{Local: "PARAM"}, ParamName: "testParam_3_2", Val: "testVal_3_2"},
			},
			DelivaryDate: "testDelDate_3",
			Deliveries: []models.ShopItemDelivery{
				{XMLName: xml.Name{Local: "DELIVERY"}, DeliveryID: "testDelId_3_1", DeliveryPrice: "testDelPrice_3_1", DeliveryPriceCOD: "testDelPriceCOD_3_1"},
				{XMLName: xml.Name{Local: "DELIVERY"}, DeliveryID: "testDelId_3_2", DeliveryPrice: "testDelPrice_3_2", DeliveryPriceCOD: "testDelPriceCOD_3_2"},
			},
			ItemGroupId: "testGroupId_3",
			Accessory:   "testAccesory_3",
			Gift:        "testGift_3",
			ExtendedWarranty: []models.ShopItemExtendedWarranty{
				{XMLName: xml.Name{Local: "EXTENDED_WARRANTY"}, Val: "testWarrantyVal_3_1", Desc: "testWarrantyDesc_3_1"},
				{XMLName: xml.Name{Local: "EXTENDED_WARRANTY"}, Val: "testWarrantyVal_3_2", Desc: "testWarrantyDesc_3_2"},
				{XMLName: xml.Name{Local: "EXTENDED_WARRANTY"}, Val: "testWarrantyVal_3_3", Desc: "testWarrantyDesc_3_3"},
			},
			SpecialService: "testSpecSvc_3",
			SalesVoucher: []models.ShopItemSalesVoucher{
				{XMLName: xml.Name{Local: "SALES_VOUCHER"}, Code: "testVoucherCode_3_1", Desc: "testVoucherDesc_3_1"},
			},
		},
		{
			XMLName:           xml.Name{Local: "SHOPITEM"},
			ItemID:            "testId_3",
			ProductName:       "testProductName_3",
			Product:           "testProduct_3",
			Description:       "testDescription_3",
			Url:               "testUrl_3",
			ImgUrl:            "testImgUrl_3",
			ImgUrlAlternative: "testImgUrlAlt_3",
			VideoUrl:          "testVideoUrl_3",
			PriceVat:          "testPriceVat_3",
			HeurekaCPC:        "testHeurekaCPC_3",
			CategoryText:      "testCategoryText_3",
			EAN:               "testEAN_3",
			ProductNo:         "testProductNo_3",
			Params: []models.ShopItemParam{
				{XMLName: xml.Name{Local: "PARAM"}, ParamName: "testParam_3_1", Val: "testVal_3_1"},
				{XMLName: xml.Name{Local: "PARAM"}, ParamName: "testParam_3_2", Val: "testVal_3_2"},
			},
			DelivaryDate: "testDelDate_3",
			Deliveries: []models.ShopItemDelivery{
				{XMLName: xml.Name{Local: "DELIVERY"}, DeliveryID: "testDelId_3_1", DeliveryPrice: "testDelPrice_3_1", DeliveryPriceCOD: "testDelPriceCOD_3_1"},
				{XMLName: xml.Name{Local: "DELIVERY"}, DeliveryID: "testDelId_3_2", DeliveryPrice: "testDelPrice_3_2", DeliveryPriceCOD: "testDelPriceCOD_3_2"},
			},
			ItemGroupId: "testGroupId_3",
			Accessory:   "testAccesory_3",
			Gift:        "testGift_3",
			ExtendedWarranty: []models.ShopItemExtendedWarranty{
				{XMLName: xml.Name{Local: "EXTENDED_WARRANTY"}, Val: "testWarrantyVal_3_1", Desc: "testWarrantyDesc_3_1"},
				{XMLName: xml.Name{Local: "EXTENDED_WARRANTY"}, Val: "testWarrantyVal_3_2", Desc: "testWarrantyDesc_3_2"},
				{XMLName: xml.Name{Local: "EXTENDED_WARRANTY"}, Val: "testWarrantyVal_3_3", Desc: "testWarrantyDesc_3_3"},
			},
			SpecialService: "testSpecSvc_3",
			SalesVoucher: []models.ShopItemSalesVoucher{
				{XMLName: xml.Name{Local: "SALES_VOUCHER"}, Code: "testVoucherCode_3_1", Desc: "testVoucherDesc_3_1"},
			},
		},
	},
}
