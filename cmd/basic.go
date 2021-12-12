package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/sdcoffey/big"
	"github.com/shipa988/techan"
)

// BasicEma is an example of how to create a basic Exponential moving average indicator
// based on the close prices of a timeseries from your exchange of choice.
func BasicEma() (techan.Indicator,*techan.TimeSeries) {
	series := techan.NewTimeSeries()

	// fetch this from your preferred exchange
	//dataset := [][]string{
	//	// Timestamp, Open, Close, High, Low, volume
	//	{"1234567", "1", "2", "3", "5", "6"},
	//}
	dataset := getDataSet()

	for _, datum := range dataset.Candles {
		//start, _ := strconv.ParseInt(datum[0], 10, 64)
		period := techan.NewTimePeriod(time.Unix(0, datum.T*int64(time.Millisecond)), time.Minute)

		candle := techan.NewCandle(period)
		candle.OpenPrice = big.NewDecimal(datum.O)
		candle.ClosePrice = big.NewDecimal(datum.C)
		candle.MaxPrice = big.NewDecimal(datum.H)
		candle.MinPrice = big.NewDecimal(datum.L)
		candle.Volume = big.NewDecimal(datum.V)
		series.AddCandle(candle)
	}

	closePrices := techan.NewClosePriceIndicator(series)
	movingAverage := techan.NewEMAIndicator(closePrices, 20) // Create an exponential moving average with a window of 10
	//volume := techan.NewVolumeIndicator(series)
	//movingAverage
	//pvtIndicator := techan.NewWindowedStandardDeviationIndicator(closePrices,  1)
	//fmt.Println("PVT - ", pvtIndicator.Calculate(n)
	//
	//pvtSignalIndicator := techan.NewPVTAndSignalIndicator(pvtIndicator, techan.NewEMAIndicator(pvtIndicator, 21))
	//fmt.Println("PVT - Signal - ", pvtSignalIndicator.Calculate(n
	return movingAverage,series
}

type responce struct {
	Candles []struct {
		T int64   `json:"t"`
		O float64 `json:"o"`
		C float64 `json:"c"`
		H float64 `json:"h"`
		L float64 `json:"l"`
		V float64 `json:"v"`
	} `json:"candles"`
}

//type Candel struct {
//	Period     int64 `json:"t"`
//	OpenPrice  big.Decimal `json:"o"`
//	ClosePrice big.Decimal `json:"c"`
//	MaxPrice   big.Decimal `json:"h"`
//	MinPrice   big.Decimal `json:"l"`
//	Volume     big.Decimal `json:"v"`
//}
//type Candles struct {
//	c []Candel `json:"candles"`
//}

func getDataSet() *responce {

	url := "https://api.exmo.me/v1.1/candles_history?symbol=BTC_USD&resolution=1&from=" +
		strconv.FormatInt(time.Now().AddDate(0, 0, -1).Unix(), 10) +
		"&to=" + strconv.FormatInt(time.Now().Unix(), 10)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(string(body))
	var candles responce
	err = json.Unmarshal(body, &candles)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &candles
}
