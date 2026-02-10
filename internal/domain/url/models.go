package url

import "time"

type Url struct {
	Id          int
	OriginalUrl string
	Alias       string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Clicks      int
}
