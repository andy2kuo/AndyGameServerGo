package socketserver

import "errors"

var ErrClientIDDuplicate error = errors.New("client id duplicate")
var ErrClientStop error = errors.New("client close connect")
var ErrConnectTimeOut error = errors.New("connection time out")
