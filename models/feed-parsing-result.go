package models

type ResultStatus string

const (
	ParsedSuccessfully ResultStatus = "PARSED_SUCCESSFULLY"
	ParsingInProgress  ResultStatus = "PARSED_IN_PROGRESS"
	ParsingErrors      ResultStatus = "PARSING_ERROR"
)

type FeedParsingResult struct {
	FeedUrl string       `json:"feedUrl"`
	Status  ResultStatus `json:"status"`
}
