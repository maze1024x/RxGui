package util

type ExitNotifier struct {
    signal  chan(struct{})
}
func MakeExitNotifier(signal chan(struct{})) ExitNotifier {
    return ExitNotifier { signal }
}
func (e ExitNotifier) Signal() <-chan(struct{}) {
    return e.signal
}
func (e ExitNotifier) Wait() {
    <- e.Signal()
}


