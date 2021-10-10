package window

type Event string

type EventListener func(Event, []interface{})

const (
	DrawEvent   Event = "draw"
	UpdateEvent Event = "update"
)
