package config

import (
	"os"
	"reflect"
	"testing"
)

func TestEnvironment_BuildDsn(t *testing.T) {
	env := Environment{
		ServerAddress: "srv",
		DBHost:        "host",
		DBPort:        "port",
		DBUser:        "user",
		DBPassword:    "pass",
		DBName:        "name",
	}
	want := "host=host port=port user=user password=pass dbname=name sslmode=disable"
	if got := env.BuildDsn(); got != want {
		t.Errorf("BuildDsn() = %v, want %v", got, want)
	}
}

func TestLoadEnvironment(t *testing.T) {
	_ = os.Setenv("SERVER_ADDRESS", "srv")
	_ = os.Setenv("POSTGRES_HOST", "host")
	_ = os.Setenv("CI", "true")
	want := Environment{
		ServerAddress: "srv",
		DBHost:        "host",
		RunningInCI:   true,
	}
	if got := LoadEnvironment(); !reflect.DeepEqual(got, want) {
		t.Errorf("LoadEnvironment() = %v, want %v", got, want)
	}
}
