// +build go1.9

package qs

import "sync"

func newSyncMap() syncMap {
	return syncMap{}
}

type syncMap struct {
	sync.Map
}
