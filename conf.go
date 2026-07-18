package confutil

import (
	// Built-in/core modules.
	"encoding"
	"encoding/json"
	"fmt"
	"os"
	"time"

	// Extended modules.
	"golang.org/x/sys/unix"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	// Third-party modules.
	"gopkg.in/yaml.v3"
)

// Function signature for a custom decoder function used with WithFile. It takes
// the file content as a byte slice and returns the decoded configuration data
// and any error encountered.
type DecoderFunc[T any] func([]byte) (*T, error)

// Structure representing a configuration object with support for JSON, YAML,
// and Protobuf files.
type Config[T any] struct {
	customDecoder     DecoderFunc[T]
	data              *T
	filePath          string
	fileType          string
	lastFileCheckTime time.Time
	lastFileModTime   time.Time
	refreshInterval   time.Duration
	tmpl              *T
}

// Returns a new Config object with the specified type template. The type
// template is used to determine the type of the configuration data to return
// (map, a specific struct type, etc.)
func NewConfig[T any](tmpl *T) *Config[T] {
	return &Config[T]{
		tmpl: tmpl,
	}
}

// Sets the JSON configuration file for the Config object and returns the Config
// object.
func (c *Config[T]) WithJSONFile(filePath string) *Config[T] {
	c.filePath = filePath
	c.fileType = "json"
	return c
}

// Sets the YAML configuration file for the Config object and returns the Config
// object. This module uses "gopkg.in/yaml.v3" to process YAML files.
func (c *Config[T]) WithYAMLFile(filePath string) *Config[T] {
	c.filePath = filePath
	c.fileType = "yaml"
	return c
}

// Sets the Proto configuration file for the Config object and returns the
// Config object.
func (c *Config[T]) WithProtoFile(filePath string) *Config[T] {
	c.filePath = filePath
	c.fileType = "proto"
	return c
}

// Sets the proto JSON configuration file for the Config object and returns the
// Config object. This is useful for Protobuf messages serialized in JSON
// format to ensure oneof fields, etc., are correctly handled.
func (c *Config[T]) WithProtoJSONFile(filePath string) *Config[T] {
	c.filePath = filePath
	c.fileType = "protojson"
	return c
}

// Sets the binary configuration file for the Config object and returns the
// Config object. This requires the encoding.BinaryUnmarshaler interface for
// unmarshaling the configuration file content, meaning the type passed to
// NewConfig must support the encoding.BinaryUnmarshaler interface.
func (c *Config[T]) WithBinaryFile(filePath string) *Config[T] {
	c.filePath = filePath
	c.fileType = "binary"
	return c
}

// Sets the text configuration file for the Config object and returns the Config
// object. This requires the encoding.TextUnmarshaler interface for
// unmarshaling the configuration file content, meaning the type passed to
// NewConfig must support the encoding.TextUnmarshaler interface.
func (c *Config[T]) WithTextFile(filePath string) *Config[T] {
	c.filePath = filePath
	c.fileType = "text"
	return c
}

// Sets configuration file path and custom decoder function for the Config
// object. Returns the Config object.
func (c *Config[T]) WithFile(
	filePath string,
	decoder DecoderFunc[T],
) *Config[T] {
	c.filePath = filePath
	c.fileType = "custom"
	c.customDecoder = decoder
	return c
}

// Sets the refresh interval for the Config object and returns the Config
// object. If the refresh interval is greater than zero, the configuration will
// be automatically reloaded when it changes. The configuration file will be
// checked after the specified refresh interval when calling GetConf.
func (c *Config[T]) WithRefresh(d time.Duration) *Config[T] {
	c.refreshInterval = d
	return c
}

// Retrieves the current configuration. If the refresh interval is set and the
// configuration file has changed since the last check, it will reload the
// configuration from the file. Returns the configuration data and any error
// encountered during the process.
func (c *Config[T]) GetConf() (*T, error) {
	if c.data != nil {
		if c.refreshInterval > 0 &&
			time.Since(c.lastFileCheckTime) >= c.refreshInterval {
			fileInfo, err := os.Stat(c.filePath)
			if err != nil {
				return nil, err
			}
			if !fileInfo.ModTime().After(c.lastFileModTime) {
				return c.data, nil
			}
		} else {
			return c.data, nil
		}
	}

	in_fh, err := os.Open(c.filePath)
	if err != nil {
		return nil, err
	}
	defer in_fh.Close()
	if err = unix.Flock(int(in_fh.Fd()), unix.LOCK_SH); err != nil {
		return nil, err
	}
	defer unix.Flock(int(in_fh.Fd()), unix.LOCK_UN)

	fileInfo, err := os.Stat(c.filePath)
	if err != nil {
		return nil, err
	}
	c.lastFileModTime = fileInfo.ModTime()
	c.lastFileCheckTime = time.Now()

	fileContent, err := os.ReadFile(c.filePath)
	if err != nil {
		return nil, err
	}

	data, err := c.unmarshal(fileContent)
	if err != nil {
		return nil, err
	}
	c.data = data

	return c.data, nil
}

func (c *Config[T]) unmarshal(fileContent []byte) (*T, error) {
	switch c.fileType {
	case "json":
		data := new(T)
		if err := json.Unmarshal(fileContent, data); err != nil {
			return nil, err
		}

		c.data = data
		return c.data, nil
	case "yaml":
		data := new(T)
		if err := yaml.Unmarshal(fileContent, data); err != nil {
			return nil, err
		}

		c.data = data
		return c.data, nil
	case "proto":
		data := new(T)
		protoData, ok := any(data).(proto.Message)
		if !ok {
			return nil, fmt.Errorf("type %T does not implement proto.Message", data)
		}
		if err := proto.Unmarshal(fileContent, protoData); err != nil {
			return nil, err
		}

		c.data = data
		return c.data, nil
	case "protojson":
		data := new(T)
		protoData, ok := any(data).(proto.Message)
		if !ok {
			return nil, fmt.Errorf("type %T does not implement proto.Message", data)
		}
		if err := protojson.Unmarshal(fileContent, protoData); err != nil {
			return nil, err
		}

		c.data = data
		return c.data, nil
	case "binary":
		data := new(T)
		binaryData, ok := any(data).(encoding.BinaryUnmarshaler)
		if !ok {
			return nil, fmt.Errorf("type %T does not implement "+
				"encoding.BinaryUnmarshaler", data)
		}
		if err := binaryData.UnmarshalBinary(fileContent); err != nil {
			return nil, err
		}

		c.data = data
		return c.data, nil
	case "text":
		data := new(T)
		textData, ok := any(data).(encoding.TextUnmarshaler)
		if !ok {
			return nil, fmt.Errorf("type %T does not implement "+
				"encoding.TextUnmarshaler", data)
		}
		if err := textData.UnmarshalText(fileContent); err != nil {
			return nil, err
		}

		c.data = data
		return c.data, nil
	case "custom":
		if c.customDecoder == nil {
			return nil, fmt.Errorf("custom decoder is not set")
		}
		data, err := c.customDecoder(fileContent)
		if err != nil {
			return nil, err
		}
		c.data = data
		return c.data, nil
	default:
		return nil, fmt.Errorf("unsupported file type: %s", c.fileType)
	}
}
