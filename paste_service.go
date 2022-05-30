package main

import (
	"log"

	"gorm.io/gorm"
)

type pasteService struct {
	logger *log.Logger
	db     *gorm.DB
}

// Create valiates the incoming Paste before persisting it to the datastore.
func (p *pasteService) Create(paste Paste) (int, error) {
	if paste.Expiration == "" {
		return 0, errorUnrecognizedExpiration{}
	}

	if !expirationIsValid(paste.Expiration) {
		return 0, errorUnrecognizedExpiration{paste.Expiration}
	}

	result := p.db.Create(&paste)

	if result.Error != nil {
		return 0, result.Error
	}

	return paste.ID, nil
}

// Read finds a Paste with the given id in the datastore. If the Paste has
// expired it is expunged and an ErrorPasteExpired error is returned to the
// caller.
func (p *pasteService) Read(id int) (Paste, error) {
	paste := Paste{}

	if err := p.db.First(&paste, id).Error; err != nil {
		return Paste{}, err
	}

	// if the paste has expired delete it and return an error to the caller
	if paste.Expired() {
		if err := p.db.Delete(&Paste{}, id).Error; err != nil {
			return Paste{}, err
		}

		return Paste{}, errorPasteExpired{ID: id}
	}

	return paste, nil
}

// Delete deletes a Paste with the given id from the datastore. If an error
// occurs it is returned to the caller.
func (p *pasteService) Delete(id int) error {
	if tx := p.db.Delete(&Paste{}, id); tx.Error != nil {
		return tx.Error
	}

	return nil
}

// List returns all the Paste's in the datastore. Any Paste's that have expired
// are deleted from the datastore and expunged from the result set.
func (p *pasteService) List() ([]Paste, error) {
	pastesRet := []Paste{}
	pastes := []Paste{}

	if tx := p.db.Find(&pastes); tx.Error != nil {
		return []Paste{}, tx.Error
	}

	for _, paste := range pastes {
		if paste.Expired() {
			if err := p.db.Delete(&Paste{}, paste.ID).Error; err != nil {
				return []Paste{}, err
			}
		} else {
			pastesRet = append(pastesRet, paste)
		}
	}

	return pastesRet, nil
}
