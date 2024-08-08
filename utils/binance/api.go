package binance

import (
	"errors"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type GetPriceResponse struct {
	Price string `json:"price"`
}

func GetPrice() (float64, error) {

	var res GetPriceResponse
	resty := resty.New()
	resp, err := resty.R().
		SetResult(&res).
		Get("https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT")
	if err != nil {
		return 0, err
	}

	if resp.StatusCode() != 200 {
		return 0, err
	}

	if res.Price == "" {
		return 0, errors.New("can't get price")
	}

	// convert string to float
	priceFloat, err := strconv.ParseFloat(res.Price, 64)
	if err != nil {
		return 0, err
	}

	return priceFloat, nil

}
