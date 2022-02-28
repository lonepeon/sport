package domain

import (
	"bytes"
	"io"
)

type ShareableMapFile struct {
	content []byte
}

func NewSharableMapFile(content []byte) ShareableMapFile {
	return ShareableMapFile{content: content}
}

func (f ShareableMapFile) File() io.Reader {
	return bytes.NewBuffer(f.content)
}
