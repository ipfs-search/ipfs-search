package crawler

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"log"
	"math/rand"

	indexTypes "github.com/ipfs-search/ipfs-search/index/types"
	t "github.com/ipfs-search/ipfs-search/types"
)

var (
	// ErrDirectoryTooLarge is returned by Ls() when a directory is larger `Config.MaxDirSize`.
	ErrDirectoryTooLarge = t.WrappedError{Err: t.ErrInvalidResource, Msg: "directory too large"}

	// errEndOfLs is an internal error to communicate the end of hte list from processNextDirEntry to processDirEntries.
	errEndOfLs = errors.New("end of list")
)

func (c *Crawler) crawlDir(ctx context.Context, r *t.AnnotatedResource, properties *indexTypes.Directory) error {
	entries := make(chan *t.AnnotatedResource, c.config.DirEntryBufferSize)

	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		return c.processDirEntries(ctx, entries, properties)
	})

	wg.Go(func() error {
		defer close(entries)
		return c.protocol.Ls(ctx, r, entries)
	})

	return wg.Wait()
}

func resourceToLinkType(r *t.AnnotatedResource) indexTypes.LinkType {
	switch r.Type {
	case t.FileType:
		return indexTypes.FileLinkType
	case t.DirectoryType:
		return indexTypes.DirectoryLinkType
	case t.UndefinedType:
		return indexTypes.UnknownLinkType
	case t.UnsupportedType:
		return indexTypes.UnsupportedLinkType
	default:
		panic("unexpected type")
	}
}

func addLink(e *t.AnnotatedResource, properties *indexTypes.Directory) {
	properties.Links = append(properties.Links, indexTypes.Link{
		Hash: e.ID,
		Name: e.Reference.Name,
		Size: e.Size,
		Type: resourceToLinkType(e),
	})
}

func (c *Crawler) processDirEntries(ctx context.Context, entries <-chan *t.AnnotatedResource, properties *indexTypes.Directory) error {
	var (
		dirCnt  uint = 0
		isLarge      = false
	)

	// Question: do we need a maximum entry cutoff point? E.g. 10^6 entries or something?
	processNextDirEntry := func() error {
		// Create (and cancel!) a new timeout context for every entry.
		ctx, cancel := context.WithTimeout(ctx, c.config.DirEntryTimeout)
		defer cancel()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case entry, ok := <-entries:
			if !ok {
				return errEndOfLs
			}

			if dirCnt%1024 == 0 {
				log.Printf("Processed %d directory entries in %v.", dirCnt, entry.Parent)
				log.Printf("Latest entry: %v", entry)
			}

			// Only add to properties up to limit (preventing oversized directory entries) - but queue entries nonetheless.
			if dirCnt == c.config.MaxDirSize {
				log.Printf("Directory %v is large: %d", entry.Parent, dirCnt)
				isLarge = true
			}

			if !isLarge {
				addLink(entry, properties)
			}

			return c.queueDirEntry(ctx, entry)
		}
	}

	var err error

	// Process entries until error.
	for err == nil {
		err = processNextDirEntry()
		dirCnt++
	}

	if errors.Is(err, errEndOfLs) {
		// Normal exit of loop, reset error condition
		err = nil

		if isLarge {
			err = ErrDirectoryTooLarge
		}
	} else {
		// Unknown error situation: fail hard
		// Prefer less over incomplete or inconsistent data.
		log.Printf("Unexpected error processing directory entries: %v", err)
	}

	return err
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
		// Index right away as invalid.
		// Rationale: as no additional protocol request is required and queue'ing returns
		// similarly fast as indexing.
		return c.indexInvalid(ctx, r, t.ErrUnsupportedType)
	default:
		panic("unexpected type")
	}
}
