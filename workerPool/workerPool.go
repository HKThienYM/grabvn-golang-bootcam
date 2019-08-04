package workerpool

//Job interface
type Job interface {
	RunJob()
}

//Worker struct
type Worker struct {
	job  chan Job
	pool chan chan Job
	quit chan int
}

//Pool struct
type Pool struct {
	workers []Worker
	job     chan Job
	pool    chan chan Job
	quit    chan int
}

func (w Worker) start() {
	for {
		w.pool <- w.job

		select {
		case value := <-w.job:
			value.RunJob()
		case <-w.quit:
			return
		}
	}
}

func (p *Pool) start() {
	for {
		select {
		case job := <-p.job:
			worker := <-p.pool
			worker <- job
		case <-p.quit:
			for _, w := range p.workers {
				w.destroy()
			}
			return
		}

	}
}

//Dispatch add new job
func (p Pool) Dispatch(job Job) {
	p.job <- job
}

//Destroy destroy the pool worker
func (p Pool) Destroy() {
	p.quit <- 0
}

//Destroy destroy the worker
func (w Worker) destroy() {
	w.quit <- 0
}

// NewWorkerPool create new worker pool
func NewWorkerPool(count int) *Pool {
	p := Pool{
		pool: make(chan chan Job, count),
		job:  make(chan Job),
		quit: make(chan int),
	}
	for i := 0; i < count; i++ {
		p.workers = append(p.workers, Worker{
			pool: p.pool,
			job:  make(chan Job),
			quit: make(chan int),
		})
	}
	for _, worker := range p.workers {
		go worker.start()
	}
	go p.start()
	return &p
}
