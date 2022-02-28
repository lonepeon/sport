package mapboxtest

import _ "embed"

//go:embed map1.png
var mapContent1 []byte

func GenerateMap() []byte {
	return mapContent1
}
