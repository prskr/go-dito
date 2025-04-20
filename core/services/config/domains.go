package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/prskr/go-dito/core/ports"
	"github.com/prskr/go-dito/core/services/parsing"
)

var _ json.Unmarshaler = (*DomainMapping)(nil)

var ErrUnknownSpecType = errors.New("unknown spec type")

type DomainMapping map[string]ports.SpecParser

// UnmarshalJSON implements json.Unmarshaler.
func (d *DomainMapping) UnmarshalJSON(data []byte) error {
	tmp := make(map[string]json.RawMessage)

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	for domain, rawSpec := range tmp {
		parsed, err := parseDomainSpec(rawSpec)
		if err != nil {
			return err
		}
		(*d)[domain] = parsed
	}

	return nil
}

func parseDomainSpec(rawSpec json.RawMessage) (ports.SpecParser, error) {
	tmp := struct {
		Type string `json:"type"`
	}{}

	if err := json.Unmarshal(rawSpec, &tmp); err != nil {
		return nil, err
	}

	switch tmp.Type {
	case "plain":
		var spec parsing.Plain
		return &spec, json.Unmarshal(rawSpec, &spec)
	case "openapi":
		var spec parsing.OpenAPI
		return &spec, json.Unmarshal(rawSpec, &spec)
	case "graphql":
		var spec parsing.GraphQL
		return &spec, json.Unmarshal(rawSpec, &spec)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownSpecType, tmp.Type)
	}
}
