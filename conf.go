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
	"google.golang.org/protobuf/proto"

	// Third-party modules.
	"gopkg.in/yaml.v3"
)

// Structure representing a configuration object with support for JSON, YAML,
// and Protobuf files.
type Config[T any] struct {
	data              *T
	filePath          string
	fileType          string
	lastFileCheckTime time.Time
	lastFileModTime   time.Time
	refreshInterval   time.Duration
	tmpl              *T
}

// Returns a new Config object with the specified type template.
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
// object.
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

// Sets the refresh interval for the Config object and returns the Config
// object. If the refresh interval is greater than zero, the configuration will
// be automatically reloaded when it changes. The configuration file will be
// checked after the specified refresh interval when calling GetConf.
func (c *Config[T]) WithRefresh(d time.Duration) *Config[T] {
	c.refreshInterval = d
	return c
}

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
	default:
		return nil, fmt.Errorf("unsupported file type: %s", c.fileType)
	}
}
