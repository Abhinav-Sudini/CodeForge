package utils

import "bytes"

type BoundBuffer struct {
	Buf bytes.Buffer
	N   int64
}


func (l *BoundBuffer) Write(p []byte) (int, error) {
	if l.N <= 0 {
		return len(p), nil // pretend we wrote it
	}

	if int64(len(p)) > l.N {
		p = p[:l.N]
	}

	n, err := l.Buf.Write(p)
	l.N -= int64(n)
	return len(p), err
}
