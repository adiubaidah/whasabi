package app

import "sync"

var ActiveRoutines = make(map[string]chan struct{})
var Mu sync.Mutex // Protects access to the map
