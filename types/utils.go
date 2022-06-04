package types

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"crypto/sha256"

	"github.com/CosmWasm/wasmd/x/wasm/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// DefaultMaxWasmCodeSize limit max bytes read to prevent gzip bombs
	DefaultMaxWasmCodeSize = 600 * 1024 * 2
)

func FindEventByType(events []abci.Event, eventType string) (abci.Event, error) {
	for _, event := range events {
		if event.Type == eventType {
			return event, nil
		}
	}

	return abci.Event{}, fmt.Errorf("no event with type %s found", eventType)
}

func FindEventsByType(events []abci.Event, eventType string) []abci.Event {
	var found []abci.Event
	for _, event := range events {
		if event.Type == eventType {
			found = append(found, event)
		}
	}

	return found
}

func FindAttributeByKey(event abci.Event, attrKey string) (abci.EventAttribute, error) {
	for _, attr := range event.Attributes {
		if string(attr.Key) == attrKey {
			return attr, nil
		}
	}

	return abci.EventAttribute{}, fmt.Errorf("no attribute with key %s found inside event with type %s", attrKey, event.Type)
}

// magic bytes to identify gzip.
// See https://www.ietf.org/rfc/rfc1952.txt
// and https://github.com/golang/go/blob/master/src/net/http/sniff.go#L186
var gzipIdent = []byte("\x1F\x8B\x08")

// uncompress returns gzip uncompressed content or given src when not gzip.
func uncompress(src []byte, limit uint64) ([]byte, error) {
	switch n := uint64(len(src)); {
	case n < 3:
		return src, nil
	case n > limit:
		return nil, types.ErrLimit
	}
	if !bytes.Equal(gzipIdent, src[0:3]) {
		return src, nil
	}
	zr, err := gzip.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}
	zr.Multistream(false)
	defer zr.Close()
	return ioutil.ReadAll(LimitReader(zr, int64(limit)))
}

// LimitReader returns a Reader that reads from r
// but stops with types.ErrLimit after n bytes.
// The underlying implementation is a *io.LimitedReader.
func LimitReader(r io.Reader, n int64) io.Reader {
	return &LimitedReader{r: &io.LimitedReader{R: r, N: n}}
}

type LimitedReader struct {
	r *io.LimitedReader
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.r.N <= 0 {
		return 0, types.ErrLimit
	}
	return l.r.Read(p)
}

func GetCodeData(input []byte) (int, string, error) {
	uncompressedCode, err := uncompress(input, DefaultMaxWasmCodeSize)
	if err != nil {
		return 0, "", err
	}

	hasher := sha256.New()
	hasher.Write(uncompressedCode)
	codeHash := hasher.Sum(nil)
	hexHashString := strings.ToUpper(hex.EncodeToString(codeHash))

	codeSize := len(uncompressedCode)

	return codeSize, hexHashString, nil
}
