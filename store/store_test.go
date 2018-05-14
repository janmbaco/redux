package store

import (
	"testing"
)

func counter(state int, action string) int {
	switch action {
	case "INC":
		return state + 1
	case "DEC":
		return state - 1
	default:
		return state
	}
}

func TestStoreNew(t *testing.T) {
	store := /*store.*/ New(counter)
	if store == nil {
		t.Error("func New fail, expected it return a pointer to store instance")
	}
}

func TestStoreDispatch(t *testing.T) {
	store := /*store.*/ New(counter)
	store.Dispatch("INC")
	if store.state["counter"].Interface() != 1 {
		t.Error("counter can not work")
	}
}

func TestStoreGetState(t *testing.T) {
	store := /*store.*/ New(counter)
	store.Dispatch("DEC")
	if store.GetState("counter") != -1 {
		t.Error("GetState should return reducer's state")
	}
}