package logging

const (
	SignalRotate = iota
	SignalShutdown
	SignalSetSigCb
	SignalSetWriter
	SignalSetFormatter
	SignalReopen
)

type LogSignal struct {
	Action  int
	Payload interface{}
}
