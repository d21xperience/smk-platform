package models

import (
	"errors"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/correspondence/events"
	"sekolah-platform/services/tu-core/internal/domain/correspondence/valueobjects"

	"github.com/google/uuid"
)

type Correspondence struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	Type             valueobjects.CorrespondenceType
	JenisSurat       string
	Kategori         string
	KodeKlasifikasi  string
	Number           string
	Date             time.Time
	Subject          string
	FromParty        string
	ToParty          string
	Description      string
	AttachmentURL    string
	FileURL          string
	TTDURL           string
	Konten           map[string]interface{}
	Status           valueobjects.CorrespondenceStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
	domainEvents     []types.DomainEvent
}

type CorrespondenceData struct {
	ID              string
	Type            valueobjects.CorrespondenceType
	JenisSurat      string
	Kategori        string
	KodeKlasifikasi string
	Number          string
	Date            time.Time
	Subject         string
	FromParty       string
	ToParty         string
	Description     string
	AttachmentURL   string
	FileURL         string
	TTDURL          string
	Konten          map[string]interface{}
}

func NewCorrespondence(data CorrespondenceData, ctx *context.OperationalContext) (*Correspondence, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	if data.Number == "" {
		return nil, errors.New("nomor surat wajib diisi")
	}
	if data.Date.IsZero() {
		return nil, errors.New("tanggal surat wajib diisi")
	}
	if data.Subject == "" {
		return nil, errors.New("subjek surat wajib diisi")
	}
	if data.FromParty == "" || data.ToParty == "" {
		return nil, errors.New("pengirim dan penerima wajib diisi")
	}
	if data.JenisSurat == "" {
		return nil, errors.New("jenis surat wajib diisi")
	}

	c := &Correspondence{
		ID:               data.ID,
		SchoolID:         ctx.SchoolID,
		AcademicPeriodID: ctx.AcademicPeriodID,
		Type:             data.Type,
		JenisSurat:       data.JenisSurat,
		Kategori:         data.Kategori,
		KodeKlasifikasi:  data.KodeKlasifikasi,
		Number:           data.Number,
		Date:             data.Date,
		Subject:          data.Subject,
		FromParty:        data.FromParty,
		ToParty:          data.ToParty,
		Description:      data.Description,
		AttachmentURL:    data.AttachmentURL,
		FileURL:          data.FileURL,
		TTDURL:           data.TTDURL,
		Konten:           data.Konten,
		Status:           valueobjects.StatusDraft,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
		domainEvents:     []types.DomainEvent{},
	}

	c.recordEvent(events.NewCorrespondenceCreated(events.CorrespondenceCreatedPayload{
		CorrespondenceID: c.ID,
		Type:             c.Type.String(),
		Number:           c.Number,
		Subject:          c.Subject,
	}, ctx))

	return c, nil
}

func (c *Correspondence) Process(ctx *context.OperationalContext) error {
	if c.Status != valueobjects.StatusDraft {
		return errors.New("hanya surat dengan status DRAFT yang dapat diproses")
	}
	c.Status = valueobjects.StatusProcessed
	c.UpdatedAt = time.Now().UTC()
	c.recordEvent(events.NewCorrespondenceProcessed(events.CorrespondenceProcessedPayload{
		CorrespondenceID: c.ID,
		Number:           c.Number,
	}, ctx))
	return nil
}

func (c *Correspondence) Archive(ctx *context.OperationalContext) error {
	if !c.Status.CanArchive() {
		return errors.New("hanya surat dengan status DRAFT atau PROCESSED yang dapat diarsipkan")
	}
	c.Status = valueobjects.StatusArchived
	c.UpdatedAt = time.Now().UTC()
	c.recordEvent(events.NewCorrespondenceArchived(events.CorrespondenceArchivedPayload{
		CorrespondenceID: c.ID,
		Number:           c.Number,
	}, ctx))
	return nil
}

func (c *Correspondence) Update(kategori, kodeKlasifikasi, subject, description, attachmentURL, fileURL, ttdURL string, konten map[string]interface{}) error {
	if c.Status != valueobjects.StatusDraft {
		return errors.New("hanya surat dengan status DRAFT yang dapat diubah")
	}
	c.Kategori = kategori
	c.KodeKlasifikasi = kodeKlasifikasi
	c.Subject = subject
	c.Description = description
	c.AttachmentURL = attachmentURL
	c.FileURL = fileURL
	c.TTDURL = ttdURL
	if konten != nil {
		c.Konten = konten
	}
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Correspondence) recordEvent(event types.DomainEvent) {
	c.domainEvents = append(c.domainEvents, event)
}

func (c *Correspondence) UncommittedEvents() []types.DomainEvent {
	return c.domainEvents
}

func (c *Correspondence) ClearEvents() {
	c.domainEvents = []types.DomainEvent{}
}
