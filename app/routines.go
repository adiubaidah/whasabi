package app

import (
	"sync"
)

var ActiveRoutines = make(map[string]chan struct{})

// var ActiveAuthenticationWaRoutine = make(map[string]<-chan whatsmeow.QRChannelItem)
var Mu sync.Mutex // Protects access to the map
