package worker

import "context"

type Job struct{
	State string
	Message string
	Cmd string
}

func NewJob(cmd string) *Job {
	return &Job{
		State: "created",
		Cmd: cmd,
	}
}

func (j *Job) Start(ctx context.Context) {

}