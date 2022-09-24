package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"

	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
)

// HTTPBodyGetter performs HTTP GET requests.
type HTTPBodyGetter interface {
	GetBody(ctx context.Context, url string, expect_status int) (io.ReadCloser, error)
}

type httpBodyGetterImpl struct {
	client *http.Client

	*instr.Instrumentation
}

func (g *httpBodyGetterImpl) get(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		// Errors here are programming errors.
		panic(fmt.Sprintf("creating request: %s", err))
	}

	return g.client.Do(req)
}

func (g *httpBodyGetterImpl) GetBody(ctx context.Context, url string, expect_status int) (io.ReadCloser, error) {
	ctx, span := g.Tracer.Start(ctx, "utils.HTTPBodyGetter.GetBody")
	defer span.End()

	resp, err := g.get(ctx, url)
	if err != nil {
		err := fmt.Errorf("%w: %v", t.ErrRequest, err)
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%w: unexpected status %s", t.ErrUnexpectedResponse, resp.Status)
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Error))
		resp.Body.Close()
		return nil, err
	}

	return resp.Body, err
}




// NewHTTPBodyGetter returns a new HTTPBodyGetter with specified client.
func NewHTTPBodyGetter(client *http.Client, instr *instr.Instrumentation) HTTPBodyGetter {
	return &httpBodyGetterImpl{
		client,
		instr,
	}
}
