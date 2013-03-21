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

const INCLUDE_KEY_TAG = "#include#"

func loadArrayInclude(arr []interface{}, dir villa.Path) {
	for _, el := range arr {
		switch vv := el.(type) {
		case map[string]interface{}:
			loadInclude(vv, dir)
		case []interface{}:
			loadArrayInclude(vv, dir)
		}
	}
}

func loadInclude(db map[string]interface{}, dir villa.Path) {
	for k, v := range db {
		if k == INCLUDE_KEY_TAG {
			switch paths := v.(type) {
			case string:
				//				fmt.Println("Including", paths, "at", dir)
				sub, err := Load(dir.Join(paths).S())
				if err == nil {
					// merge into current db
					for sk, sv := range sub.db {
						db[sk] = sv
					}
					// remove this entry
					delete(db, k)
				}
				continue
			case []interface{}:
				for _, el := range paths {
					if path, ok := el.(string); ok {
						sub, err := Load(dir.Join(path).S())
						if err == nil {
							// merge into current db
							for sk, sv := range sub.db {
								db[sk] = sv
							}
						}
					}
				}
				// remove this entry
				delete(db, k)
				continue
			} // switch
		} // if

		switch vv := v.(type) {
		case map[string]interface{}:
			loadInclude(vv, dir)
		case []interface{}:
			loadArrayInclude(vv, dir)
		}
	}
}

// Load reads configurations from a speicified file. If some error found
// during reading, it will be return, but the conf is still available.
func Load(fn string) (conf *Conf, err error) {
	path, _ := villa.Path(fn).Abs()
	conf = &Conf{
		path: path,
		db:   make(map[string]interface{}),
	}

	fin, err := path.Open()
	if err != nil {
		// if file not exists, nothing read (but configuration still usable.)
		return conf, err
	}
	func() {
		defer fin.Close()

		dec := ljson.NewDecoder(newRcReader(fin))
		dec.Decode(&conf.db)
	}()

	loadInclude(conf.db, path.Dir())

	return conf, nil
}

// fetch a value or a map[string]interface{} as an interface{},
// returns nil if not found
func (c *Conf) get(key string) interface{} {
	if key == "" {
		return c.db
	}
	parts := strings.Split(key, ".")
	var vl interface{} = c.db
	for _, p := range parts {
		mp, ok := vl.(map[string]interface{})
		if !ok {
			return nil
		}

		vl, ok = mp[p]
		if ok {
			continue
		}

		if strings.HasSuffix(p, "]") {
			// try fetch the element in an array
			idx := strings.Index(p, "[")
			if idx > 0 {
				indexes := strings.Split(p[idx+1:len(p)-1], "][")
				p = p[:idx]
				vl, ok = mp[p]
				if !ok {
					return nil
				}

				for _, sidx := range indexes {
					idx, err := strconv.ParseInt(sidx, 0, 0)
					if err != nil {
						return nil
					}

					arr, ok := vl.([]interface{})
					if !ok {
						return nil
					}

					if idx < 0 || int(idx) >= len(arr) {
						return nil
					}
					vl = arr[idx]
				}
			}
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

// floatToInt converts a float64 value into an int
func floatToInt(f float64) int {
	if f < 0 {
		return int(f - 0.5)
	}
	return int(f + 0.5)
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
		return floatToInt(v)
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
			case float64:
				e = floatToInt(et)
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
