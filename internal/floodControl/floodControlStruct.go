package floodcontrol

import (
	"fmt"
	"sync"
	"time"
)

type FloodControler struct {
	cache map[int64]UserRequest
	mu    sync.Mutex
}

type UserRequest struct {
	AmountOfRequests int
	lifetime         time.Time
}

type RequestExceeded struct {
	message string
}

func (e *RequestExceeded) Error() string {
	return fmt.Sprintf("Request exceeded: %s", e.message)
}

func RequestExceededError(userID int64) error {
	str := fmt.Sprintf("The limit of requests is exceeded  by  usedID  %d", userID)
	return &RequestExceeded{message: str}
}
