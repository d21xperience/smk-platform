package main

import (
	"sekolah-platform/services/tu-core/internal/service"
	"sekolah-platform/services/tu-core/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	classEngine "sekolah-platform/services/tu-core/internal/domain/academic/engine"
	assessmentEngine "sekolah-platform/services/tu-core/internal/domain/assessment/engine"
	correspondenceEngine "sekolah-platform/services/tu-core/internal/domain/correspondence/engine"
	financeEngine "sekolah-platform/services/tu-core/internal/domain/finance/engine"
	mutationEngine "sekolah-platform/services/tu-core/internal/domain/mutation/engine"
	registrationEngine "sekolah-platform/services/tu-core/internal/domain/registration/engine"
	studentEngine "sekolah-platform/services/tu-core/internal/domain/student/engine"
)

type AppEngines struct {
	Student        *studentEngine.StudentEngine
	Mutation       *mutationEngine.MutationEngine
	Registration   *registrationEngine.RegistrationEngine
	Assessment     *assessmentEngine.AssessmentEngine
	Finance        *financeEngine.FinanceEngine
	Correspondence *correspondenceEngine.CorrespondenceEngine
	Class          *classEngine.ClassEngine
}

type AppServices struct {
	Preference     *service.PreferenceService
	Student        *service.StudentService
	Mutation       *service.MutationService
	Registration   *service.RegistrationService
	Attendance     *service.AttendanceService
	Assessment     *service.AssessmentService
	Finance        *service.FinanceService
	Signature      *service.SignatureService
	Correspondence *service.CorrespondenceService
	Disposition    *service.DispositionService
	Class          *service.ClassService
	Stats          *service.StatisticsService
	Teacher        *service.TeacherService
	TeacherHistory *service.TeacherHistoryService
	PtkLeave       *service.PtkLeaveService
	PtkPerformance *service.PtkPerformanceService
	PtkArchive     *service.PtkArchiveService
	PtkDashboard   *service.PtkDashboardService
	BukuInduk      *service.BukuIndukService // <-- NEW
	Kurikulum      *service.KurikulumService
	Schedule       service.ScheduleService
	Jurusan        *service.JurusanService
	Enrollment     *service.EnrollmentService
	Auth           *service.AuthService
	Public         *service.PublicService
	User           *service.UserService
	School         *service.SchoolService
	Logistic       *service.LogisticService
	KIR            *service.KIRService
	LetterType     *service.LetterTypeService
	Procurement    *service.ProcurementService
	GoodsReceipt   *service.GoodsReceiptService
	Document       *service.DocumentService
	Disposal       *service.DisposalService
	Borrowing      *service.BorrowingService
	SarprasReport  *service.SarprasReportService
}

func initEngines() *AppEngines {
	return &AppEngines{
		Student:        studentEngine.NewStudentEngine(),
		Mutation:       mutationEngine.NewMutationEngine(),
		Registration:   registrationEngine.NewRegistrationEngine(),
		Assessment:     assessmentEngine.NewAssessmentEngine(),
		Finance:        financeEngine.NewFinanceEngine(),
		Correspondence: correspondenceEngine.NewCorrespondenceEngine(),
		Class:          classEngine.NewClassEngine(),
	}
}

func initServices(repos *AppRepositories, engines *AppEngines, minioClient *storage.MinioClient, eventPub service.EventPublisher, db *pgxpool.Pool) *AppServices {
	sqlDB := stdlib.OpenDBFromPool(db)
	prefService := service.NewPreferenceService(sqlDB)

	signatureSvc := service.NewSignatureService(repos.Signature)

	return &AppServices{
		Preference:     prefService,
		Student:        service.NewStudentService(repos.Student, repos.StudentDoc, minioClient, engines.Student, eventPub),
		Mutation:       service.NewMutationService(repos.Mutation, repos.Student, engines.Mutation, eventPub),
		Registration:   service.NewRegistrationService(repos.Registration, repos.Student, engines.Registration, eventPub),
		Attendance:     service.NewAttendanceService(repos.Attendance, repos.Student, eventPub),
		Assessment:     service.NewAssessmentService(repos.Assessment, repos.Student, engines.Assessment, eventPub),
		Finance:        service.NewFinanceService(repos.Finance, repos.Student, engines.Finance, eventPub),
		Signature:      signatureSvc,
		Correspondence: service.NewCorrespondenceService(repos.Correspondence, signatureSvc, engines.Correspondence, eventPub),
		Disposition:    service.NewDispositionService(repos.Disposition),
		Class:          service.NewClassService(repos.Class, engines.Class, eventPub),
		Stats:          service.NewStatisticsService(repos.Stats),
		Teacher:        service.NewTeacherService(repos.Teacher),
		TeacherHistory: service.NewTeacherHistoryService(repos.TeacherHistory),
		PtkLeave:       service.NewPtkLeaveService(repos.PtkLeave),
		PtkPerformance: service.NewPtkPerformanceService(repos.PtkPerformance),
		PtkArchive:     service.NewPtkArchiveService(repos.PtkArchive),
		PtkDashboard:   service.NewPtkDashboardService(repos.PtkDashboard),
		BukuInduk:      service.NewBukuIndukService(repos.BukuInduk), // <-- NEW
		Kurikulum:      service.NewKurikulumService(repos.Subject, repos.Pembelajaran, repos.KurikulumStat),
		Schedule:       service.NewScheduleService(repos.Schedule, db),
		Jurusan:        service.NewJurusanService(repos.Jurusan),
		Enrollment:     service.NewEnrollmentService(db, repos.Enrollment, repos.Class, repos.Student, engines.Student, eventPub),
		Auth:           service.NewAuthService(repos.User, repos.Academic),
		Public:         service.NewPublicService(repos.Public),
		User:           service.NewUserService(repos.User),
		School:         service.NewSchoolService(repos.School),
		Logistic:       service.NewLogisticService(repos.Logistic),
		KIR:            service.NewKIRService(db, repos.KirItem),
		LetterType:     service.NewLetterTypeService(repos.LetterType),
		Procurement:    service.NewProcurementService(repos.Procurement),
		GoodsReceipt:   service.NewGoodsReceiptService(repos.GoodsReceipt),
		Document:       service.NewDocumentService(repos.Document, minioClient),
		Disposal:       service.NewDisposalService(repos.Disposal),
		Borrowing:      service.NewBorrowingService(repos.Borrowing, repos.KirItem),
		SarprasReport:  service.NewSarprasReportService(db),
	}
}
