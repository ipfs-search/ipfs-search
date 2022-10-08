package types

import (
	"bytes"

	cbor "github.com/fxamacker/cbor/v2"
	lz4 "github.com/pierrec/lz4/v4"
)

// Reference represents a named reference to a Document.
type Reference struct {
	_          struct{} `cbor:",toarray"`
	ParentHash string   `json:"parent_hash"`
	Name       string   `json:"name"`
}

// References is a collection of references to a Document.
type References []Reference

// MarshalBinary marshalls into LZ4 compressed BSON.
func (r References) MarshalBinary() ([]byte, error) {
	data, err := cbor.Marshal([]Reference(r))
	if err != nil {
		return nil, err
	}

	compressed := new(bytes.Buffer)
	writer := lz4.NewWriter(compressed)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	return compressed.Bytes(), nil
}

// UnmarshalBinary unmarshalls from LZ4 compressed BSON.
func (r *References) UnmarshalBinary(data []byte) error {
	compressed := bytes.NewBuffer(data)
	uncompressed := new(bytes.Buffer)

	reader := lz4.NewReader(compressed)
	if _, err := reader.WriteTo(uncompressed); err != nil {
		return err
	}

	return cbor.Unmarshal(uncompressed.Bytes(), r)
}
