package pool

import (
	"net/http"

	"github.com/ipfs-search/ipfs-search/components/extractor"
	"github.com/ipfs-search/ipfs-search/components/extractor/nsfw"
	"github.com/ipfs-search/ipfs-search/components/extractor/tika"
	"github.com/ipfs-search/ipfs-search/components/protocol"
	"github.com/ipfs-search/ipfs-search/utils"
)

func (p *Pool) getExtractors(protocol protocol.Protocol) []extractor.Extractor {
	// Limited extractor connections (as resources are generally known to be available by now)
	extractorTransport := utils.GetHTTPTransport(p.dialer.DialContext, p.config.Workers.MaxExtractorConns)

	getter := utils.NewHTTPBodyGetter(&http.Client{Transport: extractorTransport}, p.Instrumentation)

	tikaExtractor := tika.New(p.config.TikaConfig(), getter, protocol, p.Instrumentation)
	nsfwExtractor := nsfw.New(p.config.NSFWConfig(), getter, p.Instrumentation)

	return []extractor.Extractor{tikaExtractor, nsfwExtractor}
}
