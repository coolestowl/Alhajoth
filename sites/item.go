package sites

import "time"

type Item interface {
	ID() string

	Name() string

	Text() string

	Images() []string

	CreatedAt() time.Time
}
