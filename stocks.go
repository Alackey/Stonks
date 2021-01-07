package main

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/z-Wind/alphavantage"
)

var stocks *StocksService
var awsSess *session.Session

// NewStocksClient initializes a connection to the stock api
func NewStocksClient() error {
	if stocks != nil {
		return nil
	}

	// AlphaVantage
	apiKey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.GetClient(apiKey)

	service, err := alphavantage.New(client)
	if err != nil {
		return err
	}

	stocks = &StocksService{service}

	// AWS Session
	awsSess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return nil
}

// StocksService is the struct for calling stock api's
type StocksService struct {
	client *alphavantage.Service
}

// Quote returns the quote object with the information about the stock
func (s *StocksService) Quote(ticker string) (*alphavantage.Quote, error) {
	quote, err := s.client.TimeSeries.QuoteEndpoint(ticker).Do()
	if err != nil {
		return nil, err
	}

	return quote, nil
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
