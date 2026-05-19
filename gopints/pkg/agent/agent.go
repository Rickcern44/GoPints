package agent

type Observer struct {
	chipName string
	pin      int
	handler  PulseHandler
}

func NewObserver(chipName string, pin int, handler PulseHandler) *Observer {
	return &Observer{
		chipName: chipName,
		pin:      pin,
		handler:  handler,
	}
}

func (o *Observer) Listen() {
}
