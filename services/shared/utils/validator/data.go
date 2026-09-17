package validator

var (
	countries map[string]struct{} = map[string]struct{}{
		"brunei":      {},
		"cambodia":    {},
		"indonesia":   {},
		"laos":        {},
		"malaysia":    {},
		"myanmar":     {},
		"philippines": {},
		"singapore":   {},
		"thailand":    {},
		"timor leste": {},
		"vietnam":     {},
	}

	currencies map[string]struct{} = map[string]struct{}{
		"bnd": {},
		"cny": {},
		"eur": {},
		"gbp": {},
		"idr": {},
		"jpy": {},
		"khr": {},
		"krw": {},
		"lak": {},
		"mmk": {},
		"myr": {},
		"php": {},
		"sgd": {},
		"thb": {},
		"usd": {},
		"vnd": {},
	}

	days map[string]struct{} = map[string]struct{}{
		"sun": {},
		"mon": {},
		"tue": {},
		"wed": {},
		"thu": {},
		"fri": {},
		"sat": {},
	}

	dosageForms map[string]struct{} = map[string]struct{}{
		"capsule":        {},
		"cream":          {},
		"drop":           {},
		"elixir":         {},
		"foam":           {},
		"gel":            {},
		"infusion":       {},
		"inhalation":     {},
		"injection":      {},
		"liniment":       {},
		"lotion":         {},
		"ointment":       {},
		"orodispersible": {},
		"paste":          {},
		"patch":          {},
		"pessary":        {},
		"powder":         {},
		"solution":       {},
		"spray":          {},
		"strip":          {},
		"suppository":    {},
		"suspension":     {},
		"syrup":          {},
		"tablet":         {},
		"others":         {},
	}

	packUnits map[string]struct{} = map[string]struct{}{
		"bottle": {},
		"box":    {},
		"strip":  {},
	}
)
