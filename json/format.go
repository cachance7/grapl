package json

import (
	"github.com/TylerBrock/colorjson"
)

// Pretty formats a json object to a pretty-printable byte array
func Pretty(obj interface{}) ([]byte, error) {
	f := colorjson.NewFormatter()
	f.Indent = 2
	s, _ := f.Marshal(obj)
	return s, nil
}
