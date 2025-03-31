package main

import (
	"errors"
	"os/exec"
	"syscall"
)

type Server struct {
	Path      string    `json:"path"`
	RunScript string    `json:"run_script"`
	Command   *exec.Cmd `json:"-"`
}

func NewServer(path string, run string) *Server {
	return &Server{
		path,
		run,
		nil,
	}
}

func (server *Server) Start() error {
	java := exec.Command(server.RunScript)
	java.Dir = server.Path
	return java.Start()
}

func (server *Server) Stop() error {
	var err error
	if server.Command == nil {
		return errors.New("command is nil")
	}
	serverProcess := server.Command.Process
	err = serverProcess.Signal(syscall.Signal(0))
	if serverProcess == nil || err != nil {
		return errors.New("process not found")
	}
	err = serverProcess.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) Reload() error {
	var err error
	if server.Command == nil {
		return errors.New("command is nil")
	}
	err = server.Stop()
	if err != nil {
		return err
	}
	err = server.Start()
	if err != nil {
		return err
	}
	return nil
}
