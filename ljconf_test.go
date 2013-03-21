package ljconf

import (
	"reflect"
	"testing"
)

func TestBasic(t *testing.T) {
	cf := Load("fortest.conf")
	// a case: ["key", "def", "exp"]
	cases := [][3]interface{}{
		// string
		{"http", "", ""},
		{"http.proxy", "", "proxy.example.com"},
		{"http.true", "", "true"},
		{"http.false", "", "false"},
		// bool		
		{"http.true", false, true},
		{"http.truestr", false, true},
		{"http.false", true, false},
		{"http.wrong", true, true},
		{"http.wrong", false, false},
		// int
		{"http.port", 1234, 8080},
		{"http.portstr", 1234, 8080},
		// float
		{"http.port", 1234., 8080.},
		{"http.portstr", 1234., 8080.},
		// list
		{"http.users", []interface{}(nil), []interface{}{"apple", "banana", "cat", "david"}},
		// string-list
		{"http.users", []string(nil), []string{"apple", "banana", "cat", "david"}},
		// int-list
		{"http.users", []int(nil), []int{0, 0, 0, 0}},
	}

	for _, c := range cases {
		key := c[0].(string)
		switch exp := c[2].(type) {
		case string:
			def := c[1].(string)
			act := cf.String(key, def)

			if act != exp {
				t.Errorf("[%s]: expected %v, but got %s", key, exp, act)
			}
		case bool:
			def := c[1].(bool)
			act := cf.Bool(key, def)

			if act != exp {
				t.Errorf("[%s]: expected %v, but got %s", key, exp, act)
			}
		case int:
			def := c[1].(int)
			act := cf.Int(key, def)

			if act != exp {
				t.Errorf("[%s]: expected %v, but got %s", key, exp, act)
			}
		case []interface{}:
			def := c[1].([]interface{})
			act := cf.List(key, def)

			if !reflect.DeepEqual(act, exp) {
				t.Errorf("[%s]: expected %v, but got %s", key, exp, act)
			}
		case []string:
			def := c[1].([]string)
			act := cf.StringList(key, def)

			if !reflect.DeepEqual(act, exp) {
				t.Errorf("[%s]: expected %v, but got %s", key, exp, act)
			}
		case []int:
			def := c[1].([]int)
			act := cf.IntList(key, def)

			if !reflect.DeepEqual(act, exp) {
				t.Errorf("[%s]: expected %v, but got %s", key, exp, act)
			}
		}
	}
}
