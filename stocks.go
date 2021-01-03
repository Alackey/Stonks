package main

import (
	"os"
	"strconv"

	"github.com/z-Wind/alphavantage"
)

var stocks *StocksService

// NewStocksClient initializes a connection to the stock api
func NewStocksClient() error {
	if stocks != nil {
		return nil
	}

	apiKey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.GetClient(apiKey)
	
	service, err := alphavantage.New(client)
	if err != nil {
		return err
	}

	stocks = &StocksService{service}
	return nil
}

// StocksService is the struct for calling stock api's
type StocksService struct {
	client *alphavantage.Service
}

// Quote returns the current market price of the ticker given
func (s *StocksService) Quote(ticker string) (string, error) {
	quote, err := s.client.TimeSeries.QuoteEndpoint(ticker).Do()
	if err != nil {
		return "", err
	}

	price := strconv.FormatFloat(quote.Price, 'f', 2, 64)
	return price, nil
}
