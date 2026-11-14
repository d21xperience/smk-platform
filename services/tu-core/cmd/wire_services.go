package main

import (
	"os"

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
	Preference        *service.PreferenceService
	Student           *service.StudentService
	Mutation          *service.MutationService
	Registration      *service.RegistrationService
	Attendance        *service.AttendanceService
	Assessment        *service.AssessmentService
	Finance           *service.FinanceService
	Signature         *service.SignatureService
	Agenda            *service.AgendaService
	Correspondence    *service.CorrespondenceService
	Disposition       *service.DispositionService
	Class             *service.ClassService
	Stats             *service.StatisticsService
	Teacher           *service.TeacherService
	TeacherHistory    *service.TeacherHistoryService
	PtkLeave          *service.PtkLeaveService
	PtkArchive        *service.PtkArchiveService
	PtkDashboard      *service.PtkDashboardService
	BukuInduk         *service.BukuIndukService
	LaporanBerkala    *service.LaporanBerkalaService
	LaporanEmail      *service.LaporanEmailService
	LaporanAggregator *service.LaporanAggregatorService
	PtkStatistic      *service.PtkStatisticService
	Kurikulum         *service.KurikulumService
	Schedule          service.ScheduleService
	PSB               *service.PSBService
	Jurusan           *service.JurusanService
	Enrollment        *service.EnrollmentService
	Auth              *service.AuthService
	Public            *service.PublicService
	User              *service.UserService
	School            *service.SchoolService
	Logistic          *service.LogisticService
	KIR               *service.KIRService
	LetterType        *service.LetterTypeService
	Procurement       *service.ProcurementService
	GoodsReceipt      *service.GoodsReceiptService
	Document          *service.DocumentService
	Disposal          *service.DisposalService
	Statistik         *service.StatistikService
	SKP               *service.SKPService
	BukuKontrol       *service.BukuKontrolService
	UsulanKenaikan    *service.UsulanKenaikanService
	Kehadiran         *service.KehadiranService
	Pelayanan         *service.PelayananService
	Borrowing         *service.BorrowingService
	Maintenance       *service.MaintenanceService
	SarprasReport     *service.SarprasReportService
	SKPembagian       *service.SKPembagianTugasService
	Expedition        *service.ExpeditionService
	Archive           *service.ArchiveService
	Dashboard         *service.DashboardService
	Letter            *service.LetterService
	Scholarship       *service.ScholarshipService
	StudentBukuInduk  *service.StudentBukuIndukService
	KesiswaanReport   *service.KesiswaanReportService
	AlumniLegalisir   *service.AlumniLegalisirService
	AlumniTracking    *service.AlumniTrackingService
	PklPlotting       *service.PklPlottingService
	Iuran             *service.IuranService
	BukuKas           *service.BukuKasService
	Pengeluaran       *service.PengeluaranService
	RAB               *service.RABService
	Rekonsiliasi      *service.RekonsiliasiService
	Penutupan         *service.PenutupanService
	Piutang           *service.PiutangService
	UserPinTTD        *service.UserPinTTDService
	Laporan           *service.LaporanService

	KeuanganMaster *service.KeuanganMasterService
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
	agendaSvc := service.NewAgendaService(repos.Agenda)

	return &AppServices{
		Preference:        prefService,
		Student:           service.NewStudentService(repos.Student, repos.StudentDoc, minioClient, engines.Student, eventPub),
		Mutation:          service.NewMutationService(repos.Mutation, repos.Student, eventPub),
		Registration:      service.NewRegistrationService(repos.Registration, repos.Student, eventPub),
		Attendance:        service.NewAttendanceService(repos.Attendance, repos.Student, eventPub),
		Assessment:        service.NewAssessmentService(repos.Assessment, repos.Student, engines.Assessment, eventPub),
		Finance:           service.NewFinanceService(repos.Finance, repos.Student, engines.Finance, eventPub),
		Signature:         signatureSvc,
		Agenda:            agendaSvc,
		Correspondence:    service.NewCorrespondenceService(repos.Correspondence, signatureSvc, agendaSvc, repos.School, engines.Correspondence, eventPub),
		Expedition:        service.NewExpeditionService(repos.Expedition, repos.Correspondence, repos.School),
		Archive:           service.NewArchiveService(repos.Archive, repos.Correspondence),
		Disposition:       service.NewDispositionService(repos.Disposition),
		Class:             service.NewClassService(repos.Class, engines.Class, eventPub),
		Stats:             service.NewStatisticsService(repos.Stats),
		Teacher:           service.NewTeacherService(repos.Teacher),
		TeacherHistory:    service.NewTeacherHistoryService(repos.TeacherHistory),
		PtkLeave:          service.NewPtkLeaveService(repos.PtkLeave),
		PtkArchive:        service.NewPtkArchiveService(repos.PtkArchive),
		PtkDashboard:      service.NewPtkDashboardService(repos.PtkDashboard),
		BukuInduk:         service.NewBukuIndukService(repos.BukuInduk),
		LaporanBerkala:    service.NewLaporanBerkalaService(repos.LaporanBerkala, repos.School, minioClient),
		LaporanEmail:      service.NewLaporanEmailService(os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"), os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"), "noreply@sekolah.sch.id", minioClient, repos.LaporanBerkala),
		LaporanAggregator: service.NewLaporanAggregatorService(repos.PtkStatistic),
		PtkStatistic:      service.NewPtkStatisticService(repos.PtkStatistic),
		Kurikulum:         service.NewKurikulumService(repos.Subject, repos.Pembelajaran, repos.KurikulumStat),
		Schedule:          service.NewScheduleService(repos.Schedule, db),
		PSB:               service.NewPSBService(repos.PSB),
		Jurusan:           service.NewJurusanService(repos.Jurusan),
		Enrollment:        service.NewEnrollmentService(db, repos.Enrollment, repos.Class, repos.Student, engines.Student, eventPub),
		Auth:              service.NewAuthService(repos.User, repos.Academic),
		Public:            service.NewPublicService(repos.Public),
		User:              service.NewUserService(repos.User),
		School:            service.NewSchoolService(repos.School),
		Logistic:          service.NewLogisticService(repos.Logistic),
		KIR:               service.NewKIRService(db, repos.KirItem),
		LetterType:        service.NewLetterTypeService(repos.LetterType),
		Procurement:       service.NewProcurementService(repos.Procurement),
		GoodsReceipt:      service.NewGoodsReceiptService(repos.GoodsReceipt),
		Document:          service.NewDocumentService(repos.Document, minioClient),
		Disposal:          service.NewDisposalService(repos.Disposal),
		Statistik:         service.NewStatistikService(repos.Statistik),
		SKP:               service.NewSKPService(repos.SKP, repos.School, minioClient),
		BukuKontrol:       service.NewBukuKontrolService(repos.BukuKontrol),
		UsulanKenaikan:    service.NewUsulanKenaikanService(repos.UsulanKenaikan),
		Kehadiran:         service.NewKehadiranService(repos.Kehadiran),
		Pelayanan:         service.NewPelayananService(repos.Pelayanan, repos.School, minioClient),
		Borrowing:         service.NewBorrowingService(repos.Borrowing, repos.KirItem),
		Maintenance:       service.NewMaintenanceService(repos.Maintenance),
		SarprasReport:     service.NewSarprasReportService(db),
		SKPembagian:       service.NewSKPembagianTugasService(repos.SKPembagian, repos.School, minioClient),
		Dashboard:         service.NewDashboardService(repos.Dashboard),
		Letter:            service.NewLetterService(repos.Letter, repos.Student, minioClient, eventPub),
		Scholarship:       service.NewScholarshipService(repos.Scholarship, repos.Student, eventPub),
		StudentBukuInduk:  service.NewStudentBukuIndukService(repos.StudentBukuInduk, repos.School),
		KesiswaanReport:   service.NewKesiswaanReportService(repos.KesiswaanReport, repos.School),
		AlumniLegalisir:   service.NewAlumniLegalisirService(repos.AlumniLegalisir),
		AlumniTracking:    service.NewAlumniTrackingService(repos.AlumniTracking),
		PklPlotting:       service.NewPklPlottingService(repos.PklPlotting),
		Iuran:             service.NewIuranService(repos.KeuanganTransaksi, repos.KeuanganMaster),
		BukuKas:           service.NewBukuKasService(repos.KeuanganTransaksi),
		Pengeluaran:       service.NewPengeluaranService(repos.KeuanganTransaksi),
		RAB:               service.NewRABService(repos.RAB, repos.KeuanganMaster),
		Rekonsiliasi:      service.NewRekonsiliasiService(repos.Rekonsiliasi),
		Penutupan:         service.NewPenutupanService(repos.Penutupan, repos.Rekonsiliasi),
		Piutang:           service.NewPiutangService(repos.Piutang),
		UserPinTTD:        service.NewUserPinTTDService(repos.UserPinTTD),
		Laporan:           service.NewLaporanService(repos.Laporan, service.NewRABService(repos.RAB, repos.KeuanganMaster), service.NewPiutangService(repos.Piutang), repos.Penutupan),

		KeuanganMaster: service.NewKeuanganMasterService(repos.KeuanganMaster)}
}
