package api

type API struct {
	options Options
}

func New(options *Options) (API, error) {
	a := API{}

	return a, nil
}
