package kraken

type LimitPurchaseAction int32

const (
	NOTHING                   LimitPurchaseAction = 0
	INCREASE_LIMIT_ADJUSTMENT LimitPurchaseAction = 1
)
