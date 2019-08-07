package brucecorebase

import "errors"

var (
	// ErrConfigNoAdaCoreServAddr - There is no AdaCoreServAddr in the configuration file
	ErrConfigNoAdaCoreServAddr = errors.New("There is no AdaCoreServAddr in the configuration file")
	// ErrConfigNoAdaCoreToken - There is no AdaCoreToken in the configuration file
	ErrConfigNoAdaCoreToken = errors.New("There is no AdaCoreToken in the configuration file")
)
