package channel

import (
	"container/list"
	"sync"
)

// ggkit:ignore
type Channel struct {
	mu    *sync.Mutex
	items *list.List
}
