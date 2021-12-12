package techan

import (
	"github.com/sdcoffey/big"
	"github.com/shipa988/techan/entity"
	"sync"
)

// Position is list of orders  before position closing and totalAmount amount==0
type Position struct {
	sync.RWMutex
	orders    []entity.Order
	nowType   entity.OrderSide
	nowAmount big.Decimal
	nowPrice  big.Decimal
	isClosed  bool
	nowProfit big.Decimal
}

// NewPosition returns a new Position with the passed-in order as the open order
func NewPosition(openOrder entity.Order) (p *Position) {
	p = new(Position)
	p.orders = make([]entity.Order, 0)
	p.orders = append(p.orders, openOrder)
	p.nowType = openOrder.Type
	p.nowAmount = openOrder.Amount
	p.nowPrice = openOrder.Price
	if openOrder.Type == entity.SELL {
		p.nowProfit = big.ZERO.Add(openOrder.Amount.Mul(openOrder.Price))
	} else {
		p.nowProfit = big.ZERO.Sub(openOrder.Amount.Mul(openOrder.Price))
	}
	return p
}

// NewPosition returns a new Position with the passed-in order as the open order
func EmptyPosition() (p *Position) {
	p = new(Position)
	p.orders = make([]entity.Order, 0)
	return p
}

// Enter sets the open order to the order passed in
func (p *Position) AddOrder(order entity.Order) bool {
	p.Lock()
	defer p.Unlock()
	if p.isClosed {
		return false
	}
	p.orders = append(p.orders, order)
	if order.Type == entity.SELL {
		p.nowProfit = p.nowProfit.Add(order.Amount.Mul(order.Price))
	} else {
		p.nowProfit = p.nowProfit.Sub(order.Amount.Mul(order.Price))
	}
	// если однонаправленные позиции
	if p.nowType == order.Type {
		p.nowAmount=p.nowAmount.Add(order.Amount)
	} else {
		if p.nowAmount.Cmp(order.Amount) > 0 {
			p.nowAmount = p.nowAmount.Sub(order.Amount)
		} else if p.nowAmount.Cmp(order.Amount) == 0 {
			p.isClosed = true
			p.nowAmount = big.ZERO
		} else if p.nowAmount.Cmp(order.Amount) < 0 {
			p.nowAmount = order.Amount.Sub(p.nowAmount)
			p.nowType = order.Type
		}
	}
	if !p.isClosed {
		p.nowPrice = p.nowProfit.Div(p.nowAmount).Abs()
	}
	return true
}

func (p *Position) GetAvgPrice() big.Decimal {
	p.RLock()
	defer p.RUnlock()
	if !p.isClosed {
		return p.nowPrice
	}
	return big.ZERO
}

func (p *Position) GetAmount() big.Decimal {
	p.RLock()
	defer p.RUnlock()
	return p.nowAmount
}

//// Exit sets the exit order to the order passed in
//func (p *Position) Exit(order dto.Order) {
//	p.orders[1] = &order
//}

// IsLong returns true if the entrance order is a buy order
func (p *Position) IsLong() bool {
	p.RLock()
	defer p.RUnlock()
	return p.nowType == entity.BUY
}

// IsShort returns true if the entrance order is a sell order
func (p *Position) IsShort() bool {
	p.RLock()
	defer p.RUnlock()
	return p.nowType == entity.SELL
}

// IsOpen returns true if there is an entrance order but no exit order
func (p *Position) IsOpen() bool {
	p.RLock()
	defer p.RUnlock()
	return p.isOpen()
}

func (p *Position) isOpen() bool {
	return len(p.orders) > 0 && !p.isClosed
}

// IsNew returns true if there is neither an entrance or exit order
func (p *Position) IsNew() bool {
	p.RLock()
	defer p.RUnlock()
	return len(p.orders) == 0
}

func (p *Position) GetOrder(id int) *entity.Order {
	p.RLock()
	defer p.RUnlock()
	if len(p.orders)-1 >= id {
		return &p.orders[id]
	}
	return nil
}

// EntranceOrder returns the entrance order of this position
func (p *Position) EntranceOrder() *entity.Order {
	p.RLock()
	defer p.RUnlock()
	if len(p.orders) > 0 {
		return &p.orders[0]
	}
	return nil
}

func (p *Position) GetOrders() []entity.Order {
	p.RLock()
	defer p.RUnlock()
	return p.orders
}

// ExitOrder returns the exit order of this position
func (p *Position) ExitOrder() *entity.Order {
	p.RLock()
	defer p.RUnlock()
	if len(p.orders) > 0 {
		return &p.orders[len(p.orders)-1]
	}
	return nil
}

// IsClosed returns true of there are both entrance and exit orders
func (p *Position) IsClosed() bool {
	p.RLock()
	defer p.RUnlock()
	return p.isClosed
}

// CostBasis returns the price to enter this order
func (p *Position) CostBasis() big.Decimal {
	p.RLock()
	defer p.RUnlock()
	if len(p.orders) > 0 {
		return p.orders[0].Amount.Mul(p.orders[0].Price)
	}
	return big.ZERO
}

// GetProfit returns the value accrued by closing the position
func (p *Position) GetProfit() big.Decimal {
	p.RLock()
	defer p.RUnlock()
	if p.isClosed {
		return p.nowProfit
	}

	return big.ZERO
}
