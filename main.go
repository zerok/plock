package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	flag "github.com/spf13/pflag"
)

func main() {
	var filePath string
	var exclusive bool
	var blocking bool
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: plock [options] file command [arguments]\n")
		flag.PrintDefaults()
	}
	flag.BoolVar(&exclusive, "exclusive", false, "Exclusive lock")
	flag.BoolVar(&blocking, "blocking", false, "Block instead of failing")
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	filePath = flag.Arg(0)
	file, err := os.OpenFile(filePath, syscall.O_CREAT|syscall.O_RDWR|syscall.O_CLOEXEC, 0600)
	if err != nil {
		log.Fatalf("Open failed: %v\n", err)
	}
	defer file.Close()
	lock := syscall.Flock_t{}
	if exclusive {
		lock.Type = syscall.F_WRLCK
	} else {
		lock.Type = syscall.F_RDLCK
	}
	op := syscall.F_SETLK
	if blocking {
		op = syscall.F_SETLKW
	}
	if err := syscall.FcntlFlock(file.Fd(), op, &lock); err != nil {
		log.Fatalf("Failed to run syscall: %s\n", err.Error())
	}
	defer unlock(file)
	args := flag.Args()
	var c *exec.Cmd
	if len(args) > 1 {
		c = exec.Command(flag.Arg(1), args[2:]...)
	} else {
		c = exec.Command(flag.Arg(1))
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		log.Fatalf("Command failed: %s", err.Error())
	}
}

func unlock(file *os.File) {
	lock := syscall.Flock_t{}
	lock.Type = syscall.F_UNLCK
	syscall.FcntlFlock(file.Fd(), syscall.F_SETLK, &lock)
}
