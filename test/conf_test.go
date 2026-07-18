package conf_test

import (
	// Built-in/core modules.
	"encoding/json"
	"os"
	"testing"
	"time"

	// Generated code.
	"github.com/cuberat-go/confutil/test/proto_stuff/my_proto_conf"

	// Third-party modules.
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	// First-party modules.
	"github.com/cuberat-go/confutil"
	"gopkg.in/yaml.v3"
)

func TestRefresh(t *testing.T) {
	type MyConf struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}
	expectedVal := &MyConf{
		Field1: "value1",
		Field2: 42,
	}

	confBytes, err := json.Marshal(expectedVal)
	if err != nil {
		t.Fatalf("failed to marshal config data: %v", err)
		return
	}
	tmpFile, err := os.CreateTemp("", "conf_test_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
		return
	}
	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)

	if _, err := tmpFile.Write(confBytes); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
		return
	}

	confObj := confutil.NewConfig(&MyConf{}).
		WithJSONFile(tmpFileName).WithRefresh(2 * time.Second)

	confData, err := confObj.GetConf()
	if err != nil {
		t.Fatalf("failed to get config: %v", err)
		return
	}
	if *confData != *expectedVal {
		t.Fatalf("config data mismatch: got %+v, want %+v", confData, expectedVal)
		return
	}

	expectedVal2 := &MyConf{
		Field1: "value2",
		Field2: 43,
	}
	confBytes, err = json.Marshal(expectedVal2)
	if err != nil {
		t.Fatalf("failed to marshal config data: %v", err)
		return
	}
	if err := os.WriteFile(tmpFileName, confBytes, 0644); err != nil {
		t.Fatalf("failed to write updated config to temp file: %v", err)
		return
	}

	confData, err = confObj.GetConf()
	if err != nil {
		t.Fatalf("failed to get updated config: %v", err)
		return
	}

	if *confData != *expectedVal {
		t.Fatalf("updated config data mismatch: got %+v, want %+v", confData,
			expectedVal)
		return
	}

	time.Sleep(3 * time.Second)

	// Now we should see the updated config.
	confData, err = confObj.GetConf()
	if err != nil {
		t.Fatalf("failed to get updated config: %v", err)
		return
	}
	if *confData != *expectedVal2 {
		t.Fatalf("updated config data mismatch: got %+v, want %+v", confData,
			expectedVal2)
		return
	}
}

func TestNoRefresh(t *testing.T) {
	type MyConf struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}
	expectedVal := &MyConf{
		Field1: "value1",
		Field2: 42,
	}

	confBytes, err := json.Marshal(expectedVal)
	if err != nil {
		t.Fatalf("failed to marshal config data: %v", err)
		return
	}
	tmpFile, err := os.CreateTemp("", "conf_test_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
		return
	}
	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)

	if _, err := tmpFile.Write(confBytes); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
		return
	}

	confObj := confutil.NewConfig(&MyConf{}).
		WithJSONFile(tmpFileName)

	confData, err := confObj.GetConf()
	if err != nil {
		t.Fatalf("failed to get config: %v", err)
		return
	}
	if *confData != *expectedVal {
		t.Fatalf("config data mismatch: got %+v, want %+v", confData, expectedVal)
		return
	}

	expectedVal2 := &MyConf{
		Field1: "value2",
		Field2: 43,
	}
	confBytes, err = json.Marshal(expectedVal2)
	if err != nil {
		t.Fatalf("failed to marshal config data: %v", err)
		return
	}
	if err := os.WriteFile(tmpFileName, confBytes, 0644); err != nil {
		t.Fatalf("failed to write updated config to temp file: %v", err)
		return
	}

	confData, err = confObj.GetConf()
	if err != nil {
		t.Fatalf("failed to get updated config: %v", err)
		return
	}

	if *confData != *expectedVal {
		t.Fatalf("updated config data mismatch: got %+v, want %+v", confData,
			expectedVal)
		return
	}
}

func TestYAMLConf(t *testing.T) {
	type MyConf struct {
		Field1 string `yaml:"field1"`
		Field2 int    `yaml:"field2"`
	}
	expectedVal := &MyConf{
		Field1: "value1",
		Field2: 42,
	}

	confBytes, err := yaml.Marshal(expectedVal)
	if err != nil {
		t.Fatalf("failed to marshal config data: %v", err)
		return
	}
	tmpFile, err := os.CreateTemp("", "conf_test_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
		return
	}
	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)

	if _, err := tmpFile.Write(confBytes); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
		return
	}

	confObj := confutil.NewConfig(&MyConf{}).
		WithYAMLFile(tmpFileName)

	confData, err := confObj.GetConf()
	if err != nil {
		t.Fatalf("failed to get config: %v", err)
		return
	}
	if *confData != *expectedVal {
		t.Fatalf("config data mismatch: got %+v, want %+v", confData, expectedVal)
		return
	}
}

