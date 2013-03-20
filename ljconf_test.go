package ljconf

import (
	"testing"
	"reflect"
)

func TestBasic(t *testing.T) {
	cf := Load("fortest.conf")
	cases := [][3]interface{}{
		{"http", "", ""},
		{"http.proxy", "", "proxy.example.com"},
		{"http.true", "", "true"},
		{"http.false", "", "false"},
		
		{"http.true", false, true},
		{"http.false", true, false},
		{"http.wrong", true, true},
		{"http.wrong", false, false},
		
		{"http.port", 1234, 8080},
		
		{"http.port", 1234., 8080.},
		
		{"http.users", []interface{}(nil), []interface{}{"apple", "banana", "cat", "david"}},
		
		{"http.users", []string(nil), []string{"apple", "banana", "cat", "david"}},
		
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
