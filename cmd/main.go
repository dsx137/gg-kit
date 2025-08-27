package main

import "github.com/dsx137/gg-kit/internal/concurrent"

func main() {
	at := concurrent.NewAtomic[int]()
	at.Store(42)
}