func TestJSONMapConf(t *testing.T) {
	type MyConf map[string]string
	expectedVal := &MyConf{
		"field1": "value1",
		"field2": "42",
	}

	confBytes, err := json.Marshal(expectedVal)
	if err != nil {
		t.Fatalf("failed to marshal config data: %v", err)
		return
	}
	tmpFile, err := os.CreateTemp("", "conf_test_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
		return
	}
	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)

	if _, err := tmpFile.Write(confBytes); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
		return
	}

	confObj := confutil.NewConfig(&MyConf{}).
		WithJSONFile(tmpFileName)

	confData, err := confObj.GetConf()
	if err != nil {
		t.Fatalf("failed to get config: %v", err)
		return
	}
	if !assert.Equal(t, expectedVal, confData) {
		t.Fatalf("config data mismatch: got %+v, want %+v", confData, expectedVal)
		return
	}
}

func TestProtoConf(t *testing.T) {
	expectedVal := &my_proto_conf.MyProtoConf{
		Field1: "value1",
		Field2: 42,
	}

	confBytes, err := proto.Marshal(expectedVal)
	if err != nil {
		t.Fatalf("failed to marshal config data: %v", err)
		return
	}
	tmpFile, err := os.CreateTemp("", "conf_test_*.proto")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
		return
	}
	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)

	if _, err := tmpFile.Write(confBytes); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
		return
	}

	confObj := confutil.NewConfig(&my_proto_conf.MyProtoConf{}).
		WithProtoFile(tmpFileName)

	confData, err := confObj.GetConf()
	if err != nil {
		t.Fatalf("failed to get config: %v", err)
		return
	}
	if !proto.Equal(confData, expectedVal) {
		t.Fatalf("config data mismatch: got %+v, want %+v", confData, expectedVal)
		return
	}
}

type myBinaryConf struct {
	protoVal *my_proto_conf.MyProtoConf
}

func (c *myBinaryConf) UnmarshalBinary(data []byte) error {
	if c.protoVal == nil {
		c.protoVal = &my_proto_conf.MyProtoConf{}
	}
	return proto.Unmarshal(data, c.protoVal)
}

func TestBinaryConf(t *testing.T) {
	protoVal := &my_proto_conf.MyProtoConf{
		Field1: "value1",
		Field2: 42,
	}

	confBytes, err := proto.Marshal(protoVal)
	if err != nil {
		t.Fatalf("failed to marshal config data: %v", err)
		return
	}
	tmpFile, err := os.CreateTemp("", "conf_test_*.bin")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
		return
	}
	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)

	if _, err := tmpFile.Write(confBytes); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
		return
	}

	expectedVal := &myBinaryConf{
		protoVal: &my_proto_conf.MyProtoConf{
			Field1: "value1",
			Field2: 42,
		},
	}

	confObj := confutil.NewConfig(&myBinaryConf{}).
		WithBinaryFile(tmpFileName)

	confData, err := confObj.GetConf()
	if err != nil {
		t.Fatalf("failed to get config: %v", err)
		return
	}
	if !proto.Equal(confData.protoVal, expectedVal.protoVal) {
		t.Fatalf("config data mismatch: got %+v, want %+v", confData, expectedVal)
		return
	}
}

type myTextConf struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func (c *myTextConf) UnmarshalText(data []byte) error {
	type foo myTextConf
	val := &foo{}
	err := json.Unmarshal(data, val)
	if err != nil {
		return err
	}
	*c = myTextConf(*val)
	return nil
}

func TestTextConf(t *testing.T) {
	jsonVal := &myTextConf{
		Field1: "value1",
		Field2: 42,
	}

	confBytes, err := json.Marshal(jsonVal)
	if err != nil {
		t.Fatalf("failed to marshal config data: %v", err)
		return
	}
	tmpFile, err := os.CreateTemp("", "conf_test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
		return
	}
	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)

	if _, err := tmpFile.Write(confBytes); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
		return
	}

	expectedVal := &myTextConf{
		Field1: "value1",
		Field2: 42,
	}

	confObj := confutil.NewConfig(&myTextConf{}).
		WithTextFile(tmpFileName)

	confData, err := confObj.GetConf()
	if err != nil {
		t.Fatalf("failed to get config: %v", err)
		return
	}
	if !assert.Equal(t, confData, expectedVal) {
		t.Fatalf("config data mismatch: got %+v, want %+v", confData, expectedVal)
		return
	}
}
