package microbin

import "context"

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
