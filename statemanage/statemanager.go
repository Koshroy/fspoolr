package statemanage

import "github.com/Koshroy/fspoolr/spoolr"

type StateManager interface {
    AddArtifact(ar *spoolr.Artifact)
    Start()
    Stop()

	ResourceChan() <-chan *Resource
	RequestChan() chan<- string // resource string as requests
    EventChan() chan<- *Event
}