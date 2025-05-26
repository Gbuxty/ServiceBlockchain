package domain

import "errors"

var ErrShutdownAlreadyCalled = errors.New("shutdown already called")