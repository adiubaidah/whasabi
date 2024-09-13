package history

import "time"

type History struct {
	ID        int
	Sender    string
	Recipient string
	Content   string
	Role      string
	CreatedAt time.Time
}
