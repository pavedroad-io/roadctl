package cmd

import (
	"net/url"

	"github.com/google/uuid"
)

// Types of data representations
const (
	Mime_Type = iota
	Avro
	Protobuf
	YAML
	WASM
)

// BlockResource
type BlockResource struct {
	// APIVersion for this object
	ApiVersion string

	// Kind/data representation
	Kind string

	// UUID unique ID for an object
	UUID uuid.UUID

	// LifeCycleEvents
	LifeCycleEvents []Event

	// Policies RBAC / ABAC including authentication
	Policies []string

	// Schema JSON or other object specific representation
	//   if the object is not self defining
	Schema string

	// SignatureType SHA, MD5, etc
	SignatureType string

	// Signature for anti-tampering
	Signature string

	// Mutable
	Mutable bool

	// Payload return payload for a resource
	Payload PayloadLoader

	// Events manage an event queue
	Events EventInteface
}

// PayloadLoader
//   Loads from either a external ULR or a pointer to the
//   data
type PayloadLoader struct {
	// ExternalURL nil if not set
	//   file:///~./mycode/artifacts/lint.txt
	externalURL url.URL

	dataPointer []byte
}

// Payload returns the contents of a resource
func (dl *PayloadLoader) Payload() ([]byte, error) {
	if dl.dataPointer != nil {
		return dl.dataPointer, nil
	}

	return nil, nil
	//	return LoadURL(externalURL)
}

type EventInteface struct {
	EventBusType string
	Brokers      []string
	ReadEvents   []Event
	WriteEvents  []Event
}

type Event struct {
	Topic string
}
