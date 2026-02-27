package session

import "io"

type Session struct{
	StdErr io.ReadWriter
	StdOut io.ReadWriter
}
