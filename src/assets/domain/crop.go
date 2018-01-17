package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Crop struct {
	UID          uuid.UUID
	BatchID      string
	InitialArea  Area
	CurrentAreas []Area
	Type         CropType
	Inventory    InventoryMaterial
	Container    CropContainer
	Notes        map[uuid.UUID]CropNote
	CreatedDate  time.Time
}

// CropType defines type of crop
type CropType interface {
	Code() string
}

// Seeding implements CropType
type Seeding struct{}

func (s Seeding) Code() string { return "seeding" }

// Growing implements CropType
type Growing struct{}

func (g Growing) Code() string { return "growing" }

// CropContainer defines the container of a crop
type CropContainer struct {
	Quantity int
	Type     CropContainerType
}

// CropContainerType defines the type of a container
type CropContainerType interface {
	Code() string
}

// Tray implements CropContainerType
type Tray struct {
	Cell int
}

func (t Tray) Code() string { return "tray" }

// Pot implements CropContainerType
type Pot struct{}

func (p Pot) Code() string { return "pot" }

type CropNote struct {
	UID         uuid.UUID `json:"uid"`
	Content     string    `json:"content"`
	CreatedDate time.Time `json:"created_date"`
}

func CreateCropBatch(area Area) (Crop, error) {
	if area.UID == (uuid.UUID{}) {
		return Crop{}, CropError{Code: CropErrorInvalidArea}
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return Crop{}, err
	}

	return Crop{
		UID:          uid,
		InitialArea:  area,
		CurrentAreas: []Area{area},
		CreatedDate:  time.Now(),
	}, nil
}

func (c *Crop) ChangeCropType(cropType CropType) error {
	err := validateCropType(cropType)
	if err != nil {
		return err
	}

	c.Type = cropType

	return nil
}

func (c *Crop) ChangeContainer(container CropContainer) error {
	err := validateCropContainer(container)
	if err != nil {
		return err
	}

	c.Container = container

	return nil
}

func (c *Crop) AddNewNote(content string) error {
	if content == "" {
		return CropError{Code: CropNoteErrorInvalidContent}
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	cropNote := CropNote{
		UID:         uid,
		Content:     content,
		CreatedDate: time.Now(),
	}

	if len(c.Notes) == 0 {
		c.Notes = make(map[uuid.UUID]CropNote)
	}

	c.Notes[uid] = cropNote

	return nil
}

func (c *Crop) RemoveNote(uid string) error {
	if uid == "" {
		return CropError{Code: CropNoteErrorNotFound}
	}

	uuid, err := uuid.FromString(uid)
	if err != nil {
		return AreaError{Code: AreaNoteErrorNotFound}
	}

	found := false
	for _, v := range c.Notes {
		if v.UID == uuid {
			delete(c.Notes, uuid)
			found = true
		}
	}

	if !found {
		return CropError{Code: CropNoteErrorNotFound}
	}

	return nil
}

// CalculateDaysSinceSeeding will find how long since its been seeded
// It basically tell use the days since this crop is created.
func (c Crop) CalculateDaysSinceSeeding() int {
	now := time.Now()

	diff := now.Sub(c.CreatedDate)

	days := int(diff.Hours()) / 24

	return days
}

func validateCropType(cropType CropType) error {
	switch cropType.(type) {
	case Seeding:
	case Growing:
	default:
		return CropError{Code: CropErrorInvalidCropType}
	}

	return nil
}

func validateCropContainer(container CropContainer) error {
	switch container.Type.(type) {
	case Tray:
	case Pot:
	default:
		return CropError{Code: CropContainerErrorInvalidType}
	}

	return nil
}
