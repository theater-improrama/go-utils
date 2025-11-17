package query

type Order int

const (
	OrderAscending Order = iota
	OrderDescending
)

func OrderFromBool(b bool) Order {
	if b {
		return OrderDescending
	}

	return OrderAscending
}
