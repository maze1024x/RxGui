package core

import (
	"time"
	"rxgui/standalone/qt"
	"rxgui/standalone/ctn"
)


type EventLoop struct {
	tasks   chan func()
	timers  *timerHeap
}
func CreateEventLoop() *EventLoop {
	var e = &EventLoop {
		tasks:  make(chan func(), 256),
		timers: createTimerHeap(),
	}
	go (func() {
		for {
			select {
			case k := <- e.tasks:
				qt.Schedule(k)
			default:
				select {
				case k := <- e.tasks:
					qt.Schedule(k)
				case <- e.timers.notify:
					if k, ok := e.timers.shift(); ok {
						qt.Schedule(k)
					}
				case f := <- e.timers.modify:
					f()
				}
			}
		}
	})()
	return e
}
func (e *EventLoop) addTask(k func()) {
	e.tasks <- k
}
func (e *EventLoop) addTimer(ms int, n int, ctx *context, k func()) {
	e.timers.modify <- func() {
		var dur = (time.Millisecond * time.Duration(ms))
		e.timers.insert(dur, n, ctx, k)
	}
}

type timerHeap struct {
	data    ctn.MutHeap[*timer]
	notify  <- chan time.Time
	modify  chan func()
}
type timer struct {
	baseTime  time.Time
	duration  time.Duration
	repeat    int
	context   *context
	callback  func()
}
func timerDue(a *timer) (bool,time.Duration) {
	var delta = time.Now().Sub(a.baseTime)
	return (delta >= a.duration), (a.duration - delta)
}
func timerLessThan(a *timer, b *timer) bool {
	return (a.baseTime.Sub(b.baseTime) < (b.duration - a.duration))
}
func createTimerHeap() *timerHeap {
	return &timerHeap {
		data:   ctn.MakeMutHeap(timerLessThan),
		notify: nil,
		modify: make(chan func(), 256),
	}
}
func (h *timerHeap) shift() (func(), bool) {
	var first, ok = h.data.First()
	if !(ok) {
		return nil, false
	}
	{ var due, dur = timerDue(first)
	if !(due) {
		h.notify = time.After(dur)
		return nil, false
	}}
	var this = first
	h.data.Shift()
	defer (func() {
		if first, ok := h.data.First(); ok {
			var due, dur = timerDue(first)
			if due {
				h.notify = time.After(0)
			} else {
				h.notify = time.After(dur)
			}
		}
	})()
	if this.context.isDisposed() {
		return nil, false
	}
	if this.repeat != 0 {
		var new_repeat_count int
		if this.repeat < 0 {
			// infinite
			new_repeat_count = this.repeat
		} else {
			new_repeat_count = (this.repeat - 1)
		}
		var repeat = &timer {
			baseTime: time.Now(),
			duration: this.duration,
			repeat:   new_repeat_count,
			context:  this.context,
			callback: this.callback,
		}
		h.data.Insert(repeat)
	}
	return this.callback, true
}
func (h *timerHeap) insert(dur time.Duration, n int, ctx *context, k func()) {
	if n == 0 {
		return
	}
	if dur < 0 {
		dur = 0
	}
	var a = &timer {
		baseTime: time.Now(),
		duration: dur,
		repeat:   (n - 1),
		context:  ctx,
		callback: k,
	}
	var first, has_first_before = h.data.First()
	h.data.Insert(a)
	if !(has_first_before) || (has_first_before && timerLessThan(a, first)) {
		h.notify = time.After(dur)
	}
}


