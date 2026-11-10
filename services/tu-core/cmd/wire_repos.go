package main

import (
	"sekolah-platform/services/tu-core/internal/notification"
	notifHandlers "sekolah-platform/services/tu-core/internal/notification/handlers"
	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppRepositories struct {
	Student        repository.StudentRepository
	StudentDoc     repository.StudentDocumentRepository
	Mutation       repository.MutationRepository
	Registration   repository.RegistrationRepository
	Attendance     repository.AttendanceRepository
	Assessment     repository.AssessmentRepository
	Finance        repository.FinanceRepository
	Correspondence repository.CorrespondenceRepository
	Disposition    repository.DispositionRepository
	Signature      repository.SignatureRepository
	Class          repository.ClassRepository
	Stats          repository.StatisticsRepository
	Teacher        repository.TeacherRepository
	TeacherHistory repository.TeacherHistoryRepository
	PtkLeave       repository.PtkLeaveRepository
	PtkPerformance repository.PtkPerformanceRepository
	PtkArchive     repository.PtkArchiveRepository
	PtkDashboard   repository.PtkDashboardRepository
	BukuInduk      repository.BukuIndukRepository // <-- NEW
	Enrollment     repository.EnrollmentRepository
	User           repository.UserRepository
	Public         repository.PublicRepository
	Academic       repository.AcademicRepository
	School         repository.SchoolRepository
	Subject        repository.SubjectRepository
	Pembelajaran   repository.PembelajaranRepository
	KurikulumStat  repository.KurikulumStatRepository
	Schedule       repository.ScheduleRepository
	Jurusan        repository.JurusanRepository
	Logistic       repository.LogisticRepository
	KirItem        repository.InventoryItemRepository
	Notification   notifHandlers.InvoiceRepository
	LetterType     repository.LetterTypeRepository
	Procurement    repository.ProcurementRepository
	GoodsReceipt   repository.GoodsReceiptRepository
	Document       repository.DocumentRepository
	Disposal       repository.DisposalRepository
	Borrowing      repository.BorrowingRepository
}

func initRepositories(db *pgxpool.Pool) *AppRepositories {
	return &AppRepositories{
		Student:        repository.NewStudentRepository(db),
		StudentDoc:     repository.NewStudentDocumentRepository(db),
		Mutation:       repository.NewMutationRepository(db),
		Registration:   repository.NewRegistrationRepository(db),
		Attendance:     repository.NewAttendanceRepository(db),
		Assessment:     repository.NewAssessmentRepository(db),
		Finance:        repository.NewFinanceRepository(db),
		Correspondence: repository.NewCorrespondenceRepository(db),
		Disposition:    repository.NewDispositionRepository(db),
		Signature:      repository.NewSignatureRepository(db),
		Class:          repository.NewClassRepository(db),
		Stats:          repository.NewStatisticsRepository(db),
		Teacher:        repository.NewTeacherRepository(db),
		TeacherHistory: repository.NewTeacherHistoryRepository(db),
		PtkLeave:       repository.NewPtkLeaveRepository(db),
		PtkPerformance: repository.NewPtkPerformanceRepository(db),
		PtkArchive:     repository.NewPtkArchiveRepository(db),
		PtkDashboard:   repository.NewPtkDashboardRepository(db),
		BukuInduk:      repository.NewBukuIndukRepository(db), // <-- NEW
		Enrollment:     repository.NewEnrollmentRepository(db),
		User:           repository.NewUserRepository(db),
		Public:         repository.NewPublicRepository(db),
		Academic:       repository.NewAcademicRepository(db),
		School:         repository.NewSchoolRepository(db),
		Subject:        repository.NewSubjectRepository(db),
		Pembelajaran:   repository.NewPembelajaranRepository(db),
		KurikulumStat:  repository.NewKurikulumStatRepository(db),
		Schedule:       repository.NewScheduleRepository(db),
		Jurusan:        repository.NewJurusanRepository(db),
		Logistic:       repository.NewLogisticRepository(db),
		KirItem:        repository.NewInventoryItemRepository(db),
		Notification:   notification.NewNotificationRepositoryImpl(db),
		LetterType:     repository.NewLetterTypeRepository(db),
		Procurement:    repository.NewProcurementRepository(db),
		GoodsReceipt:   repository.NewGoodsReceiptRepository(db),
		Document:       repository.NewDocumentRepository(db),
		Disposal:       repository.NewDisposalRepository(db),
		Borrowing:      repository.NewBorrowingRepository(db),
	}
}
