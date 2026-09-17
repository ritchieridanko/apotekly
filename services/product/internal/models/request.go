package models

type (
	CreateProductReq struct {
		BrandName      string
		GenericName    string
		Description    *string
		RequiresRX     bool
		DosageForm     string
		Strength       *string
		PackUnit       string
		PackSize       int
		HeightCM       float64
		LengthCM       float64
		WidthCM        float64
		WeightG        float64
		Price          string
		Currency       string
		Quantity       int
		IsActive       bool
		ManufacturedBy string
		ManufacturedIn string
		RegAuthority   string
		RegIdentifier  string
	}
)
