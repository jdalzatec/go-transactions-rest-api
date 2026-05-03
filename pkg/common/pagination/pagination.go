package pagination

import (
	"github.com/oklog/ulid/v2"
)

type Pagination struct {
	Cursor  *ulid.ULID `json:"cursor"`
	HasMore bool       `json:"has_more"`
}

type Paginated[T any] struct {
	Items      []*T       `json:"items"`
	Pagination Pagination `json:"pagination"`
}
