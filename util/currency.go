package util

const (
	AUD = "AUD"
	JPY = "JPY"
	RMB = "RMB"
	CHF = "CHF"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case AUD, JPY, RMB, CHF:
		return true
	}
	return false
}
