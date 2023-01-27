package bvh

import (
	"fmt"
	"sync"

	"github.com/downflux/go-bvh/bvh"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-geometry/2d/hyperrectangle"
	"golang.org/x/sync/errgroup"

	hnd "github.com/downflux/go-geometry/nd/hyperrectangle"
)

type RO interface {
	BroadPhase(l Layer, q hyperrectangle.R) []id.ID
}

type O bvh.O

// Layer is a bitmask representing the different collision layers an object
// exists within.
type Layer uint16

type BVH struct {
	lookup map[id.ID]Layer
	layers [16]*bvh.T
}

func New(o O) *BVH {
	t := &BVH{
		lookup: make(map[id.ID]Layer, 256),
	}

	for i := 0; i < 16; i++ {
		t.layers[i] = bvh.New(bvh.O(o))
	}

	return t
}

func (bvh *BVH) Insert(x id.ID, l Layer, aabb hyperrectangle.R) {
	if _, ok := bvh.lookup[x]; ok {
		panic(fmt.Sprintf("cannot insert duplicate node: %v", x))
	}

	bvh.lookup[x] = l

	var wg errgroup.Group
	for i := 0; i < 16; i++ {
		if m := Layer(1 << i); m&l != 0 {
			wg.Go(func() error { return bvh.layers[i].Insert(x, hnd.R(aabb)) })
		}
	}

	if err := wg.Wait(); err != nil {
		panic(fmt.Sprintf("cannot insert node: %v", err))
	}
}

func (bvh *BVH) Remove(x id.ID) {
	l, ok := bvh.lookup[x]
	if !ok {
		panic(fmt.Sprintf("cannot remove non-existent node: %v", x))
	}

	var wg errgroup.Group
	for i := 0; i < 16; i++ {
		if m := Layer(1 << i); m&l != 0 {
			wg.Go(func() error { return bvh.layers[i].Remove(x) })
		}
	}

	if err := wg.Wait(); err != nil {
		panic(fmt.Sprintf("cannot remove node: %v", err))
	}

	delete(bvh.lookup, x)
}

func (bvh *BVH) Update(x id.ID, aabb hyperrectangle.R) {
	l, ok := bvh.lookup[x]
	if !ok {
		panic(fmt.Sprintf("cannot update non-existent node: %v", x))
	}

	var wg errgroup.Group
	for i := 0; i < 16; i++ {
		if m := Layer(1 << i); m&l != 0 {
			wg.Go(func() error { return bvh.layers[i].Update(x, hnd.R(aabb)) })
		}
	}

	if err := wg.Wait(); err != nil {
		panic(fmt.Sprintf("cannot update node: %v", err))
	}
}

func (bvh *BVH) BroadPhase(l Layer, q hyperrectangle.R) []id.ID {
	ch := make(chan id.ID, 256)

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		if m := Layer(1 << i); m&l != 0 {
			wg.Add(1)
			go func(ch chan<- id.ID) {
				defer wg.Done()
				for _, x := range bvh.layers[i].BroadPhase(hnd.R(q)) {
					ch <- x
				}
			}(ch)
		}
	}

	go func(ch chan id.ID) {
		defer close(ch)
		wg.Wait()
	}(ch)

	ids := make([]id.ID, 256)
	for _, x := range ids {
		ids = append(ids, x)
	}
	return ids
}
