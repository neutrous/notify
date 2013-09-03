// 
// Abstraction of communication environment.
// 
// Author: neutrous
// 
package entity

import (
	"os"
)

// Abstraction of communication environment.
type CommEnv interface {
	addEntity(obj *Endpoint) error
	removeEntity(obj *Endpoint)
	
	handleSignal(sig chan os.Signal)
	Close()
}



















