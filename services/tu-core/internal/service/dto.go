package service

import "time"

// === STUDENT DTOs ===

// CreateStudentRequest adalah request untuk membuat siswa baru.
type CreateStudentRequest struct {
	NIS         string    `json:"nis"`
	NISN        string    `json:"nisn"`
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	BirthDate   time.Time `json:"birthDate"`
	BirthPlace  string    `json:"birthPlace"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	ParentName  string    `json:"parentName"`
	ParentPhone string    `json:"parentPhone"`
	ParentEmail string    `json:"parentEmail"`
	ClassID     string    `json:"classId"`
}

// UpdateStudentRequest adalah request untuk update siswa.
type UpdateStudentRequest struct {
	Name        string    `json:"name"`
	NIS         string    `json:"nis"`
	NISN        string    `json:"nisn"`
	ClassID     string    `json:"classId"`
	Gender      string    `json:"gender"`
	BirthDate   time.Time `json:"birthDate"`
	BirthPlace  string    `json:"birthPlace"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	ParentName  string    `json:"parentName"`
	ParentPhone string    `json:"parentPhone"`
	ParentEmail string    `json:"parentEmail"`
}

// StudentResponse adalah response untuk data siswa.
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

// StudentListResponse adalah response untuk list siswa dengan pagination.
type StudentListResponse struct {
	Data    []StudentResponse `json:"data"`
	Total   int               `json:"total"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
	HasMore bool              `json:"hasMore"`
}

// StudentFilterRequest adalah request filter untuk list siswa.
type StudentFilterRequest struct {
	ClassID string `json:"classId"`
	Status  string `json:"status"`
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

// === TEACHER DTOs ===

// CreateTeacherRequest adalah request untuk membuat guru baru.
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

// TeacherResponse adalah response untuk data guru.
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

// SubmitAttendanceRequest adalah request untuk submit absensi.
type SubmitAttendanceRequest struct {
	ClassID string                    `json:"classId"`
	Date    time.Time                 `json:"date"`
	Records []AttendanceRecordRequest `json:"records"`
}

// AttendanceRecordRequest adalah data absensi per siswa.
type AttendanceRecordRequest struct {
	StudentID string `json:"studentId"`
	Status    string `json:"status"` // PRESENT, ABSENT, LATE, SICK, PERMISSION
	Note      string `json:"note"`
}

// AttendanceResponse adalah response untuk data absensi.
type AttendanceResponse struct {
	ID        string    `json:"id"`
	StudentID string    `json:"studentId"`
	ClassID   string    `json:"classId"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
}

// === ASSESSMENT DTOs===
// CreateAssessmentRequest adalah request untuk membuat nilai.
type CreateAssessmentRequest struct {
	StudentID      string    `json:"studentId"`
	SubjectID      string    `json:"subjectId"`
	AssessmentType string    `json:"assessmentType"` // DAILY, MIDTERM, FINAL
	Score          float64   `json:"score"`
	MaxScore       float64   `json:"maxScore"`
	AssessmentDate time.Time `json:"assessmentDate"`
	Semester       int       `json:"semester"` // 1 atau 2
	Notes          string    `json:"notes"`
}

// AssessmentResponse adalah response untuk data nilai.
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

// SubjectAverageResponse adalah response untuk rata-rata per mata pelajaran.
type SubjectAverageResponse struct {
	SubjectID        string  `json:"subjectId"`
	SubjectName      string  `json:"subjectName"`
	TotalScore       float64 `json:"totalScore"`
	TotalAssessments int     `json:"totalAssessments"`
	Average          float64 `json:"average"`
	Grade            string  `json:"grade"`
}

// === FINANCE DTOs ===

// CreateInvoiceRequest adalah request untuk membuat tagihan.
type CreateInvoiceRequest struct {
	StudentID     string    `json:"studentId"`
	ComponentType string    `json:"componentType"` // SPP, UNIFORM, BOOK
	Month         int       `json:"month"`
	Year          int       `json:"year"`
	Amount        float64   `json:"amount"`
	DueDate       time.Time `json:"dueDate"`
}

// InvoiceResponse adalah response untuk data tagihan.
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

// CreatePaymentRequest adalah request untuk membuat pembayaran.
type CreatePaymentRequest struct {
	InvoiceID       string    `json:"invoiceId"`
	Amount          float64   `json:"amount"`
	PaymentMethod   string    `json:"paymentMethod"` // CASH, TRANSFER, VA, QRIS
	PaymentDate     time.Time `json:"paymentDate"`
	ReferenceNumber string    `json:"referenceNumber"`
}

// PaymentResponse adalah response untuk data pembayaran.
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

// CreateCorrespondenceRequest adalah request untuk membuat surat.
type CreateCorrespondenceRequest struct {
	Type          string    `json:"type"` // INCOMING, OUTGOING
	Number        string    `json:"number"`
	Date          time.Time `json:"date"`
	Subject       string    `json:"subject"`
	From          string    `json:"from"`
	To            string    `json:"to"`
	Description   string    `json:"description"`
	AttachmentURL string    `json:"attachmentUrl"`
}

// CorrespondenceResponse adalah response untuk data surat.
type CorrespondenceResponse struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	Number        string    `json:"number"`
	Date          time.Time `json:"date"`
	Subject       string    `json:"subject"`
	From          string    `json:"from"`
	To            string    `json:"to"`
	Description   string    `json:"description"`
	AttachmentURL string    `json:"attachmentUrl"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
}
