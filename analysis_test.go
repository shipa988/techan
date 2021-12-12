package techan

import (
	"github.com/shipa988/techan/entity"
	"testing"
	"time"

	"bytes"

	"bufio"

	"fmt"

	"github.com/sdcoffey/big"
	"github.com/stretchr/testify/assert"
)

const example = "EXM"

func TestTotalProfitAnalysis(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		record := NewTradingRecord()
		tpa := TotalProfitAnalysis{}

		orders := []entity.Order{
			{
				Type:    entity.BUY,
				Amount:  big.NewDecimal(2),
				Price:   big.NewDecimal(1),
				Pair:    example,
				Created: time.Now(),
			},
			{
				Type:    entity.SELL,
				Amount:  big.NewDecimal(2),
				Price:   big.NewDecimal(1),
				Pair:    example,
				Created: time.Now(),
			},
			{
				Type:    entity.SELL,
				Amount:  big.NewDecimal(1),
				Price:   big.NewDecimal(2),
				Pair:    example,
				Created: time.Now(),
			},
			{
				Type:    entity.BUY,
				Amount:  big.NewDecimal(1),
				Price:   big.NewDecimal(1),
				Pair:    example,
				Created: time.Now(),
			},
		}

		for _, order := range orders {
			record.Operate(order)
		}

		assert.EqualValues(t, 1.0, tpa.Analyze(record))

		record.Operate(entity.Order{
			Type:    entity.BUY,
			Amount:  big.ONE,
			Price:   big.ONE,
			Pair:    example,
			Created: time.Now(),
		})

		record.Operate(entity.Order{
			Type:    entity.SELL,
			Amount:  big.NewFromString("0.5"),
			Price:   big.ONE,
			Pair:    example,
			Created: time.Now(),
		})
		record.Operate(entity.Order{
			Type:    entity.SELL,
			Amount:  big.NewFromString("0.5"),
			Price:   big.ONE,
			Pair:    example,
			Created: time.Now(),
		})
		assert.EqualValues(t, 1, tpa.Analyze(record))

		record.Operate(entity.Order{
			Type:    entity.BUY,
			Amount:  big.ONE,
			Price:   big.NewFromString("2"),
			Pair:    example,
			Created: time.Now(),
		})

		record.Operate(entity.Order{
			Type:    entity.SELL,
			Amount:  big.NewFromString("0.5"),
			Price:   big.ONE,
			Pair:    example,
			Created: time.Now(),
		})
		record.Operate(entity.Order{
			Type:    entity.SELL,
			Amount:  big.NewFromString("0.5"),
			Price:   big.ONE,
			Pair:    example,
			Created: time.Now(),
		})
		assert.EqualValues(t, 0, tpa.Analyze(record))
	})
}

