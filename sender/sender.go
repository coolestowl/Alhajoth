package sender

import "github.com/coolestowl/Alhajoth/sites"

type Sender interface {
	Send(sites.Item) error
}
