package ljconf

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	cf, _ := Load("testdata/fortest.conf")
	t.Logf("Path: %v", cf.ConfPath())
	js, _ := json.MarshalIndent(cf.Object("", nil), "", "    ")
	t.Logf("Loaded: %v", string(js))
	// a case: ["key", "def", "exp"]
	cases := [][3]interface{}{
		// string
		{"http", "", ""},
		{"http.proxy", "", "proxy.example.com"},
		{"http.true", "", "true"},
		{"http.false", "", "false"},
		{"http.users[1]", "", "banana"},
		{"entries[0].apple.name", "<notfound>", "Apple"},
		{"entries[2][0]", "<notfound>", "apple"},
		{"entries[2][3][0]", "<notfound>", "david"},
		{"entries[2][3][1].apple.name", "<notfound>", "Apple"},
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
		{"http.nums", []int(nil), []int{1, -2, 3}},
		{"http.users", []int(nil), []int{0, 0, 0, 0}},
		// duration
		{"http.gap", "", "1m2s"},
		{"http.gap", time.Duration(0), 1*time.Minute + 2*time.Second},
		// time
		{"http.start", "", "2013-7-10 17:39:25"},
		{"http.start", time.Now(), time.Date(2013, 7, 10, 17, 39, 25, 0, time.UTC)},

		// included
		{"sub.value", "", "hello"},
		{"http.sub.value", "", "hello"},
	}

	for _, c := range cases {
		key := c[0].(string)
		switch exp := c[2].(type) {
		case string:
			def := c[1].(string)
			act := cf.String(key, def)

			if act != exp {
				t.Errorf("[%s]: expected %v, but got %v", key, exp, act)
			}
		case bool:
			def := c[1].(bool)
			act := cf.Bool(key, def)

			if act != exp {
				t.Errorf("[%s]: expected %v, but got %v", key, exp, act)
			}
		case int:
			def := c[1].(int)
			act := cf.Int(key, def)

			if act != exp {
				t.Errorf("[%s]: expected %v, but got %v", key, exp, act)
			}
		case float64:
			def := c[1].(float64)
			act := cf.Float(key, def)

			if act != exp {
				t.Errorf("[%s]: expected %v, but got %v", key, exp, act)
			}
		case []interface{}:
			def := c[1].([]interface{})
			act := cf.List(key, def)

			if !reflect.DeepEqual(act, exp) {
				t.Errorf("[%s]: expected %v, but got %v", key, exp, act)
			}
		case []string:
			def := c[1].([]string)
			act := cf.StringList(key, def)

			if !reflect.DeepEqual(act, exp) {
				t.Errorf("[%s]: expected %v, but got %v", key, exp, act)
			}
		case []int:
			def := c[1].([]int)
			act := cf.IntList(key, def)

			if !reflect.DeepEqual(act, exp) {
				t.Errorf("[%s]: expected %v, but got %v", key, exp, act)
			}
		case time.Duration:
			def := c[1].(time.Duration)
			act := cf.Duration(key, def)

			if act != exp {
				t.Errorf("[%s]: expected %v, but got %v", key, exp, act)
			}
		case time.Time:
			def := c[1].(time.Time)
			act := cf.Time(key, "2006-1-2 15:04:05", def)

			if act != exp {
				t.Errorf("[%s]: expected %v, but got %v", key, exp, act)
			}
		default:
			t.Errorf("Unknown type of %v", c[2])
		}
	}
}

func TestPath(t *testing.T) {
	cf, _ := Load("fortest.conf")
	t.Logf("Path: %v", cf.ConfPath())
	js, _ := json.MarshalIndent(cf.Object("", nil), "", "    ")
	t.Logf("Loaded: %v", string(js))
}

func TestFormatError(t *testing.T) {
	_, err := Load("testdata/wrongfmt.conf")
	if err == nil {
		t.Errorf("Wrong format error not reported")
	}
}

func TestConfNotExists(t *testing.T) {
	_, err := Load("testdata/nonexist.conf")
	if err != nil {
		t.Error(err)
	}
}

func TestDecode(t *testing.T) {
	type Apple struct {
		Name   string   `json:"name"`
		Weight string   `json:"weight"`
		Colors []string `json:"colors"`
	}

	var val Apple
	cf, _ := Load("testdata/fortest.conf")
	key := "entries[0].apple"
	cf.Decode(key, &val)
	t.Logf("%#v\n", val)
	if val.Name != "Apple" {
		t.Errorf("[%s]name: expected %v, got %v", key, "Apple", val.Name)
	}
	if val.Weight != "10kg" {
		t.Errorf("[%s]weight: expected %v, got %v", key, "10kg", val.Weight)
	}
	if val.Colors[1] != "green" {
		t.Errorf("[%s]color[1]: expected %v, got %v", key, "green", val.Colors[1])
	}

}

func TestSection(t *testing.T) {
	cf, _ := Load("testdata/fortest.conf")
	sec, _ := cf.Section("http")
	if sec.String("proxy", "") != "proxy.example.com" {
		t.Errorf("expected %v, got %v", "proxy.example.com", sec.String("proxy", ""))
	}

	// test hierarchy key name
	sec, _ = cf.Section("entries[1]")
	t.Logf("%#v", sec)
	if sec.String("sub.value", "") != "hello" {
		t.Errorf("expected %v, got %v", "hello", sec.String("sub.value", ""))
	}
}
