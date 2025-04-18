package worker

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/tobiashort/ansi"
	"github.com/tobiashort/isatty"
)

type Pool interface {
	GetWorker() Worker
	Wait()
}

type pool struct {
	cap     int
	workers []*worker
	mutex   sync.Mutex
}

func NewPool(cap int) Pool {
	pool := &pool{}
	pool.cap = cap
	pool.workers = make([]*worker, cap)
	pool.mutex = sync.Mutex{}
	if isatty.IsTerminal(os.Stdout) {
		for range cap {
			fmt.Print(ansi.EraseEntireLine)
			fmt.Println()
		}
		fmt.Print(ansi.MoveCursorUp(pool.cap))
	}
	for i := range cap {
		pool.workers[i] = &worker{
			num:  i + 1,
			pool: pool,
			done: true,
			msg:  "",
		}
		pool.print(pool.workers[i])
	}
	return pool
}

func (p *pool) GetWorker() Worker {
	i := 0
	for {
		mod := i % p.cap
		worker := p.workers[mod]
		if worker.done {
			worker.done = false
			return worker
		}
		i++
	}
}

func (p *pool) Wait() {
	for {
		allDone := true
		for _, w := range p.workers {
			allDone = allDone && w.done
		}
		if allDone {
			break
		}
	}
	if isatty.IsTerminal(os.Stdout) {
		fmt.Print(ansi.EraseFromCursorToEndOfScreen)
	}
}

func (p *pool) print(w *worker) {
	if isatty.IsTerminal(os.Stdout) {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		capLen := len(strconv.Itoa(p.cap))
		prefixFmt := fmt.Sprintf("[Worker %%%dd] ", capLen)
		fmt.Print(ansi.MoveCursorDown(w.num))
		fmt.Print(ansi.EraseEntireLine)
		fmt.Printf(prefixFmt, w.num)
		fmt.Print(w.msg)
		fmt.Print(ansi.MoveCursorToColumn(0))
		fmt.Print(ansi.MoveCursorUp(w.num))
	}
}

func (p *pool) log(w *worker) {
	p.mutex.Lock()
	if isatty.IsTerminal(os.Stdout) {
		fmt.Print(ansi.EraseEntireLine)
	}
	fmt.Printf("%s\n", w.msg)
	w.msg = ""
	if isatty.IsTerminal(os.Stdout) {
		for range p.cap {
			fmt.Print(ansi.EraseEntireLine)
			fmt.Println()
		}
		fmt.Print(ansi.MoveCursorUp(p.cap))
	}
	p.mutex.Unlock()
	for _, w := range p.workers {
		p.print(w)
	}
}
