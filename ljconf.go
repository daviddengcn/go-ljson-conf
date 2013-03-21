package ljconf

import (
	"fmt"
	"github.com/daviddengcn/go-villa"
	"github.com/daviddengcn/ljson"
	"strconv"
	"strings"
)

type Conf struct {
	path villa.Path
	db   map[string]interface{}
}

// Load reads configurations from a speicified file
func Load(path string) (conf *Conf) {
	conf = &Conf{
		path: villa.Path(path),
		db:   make(map[string]interface{}),
	}

	fin, err := conf.path.Open()
	if err != nil {
		// if file not exists, nothing read (but configuration still usable.
		return
	}
	defer fin.Close()

	dec := ljson.NewDecoder(fin)
	dec.Decode(&conf.db)

	return
}

// fetch a value or a map[string]interface{} as an interface{},
// returns nil if not found
func (c *Conf) get(key string) interface{} {
	parts := strings.Split(key, ".")
	var vl interface{} = c.db
	for _, p := range parts {
		mp, ok := vl.(map[string]interface{})
		if !ok {
			return nil
		}

		vl, ok = mp[p]
		if !ok {
			return nil
		}
	}

	return vl
}

// Interface retrieves a value as an interface{} of the key. def is returned
// if the value does not exist.
func (c *Conf) Interface(key string, def interface{}) interface{} {
	vl := c.get(key)
	if vl == nil {
		return def
	}

	return vl
}

// String retrieves a value as a string of the key. def is returned
// if the value does not exist or cannot be converted to a string(e.g. is an
// object).
func (c *Conf) String(key, def string) string {
	vl := c.get(key)
	if vl == nil {
		return def
	}

	switch vl.(type) {
	case string, float64, bool:
		return fmt.Sprint(vl)
	}

	return def
}

// Bool retrieves a value as a bool of the key. def is returned
// if the value does not exist or is not a bool. A string will be converted
// using strconv.ParseBool.
func (c *Conf) Bool(key string, def bool) bool {
	vl := c.get(key)
	if vl == nil {
		return def
	}

	switch v := vl.(type) {
	case bool:
		return v
	case string:
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}

	return def
}

// Int retrieves a value as a string of the key. def is returned
// if the value does not exist or is not a number. A float number will be
// round up to the closest interger. A string will be converted using
// strconv.ParseInt.
func (c *Conf) Int(key string, def int) int {
	vl := c.get(key)
	if vl == nil {
		return def
	}

	switch v := vl.(type) {
	case float64:
		return int(v + 0.5)
	case string:
		i, err := strconv.ParseInt(v, 0, 0)
		if err == nil {
			return int(i)
		}
	}

	return def
}

// Float retrieves a value as a float64 of the key. def is returned
// if the value does not exist or is not a number. A string will be converted
// using strconv.ParseFloat.
func (c *Conf) Float(key string, def float64) float64 {
	vl := c.get(key)
	if vl == nil {
		return def
	}

	switch v := vl.(type) {
	case float64:
		return v
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return f
		}
	}

	return def
}

// Object retrieves a value as a map[string]interface{} of the key. def is returned
// if the value does not exist or is not an object.
func (c *Conf) Object(key string, def map[string]interface{}) map[string]interface{} {
	vl := c.get(key)
	if vl == nil {
		return def
	}

	switch v := vl.(type) {
	case map[string]interface{}:
		return v
	}

	return def
}

// List retrieves a value as a slice of interface{} of the key. def is returned
// if the value does not exist or is not an array.
func (c *Conf) List(key string, def []interface{}) []interface{} {
	vl := c.get(key)
	if vl == nil {
		return def
	}

	switch v := vl.(type) {
	case []interface{}:
		return v
	}

	return def
}

// StringList retrieves a value as a slice of string of the key. def is returned
// if the value does not exist or is not an array. Elements of the array are
// converted to strings using fmt.Sprint.
func (c *Conf) StringList(key string, def []string) []string {
	vl := c.get(key)
	if vl == nil {
		return def
	}

	switch v := vl.(type) {
	case []interface{}:
		res := make([]string, 0, len(v))
		for _, el := range v {
			res = append(res, fmt.Sprint(el))
		}
		return res
	}

	return def
}

// IntList retrieves a value as a slice of int of the key. def is returned
// if the value does not exist or is not an array. Elements of the array are
// converted to int. Zero is used when converting failed.
func (c *Conf) IntList(key string, def []int) []int {
	vl := c.get(key)
	if vl == nil {
		return def
	}

	switch v := vl.(type) {
	case []interface{}:
		res := make([]int, 0, len(v))
		for _, el := range v {
			var e int
			switch et := el.(type) {
			case int:
				e = et
			case string:
				i, _ := strconv.ParseInt(et, 0, 0)
				e = int(i)
			case bool:
				if et {
					e = 1
				} else {
					e = 0
				}
			}
			res = append(res, e)
		}
		return res
	}

	return def
}