func TestPercentGainAnalysis(t *testing.T) {
	t.Run("Zero", func(t *testing.T) {
		record := NewTradingRecord()

		pga := PercentGainAnalysis{}

		assert.EqualValues(t, 0, pga.Analyze(record))
	})

	t.Run("Simple gain", func(t *testing.T) {
		record := NewTradingRecord()

		pga := PercentGainAnalysis{}

		orders := []entity.Order{
			{
				Type:    entity.BUY,
				Amount:  big.NewDecimal(1),
				Price:   big.NewDecimal(1),
				Pair:    example,
				Created: time.Now(),
			},
			{
				Type:    entity.SELL,
				Amount:  big.NewDecimal(2),
				Price:   big.NewDecimal(1),
				Pair:    example,
				Created: time.Now(),
			},
		}

		for _, order := range orders {
			record.Operate(order)
		}

		gain := pga.Analyze(record)
		assert.EqualValues(t, 1, gain)
	})

	t.Run("Simple loss", func(t *testing.T) {
		record := NewTradingRecord()

		pga := PercentGainAnalysis{}

		orders := []entity.Order{
			{
				Type:    entity.BUY,
				Amount:  big.NewDecimal(2),
				Price:   big.NewDecimal(1),
				Pair:    example,
				Created: time.Now(),
			},
			{
				Type:    entity.SELL,
				Amount:  big.NewDecimal(1),
				Price:   big.NewDecimal(1),
				Pair:    example,
				Created: time.Now(),
			},
		}

		for _, order := range orders {
			record.Operate(order)
		}

		gain := pga.Analyze(record)
		assert.EqualValues(t, -.5, gain)
	})

	t.Run("Small loss and gain", func(t *testing.T) {
		record := NewTradingRecord()

		pga := PercentGainAnalysis{}

		orders := []entity.Order{
			{
				Type:    entity.BUY,
				Amount:  big.NewDecimal(2),
				Price:   big.NewDecimal(1),
				Pair:    example,
				Created: time.Now(),
			},
			{
				Type:    entity.SELL,
				Amount:  big.NewDecimal(1),
				Price:   big.NewDecimal(1),
				Pair:    example,
				Created: time.Now(),
			},
			{
				Type:    entity.BUY,
				Amount:  big.NewDecimal(1),
				Price:   big.NewDecimal(1),
				Pair:    example,
				Created: time.Now(),
			},
			{
				Type:    entity.SELL,
				Amount:  big.NewDecimal(1),
				Price:   big.NewDecimal(1.25),
				Pair:    example,
				Created: time.Now(),
			},
		}

		for _, order := range orders {
			record.Operate(order)
		}

		gain := pga.Analyze(record)
		assert.EqualValues(t, -.375, gain)
	})
}

func TestNumTradesAnalysis(t *testing.T) {
	record := NewTradingRecord()

	var nta NumTradesAnalysis

	assert.EqualValues(t, 0, nta.Analyze(record))
}

func TestLogTradesAnalysis(t *testing.T) {
	buffer := bytes.NewBufferString("")

	logger := LogTradesAnalysis{
		Writer: buffer,
	}

	record := NewTradingRecord()

	now := time.Now().UTC()
	dates := []time.Time{
		now,
		now.AddDate(0, 0, 1),
		now.AddDate(0, 0, 2),
		now.AddDate(0, 0, 3),
	}

	orders := []entity.Order{
		{
			Type:    entity.BUY,
			Amount:  big.NewDecimal(1),
			Price:   big.NewDecimal(2),
			Pair:    example,
			Created: dates[0],
		},
		{
			Type:    entity.SELL,
			Amount:  big.NewDecimal(1),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: dates[1],
		},
		{
			Type:    entity.BUY,
			Amount:  big.NewDecimal(1),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: dates[2],
		},
		{
			Type:    entity.SELL,
			Amount:  big.NewDecimal(1),
			Price:   big.NewDecimal(1.25),
			Pair:    example,
			Created: dates[3],
		},
	}

	for _, order := range orders {
		record.Operate(order)
	}

	val := logger.Analyze(record)
	assert.EqualValues(t, 0, val)

	scanner := bufio.NewScanner(buffer)

	var i int
	for scanner.Scan() {
		text := scanner.Text()

		var expected string
		switch i {
		case 0:
			expected = fmt.Sprintf("%s - enter with buy EXM (1 @ $2)", dates[0].Format(time.RFC822))
		case 1:
			expected = fmt.Sprintf("%s - exit with sell EXM (1 @ $1)", dates[1].Format(time.RFC822))
		case 2:
			expected = "Profit: $-1"
		case 3:
			expected = fmt.Sprintf("%s - enter with buy EXM (1 @ $1)", dates[2].Format(time.RFC822))
		case 4:
			expected = fmt.Sprintf("%s - exit with sell EXM (1 @ $1.25)", dates[3].Format(time.RFC822))
		case 5:
			expected = "Profit: $0.25"
		}

		assert.EqualValues(t, expected, text)
		i++
	}
}

