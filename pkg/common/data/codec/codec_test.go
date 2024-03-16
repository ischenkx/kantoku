package codec

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testFromStruct struct {
	A string
	B int
	C map[string]int
}

func TestCodec(t *testing.T) {
	type testCase struct {
		name   string
		codec  Codec[testFromStruct, []byte]
		sample testFromStruct
	}
	tests := []testCase{
		{
			name:  "1",
			codec: JSON[testFromStruct](),
			sample: testFromStruct{
				A: "123",
				B: 42,
				C: map[string]int{
					"x": 46,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encodedSample, err := tt.codec.Encode(tt.sample)
			if err != nil {
				t.Fatal("failed to encode the sample:", err)
			}

			decodedSample, err := tt.codec.Decode(encodedSample)
			if err != nil {
				t.Fatal("failed to decode the encoded sample:", err)
			}

			assert.Equal(t, tt.sample, decodedSample)
		})
	}
}
