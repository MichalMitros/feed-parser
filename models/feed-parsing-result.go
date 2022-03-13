package models

type ResultStatus string

const (
	ParsedSuccessfully ResultStatus = "PARSED_SUCCESSFULLY"
	ParsingErrors      ResultStatus = "PARSING_ERROR"
)

type FeedParsingResult struct {
	FeedUrl string
	Status  ResultStatus
}
