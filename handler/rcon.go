package handler

type RCon interface {
	Command(cmd string) error
}
