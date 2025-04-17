package worker

import (
	"strings"

	"github.com/tobiashort/cfmt"
)

type Worker interface {
	Done()
	Printf(format string, args ...any)
	Logf(format string, args ...any)
}

type worker struct {
	num  int
	pool *pool
	done bool
	msg  string
}

func (w *worker) Done() {
	w.done = true
}

func (w *worker) Printf(format string, args ...any) {
	w.msg = cfmt.Sprintf(format, args...)
	w.msg = strings.TrimSpace(w.msg)
	w.pool.print(w)
}

func (w *worker) Logf(format string, args ...any) {
	w.msg = cfmt.Sprintf(format, args...)
	w.msg = strings.TrimSpace(w.msg)
	w.pool.log(w)
}
