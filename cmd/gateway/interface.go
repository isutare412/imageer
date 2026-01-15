package main

type leaderElector interface {
	Run()
	Shutdown()
}
