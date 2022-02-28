package domain

// GPXFilePath represents the path to a GPX file
type GPXFilePath string

// String implements Stringer interface
func (g GPXFilePath) String() string {
	return string(g)
}
