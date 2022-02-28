package domain

import (
	"bytes"
	"io"
)

type MapFile struct {
	content []byte
}

func NewMapFile(content []byte) MapFile {
	return MapFile{content: content}
}

func (f MapFile) File() io.Reader {
	return bytes.NewBuffer(f.content)
}
