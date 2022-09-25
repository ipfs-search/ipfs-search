package nsfw

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/ipfs-search/ipfs-search/components/extractor"

	indexTypes "github.com/ipfs-search/ipfs-search/components/index/types"
	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
	"github.com/ipfs-search/ipfs-search/utils"
)

// Extractor extracts metadata using the nsfw-server.
type Extractor struct {
	config *Config
	getter utils.HTTPBodyGetter

	*instr.Instrumentation
}

func (e *Extractor) getExtractURL(r *t.AnnotatedResource) string {
	return fmt.Sprintf("%s/classify/%s", e.config.NSFWServerURL, r.ID)
}

func matchOne(haystack string, needles []*regexp.Regexp) bool {
	for _, needle := range needles {
		if needle.MatchString(haystack) {
			return true
		}
	}

	return false
}

func getFileStringField(f *indexTypes.File, field string) string {
	v := f.Metadata[field] // Return zero-value if field not available.
	if v == nil {
		return "" // Return zero-value if field is not available.
	}

	switch x := v.(type) {
	case []interface{}:
		return x[0].(string)
	case string:
		return x
	default:
		// Panic if we're not a string
		panic(fmt.Sprintf("invalid type %T for field %s of metadata", v, field))
	}
}

var compatibleMimes = []*regexp.Regexp{
	regexp.MustCompile("^image/jpeg"),
	regexp.MustCompile("^image/png"),
	regexp.MustCompile("^image/gif"),
	regexp.MustCompile("^image/bmp"),
}

func isCompatible(r *t.AnnotatedResource, f *indexTypes.File) bool {
	// Check compatible protocol
	if r.Protocol != t.IPFSProtocol {
		return false
	}

	contentType := getFileStringField(f, "Content-Type")
	// log.Printf("Found Content-Type: %s", contentType)
	if contentType == "" {
		// No content-type set, assume we're not compatible
		return false
	}

	return matchOne(contentType, compatibleMimes)
}

// Extract metadata from a (potentially) referenced resource, updating
// Metadata or returning an error.
func (e *Extractor) Extract(ctx context.Context, r *t.AnnotatedResource, m interface{}) error {
	ctx, span := e.Tracer.Start(ctx, "extractor.nsfw_server.Extract")
	defer span.End()

	// Timeout if extraction hasn't fully completed within this time.
	ctx, cancel := context.WithTimeout(ctx, e.config.RequestTimeout)
	defer cancel()

	file := m.(*indexTypes.File) // Panics if we're not a File.

	if !isCompatible(r, file) {
		// log.Printf("Not extracting NSFW for incompatible %s", r)
		return nil
	}

	if r.Size > uint64(e.config.MaxFileSize) {
		err := fmt.Errorf("%w: %d", extractor.ErrFileTooLarge, r.Size)
		span.RecordError(extractor.ErrFileTooLarge) // TODO: Enable after otel upgrade.
		// attribute.Int64("file.size", r.Size),

		return err
	}

	body, err := e.getter.GetBody(ctx, e.getExtractURL(r), 200)
	if err != nil {
		return err
	}
	defer body.Close()

	var nsfwData indexTypes.NSFW
	if err := json.NewDecoder(body).Decode(&nsfwData); err != nil {
		err := fmt.Errorf("%w: decoding error %s", t.ErrUnexpectedResponse, err)
		span.RecordError(err)
		return err
	}

	// Success, update NSFW data.
	file.NSFW = &nsfwData

	log.Printf("Got nsfw metadata metadata for '%v'", r)
	return nil
}

// New returns a new nsfw-server extractor.
func New(config *Config, getter utils.HTTPBodyGetter, instr *instr.Instrumentation) extractor.Extractor {
	return &Extractor{
		config,
		getter,
		instr,
	}
}

// Compile-time assurance that implementation satisfies interface.
var _ extractor.Extractor = &Extractor{}
