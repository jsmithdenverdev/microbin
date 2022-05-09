package microbin

import (
	"log"

	"gorm.io/gorm"
)

type PasteService struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	db       *gorm.DB
}

func NewPasteService(db *gorm.DB, infoLog, errorLog *log.Logger) *PasteService {
	return &PasteService{
		db:       db,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

func (p *PasteService) Create(paste Paste) (int, error) {
	if paste.Expiration == "" {
		return 0, ErrorUnrecognizedExpiration{}
	}

	if !ValidExpiration(paste.Expiration) {
		return 0, ErrorUnrecognizedExpiration{paste.Expiration}
	}

	result := p.db.Create(&paste)

	if result.Error != nil {
		return 0, result.Error
	}

	return paste.ID, nil
}

func (p *PasteService) Read(id int) (Paste, error) {
	paste := Paste{}

	if err := p.db.First(&paste, id).Error; err != nil {
		return Paste{}, err
	}

	// if the paste has expired delete it and return an error to the caller
	if paste.Expired() {
		if err := p.db.Delete(&Paste{}, id).Error; err != nil {
			return Paste{}, err
		}

		return Paste{}, ErrorPasteExpired{ID: id}
	}

	return paste, nil
}

func (p *PasteService) Delete(id int) error {
	if tx := p.db.Delete(&Paste{}, id); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (p *PasteService) List() ([]Paste, error) {
	pastes := []Paste{}

	if tx := p.db.Find(&pastes); tx.Error != nil {
		return []Paste{}, tx.Error
	}

	for _, paste := range pastes {
		if paste.Expired() {
			if tx := p.db.Delete(&Paste{}, paste.ID); tx.Error != nil {
				return []Paste{}, tx.Error
			}
		}
	}

	return pastes, nil
}
