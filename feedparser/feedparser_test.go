package feedparser

import (
	"encoding/xml"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/MichalMitros/feed-parser/filefetcher/httpfilefetcher"
	"github.com/MichalMitros/feed-parser/fileparser/xmlparser"
	"github.com/MichalMitros/feed-parser/models"
)

func TestFeedParserFunctionsCalling(t *testing.T) {
	// Prepare mocked data
	mockedFetcher := MockedFileFetcher{}
	mockedFileParser := MockedFileParser{}
	mockedWriter := MockedQueueWriter{}
	mockedFeedParser := NewFeedParser(
		&mockedFetcher,
		&mockedFileParser,
		&mockedWriter,
	)
	testUrls := []string{"test_url_1", "test_url_2"}

	// Use ParseFeed function
	mockedFeedParser.ParseFeeds(testUrls)

	// Check if all parser's building blocks has been called
	if mockedFetcher.NumOfFuncCalls != len(testUrls) {
		t.Fatalf(
			`FeedParser.ParseFeedsAsync(testUrls), number of file fetcher calls = %d, want %d`,
			mockedFetcher.NumOfFuncCalls,
			len(testUrls),
		)
	}
	if mockedFileParser.NumOfFuncCalls != len(testUrls) {
		t.Fatalf(
			`FeedParser.ParseFeedsAsync(testUrls), number of file parser calls = %d, want %d`,
			mockedFileParser.NumOfFuncCalls,
			len(testUrls),
		)
	}
	if mockedWriter.NumOfFuncCalls != 2*len(testUrls) {
		t.Fatalf(
			`FeedParser.ParseFeedsAsync(testUrls), number of queue writer calls = %d, want %d`,
			mockedWriter.NumOfFuncCalls,
			2*len(testUrls),
		)
	}
}

func TestFeedParserResults(t *testing.T) {
	// Prepare mocked data
	mockedFetcher := httpfilefetcher.NewHttpFileFetcher(
		MockedHttpClient{},
	)
	mockedFileParser := xmlparser.NewXmlFeedParser()
	mockedWriter := NewMockedQueueWriter()
	mockedFeedParser := NewFeedParser(
		mockedFetcher,
		mockedFileParser,
		mockedWriter,
	)
	testUrls := []string{"test_url_1"}

	// Prepare results
	expectedBiddingItems := []models.ShopItem{}
	expectedAllItems := []models.ShopItem{}
	for _, item := range mockedCorrectShop.ShopItems {
		if len(item.HeurekaCPC) > 0 {
			expectedBiddingItems = append(expectedBiddingItems, item)
		}
		expectedAllItems = append(expectedAllItems, item)
	}

	// Use ParseFeed function
	mockedFeedParser.ParseFeeds(testUrls)

	// Check if mockedQueueWriter has proper queues
	isBiddingItemsQueueCreated := false
	isAllItemsQueueCreated := false
	for k := range mockedWriter.queues {
		switch k {
		case "shop_items":
			isAllItemsQueueCreated = true
		case "shop_items_bidding":
			isBiddingItemsQueueCreated = true
		}
	}
	if !isAllItemsQueueCreated {
		t.Fatalf(
			`FeedParser.ParseFeedsAsync(testUrls), "shop_items" queue not created, but expected`,
		)
	}
	if !isBiddingItemsQueueCreated {
		t.Fatalf(
			`FeedParser.ParseFeedsAsync(testUrls), "shop_items" queue not created, but expected`,
		)
	}

	// Check if all bidding items are stored in the queue
	allItemsResult := mockedWriter.queues["shop_items"]
	if !reflect.DeepEqual(allItemsResult, expectedAllItems) {
		t.Fatalf(
			"FeedParser.ParseFeedsAsync(testUrls), \"shop_items\" queue contains \n%v\n wanted\n%v\n",
			allItemsResult,
			expectedAllItems,
		)
	}

	// Check if all bidding items are stored in the queue
	biddingItemsResult := mockedWriter.queues["shop_items_bidding"]
	if !reflect.DeepEqual(biddingItemsResult, expectedBiddingItems) {
		t.Fatalf(
			"FeedParser.ParseFeedsAsync(testUrls), \"shop_items_bidding\" queue contains \n%v\n wanted\n%v\n",
			biddingItemsResult,
			expectedBiddingItems,
		)
	}
}

// MOCKED DATA

// Mocked QueueWriter with HasBeenCalled value for checking functions calling
type MockedQueueWriter struct {
	queues         map[string][]models.ShopItem
	NumOfFuncCalls int
}

func NewMockedQueueWriter() *MockedQueueWriter {
	return &MockedQueueWriter{
		queues:         make(map[string][]models.ShopItem),
		NumOfFuncCalls: 0,
	}
}

func (w *MockedQueueWriter) WriteToQueue(
	queueName string,
	shopItems chan models.ShopItem,
) error {
	w.NumOfFuncCalls++
	for item := range shopItems {
		queueItems := w.queues[queueName]
		if queueItems == nil {
			queueItems = []models.ShopItem{}
		}
		w.queues[queueName] = append(queueItems, item)
	}
	return nil
}

// Mocked FileParser with HasBeenCalled value for checking functions calling
type MockedFileFetcher struct {
	NumOfFuncCalls int
}

func (f *MockedFileFetcher) FetchFile(url string) (io.ReadCloser, string, error) {
	f.NumOfFuncCalls++
	return nil, "", nil
}

// Mocked FileParser with HasBeenCalled value for checking functions calling
type MockedFileParser struct {
	NumOfFuncCalls int
}

func (p *MockedFileParser) ParseFile(
	feedFile io.ReadCloser,
	shopItemsOutput chan models.ShopItem,
) error {
	defer close(shopItemsOutput)
	p.NumOfFuncCalls++
	return nil
}

// Mocked http.Client as struct implementing FileFetcher interface
type MockedHttpClient struct{}

func (c MockedHttpClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       mockedReadCloser,
	}, nil
}

// Mocked ErrorsCollector with HasBeenCalled value for checking functions calling
type MockedErrorsCollector struct {
	CollectedErrors []error
	NumOfCalls      int
}

func NewMockedErrorsCollector() *MockedErrorsCollector {
	return &MockedErrorsCollector{
		CollectedErrors: []error{},
		NumOfCalls:      0,
	}
}

func (e *MockedErrorsCollector) HandleErrors(
	feedUrl string,
	stageName string,
) (errorsInput chan error, err error) {
	e.NumOfCalls++

	// Create channel for errors collecting
	errorsInput = make(chan error)

	// Start new go routine for collecting errors
	go func(feedUrl string, stageName string, errorsInput chan error) {
		for er := range errorsInput {
			e.CollectedErrors = append(e.CollectedErrors, er)
		}
	}(feedUrl, stageName, errorsInput)

	return errorsInput, nil
}

// io.ReadCloser with mockedCorrectShop
var mockedXmlFileBytes, _ = xml.Marshal(mockedCorrectShop)
var mockedReadCloser = io.NopCloser(strings.NewReader(string(mockedXmlFileBytes)))

// Correct shop with 3 items (item No. 2 has no HeurekaCPC)
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
