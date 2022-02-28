package domain

// MapFilePath represents the path to a map file
type MapFilePath string

// String implements Stringer interface
func (m MapFilePath) String() string {
	return string(m)
}

// ShareableMapFilePath represents the path to a shareable map file
type ShareableMapFilePath string

// String implements Stringer interface
func (s ShareableMapFilePath) String() string {
	return string(s)
}
