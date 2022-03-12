package errorscollector

import (
	"fmt"

	"go.uber.org/zap"
)

// Collects errors from many feed processing async stages and handles them
type ErrorsCollectorInterface interface {
	HandleErrors(feedUrl string, stageName string) (errorsInput chan error, err error)
}

// Collects errors from many feed processing async stages and logs them
type ErrorsCollector struct{}

func NewErrorsCollector() *ErrorsCollector {
	return &ErrorsCollector{}
}

// Runs new go routine collecting all errors from errorsInput
func (e ErrorsCollector) HandleErrors(
	feedUrl string,
	stageName string,
) (errorsInput chan error, err error) {
	defer zap.L().Sync()

	// Create channel for errors collecting
	errorsInput = make(chan error)

	// Start new go routine for collecting errors
	go func(feedUrl string, stageName string, errorsInput chan error) {
		for e := range errorsInput {
			// Print all incoming errors to logs
			zap.L().Error(
				fmt.Sprintf("Error during processing feed file from %s", feedUrl),
				zap.String("stage", stageName),
				zap.Error(e),
			)
			zap.L().Sync()
		}
	}(feedUrl, stageName, errorsInput)
	zap.L().Info(
		fmt.Sprintf("Started collecting errors of stage %s during processing feed from %s", stageName, feedUrl),
		zap.String("feedUrl", feedUrl),
		zap.String("stage", stageName),
	)
	return errorsInput, nil
}
