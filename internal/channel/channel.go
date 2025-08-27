package channel

import (
	"container/list"
	"sync"
)

type Channel struct {
	mu    *sync.Mutex
	items *list.List
}
