package techan

//lint:file-ignore S1038 prefer Fprintln

import (
	"fmt"
	"github.com/shipa988/techan/entity"
	"io"
	"time"

	"github.com/sdcoffey/big"
)

// Analysis is an interface that describes a methodology for taking a TradingRecord as input,
// and giving back some float value that describes it's performance with respect to that methodology.
type Analysis interface {
	Analyze(*TradingRecord) float64
}

// TotalProfitAnalysis analyzes the trading record for total profit.
type TotalProfitAnalysis struct{}

// Analyze analyzes the trading record for total profit.
func (tps TotalProfitAnalysis) Analyze(record *TradingRecord) float64 {
	totalProfit := big.NewDecimal(0)
	for _, trade := range record.Trades {
		if trade.IsClosed() {
			totalProfit = totalProfit.Add(trade.GetProfit())
		}
	}

	return totalProfit.Float()
}

// PercentGainAnalysis analyzes the trading record for the percentage profit gained relative to start
type PercentGainAnalysis struct{}

// Analyze analyzes the trading record for the percentage profit gained relative to start
func (pga PercentGainAnalysis) Analyze(record *TradingRecord) float64 {
	if len(record.Trades) > 0 && record.Trades[0].IsClosed() {
		return (record.Trades[len(record.Trades)-1].GetProfit().Div(record.Trades[0].CostBasis())).Sub(big.NewDecimal(1)).Float()
	}

	return 0
}

// NumTradesAnalysis analyzes the trading record for the number of trades executed
type NumTradesAnalysis string

// Analyze analyzes the trading record for the number of trades executed
func (nta NumTradesAnalysis) Analyze(record *TradingRecord) float64 {
	return float64(len(record.Trades))
}

// LogTradesAnalysis is a wrapper around an io.Writer, which logs every trade executed to that writer
type LogTradesAnalysis struct {
	io.Writer
}

// Analyze logs trades to provided io.Writer
func (lta LogTradesAnalysis) Analyze(record *TradingRecord) float64 {
	logOrder := func(trade *Position) {
		for _, order := range trade.GetOrders() {
			if order.Type == entity.SELL {
				fmt.Fprintln(lta.Writer, fmt.Sprintf("%s - order sell %s (%s @ $%s)", order.Created.UTC().Format(time.RFC822), order.Pair, order.Amount, order.Price))
			} else {
				fmt.Fprintln(lta.Writer, fmt.Sprintf("%s - order buy %s (%s @ $%s)", order.Created.UTC().Format(time.RFC822), order.Pair, order.Amount, order.Price))
			}
		}

		profit := trade.GetProfit().Sub(trade.CostBasis())
		fmt.Fprintln(lta.Writer, fmt.Sprintf("Profit: $%s", profit))
	}

	for _, trade := range record.Trades {
		if trade.IsClosed() {
			logOrder(trade)
		}
	}
	return 0.0
}

// PeriodProfitAnalysis analyzes the trading record for the average profit based on the time period provided.
// i.e., if the trading record spans a year of trading, and PeriodProfitAnalysis wraps one month, Analyze will return
// the total profit for the whole time period divided by 12.
type PeriodProfitAnalysis struct {
	Period time.Duration
}

// Analyze returns the average profit for the trading record based on the given duration
func (ppa PeriodProfitAnalysis) Analyze(record *TradingRecord) float64 {
	var tp TotalProfitAnalysis
	totalProfit := tp.Analyze(record)

	periods := record.Trades[len(record.Trades)-1].ExitOrder().Created.Sub(record.Trades[0].EntranceOrder().Created) / ppa.Period
	return totalProfit / float64(periods)
}

// ProfitableTradesAnalysis analyzes the trading record for the number of profitable trades
type ProfitableTradesAnalysis struct{}

// Analyze returns the number of profitable trades in a trading record
func (pta ProfitableTradesAnalysis) Analyze(record *TradingRecord) float64 {
	var profitableTrades int
	for _, trade := range record.Trades {
		if trade.nowProfit.Cmp(big.ZERO)>0 {
			profitableTrades++
		}
	}

	return float64(profitableTrades)
}

// AverageProfitAnalysis returns the average profit for the trading record. Average profit is represented as the total
// profit divided by the number of trades executed.
type AverageProfitAnalysis struct{}

// Analyze returns the average profit of the trading record
func (apa AverageProfitAnalysis) Analyze(record *TradingRecord) float64 {
	var tp TotalProfitAnalysis
	totalProft := tp.Analyze(record)

	return totalProft / float64(len(record.Trades))
}

// BuyAndHoldAnalysis returns the profit based on a hypothetical where a purchase order was made on the first period available
// and held until the date on the last trade of the trading record. It's useful for comparing the performance of your strategy
// against a simple long position.
type BuyAndHoldAnalysis struct {
	TimeSeries    *TimeSeries
	StartingMoney float64
}

// Analyze returns the profit based on a simple buy and hold strategy
func (baha BuyAndHoldAnalysis) Analyze(record *TradingRecord) float64 {
	if len(record.Trades) == 0 {
		return 0
	}

	openOrder := entity.Order{
		Type:   entity.BUY,
		Amount: big.NewDecimal(baha.StartingMoney).Div(baha.TimeSeries.Candles[0].ClosePrice),
		Price:  baha.TimeSeries.Candles[0].ClosePrice,
	}

	closeOrder := entity.Order{
		Type:   entity.SELL,
		Amount: openOrder.Amount,
		Price:  baha.TimeSeries.Candles[len(baha.TimeSeries.Candles)-1].ClosePrice,
	}

	pos := NewPosition(openOrder)
	pos.AddOrder(closeOrder)

	return pos.GetProfit().Float()
}
