package crawler

import (
	"context"
	"errors"
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

func resourceToLinkType(r *t.AnnotatedResource) (indexTypes.LinkType, error) {
	switch r.Type {
	case t.FileType:
		return indexTypes.FileLinkType, nil
	case t.DirectoryType:
		return indexTypes.DirectoryLinkType, nil
	case t.UndefinedType:
		return indexTypes.UnknownLinkType, nil
	case t.UnsupportedType:
		return indexTypes.UnsupportedLinkType, nil
	default:
		return "", UnexpectedTypeError{r.Type}
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
			linkType, err := resourceToLinkType(entry)
			if err != nil {
				return err
			}

			properties.Links = append(properties.Links, indexTypes.Link{
				Hash: entry.ID,
				Name: entry.Reference.Name,
				Size: entry.Size,
				Type: linkType,
			})

			// Queue directory entry. Fail hard on error; prefer less over incomplete
			// or inconsistent data.
			// TODO: Consider doing this in a separate Goroutine, as it's blocking.
			if err := c.queueDirEntry(ctx, entry); err != nil {
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
	case t.UndefinedType:
		return c.queues.Hashes.Publish(ctx, r, priority)
	case t.FileType:
		return c.queues.Files.Publish(ctx, r, priority)
	case t.DirectoryType:
		return c.queues.Directories.Publish(ctx, r, priority)
	case t.UnsupportedType:
		return c.indexInvalid(ctx, r, errors.New(indexTypes.UnsupportedTypeError))
	default:
		return UnexpectedTypeError{r.Type}
	}
}
