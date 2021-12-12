package techan

import "github.com/shipa988/techan/entity"

// TradingRecord is an object describing a series of trades made and a current position
/**
 * A history/record of a trading session.
 * <p>
 * Holds the full trading record when running a {@link Strategy strategy}.
 * It is used to:
 * <ul>
 * <li>check to satisfaction of some trading rules (when running a strategy)
 * <li>analyze the performance of a trading strategy
 * </ul>
 */
type TradingRecord struct {
	Trades          []*Position
	currentPosition *Position
}

// NewTradingRecord returns a new TradingRecord
func NewTradingRecord() (t *TradingRecord) {
	t = new(TradingRecord)
	t.Trades = make([]*Position, 0)
	t.currentPosition=EmptyPosition()
	return t
}

// CurrentPosition returns the current position in this record
func (tr *TradingRecord) CurrentPosition() *Position {
	return tr.currentPosition
}

// LastTrade returns the last trade executed in this record
func (tr *TradingRecord) LastTrade() *Position {
	if len(tr.Trades) == 0 {
		return nil
	}

	return tr.Trades[len(tr.Trades)-1]
}

// Operate takes an order and adds it to the current TradingRecord. It will only add the order if:
// - The current position is open and the passed order was executed after the entrance order
// - The current position is new and the passed order was executed after the last exit order
func (tr *TradingRecord) Operate(order entity.Order) {
	if tr.currentPosition.IsOpen() {
		if order.Created.Before(tr.CurrentPosition().EntranceOrder().Created) {
			return
		}

		tr.currentPosition.AddOrder(order)
		if tr.currentPosition.IsClosed(){
			tr.Trades = append(tr.Trades, tr.currentPosition)
			tr.currentPosition = EmptyPosition()
		}
	} else if tr.currentPosition.IsNew() {
		if tr.LastTrade() != nil && order.Created.Before(tr.LastTrade().ExitOrder().Created) {
			return
		}

		tr.currentPosition = NewPosition(order)
	}
}
