package statemanage

const (
	EV_SPACKCHANGE = iota
    EV_CHANGE = iota // create and change
    EV_REBUILD = iota
    EV_GLOBALCHANGE = iota
)

type Event struct {
	Type int
	Data string
}

func NewEvent(evType int, evData string) *Event {
	return &Event{Type: evType, Data: evData}
}

func (e Event) String() string {
	var typeStr string

	switch(e.Type) {
	case EV_SPACKCHANGE:
		typeStr = "SPACKCHANGE"
	case EV_CHANGE:
		typeStr = "CHANGE"
	case EV_REBUILD:
		typeStr = "REBUILD"
	case EV_GLOBALCHANGE:
		typeStr = "GLOBALCHANGE"
	default:
		typeStr = ""
	}

	return "Event type: " + typeStr + " data: " + e.Data
}