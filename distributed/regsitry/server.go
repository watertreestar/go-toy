package regsitry

import "sync"

type registry struct {
	registrations []Registration
	mutex         sync.Mutex
}
