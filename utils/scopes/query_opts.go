package scopes

type QueryOpts struct {
	Page     int
	PageSize int
	Order    string
	OrderDir Direction
}
