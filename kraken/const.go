package kraken

type LimitPurchaseAction int32

const (
	NONE                      LimitPurchaseAction = 0
	INCREASE_LIMIT_ADJUSTMENT LimitPurchaseAction = 1
)
