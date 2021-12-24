package types_test

import (
	"github.com/stafiprotocol/stafi-types"
	"testing"
)

type mockSarpc struct{}

func (m *mockSarpc) RegCustomTypes(content []byte) {}
func (m *mockSarpc) GetSystemChain() (string, error) {
	return "Development", nil
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, ctx ...interface{})  {}
func (m *mockLogger) Warn(msg string, ctx ...interface{})  {}
func (m *mockLogger) Error(msg string, ctx ...interface{}) {}

func TestNewTypes(t *testing.T) {
	typesBts, err := types.NewTypes(&mockSarpc{}, &mockLogger{}, 5, "http://127.0.0.1:9944")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(typesBts.GetStafiJsonTypes()))
	typesBts.StartMonitor()
	select {}
}
