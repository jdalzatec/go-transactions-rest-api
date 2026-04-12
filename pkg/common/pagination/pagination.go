package pagination

import "github.com/google/uuid"

type Pagination struct {
	Cursor  uuid.UUID `json:"cursor"`
	HasMore bool      `json:"has_more"`
}

type Paginated[T any] struct {
	Items      []T        `json:"items"`
	Pagination Pagination `json:"pagination"`
}
