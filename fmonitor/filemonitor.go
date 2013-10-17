package fmonitor

import (
	"log"
	"github.com/Koshroy/fspoolr/spoolr"
	"github.com/Koshroy/fspoolr/statemanage"
	"github.com/howeyc/fsnotify"
	"path/filepath"
)


type filemonitor struct {
	artifactSet map[spoolr.Artifact]bool
	watcher *fsnotify.Watcher
	started bool
	eventChan chan<- *statemanage.Event
}

func NewFilemonitor() (*filemonitor, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &filemonitor{artifactSet: make(map[spoolr.Artifact]bool), 
		watcher: w, started: false, eventChan: nil }, nil
}

func (f *filemonitor) AddArtifact(a spoolr.Artifact) {
	log.Println("adding artifact", a)
	if _, ok := f.artifactSet[a]; !ok {
		f.artifactSet[a] = true
		if f.started {
			f.watcher.Watch(a.RootDir())
		}
	}
}

func (f *filemonitor) EventChan() (chan<- *statemanage.Event) {
	return f.eventChan
}

func (f *filemonitor) SetEventChan(c chan<- *statemanage.Event)  {
	f.eventChan = c
}

func (f *filemonitor) Start() {
	if !f.started {
		log.Println("starting filemonitor")
		for a, _ := range f.artifactSet {
			log.Println("watching", a.RootDir())
			f.watcher.Watch(a.RootDir())
		}
		go f.process()
	}
}

func (f *filemonitor) Restart() {
	if f.started {
		f.watcher.Close()
	}
	f.Start()
}

func (f *filemonitor) Stop() {
	if f.started {
		f.watcher.Close()
	}
}

func (f *filemonitor) process() {
	for {
		select {
		case ev := <-f.watcher.Event:
			log.Println("event ocurred", ev)
			if ev.IsCreate() || ev.IsModify() || ev.IsRename() || ev.IsDelete() {
				log.Println("will resolve", ev.Name, "into absolute path")
				absEvPath, err := filepath.Abs(ev.Name)
				if err != nil {
					log.Println("could not find absolute path for path", ev.Name)
					continue
				}
				smEv := statemanage.NewEvent(statemanage.EV_CHANGE, absEvPath)
				if f.eventChan != nil {
					f.eventChan <- smEv
				}
			}
		case err := <-f.watcher.Error:
			log.Println("error ocurred", err)
		}
	}
}
