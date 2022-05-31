package microbin

import "context"

type PasteStore interface {
	Create(context.Context, Paste) error
	Read(ctx context.Context, id int) (Paste, error)
	List(context.Context) ([]Paste, error)
	Delete(ctx context.Context, id int) error
}
