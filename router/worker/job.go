package worker

import (
	"context"
	"io"
	"os/exec"
)

type Job struct {
	Cmd    *exec.Cmd
	Status string

	Output string

	cancel context.CancelFunc
}

func SetupJob(chunks []string, stderr, stdout io.ReadWriter) *Job {
	command, args := chunks[0], chunks[0:]

	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, command, args...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return &Job{
		Cmd:    cmd,
		Status: "created",
		cancel: cancel,
	}
}

func (j *Job) Run() error {
	j.Status = "started"

	if err := j.Cmd.Run(); err != nil {
		j.Status = "error"

		return err
	}

	j.Status = "finished"

	return nil
}
