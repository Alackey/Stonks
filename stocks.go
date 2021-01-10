package main

import (
	"bytes"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spacecodewor/fmpcloud-go"
	"github.com/spacecodewor/fmpcloud-go/objects"
)

var stocks *StocksService
var awsSess *session.Session

// NewStocksClient initializes a connection to the stock api
func NewStocksClient() error {
	if stocks != nil {
		return nil
	}

	// Financial Modeling Prep
	fmpAPIKey := os.Getenv("FMP_API_KEY")

	client, err := fmpcloud.NewAPIClient(fmpcloud.Config{APIKey: fmpAPIKey})
	if err != nil {
		return err
	}

	stocks = &StocksService{client}

	// AWS Session
	awsSess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return nil
}

// StocksService is the struct for calling stock api's
type StocksService struct {
	client *fmpcloud.APIClient
}

// Quote returns the quote object with the information about the stock
func (s *StocksService) Quote(symbol string) (objects.StockQuote, error) {
	quote, err := s.client.Stock.Quote(symbol)
	if err != nil {
		return objects.StockQuote{}, err
	}
	if len(quote) < 1 {
		return objects.StockQuote{}, errors.New("No quote found for symbol: " + symbol)
	}

	return quote[0], nil
}

// Market returns the market heatmap
func (s *StocksService) Market() (*bytes.Reader, error) {
	var buf []byte
	awsBuf := aws.NewWriteAtBuffer(buf)

	input := &s3.GetObjectInput{
		Bucket: aws.String("stockbot-heatmap"),
		Key:    aws.String("marketHeatmap.png"),
	}

	downloader := s3manager.NewDownloader(awsSess)

	_, err := downloader.Download(awsBuf, input)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(awsBuf.Bytes()), nil
}
