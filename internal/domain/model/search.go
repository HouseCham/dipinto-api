package model

type ClientSearchParams struct {
	SearchValue string
	CategoryID  uint64
	MinPrice    float64
	MaxPrice    float64
	OrderBy     string
	Offset      int
	Limit       int
}
