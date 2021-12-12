package techan

import "fmt"

// NewCrossUpIndicatorRule returns a new rule that is satisfied when the lower indicator has crossed above the upper
// indicator.
func NewCrossUpIndicatorRule(upper, lower Indicator) Rule {
	return crossRule{
		upper: upper,
		lower: lower,
		cmp:   1,
	}
}

// NewCrossDownIndicatorRule returns a new rule that is satisfied when the upper indicator has crossed below the lower
// indicator.
func NewCrossDownIndicatorRule(upper, lower Indicator) Rule {
	return crossRule{
		upper: lower,
		lower: upper,
		cmp:   -1,
	}
}

type crossRule struct {
	upper Indicator
	lower Indicator
	cmp   int
}

func (cr crossRule) IsSatisfied(index int, record *TradingRecord) bool {
	i := index

	if i == 0 {
		return false
	}
	fmt.Println("l",cr.lower.Calculate(i).String(),"u",cr.upper.Calculate(i))
	if cmp := cr.lower.Calculate(i).Cmp(cr.upper.Calculate(i)); cmp == 0 || cmp == cr.cmp {
		for ; i >= 0; i-- {
			fmt.Println("l",cr.lower.Calculate(i).String(),"u",cr.upper.Calculate(i))
			if cmp = cr.lower.Calculate(i).Cmp(cr.upper.Calculate(i)); cmp == 0 || cmp == -cr.cmp {
				return true
			}
			//https://github.com/ta4j/ta4j/blob/8994ecbc9e644f7cfc36216c7de9eef4b65a4b22/ta4j-core/src/main/java/org/ta4j/core/indicators/helpers/CrossIndicator.java#L68
		}
	}

	return false
}
