package terminux

type EventType int

const (
	MousePressed EventType = iota
	MouseReleased
	KeyPressed
	Error
)

type Event struct {
	Type EventType
	X, Y int
	Key  string
}

func newEventError() *Event {
	return newEvent(Error, -1, -1, "")
}

func newEvent(t EventType, x, y int, k string) *Event {
	return &Event{
		Type: t,
		X:    x,
		Y:    y,
		Key:  k,
	}
}

func (s *Screen) ListenForEvents() {
	buf := make([]byte, 64)

	for {
		n, _ := s.inputReader.Read(buf)
		if buf[0] == 0x1b {
			ev, err := parseANSIEvent(buf[:n])
			if err != nil {
				ev = newEventError()
				s.Debug(err, 0, 1)
			}
			s.events <- ev
			continue
		}
		s.events <- newEvent(KeyPressed, -1, -1, string(buf[0]))
	}
}

func (s *Screen) Events() <-chan *Event {
	return s.events
}

var eventRegexMap = make(map[EventType]string)

func initEventRegexMap() {
	eventRegexMap[MousePressed] = `\x1b\[<0;([0-9]+);([0-9]+)m`
	eventRegexMap[MouseReleased] = `\x1b\[<0;([0-9]+);([0-9]+)M`
}
