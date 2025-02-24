package internal

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
)

// Values represents a collection of chart values.
type Values map[string]interface{}

// ReadValues will parse YAML byte data into a Values.
func ReadValues(data []byte) (vals Values, err error) {
	err = yaml.Unmarshal(data, &vals)
	if len(vals) == 0 {
		vals = Values{}
	}
	return
}

// ReadValuesFile will parse a YAML file into a map of values.
func ReadValuesFile(filename string) (Values, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return ReadValues(data)
}

// istable is a special-purpose function to see if the present thing matches the definition of a YAML table.
func istable(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}

// CoalesceValues merges Values, dst is the authoritative values and takes precedence
func CoalesceValues(dst, src Values) Values {
	return Values(coalesceTables(dst, src))
}

func coalesceTables(dst, src map[string]interface{}) map[string]interface{} {
	// Because dest has higher precedence than src, dest values override src
	// values.
	for key, val := range src {
		if istable(val) {
			if innerdst, ok := dst[key]; !ok {
				dst[key] = val
			} else if istable(innerdst) {
				coalesceTables(innerdst.(map[string]interface{}), val.(map[string]interface{}))
			} else {
				fmt.Fprintf(os.Stderr, "Warning: Attempt to override Map value [%s] with String value [%s] for key [%s].\n", val, innerdst, key)
			}
			continue
		} else if dv, ok := dst[key]; ok && istable(dv) {
			fmt.Fprintf(os.Stderr, "Warning: Attempt to override String value [%s] with Map values [%s] for key [%s].\n", val, dv, key)
			continue
		} else if !ok { // <- ok is still in scope from preceding conditional.
			dst[key] = val
			continue
		}
	}
	return dst
}
