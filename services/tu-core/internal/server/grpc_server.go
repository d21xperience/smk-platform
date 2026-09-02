package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	gRPCHandler "sekolah-platform/services/tu-core/internal/handler/grpc"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
	pbAssessment "sekolah-platform/services/tu-core/proto/assessment/v1"
	pbAttendance "sekolah-platform/services/tu-core/proto/attendance/v1"
	pbCorrespondence "sekolah-platform/services/tu-core/proto/correspondence/v1"
	pbFinance "sekolah-platform/services/tu-core/proto/finance/v1"
	pbMutation "sekolah-platform/services/tu-core/proto/mutation/v1"
	pbRegistration "sekolah-platform/services/tu-core/proto/registration/v1"
	pb "sekolah-platform/services/tu-core/proto/student/v1"
)

type Config struct {
	Port int
}

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
	config   Config
}

func NewGRPCServer(
	config Config,
	studentService *service.StudentService,
	mutationService *service.MutationService,
	registrationService *service.RegistrationService,
	attendanceService *service.AttendanceService,
	assessmentService *service.AssessmentService,
	financeService *service.FinanceService,
	correspondenceService *service.CorrespondenceService,
	logger middleware.Logger,
) (*GRPCServer, error) {

	chain := middleware.ChainUnaryInterceptor(
		middleware.RecoveryUnaryInterceptor(logger),
		middleware.LoggingUnaryInterceptor(logger),
		middleware.OperationalContextUnaryInterceptor(),
	)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(chain),
	)

	// Register services
	studentHandler := gRPCHandler.NewStudentGRPCHandler(studentService)
	pb.RegisterStudentServiceServer(server, studentHandler)

	mutationHandler := gRPCHandler.NewMutationGRPCHandler(mutationService)
	pbMutation.RegisterMutationServiceServer(server, mutationHandler)

	registrationHandler := gRPCHandler.NewRegistrationGRPCHandler(registrationService)
	pbRegistration.RegisterRegistrationServiceServer(server, registrationHandler)

	attendanceHandler := gRPCHandler.NewAttendanceGRPCHandler(attendanceService)
	pbAttendance.RegisterAttendanceServiceServer(server, attendanceHandler)

	assessmentHandler := gRPCHandler.NewAssessmentGRPCHandler(assessmentService)
	pbAssessment.RegisterAssessmentServiceServer(server, assessmentHandler)

	financeHandler := gRPCHandler.NewFinanceGRPCHandler(financeService)
	pbFinance.RegisterFinanceServiceServer(server, financeHandler)

	correspondenceHandler := gRPCHandler.NewCorrespondenceGRPCHandler(correspondenceService)
	pbCorrespondence.RegisterCorrespondenceServiceServer(server, correspondenceHandler)

	reflection.Register(server)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return nil, fmt.Errorf("gagal listen pada port %d: %w", config.Port, err)
	}

	return &GRPCServer{
		server:   server,
		listener: listener,
		config:   config,
	}, nil
}

func (s *GRPCServer) Start() error {
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

func (s *GRPCServer) ForceStop() {
	s.server.Stop()
}

func (s *GRPCServer) Port() int {
	return s.config.Port
}
