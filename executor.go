package geektime

import "sync"

type executor struct {
	tasks       chan task
	workerCount int

	exitChan chan struct{}
	wg       sync.WaitGroup
}

func newExecutor(workerCount, taskBufSize int) *executor {
	return &executor{
		tasks:       make(chan task, taskBufSize),
		workerCount: workerCount,
		exitChan:    make(chan struct{}),
	}
}

func (e *executor) start() {
	for i := 0; i < e.workerCount; i++ {
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()

			for {
				select {
				case t, ok := <-e.tasks:
					if !ok {
						return
					}
					t.run()
				case <-e.exitChan:
					return
				}
			}
		}()
	}
}

func (e *executor) execute(t task) {
	e.tasks <- t
}

func (e *executor) shutdown() {
	close(e.exitChan)
	e.wg.Wait()
}
