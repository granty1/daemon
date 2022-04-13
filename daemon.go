package daemon

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	// Daemon process env
	DaemonEnv = "DAEMON_PROCESS=true"
)

var (
	// --daemon flag
	daemonFlag = flag.Bool("d", false, "run as daemon process")
)

func Run() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if !(*daemonFlag) {
		return
	}

	Start()
}

func Start() {
	var child_proc bool
	envs := os.Environ()
	for _, env := range envs {
		if env == DaemonEnv {
			child_proc = true
		}
	}

	if !child_proc {
		log.Println("boot proc:", os.Getpid())
		// start fork
		cmd := exec.Command(os.Args[0], os.Args[1:]...)

		envs = append(envs, DaemonEnv)
		fd, err := os.OpenFile("daemon.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			error(err.Error())
		}

		cmd.Stderr = fd
		cmd.Stdout = fd
		cmd.Env = append(cmd.Env, envs...)
		if err := cmd.Start(); err != nil {
			error(err.Error())
		}
		log.Println("start daemon proc:", cmd.Process.Pid)
		log.Println("echo log >> daemon.log")
		// parent process exit success
		os.Exit(0)
	}
}

func error(content ...string) {
	fmt.Fprintln(os.Stdout, content)
	os.Exit(1)
}
