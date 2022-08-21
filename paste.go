package microbin

import (
	"time"

	"gorm.io/gorm"
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
	ID int `json:"id" gorm:"primaryKey"`
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
	now := time.Now()
	duration := expirationDuration[p.Expiration]

	return p.CreatedAt.Add(duration).Before(now)
}

type PasteService struct {
	DB *gorm.DB
}

// Create valiates the incoming Paste before persisting it to the datastore.
func (p *PasteService) Create(paste Paste) (int, error) {
	if paste.Expiration == "" {
		return 0, ErrorUnrecognizedExpiration{}
	}

	if !paste.Expiration.IsValid() {
		return 0, ErrorUnrecognizedExpiration{
			paste.Expiration,
		}
	}

	result := p.DB.Create(&paste)

	if result.Error != nil {
		return 0, result.Error
	}

	return paste.ID, nil
}

// Read finds a Paste with the given id in the datastore. If the Paste has
// expired it is expunged and an ErrorPasteExpired error is returned to the
// caller.
func (p *PasteService) Read(id int) (Paste, error) {
	paste := Paste{}

	if err := p.DB.First(&paste, id).Error; err != nil {
		return Paste{}, err
	}

	// if the paste has expired delete it and return an error to the caller
	if paste.Expired() {
		if err := p.DB.Delete(&Paste{}, id).Error; err != nil {
			return Paste{}, err
		}

		return Paste{}, ErrorPasteExpired{ID: id}
	}

	return paste, nil
}

// Delete deletes a Paste with the given id from the datastore. If an error
// occurs it is returned to the caller.
func (p *PasteService) Delete(id int) error {
	if tx := p.DB.Delete(&Paste{}, id); tx.Error != nil {
		return tx.Error
	}

	return nil
}

// List returns all the Paste's in the datastore. Any Paste's that have expired
// are deleted from the datastore and expunged from the result set.
func (p *PasteService) List() ([]Paste, error) {
	pastesRet := []Paste{}
	pastes := []Paste{}

	if tx := p.DB.Find(&pastes); tx.Error != nil {
		return []Paste{}, tx.Error
	}

	for _, paste := range pastes {
		if paste.Expired() {
			if err := p.DB.Delete(&Paste{}, paste.ID).Error; err != nil {
				return []Paste{}, err
			}
		} else {
			pastesRet = append(pastesRet, paste)
		}
	}

	return pastesRet, nil
}
