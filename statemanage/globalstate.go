package statemanage

import (
	"github.com/Koshroy/fspoolr/spoolr"
	"log"
	"path"
	"strings"
	"os/exec"
)

type globalState struct {
	resources []string // list of all resources
	resourceMap map[string]*Resource // maps resource name to Resource
	pathMap map[string]string // maps resource name to package.json path
	buildMap map[string]string // maps filename to build command

	resChan chan *Resource // w-only to satisfy resource requests
	reqChan chan string // r-only to receive resource requests
	eventChan chan *Event // r-only to receive events

	started bool // do not allow synchronous state changes
	shutdown bool // stop channel listening
}

func NewGlobalState() *globalState {
	return &globalState{
		resources: make([]string, 0), resourceMap: make(map[string]*Resource),
		pathMap: make(map[string]string), buildMap: make(map[string]string),
		resChan: make(chan *Resource), reqChan: make(chan string), eventChan: make(chan *Event),
		started: false, shutdown: false }
}

func runBuildCommand(buildCmd string) {
	buildCmdSplit := strings.Split(buildCmd, " ")
	cmd := exec.Command(buildCmdSplit[0], buildCmdSplit[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Println("error running [", buildCmd, "]")
		log.Println(err)
	}
}

func (g* globalState) processRequests() {
	for {
		if g.shutdown {
			break
		}
		_, ok := <-g.reqChan
		if !ok {
			log.Fatalln("could not read from requests channel")
			break
		}
	}
}

func (g* globalState) processEvents() {
	for {
		if g.shutdown {
			break
		}
		ev := <-g.eventChan
		log.Println("received event", ev)
		buildCmd, ok := g.buildMap[ev.Data]
		if !ok {
			log.Println("could not find event", ev, "in build map")
			continue
		}
		switch ev.Type {
		case EV_CHANGE:
			log.Println("running build command", buildCmd)
			runBuildCommand(g.buildMap[ev.Data])
		case EV_REBUILD:
			log.Println("running build command", buildCmd)
			runBuildCommand(g.buildMap[ev.Data])
		}
	}
}

func (g *globalState) processSettings() {
	return
}

func (g* globalState) flush() {
	g.resources = make([]string, 0)
	g.resourceMap = make(map[string]*Resource)
	g.pathMap = make(map[string]string)
	g.buildMap = make(map[string]string)
}

func (g* globalState) insertArtifact(ar spoolr.Artifact) {
	g.resources = append(g.resources, ar.Name())
	g.resourceMap[ar.Name()] = NewResource(ar.Target().File, ar.Target().MimeType)
	g.pathMap[ar.Name()] = path.Join(ar.RootDir(), "package.json")
	log.Println("ar.Files()", ar.Files())
	for _, elem := range ar.Files() {
		g.buildMap[elem] = ar.BuildCmd()
	}
}

func (g* globalState) AddArtifact(ar spoolr.Artifact) {
	if g.started == false {
		g.insertArtifact(ar)
	}
}

func (g* globalState) Start() {
	g.started = true
	go g.processEvents()
	go g.processRequests()
}

func (g* globalState) Stop() {
	g.shutdown = true
}

func (g *globalState) ResourceChan() <-chan *Resource {
	return g.resChan
}

func (g *globalState) RequestChan() chan<- string {
	return g.reqChan
}

func (g *globalState) EventChan() chan<- *Event {
	return g.eventChan
}