package fastrand

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"sync"
)

const numRand = 100

type threadSafeSource struct {
	lk  sync.Mutex
	src rand.Source
}

func newThreadSafeSource(src rand.Source) *threadSafeSource {
	return &threadSafeSource{
		src: src,
	}
}
func (s *threadSafeSource) Int63() (n int64) {
	s.lk.Lock()
	n = s.src.Int63()
	s.lk.Unlock()
	return
}
func (s *threadSafeSource) Seed(seed int64) {
	s.lk.Lock()
	s.src.Seed(seed)
	s.lk.Unlock()
}

// The last *rand.Rand is for func selectRand.
var rr [numRand + 1]*rand.Rand

func seedSecurely() (randInstance *rand.Rand, err error) {
	n, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	return rand.New(newThreadSafeSource(rand.NewSource(n.Int64()))), nil
}

func init() {
	var err error
	for i := range rr {
		rr[i], err = seedSecurely()
		if err != nil {
			panic(err)
		}
	}
}

func selectRand() *rand.Rand {
	return rr[rr[numRand].Int()%numRand]
}

func Intn(n int) int                   { return selectRand().Intn(n) }
func Int63n(n int64) int64             { return selectRand().Int63n(n) }
func Float64() float64                 { return selectRand().Float64() }
func Read(p []byte) (n int, err error) { return selectRand().Read(p) }
func Rand() *rand.Rand                 { return selectRand() }
