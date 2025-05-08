package mapping

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

const (
	nullTag      = "!!null"
	boolTag      = "!!bool"
	strTag       = "!!str"
	intTag       = "!!int"
	floatTag     = "!!float"
	timestampTag = "!!timestamp"
	seqTag       = "!!seq"
	mapTag       = "!!map"
	binaryTag    = "!!binary"
	mergeTag     = "!!merge"
)

var (
	ErrExpectedDocument   = errors.New("expected document node to a single document")
	ErrOddMapLength       = errors.New("expected mapping node to have an even number of elements")
	ErrUnexpectedNodeKind = errors.New("unexpected node type")
)

func YamlToJson(node *yaml.Node) ([]byte, error) {
	val, err := yamlNodeToValue(node)
	if err != nil {
		return nil, err
	}

	return json.Marshal(val)
}

func yamlNodeToValue(node *yaml.Node) (any, error) {
	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) == 1 {
			return yamlNodeToValue(node.Content[0])
		} else {
			return nil, ErrExpectedDocument
		}
	case yaml.ScalarNode:
		switch node.Tag {
		case nullTag:
			return nil, nil
		case strTag:
			return node.Value, nil
		case intTag:
			number, err := strconv.ParseInt(node.Value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parsing content node value as int: %w", err)
			}

			return number, err
		case floatTag:
			number, err := strconv.ParseFloat(node.Value, 64)
			if err != nil {
				return nil, fmt.Errorf("parsing content node value as float: %w", err)
			}

			return number, err
		case boolTag:
			boolValue, err := strconv.ParseBool(node.Value)
			if err != nil {
				return nil, fmt.Errorf("parsing content node value as boolean: %w", err)
			}

			return boolValue, nil
		}
	case yaml.SequenceNode:
		items := make([]any, len(node.Content))
		for idx, item := range node.Content {
			mappedVal, err := yamlNodeToValue(item)
			if err != nil {
				return nil, err
			}

			items[idx] = mappedVal
		}
		return items, nil
	case yaml.MappingNode:
		contentLength := len(node.Content)
		if contentLength%2 != 0 {
			return nil, ErrOddMapLength
		}

		itemMap := make(map[string]any, contentLength/2)

		for idx := 0; idx < contentLength; idx += 2 {
			val, err := yamlNodeToValue(node.Content[idx+1])
			if err != nil {
				return nil, err
			}

			itemMap[node.Content[idx].Value] = val
		}
		return itemMap, nil
	}

	return nil, fmt.Errorf("%w: %v", ErrUnexpectedNodeKind, node.Kind)
}
