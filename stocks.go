package main

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/alackey/go-tdameritrade"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/spacecodewor/fmpcloud-go"
	"github.com/spacecodewor/fmpcloud-go/objects"
	"golang.org/x/oauth2"
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
	if fmpAPIKey == "" {
		return errors.New("Must set environment variable FMP_API_KEY")
	}

	fmpClient, err := fmpcloud.NewAPIClient(fmpcloud.Config{APIKey: fmpAPIKey})
	if err != nil {
		return err
	}

	// TD Ameritrade
	tdaAPIKey := os.Getenv("TDAMERITRADE_API_KEY")
	if tdaAPIKey == "" {
		return errors.New("must set environment variable TDAMERITRADE_API_KEY")
	}

	tdaRefreshToken := os.Getenv("TDAMERITRADE_REFRESH_TOKEN")
	if tdaRefreshToken == "" {
		return errors.New("must set environment variable TDAMERITRADE_REFRESH_TOKEN")
	}

	oauthConfig := oauth2.Config{
		ClientID: tdaAPIKey,
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://api.tdameritrade.com/v1/oauth2/token",
		},
		RedirectURL: "http://localhost",
	}

	oauthToken := &oauth2.Token{
		RefreshToken: tdaRefreshToken,
	}

	ctx := context.Background()
	httpClient := oauthConfig.Client(ctx, oauthToken)

	tdaClient, err := tdameritrade.NewClient(httpClient)
	if err != nil {
		return err
	}

	stocks = &StocksService{fmp: fmpClient, tdameritrade: tdaClient}

	// AWS Session
	awsSess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return nil
}

// FuturesData contains the futures and their data
type FuturesData struct {
	ES *tdameritrade.Quote
	NQ *tdameritrade.Quote
}

// StocksService is the struct for calling stock api's
type StocksService struct {
	fmp          *fmpcloud.APIClient
	tdameritrade *tdameritrade.Client
}

// Quote returns the quote object with the information about the stock
func (s *StocksService) Quote(symbol string) (objects.StockQuote, error) {
	quote, err := s.fmp.Stock.Quote(strings.ToUpper(symbol))
	if err != nil {
		return objects.StockQuote{}, err
	}
	if len(quote) < 1 {
		return objects.StockQuote{}, errors.New("no quote found for symbol: " + symbol)
	}

	return quote[0], nil
}

// Futures returns data for specific futures
func (s *StocksService) Futures() (*FuturesData, error) {
	ctx := context.Background()

	// S&P 500
	symbol := "/ES"
	quotesES, _, err := s.tdameritrade.Quotes.GetQuotes(ctx, symbol)
	if err != nil {
		return nil, err
	}
	if len(*quotesES) <= 0 {
		return nil, errors.New("no quote found for future " + symbol)
	}
	es := (*quotesES)[symbol]

	// NASDAQ
	symbol = "/NQ"
	quotesNQ, _, err := s.tdameritrade.Quotes.GetQuotes(ctx, symbol)
	if err != nil {
		return nil, err
	}
	if len(*quotesNQ) <= 0 {
		return nil, errors.New("no quote found for future " + symbol)
	}
	nq := (*quotesNQ)[symbol]

	return &FuturesData{ES: es, NQ: nq}, nil
}

// Market returns the market heatmap
func (s *StocksService) Market(key string) (string, error) {
	svc := dynamodb.New(awsSess)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(key),
			},
		},
		TableName: aws.String("Stonks_Heatmaps"),
	}

	result, err := svc.GetItem(input)
	if err != nil {
		return "", err
	}

	return *result.Item["Image"].S, nil
}

// News returns the most recent news about the stock symbol
func (s *StocksService) News(symbol string) ([]objects.StockNews, error) {
	request := objects.RequestStockNews{
		SymbolList: []string{symbol},
		Limit:      5,
	}

	news, err := s.fmp.CompanyValuation.StockNews(request)
	if err != nil {
		return nil, err
	}
	return news, nil
}
