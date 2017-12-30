package store

import (
	"encoding/json"
	"fmt"
	"github.com/dannypsnl/redux/action"
	"strings"
	"testing"
)

var users = map[string]string{
	"danny": "1234",
}

func login(state interface{}, action action.Action) interface{} {
	if state == nil {
		return "Guest"
	}
	switch action.Type {
	case "login":
		if users[action.Args["user"].(string)] == action.Args["password"] {
			return action.Args["user"].(string) + " Login"
		}
		return state
	default:
		return state
	}
}

func TestSerialize(t *testing.T) {
	type StoreByCounterLogin struct {
		Counter int    `json:"counter"`
		Login   string `json:"login"`
	}
	store := New(counter, login)
	store.Subscribe(func() {
		// Use Subscribe to test is the most easy way to reduce lots of code
		expected := []byte(store.Marshal())
		expectedStore := &StoreByCounterLogin{}
		json.Unmarshal(expected, expectedStore)
		if store.GetState("counter") != expectedStore.Counter || store.GetState("login") != expectedStore.Login {
			t.Errorf("serialized result is not expected, expected:\n`%s`, actual:\n`%s`", expected, store.Marshal())
		}
	})
	store.Dispatch(action.New("INC"))
	store.Dispatch(action.New("INC"))
	store.Dispatch(
		action.New("login").
			Arg("user", "danny").
			Arg("password", 1234))
}

type file struct {
	// Ext mean extension
	Ext string `json:"ext"`
	// two mod: +x, +d
	Mod string `json:"mod"`
}

func fileUpdator(state interface{}, act action.Action) interface{} {
	if state == nil {
		return &file{
			Ext: "elz",
			Mod: "+x",
		}
	}
	switch act.Type {
	case "chmod":
		if act.Args["mod"] == "+x" && !strings.Contains(state.(file).Mod, "+x") {
			return &file{
				Ext: state.(file).Ext,
				Mod: state.(file).Mod + "+x",
			}
		} else if act.Args["mod"] == "+d" && !strings.Contains(state.(file).Mod, "+d") {
			return &file{
				Ext: state.(file).Ext,
				Mod: state.(file).Mod + "+d",
			}
		} else if act.Args["mod"] == "-d" {
			strings.Replace(state.(file).Mod, "+d", "", -1)
		} else if act.Args["mod"] == "-x" {
			strings.Replace(state.(file).Mod, "+x", "", -1)
		} else {
			return state
		}
	case "new ext":
		if act.Args["Ext"] != nil {
			return &file{
				Ext: act.Args["Ext"].(string),
				Mod: state.(file).Mod,
			}
		}
		return state
	default:
		return state
	}
	return state
}

// Make sure Marshal can work with Struct Type
func TestSerializeStruct(t *testing.T) {
	store := New(fileUpdator)
	expected := `{"fileUpdator":{"ext":"elz","mod":"+x"}}`
	if store.Marshal() != expected {
		t.Errorf("expected: %s, actual: %s", expected, store.Marshal())
	}
	store.Dispatch(
		action.New("new ext").
			Arg("ext", "cpp"))
}

func foo(state interface{}, act action.Action) interface{} {
	if state == nil {
		return func() {}
	}
	return state
}

func TestMarshalShouldPanic(t *testing.T) {
	store := New(foo)
	defer func() {
		if r := recover(); r == nil {
			t.Error(`Marshal didn't panic when data can't Marshal`)
		}
	}()
	store.Marshal()
}

// !!! If nest the struct, then the Marshal can't get the correct key
func TestJson(t *testing.T) {
	type NestData struct {
		I int `json:"integer"`
	}

	type TestData struct {
		Nest NestData `json:"nest"`
	}
	td := TestData{
		Nest: NestData{I: 10},
	}
	b, _ := json.Marshal(td)
	formatString := fmt.Sprintf("%s", b)
	expected := `{"nest":{"integer":10}}`
	if expected != formatString {
		t.Errorf("expected: %s, actual: %s", expected, formatString)
	}
}