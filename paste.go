package microbin

import (
	"context"
	"time"
)

type PasteType int

const (
	PasteTypeText PasteType = iota
	PasteTypeFile
	PasteTypeURL
)

// Paste represents a piece of content that will be persisted for a duration or
// indefinitely. It can be text or binary content.
type Paste struct {
	ID int `json:"id"`
	// Content represents the text content uploaded as a paste.
	Content string `json:"content"`
	// BinaryContent represents a file that has been uploaded as a paste. It is omitted from JSON responses.
	BinaryContent []byte     `json:"-"`
	File          string     `json:"file"`
	Expiration    Expiration `json:"expiration"`
	Type          PasteType  `json:"type"`
	CreatedAt     time.Time  `json:"createdAt"`
}

func (p *Paste) Expired() bool {
	var (
		now      = time.Now()
		duration = p.Expiration.ToDuration()
	)

	return p.CreatedAt.Add(duration).Before(now)
}

func CreatePaste(store PasteStore) func(context.Context, Paste) error {
	return func(ctx context.Context, p Paste) error {
		if !p.Expiration.IsValid() {
			return ErrInvalidExpiration
		}

		return store.Create(ctx, p)
	}
}

func ReadPaste(store PasteStore) func(ctx context.Context, id int) (Paste, error) {
	return func(ctx context.Context, id int) (Paste, error) {
		paste, err := store.Read(ctx, id)

		if err != nil {
			return Paste{}, err
		}

		if paste.Expired() {
			err = store.Delete(ctx, id)

			if err != nil {
				return Paste{}, err
			}

			return Paste{}, ErrExpiredPaste(id)
		}

		return paste, nil
	}
}

func DeletePaste(store PasteStore) func(ctx context.Context, id int) error {
	return func(ctx context.Context, id int) error {
		return store.Delete(ctx, id)
	}
}

func ListPastes(store PasteStore) func(ctx context.Context) ([]Paste, error) {
	return func(ctx context.Context) ([]Paste, error) {
		pastes, err := store.List(ctx)

		if err != nil {
			return []Paste{}, err
		}

		delPastes := []Paste{}
		retPastes := []Paste{}

		for _, paste := range pastes {
			if !paste.Expired() {
				retPastes = append(retPastes, paste)
			} else {
				delPastes = append(delPastes, paste)
			}
		}

		if len(delPastes) > 0 {
			for _, paste := range delPastes {
				if err := store.Delete(ctx, paste.ID); err != nil {
					return []Paste{}, err
				}
			}
		}

		return retPastes, nil
	}
}
