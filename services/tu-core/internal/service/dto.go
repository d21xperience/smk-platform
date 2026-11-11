package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// FlexString adalah tipe kustom yang bisa menerima string maupun number dari JSON.
type FlexString string

func (s *FlexString) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*s = ""
		return nil
	}
	if data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		*s = FlexString(str)
		return nil
	}
	if data[0] == '-' || (data[0] >= '0' && data[0] <= '9') {
		var num json.Number
		if err := json.Unmarshal(data, &num); err != nil {
			return err
		}
		*s = FlexString(num.String())
		return nil
	}
	return fmt.Errorf("cannot unmarshal %s into FlexString", string(data))
}

// === STUDENT DTOs ===
type CreateStudentRequest struct {
	NIS         FlexString `json:"nis"`
	NISN        FlexString `json:"nisn"`
	Name        string     `json:"name"`
	Gender      string     `json:"gender"`
	BirthDate   *time.Time `json:"birthDate"`
	BirthPlace  string     `json:"birthPlace"`
	Address     string     `json:"address"`
	Phone       FlexString `json:"phone"`
	Email       string     `json:"email"`
	ParentName  string     `json:"parentName"`
	ParentPhone FlexString `json:"parentPhone"`
	ParentEmail string     `json:"parentEmail"`
	ClassID     string     `json:"classId"`
}

type UpdateStudentRequest struct {
	Name        string     `json:"name"`
	NIS         FlexString `json:"nis"`
	NISN        FlexString `json:"nisn"`
	ClassID     string     `json:"classId"`
	Gender      string     `json:"gender"`
	BirthDate   *time.Time `json:"birthDate"`
	BirthPlace  string     `json:"birthPlace"`
	Address     string     `json:"address"`
	Phone       FlexString `json:"phone"`
	Email       string     `json:"email"`
	ParentName  string     `json:"parentName"`
	ParentPhone FlexString `json:"parentPhone"`
	ParentEmail string     `json:"parentEmail"`
}

type StudentResponse struct {
	ID               string    `json:"id"`
	SchoolID         string    `json:"schoolId"`
	AcademicPeriodID string    `json:"academicPeriodId"`
	NIS              string    `json:"nis"`
	NISN             string    `json:"nisn"`
	Name             string    `json:"name"`
	Gender           string    `json:"gender"`
	BirthDate        time.Time `json:"birthDate"`
	BirthPlace       string    `json:"birthPlace"`
	Address          string    `json:"address"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email"`
	ParentName       string    `json:"parentName"`
	ParentPhone      string    `json:"parentPhone"`
	ParentEmail      string    `json:"parentEmail"`
	ClassName        string    `json:"className"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type StudentListResponse struct {
	Data    []StudentResponse `json:"data"`
	Total   int               `json:"total"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
	HasMore bool              `json:"hasMore"`
}

