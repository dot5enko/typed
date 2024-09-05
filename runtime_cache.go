package typed

import (
	"fmt"
	"sync"
	"time"
)

var (
	ErrWrongKeyFromGeneratedData = fmt.Errorf("")
)

type cacheEntry[T any] struct {
	Created time.Time
	V       T
}

type RuntimeCache[T any, K comparable] struct {
	generator func(key K) (T, error)

	allocator map[K]*cacheEntry[T]

	options RuntimeCacheOptions

	syncMutex *sync.Mutex
}

type RuntimeCacheOptions struct {
	ThreadSafe bool
	Expiry     time.Duration
}

func NewRuntimeCache[T any, K comparable](
	gen func(key K) (T, error),
	options *RuntimeCacheOptions,
) RuntimeCache[T, K] {

	result := RuntimeCache[T, K]{
		allocator: make(map[K]*cacheEntry[T]),
		generator: gen,
	}

	if options != nil {

		result.options = *options

		if options.ThreadSafe {
			result.syncMutex = &sync.Mutex{}
		}
	}

	return result
}

func (rc RuntimeCache[T, K]) Get(key K) (result *T, topErr error) {

	defer RecoverPanic(func(pe *PanicError) {
		topErr = pe
	})

	if rc.options.ThreadSafe {
		rc.syncMutex.Lock()
		defer rc.syncMutex.Unlock()
	}

	val, ok := rc.allocator[key]

	expiryIn := time.Minute

	if rc.options.Expiry != 0 {
		expiryIn = rc.options.Expiry
	}

	if !ok || time.Since(val.Created).Milliseconds() > expiryIn.Milliseconds() {

		nval, errGen := rc.generator(key)
		if errGen != nil {
			topErr = errGen
			return
		}

		cachedEntryData := &cacheEntry[T]{
			V:       nval,
			Created: time.Now(),
		}

		rc.allocator[key] = cachedEntryData

		result = &cachedEntryData.V

		return

	} else {
		return &val.V, nil
	}

}
