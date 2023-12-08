package metricsjobs

import "sync"

type CreatingCountersJob struct {
	creatingCounters sync.WaitGroup
	jobs             int
}

func (job *CreatingCountersJob) Add(n int) {
	job.jobs += n
	job.creatingCounters.Add(n)
}

func (job *CreatingCountersJob) Done() {
	job.jobs -= 1
	job.creatingCounters.Done()
}

func (job *CreatingCountersJob) Wait() {
	job.creatingCounters.Wait()
}

func (job *CreatingCountersJob) IsActive() bool {
	return job.jobs > 0
}
