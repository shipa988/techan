package techan

import (
	"github.com/shipa988/techan/entity"
	"testing"
	"time"

	"github.com/sdcoffey/big"
	"github.com/stretchr/testify/assert"
)

func TestPosition_NoOrders_IsNew(t *testing.T) {
	position := new(Position)

	assert.True(t, position.IsNew())
}

func TestPosition_NewPosition_IsOpen(t *testing.T) {
	order := entity.Order{
		Type:   entity.BUY,
		Amount: big.ONE,
		Price:  big.NewFromString("2"),
	}

	position := NewPosition(order)
	assert.True(t, position.IsOpen())
	assert.False(t, position.IsNew())
	assert.False(t, position.IsClosed())
	assert.EqualValues(t, 2.0,position.GetAvgPrice().Float())
}

func TestNewPosition_WithBuy_IsLong(t *testing.T) {
	order := entity.Order{
		Type:   entity.BUY,
		Amount: big.ONE,
		Price:  big.NewFromString("2"),
	}

	position := NewPosition(order)
	assert.True(t, position.IsLong())
}

func TestNewPosition_WithSell_IsShort(t *testing.T) {
	order := entity.Order{
		Type:   entity.SELL,
		Amount: big.ONE,
		Price:  big.NewFromString("2"),
	}

	position := NewPosition(order)
	assert.True(t, position.IsShort())
}

func TestPosition_Enter(t *testing.T) {
	order := entity.Order{
		Type:   entity.BUY,
		Amount: big.ONE,
		Price:  big.NewFromString("2"),
	}
	position := NewPosition(order)

	assert.True(t, position.IsOpen())
	assert.EqualValues(t, order.Amount, position.GetOrder(0).Amount)
	assert.EqualValues(t, order.Price, position.GetOrder(0).Price)
	assert.EqualValues(t, order.Created, position.GetOrder(0).Created)
}

func TestPosition_AvgPrice(t *testing.T) {
	order := entity.Order{
		Type:   entity.BUY,
		Amount: big.ONE,
		Price:  big.NewFromString("2"),
	}
	position := NewPosition(order)
	assert.EqualValues(t, 2.0,position.GetAvgPrice().Float())
	position.AddOrder(order)
	assert.EqualValues(t, 2.0,position.GetAvgPrice().Float())
	order3 := entity.Order{
		Type:   entity.BUY,
		Amount: big.NewFromString("2"),
		Price:  big.NewFromString("4"),
	}
	position.AddOrder(order3)
	assert.EqualValues(t, 3.0,position.GetAvgPrice().Float())

	assert.True(t, position.IsOpen())
	assert.EqualValues(t, 4.0, position.GetAmount().Float())

	orderRevers := entity.Order{
		Type:   entity.SELL,
		Amount: big.NewFromString("4"),
		Price:  big.NewFromString("3"),
	}
	position.AddOrder(orderRevers)
	assert.True(t, position.IsClosed())
	assert.EqualValues(t, 0.0,position.GetProfit().Float())
}

func TestPosition_AvgPriceSell(t *testing.T) {
	order := entity.Order{
		Type:   entity.SELL,
		Amount: big.ONE,
		Price:  big.NewFromString("2"),
	}
	position := NewPosition(order)
	assert.EqualValues(t, 2.0,position.GetAvgPrice().Float())
	position.AddOrder(order)
	assert.EqualValues(t, 2.0,position.GetAvgPrice().Float())
	order3 := entity.Order{
		Type:   entity.SELL,
		Amount: big.NewFromString("2"),
		Price:  big.NewFromString("4"),
	}
	position.AddOrder(order3)
	assert.EqualValues(t, 3.0,position.GetAvgPrice().Float())

	assert.True(t, position.IsOpen())
	assert.EqualValues(t, 4.0, position.GetAmount().Float())
	orderRevers := entity.Order{
		Type:   entity.BUY,
		Amount: big.NewFromString("4"),
		Price:  big.NewFromString("3"),
	}
	position.AddOrder(orderRevers)
	assert.True(t, position.IsClosed())
	assert.EqualValues(t, 0.0,position.GetProfit().Float())
}


func TestPosition_Close(t *testing.T) {
	entranceOrder := entity.Order{
		Type:   entity.BUY,
		Amount: big.ONE,
		Price:  big.NewFromString("2"),
	}

	position := NewPosition(entranceOrder)

	assert.True(t, position.IsOpen())
	assert.EqualValues(t, entranceOrder.Amount, position.GetOrder(0).Amount)
	assert.EqualValues(t, entranceOrder.Price, position.GetOrder(0).Price)
	assert.EqualValues(t, entranceOrder.Created, position.GetOrder(0).Created)

	exitOrder := entity.Order{
		Type:    entity.SELL,
		Amount:  big.ONE,
		Price:   big.NewFromString("4"),
		Created: time.Now(),
	}

	position.AddOrder(exitOrder)

	assert.True(t, position.IsClosed())

	assert.EqualValues(t, exitOrder.Amount, position.GetOrder(1).Amount)
	assert.EqualValues(t, exitOrder.Price, position.GetOrder(1).Price)
	assert.EqualValues(t, exitOrder.Created, position.GetOrder(1).Created)
}

func TestPosition_CostBasis(t *testing.T) {
	t.Run("When entrance order nil, returns 0", func(t *testing.T) {
		p := new(Position)
		assert.EqualValues(t, "0", p.CostBasis().String())
	})

	t.Run("When entracne order not nil, returns cost basis", func(t *testing.T) {

		order := entity.Order{
			Type:   entity.BUY,
			Amount: big.ONE,
			Price:  big.NewFromString("2"),
		}

		p := NewPosition(order)

		assert.EqualValues(t, "2.00", p.CostBasis().FormattedString(2))
	})
}

func TestPosition_ExitValue(t *testing.T) {
	t.Run("when not closed, returns 0", func(t *testing.T) {
		order := entity.Order{
			Type:   entity.BUY,
			Amount: big.ONE,
			Price:  big.NewFromString("2"),
		}

		p := NewPosition(order)

		assert.EqualValues(t, "0.00", p.GetProfit().FormattedString(2))
	})
	//
	//t.Run("when closed, returns exit value", func(t *testing.T) {
	//
	//	order := entity.Order{
	//		Type:   entity.BUY,
	//		Amount: big.ONE,
	//		Price:  big.NewFromString("2"),
	//	}
	//
	//	p := NewPosition(order)
	//	order = entity.Order{
	//		Type:   entity.SELL,
	//		Amount: big.ONE,
	//		Price:  big.NewFromString("12"),
	//	}
	//
	//	p.AddOrder(order)
	//
	//	assert.EqualValues(t, "12.00", p.GetProfit().FormattedString(2))
	//})
}