type StudentFilterRequest struct {
	ClassID string `json:"classId"`
	Status  string `json:"status"`
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

// === TEACHER DTOs ===
type CreateTeacherRequest struct {
	NIP       string    `json:"nip"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birthDate"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	SubjectID string    `json:"subjectId"`
}

type TeacherResponse struct {
	ID               string    `json:"id"`
	SchoolID         string    `json:"schoolId"`
	AcademicPeriodID string    `json:"academicPeriodId"`
	NIP              string    `json:"nip"`
	Name             string    `json:"name"`
	Gender           string    `json:"gender"`
	BirthDate        time.Time `json:"birthDate"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	Address          string    `json:"address"`
	SubjectID        string    `json:"subjectId"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
}

// === ATTENDANCE DTOs ===
type SubmitAttendanceRequest struct {
	ClassID string                    `json:"classId"`
	Date    time.Time                 `json:"date"`
	Records []AttendanceRecordRequest `json:"records"`
}

type AttendanceRecordRequest struct {
	StudentID string `json:"studentId"`
	Status    string `json:"status"`
	Note      string `json:"note"`
}

type AttendanceResponse struct {
	ID        string    `json:"id"`
	StudentID string    `json:"studentId"`
	ClassID   string    `json:"classId"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
}

// === ASSESSMENT DTOs ===
type CreateAssessmentRequest struct {
	StudentID      string    `json:"studentId"`
	SubjectID      string    `json:"subjectId"`
	AssessmentType string    `json:"assessmentType"`
	Score          float64   `json:"score"`
	MaxScore       float64   `json:"maxScore"`
	AssessmentDate time.Time `json:"assessmentDate"`
	Semester       int       `json:"semester"`
	Notes          string    `json:"notes"`
}

type AssessmentResponse struct {
	ID               string    `json:"id"`
	SchoolID         string    `json:"schoolId"`
	AcademicPeriodID string    `json:"academicPeriodId"`
	StudentID        string    `json:"studentId"`
	SubjectID        string    `json:"subjectId"`
	AssessmentType   string    `json:"assessmentType"`
	Score            float64   `json:"score"`
	MaxScore         float64   `json:"maxScore"`
	Percentage       float64   `json:"percentage"`
	AssessmentDate   time.Time `json:"assessmentDate"`
	Semester         int       `json:"semester"`
	Average          float64   `json:"average"`
	Notes            string    `json:"notes"`
	CreatedBy        string    `json:"createdBy"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type SubjectAverageResponse struct {
	SubjectID        string  `json:"subjectId"`
	SubjectName      string  `json:"subjectName"`
	TotalScore       float64 `json:"totalScore"`
	TotalAssessments int     `json:"totalAssessments"`
	Average          float64 `json:"average"`
	Grade            string  `json:"grade"`
}

// === FINANCE DTOs ===
type CreateInvoiceRequest struct {
	StudentID     string    `json:"studentId"`
	ComponentType string    `json:"componentType"`
	Month         int       `json:"month"`
	Year          int       `json:"year"`
	Amount        float64   `json:"amount"`
	DueDate       time.Time `json:"dueDate"`
}

type InvoiceResponse struct {
	ID               string     `json:"id"`
	SchoolID         string     `json:"schoolId"`
	AcademicPeriodID string     `json:"academicPeriodId"`
	StudentID        string     `json:"studentId"`
	ComponentType    string     `json:"componentType"`
	Month            int        `json:"month"`
	Year             int        `json:"year"`
	Amount           float64    `json:"amount"`
	DueDate          time.Time  `json:"dueDate"`
	Status           string     `json:"status"`
	PaidAmount       float64    `json:"paidAmount"`
	PaidAt           *time.Time `json:"paidAt"`
	CreatedBy        string     `json:"cretedBy"`
	Notes            string     `json:"notes"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

type CreatePaymentRequest struct {
	InvoiceID       string    `json:"invoiceId"`
	Amount          float64   `json:"amount"`
	PaymentMethod   string    `json:"paymentMethod"`
	PaymentDate     time.Time `json:"paymentDate"`
	ReferenceNumber string    `json:"referenceNumber"`
}

type PaymentResponse struct {
	ID               string    `json:"id"`
	SchoolID         string    `json:"schoolId"`
	AcademicPeriodID string    `json:"academicPeriodId"`
	InvoiceID        string    `json:"invoiceId"`
	Amount           float64   `json:"amount"`
	PaymentMethod    string    `json:"paymentMethod"`
	PaymentDate      time.Time `json:"paymentDate"`
	ReferenceNumber  string    `json:"referenceNumber"`
	CreatedBy        string    `json:"cretedBy"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// === CORRESPONDENCE DTOs ===
type CreateCorrespondenceRequest struct {
	Type            string                 `json:"type"`
	JenisSurat      string                 `json:"jenis_surat"`
	Kategori        string                 `json:"kategori"`
	KodeKlasifikasi string                 `json:"kode_klasifikasi"`
	Number          string                 `json:"number"`
	Date            time.Time              `json:"date"`
	Subject         string                 `json:"subject"`
	From            string                 `json:"from"`
	To              string                 `json:"to"`
	Description     string                 `json:"description"`
	AttachmentURL   string                 `json:"attachment_url"`
	FileURL         string                 `json:"file_url"`
	TTDURL          string                 `json:"ttd_url"`
	Konten          map[string]interface{} `json:"konten"`
}

type UpdateCorrespondenceRequest struct {
	Kategori        string                 `json:"kategori"`
	KodeKlasifikasi string                 `json:"kode_klasifikasi"`
	Subject         string                 `json:"subject"`
	Description     string                 `json:"description"`
	AttachmentURL   string                 `json:"attachment_url"`
	FileURL         string                 `json:"file_url"`
	TTDURL          string                 `json:"ttd_url"`
	Konten          map[string]interface{} `json:"konten"`
}

type CorrespondenceResponse struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"`
	JenisSurat      string                 `json:"jenis_surat"`
	Kategori        string                 `json:"kategori"`
	KodeKlasifikasi string                 `json:"kode_klasifikasi"`
	Number          string                 `json:"number"`
	NomorAgenda     string                 `json:"nomor_agenda"`
	Date            time.Time              `json:"date"`
	Subject         string                 `json:"subject"`
	From            string                 `json:"from"`
	To              string                 `json:"to"`
	Description     string                 `json:"description"`
	AttachmentURL   string                 `json:"attachment_url"`
	FileURL         string                 `json:"file_url"`
	TTDURL          string                 `json:"ttd_url"`
	Konten          map[string]interface{} `json:"konten"`
	Status          string                 `json:"status"`
	CreatedAt       time.Time              `json:"createdAt"`
}
