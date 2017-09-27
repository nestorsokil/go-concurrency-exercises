//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"os"
	"os/signal"
)

func onSigint(terminationFunc func(), stop <-chan struct{}) {
	go func() {
		interruptChan := make(chan os.Signal, 1)
		defer close(interruptChan)
		signal.Notify(interruptChan, os.Interrupt)
		for {
			select {
			case <-interruptChan:
				killswitch := func() { os.Exit(1) }
				onSigint(killswitch, stop)
				terminationFunc()
				return
			case <-stop:
				return
			}
		}
	}()
}

func main() {
	// Create a process
	proc := MockProcess{}
	stop := make(chan struct{})
	defer close(stop)
	onSigint(proc.Stop, stop)
	// Run the process (blocking)
	proc.Run()
}
