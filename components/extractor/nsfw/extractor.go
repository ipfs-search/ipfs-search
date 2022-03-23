package nsfw

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"

	"github.com/ipfs-search/ipfs-search/components/extractor"

	indexTypes "github.com/ipfs-search/ipfs-search/components/index/types"
	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
)

// Extractor extracts metadata using the nsfw-server.
type Extractor struct {
	config *Config
	client *http.Client

	*instr.Instrumentation
}

func (e *Extractor) get(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		// Errors here are programming errors.
		panic(fmt.Sprintf("creating request: %s", err))
	}

	return e.client.Do(req)
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
		span.RecordError(
			ctx, extractor.ErrFileTooLarge, trace.WithErrorStatus(codes.Error),
			// TODO: Enable after otel upgrade.
			// label.Int64("file.size", r.Size),
		)
		return err
	}

	resp, err := e.get(ctx, e.getExtractURL(r))
	if err != nil {
		err := fmt.Errorf("%w: %v", extractor.ErrRequest, err)
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("%w: unexpected status %s", extractor.ErrUnexpectedResponse, resp.Status)
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return err
	}

	var nsfwData indexTypes.NSFW
	if err := json.NewDecoder(resp.Body).Decode(&nsfwData); err != nil {
		err := fmt.Errorf("%w: decoding error %s", extractor.ErrUnexpectedResponse, err)
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return err
	}

	file.NSFW = &nsfwData

	log.Printf("Got nsfw metadata metadata for '%v'", r)
	return nil
}

// New returns a new nsfw-server extractor.
func New(config *Config, client *http.Client, instr *instr.Instrumentation) extractor.Extractor {
	return &Extractor{
		config,
		client,
		instr,
	}
}

// Compile-time assurance that implementation satisfies interface.
var _ extractor.Extractor = &Extractor{}
