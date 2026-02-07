package worker

import (
	"context"
	"os"
	"os/exec"
)

type Job struct {
	Cmd     *exec.Cmd
	Status  string

	Output string
}

func SetupJob(chunks []string) (*Job, context.CancelFunc) {
	command, args := chunks[0], chunks[0:]

	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	return &Job{
		Cmd: cmd,
		Status: "created",
	}, cancel
}

func (j *Job) Run() error {
	j.Status = "started"

	if err := j.Cmd.Run();err != nil{
		j.Status = "error"	

		return err
	}

	j.Status = "finished"

	return nil
}