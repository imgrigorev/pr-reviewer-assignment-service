package errors

import "errors"

var ErrTeamExists = errors.New("TEAM EXISTS")
var ErrTeamNotFound = errors.New("TEAM NOT FOUND")
var ErrUserNotFound = errors.New("USER NOT FOUND")
var ErrPrNotFound = errors.New("PR NOT FOUND")
var ErrPRAlreadyExists = errors.New("PR id already exists")
