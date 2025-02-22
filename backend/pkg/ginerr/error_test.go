package ginerr

import (
	"reflect"
	"testing"
)

func TestBuild(t *testing.T) {
	if !reflect.DeepEqual(Build(""), ErrorResp{Status: "error", Message: ""}) {
		t.Error("Build() failed on empty string")
	}

	if !reflect.DeepEqual(Build("not found"), ErrorResp{Status: "error", Message: "not found"}) {
		t.Error("Build() failed on not found")
	}
}
