/*
	Copyright 2017 - 2018 OneLedger
*/

package action

type CommandType int

const (
	NOOP CommandType = iota
	SUBMIT_TRANSACTION
	CREATE_LOCKBOX
	SIGN_LOCKBOX
	VERIFY_LOCKBOX
	SEND_KEY
	READ_CHAIN
	OPEN_LOCKBOX
	WAIT_FOR_CHAIN
)

// A command to execute again a chain, needs to be polymorphic
type Command struct {
	Function CommandType
	Data     map[string]string
}

type Commands []Command

func (commands Commands) Count() int {
	return len(commands)
}