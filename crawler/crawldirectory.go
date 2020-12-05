package crawler

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/rand"

	indexTypes "github.com/ipfs-search/ipfs-search/index/types"
	t "github.com/ipfs-search/ipfs-search/types"
)

const entryBufferSize = 256 // Size of buffer for processing channels. TODO: Make configurable.

func (c *Crawler) crawlDir(ctx context.Context, r *t.AnnotatedResource, properties *indexTypes.Directory) error {
	entries := make(chan *t.AnnotatedResource, entryBufferSize)

	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		err := c.processDirEntries(ctx, entries, properties)
		return err
	})

	wg.Go(func() error {
		err := c.protocol.Ls(ctx, r, entries)
		close(entries)
		return err
	})

	return wg.Wait()
}

func resourceIndexType(r *t.AnnotatedResource) indexTypes.LinkType {
	switch r.Type {
	case t.FileType:
		return indexTypes.FileLinkType
	case t.DirectoryType:
		return indexTypes.DirectoryLinkType
	default:
		// Behaviour for other types not (yet) defined.
		panic(fmt.Sprintf("Unsupported type returned from listing: %s", r.Type))
	}
}

func (c *Crawler) processDirEntries(ctx context.Context, entries <-chan *t.AnnotatedResource, properties *indexTypes.Directory) error {
	// Question: do we need a maximum entry cutoff point? E.g. 10^6 entries or something?
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case entry, ok := <-entries:
			if !ok {
				// Channel closed, we're done
				return nil
			}

			// TODO: Implement timeout waiting for new directory entries here.

			// Add link
			properties.Links = append(properties.Links, indexTypes.Link{
				Hash: entry.ID,
				Name: entry.Reference.Name,
				Size: entry.Size,
				Type: resourceIndexType(entry),
			})

			// Queue directory entry. Fail hard on error; prefer less over incomplete
			// or inconsistent data.
			// TODO: Consider doing this in a separate Goroutine, as it's blocking.
			err := c.queueDirEntry(ctx, entry)
			if err != nil {
				return err
			}
		}

	}
}

func (c *Crawler) queueDirEntry(ctx context.Context, r *t.AnnotatedResource) error {
	// Generate random lower priority for items in this directory
	// Rationale; directories might have different availability but
	// within a directory, items are likely to have similar availability.
	// We want consumers to get a varied mixture of availability, for
	// consistent overall indexing load.
	priority := uint8(1 + rand.Intn(7))

	switch r.Type {
	case t.FileType:
		return c.queues.Files.Publish(ctx, r, priority)
	case t.DirectoryType:
		return c.queues.Directories.Publish(ctx, r, priority)
	default:
		// Behaviour for other types not (yet) defined.
		panic(fmt.Sprintf("Unsupported type returned from listing: %s", r.Type))
	}
}