func TestPeriodProfitAnalysis(t *testing.T) {
	record := NewTradingRecord()

	now := time.Now().Add(-time.Minute * 5)

	orders := []entity.Order{
		{
			Type:    entity.BUY,
			Amount:  big.NewDecimal(1),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: now,
		},
		{
			Type:    entity.SELL,
			Amount:  big.NewDecimal(2),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: now.Add(time.Minute),
		},
		{
			Type:    entity.BUY,
			Amount:  big.NewDecimal(2),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: now.Add(time.Minute * 2),
		},
		{
			Type:    entity.SELL,
			Amount:  big.NewDecimal(3),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: now.Add(time.Minute * 3),
		},
	}

	for _, order := range orders {
		record.Operate(order)
	}

	ppa := PeriodProfitAnalysis{
		Period: time.Minute * 2,
	}

	assert.EqualValues(t, 2, ppa.Analyze(record))
}

func TestProfitableTradesAnalysis(t *testing.T) {
	record := NewTradingRecord()

	orders := []entity.Order{
		{
			Type:    entity.BUY,
			Amount:  big.NewDecimal(1),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: time.Now(),
		},
		{
			Type:    entity.SELL,
			Amount:  big.NewDecimal(2),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: time.Now(),
		},
		{
			Type:    entity.BUY,
			Amount:  big.NewDecimal(2),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: time.Now(),
		},
		{
			Type:    entity.SELL,
			Amount:  big.NewDecimal(1),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: time.Now(),
		},
	}

	for _, order := range orders {
		record.Operate(order)
	}

	pta := ProfitableTradesAnalysis{}

	assert.EqualValues(t, 1, pta.Analyze(record))
}

func TestAverageProfitAnalysis(t *testing.T) {
	record := NewTradingRecord()

	orders := []entity.Order{
		{
			Type:    entity.BUY,
			Amount:  big.NewDecimal(1),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: time.Now(),
		},
		{
			Type:    entity.SELL,
			Amount:  big.NewDecimal(2),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: time.Now(),
		},
		{
			Type:    entity.BUY,
			Amount:  big.NewDecimal(2),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: time.Now(),
		},
		{
			Type:    entity.SELL,
			Amount:  big.NewDecimal(5),
			Price:   big.NewDecimal(1),
			Pair:    example,
			Created: time.Now(),
		},
	}

	for _, order := range orders {
		record.Operate(order)
	}

	pta := AverageProfitAnalysis{}

	assert.EqualValues(t, 2, pta.Analyze(record))
}

func TestBuyAndHoldAnalysis(t *testing.T) {
	series := mockTimeSeries("1", "2", "3", "2", "6")
	record := NewTradingRecord()

	t.Run("== 0 trades returns zero", func(t *testing.T) {
		buyAndHoldAnalysis := BuyAndHoldAnalysis{
			TimeSeries:    series,
			StartingMoney: 1,
		}

		assert.EqualValues(t, 0, buyAndHoldAnalysis.Analyze(record))
	})

	t.Run("> 0 trades", func(t *testing.T) {
		orders := []entity.Order{
			{
				Type:    entity.BUY,
				Amount:  big.NewDecimal(1),
				Price:   big.NewDecimal(1),
				Pair:    example,
				Created: time.Now(),
			},
			{
				Type:    entity.SELL,
				Amount:  big.NewDecimal(1),
				Price:   big.NewDecimal(2),
				Pair:    example,
				Created: time.Now(),
			},
			{
				Type:    entity.BUY,
				Amount:  big.NewDecimal(1),
				Price:   big.NewDecimal(2),
				Pair:    example,
				Created: time.Now(),
			},
			{
				Type:    entity.SELL,
				Amount:  big.NewDecimal(1),
				Price:   big.NewDecimal(6),
				Pair:    example,
				Created: time.Now(),
			},
		}

		for _, order := range orders {
			record.Operate(order)
		}

		buyAndHoldAnalysis := BuyAndHoldAnalysis{
			TimeSeries:    series,
			StartingMoney: 1,
		}

		assert.EqualValues(t, 5.0, buyAndHoldAnalysis.Analyze(record))
	})
}
