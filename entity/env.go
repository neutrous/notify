// 
// Abstraction of communication environment.
// 
// Author: neutrous
// 
package entity

import (
	"os"
)

// CommEnv abstracts the communication environment.
type CommEnv interface {
	addEntity(obj *endpoint) error
	removeEntity(obj *endpoint)
	
	handleSignal(sig chan os.Signal)
	Close()
}

