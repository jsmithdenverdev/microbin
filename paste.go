package microbin

import (
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
