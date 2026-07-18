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
)

// Tests JSON files where the configuration is represented as a struct and
// refresh is enabled.
func TestRefresh(t *testing.T) {
	type MyConf struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	expectedVal := &MyConf{
		Field1: "value1",
		Field2: 42,
	}

	confBytes := []byte(`{"field1":"value1","field2":42}`)
	tmpFile, err := os.CreateTemp("", "conf_test_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
		return
	}
	confFileName := tmpFile.Name()
	defer os.Remove(confFileName)

	if _, err := tmpFile.Write(confBytes); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
		return
	}

	confObj := confutil.NewConfig(&MyConf{}).
		WithJSONFile(confFileName).WithRefresh(2 * time.Second)

	confData, err := confObj.GetConf()
	if err != nil {
		t.Fatalf("failed to get config: %v", err)
		return
	}
	t.Logf("confData: %+v", confData)
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
	if err := os.WriteFile(confFileName, confBytes, 0644); err != nil {
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

// Tests configuration without refresh enabled.
func TestNoRefresh(t *testing.T) {
	type MyConf struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}
	expectedVal := &MyConf{
		Field1: "value1",
		Field2: 42,
	}

	confBytes := []byte(`{"field1":"value1","field2":42}`)
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

// Tests YAML files where the configuration is represented as a struct.
func TestYAMLConf(t *testing.T) {
	type MyConf struct {
		Field1 string `yaml:"field1"`
		Field2 int    `yaml:"field2"`
	}
	expectedVal := &MyConf{
		Field1: "value1",
		Field2: 42,
	}

	confBytes := []byte(`field1: "value1"
field2: 42
`)

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

// Tests JSON files where the configuration is represented as a map.
func TestJSONMapConf(t *testing.T) {
	type MyConf map[string]string
	expectedVal := &MyConf{
		"field1": "value1",
		"field2": "42",
	}

	confBytes := []byte(`{"field1":"value1","field2":"42"}`)
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

// Tests protobuf binary files.
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

// Tests protojson files (protobufs marshalled to JSON). The standard JSON
// decoder doesn't handle protobuf oneof fields "correctly", because Go uses
// interfaces to enforce oneof constraints, so WithProtoJSONFile indicates we
// should use the protojson package.
func TestProtoJSONConf(t *testing.T) {
	confText := `{
	"top_field": "value1",
	"field2": 42
}`

	// Go uses interfaces to enforce oneof constraints, so we need to set the
	// appropriate oneof field.
	expectedVal := &my_proto_conf.ProtoWithOneof{
		TopField: "value1",
		TestOneof: &my_proto_conf.ProtoWithOneof_Field2{
			Field2: 42,
		},
	}

	confBytes := []byte(confText)
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

	confObj := confutil.NewConfig(&my_proto_conf.ProtoWithOneof{}).
		WithProtoJSONFile(tmpFileName)

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

// Implements the BinaryUnmarshaller interface.
func (c *myBinaryConf) UnmarshalBinary(data []byte) error {
	if c.protoVal == nil {
		c.protoVal = &my_proto_conf.MyProtoConf{}
	}
	return proto.Unmarshal(data, c.protoVal)
}

// Tests binary files (parsed using the BinaryUnmarshaller interface). See the
// `WithBinaryFile` method.
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

// Implements the TextUnmarshaller interface.
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

// Tests text files (parsed using the TextUnmarshaller interface). See the
// `WithTextFile` method.
func TestTextConf(t *testing.T) {
	confBytes := []byte(`{"field1":"value1","field2":42}`)
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

// Tests custom decoding of configuration files with a custom decoder. See the
// `WithFile` method.
func TestCustomConf(t *testing.T) {
	type myCustomConf struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	decoder := func(data []byte) (*myCustomConf, error) {
		val := &myCustomConf{}
		if err := json.Unmarshal(data, val); err != nil {
			return nil, err
		}
		return val, nil
	}

	confBytes := []byte(`{"field1":"value1","field2":42}`)
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

	expectedVal := &myCustomConf{
		Field1: "value1",
		Field2: 42,
	}

	confObj := confutil.NewConfig(&myCustomConf{}).
		WithFile(tmpFileName, decoder)

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
