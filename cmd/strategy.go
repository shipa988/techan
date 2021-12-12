package main

import (
	"fmt"
	"github.com/sdcoffey/big"
	"github.com/shipa988/techan"
	"github.com/shipa988/techan/entity"
	"time"
)

// StrategyExample shows how to create a simple trading strategy. In this example, a position should
// be opened if the price moves above 70, and the position should be closed if a position moves below 30.
func StrategyExample() {
	indicator, candles := BasicEma() // from basic.go
	cci := techan.NewCCIIndicator(candles, 20)
	// record trades on this object
	record := techan.NewTradingRecord()

	entryConstant := techan.NewConstantIndicator(57400)
	exitConstant := techan.NewConstantIndicator(10)

	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(entryConstant, indicator),
		techan.PositionNewRule{}) // Is satisfied when the price ema moves above 30 and the current position is new

	exitRule := techan.And(
		techan.NewCrossDownIndicatorRule(indicator, exitConstant),
		techan.PositionOpenRule{}) // Is satisfied when the price ema moves below 10 and the current position is open

	strategy := techan.RuleStrategy{
		UnstablePeriod: 20,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
	candles.LastCandle()

	if strategy.ShouldEnter(1000, record) {
		// иду в эксмо пробую сделать оордер если удается-то заношу в рекорд
		record.Operate(entity.Order{
			Type:    0,
			Pair:    "",
			Price:   big.Decimal{},
			Amount:  big.Decimal{},
			Created: time.Time{},
		})
		fmt.Println("enter", 1000)
	}
	if strategy.ShouldExit(1000, record) {
		fmt.Println("exit", 1000)
	}
	for i := 0; i < len(candles.Candles); i++ {

		fmt.Println(indicator.Calculate(i), ",", cci.Calculate(i), ",", candles.Candles[i].ClosePrice)
		//candles.Candles[i].AddTrade()
	}
	//разобраться с тикером и добавлением новых свечей
	//candles.AddCandle(techan.Candle{
	//	Period:     techan.TimePeriod{},
	//	OpenPrice:  big.Decimal{},
	//	ClosePrice: big.Decimal{},
	//	MaxPrice:   big.Decimal{},
	//	MinPrice:   big.Decimal{},
	//	Volume:     big.Decimal{},
	//	TradeCount: 0,
	//})
	record.Operate(entity.Order{})

}
//https://programtalk.com/vs/ta4j/ta4j-examples/src/main/java/ta4jexamples/bots/TradingBotOnMovingTimeSeries.java/

//https://github.com/ta4j/ta4j/blob/master/ta4j-examples/src/main/java/ta4jexamples/strategies/MovingMomentumStrategy.java
/*
func tickerUpdate(final TickerDTO ticker) {
	getLastTickers().put(ticker.getCurrencyPair(), ticker)
	// If there is no bar or if the duration between the last bar and the ticker is enough.
	if lastAddedBarTimestamp == null || ticker.getTimestamp().isEqual(lastAddedBarTimestamp.plus(getDelayBetweenTwoBars())) || ticker.getTimestamp().isAfter(lastAddedBarTimestamp.plus(getDelayBetweenTwoBars())) {
		// Add the ticker to the series.
		Number
		openPrice = MoreObjects.firstNonNull(ticker.getOpen(), 0)
		Number
		highPrice = MoreObjects.firstNonNull(ticker.getHigh(), 0)
		Number
		lowPrice = MoreObjects.firstNonNull(ticker.getLow(), 0)
		Number
		closePrice = MoreObjects.firstNonNull(ticker.getLast(), 0)
		Number
		volume = MoreObjects.firstNonNull(ticker.getVolume(), 0)
		series.addBar(ticker.getTimestamp(), openPrice, highPrice, lowPrice, closePrice, volume)
		lastAddedBarTimestamp = ticker.getTimestamp()
		// Ask what to do to the strategy.
		int
		endIndex = series.getEndIndex()
		if strategy.shouldEnter(endIndex) {
			// Our strategy should enter.
			shouldEnter()
		} else if strategy.shouldExit(endIndex) {
			// Our strategy should exit.
			shouldExit()
		}
	}
	onTickerUpdate(ticker)
}

*/

/*
public static Strategy buildStrategy(TimeSeries series) {
if (series == null) {
throw new IllegalArgumentException("Series cannot be null");
}
CCIIndicator longCci = new CCIIndicator(series, 200);
CCIIndicator shortCci = new CCIIndicator(series, 5);
Decimal plus100 = Decimal.HUNDRED;
Decimal minus100 = Decimal.valueOf(-100);
Rule entryRule = // Bull trend
new OverIndicatorRule(longCci, plus100).and(// Signal
new UnderIndicatorRule(shortCci, minus100));
Rule exitRule = // Bear trend
new UnderIndicatorRule(longCci, minus100).and(// Signal
new OverIndicatorRule(shortCci, plus100));
Strategy strategy = new Strategy(entryRule, exitRule);
strategy.setUnstablePeriod(5);
return strategy;
}



public class RSI2Strategy {


public static Strategy buildStrategy(TimeSeries series) {
if (series == null) {
throw new IllegalArgumentException("Series cannot be null");
}

ClosePriceIndicator closePrice = new ClosePriceIndicator(series);
SMAIndicator shortSma = new SMAIndicator(closePrice, 5);
SMAIndicator longSma = new SMAIndicator(closePrice, 200);

// We use a 2-period RSI indicator to identify buying
// or selling opportunities within the bigger trend.
RSIIndicator rsi = new RSIIndicator(closePrice, 2);

// Entry rule
// The long-term trend is up when a security is above its 200-period SMA.
Rule entryRule = new OverIndicatorRule(shortSma, longSma) // Trend
.and(new CrossedDownIndicatorRule(rsi, Decimal.valueOf(5))) // Signal 1
.and(new OverIndicatorRule(shortSma, closePrice)); // Signal 2

// Exit rule
// The long-term trend is down when a security is below its 200-period SMA.
Rule exitRule = new UnderIndicatorRule(shortSma, longSma) // Trend
.and(new CrossedUpIndicatorRule(rsi, Decimal.valueOf(95))) // Signal 1
.and(new UnderIndicatorRule(shortSma, closePrice)); // Signal 2

// TODO: Finalize the strategy

return new Strategy(entryRule, exitRule);
}

public static void main(String[] args) {

// Getting the time series
TimeSeries series = CsvTradesLoader.loadBitstampSeries();

// Building the trading strategy
Strategy strategy = buildStrategy(series);

// Running the strategy
TradingRecord tradingRecord = series.run(strategy);
System.out.println("Number of trades for the strategy: " + tradingRecord.getTradeCount());

// Analysis
System.out.println("Total profit for the strategy: " + new TotalProfitCriterion().calculate(series, tradingRecord));
}

}


*/
