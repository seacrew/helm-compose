package compose

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
)

var b64 = base64.StdEncoding

var magicGzip = []byte{0x1f, 0x8b, 0x08}

// encodeComposeConfig encodes the config file returning a base64 encoded
// gzipped string representation, or error.
func encodeComposeConfig(config *Config) (string, error) {
	b, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	if _, err = w.Write(b); err != nil {
		return "", err
	}
	w.Close()

	return b64.EncodeToString(buf.Bytes()), nil
}

// decodeComposeConfig decodes the bytes of data into a compose
// config. Data must contain a base64 encoded gzipped string of a
// valid release, otherwise an error is returned.
func decodeComposeConfig(data string) (*Config, error) {
	// base64 decode string
	b, err := b64.DecodeString(data)
	if err != nil {
		return nil, err
	}

	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	b2, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	b = b2

	var config Config
	// unmarshal release object bytes
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
