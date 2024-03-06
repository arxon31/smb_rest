package sessions

import "errors"

var errNoSessionAvailable = errors.New("no session available")
var errUnableConnectToHost = errors.New("unable to connect to host")
