package main

import (
	"context"
	"log"

	"sekolah-platform/services/tu-core/internal/handler/rest"
	"sekolah-platform/services/tu-core/internal/health"
	"sekolah-platform/services/tu-core/internal/server"
)

type AppServers struct {
	GRPC *server.GRPCServer
	HTTP *server.HTTPServer
	REST *server.RESTServer
}

func (s *AppServers) Start() {
	appServer := server.NewAppServer(s.GRPC, s.HTTP, s.REST)
	appServer.Start()
}

func (s *AppServers) Stop(ctx context.Context) {
	appServer := server.NewAppServer(s.GRPC, s.HTTP, s.REST)
	appServer.Stop(ctx)
}

func setupServers(services *AppServices, repos *AppRepositories, infra *Infrastructure) *AppServers {
	grpcServer, err := server.NewGRPCServer(
		server.Config{Port: getEnvAsInt("GRPC_PORT", 50051)},
		services.Student, services.Mutation, services.Registration,
		services.Attendance, services.Assessment, services.Finance,
		services.Correspondence, &dummyLogger{},
	)
	if err != nil {
		log.Fatalf("❌ Gagal membuat gRPC Server: %v", err)
	}

	httpServer := server.NewHTTPServer(
		server.HTTPServerConfig{Port: getEnvAsInt("HTTP_PORT", 8080)},
		health.NewHealthHandler(infra.DB, infra.NC, infra.JS),
	)

	scheduleHandler := rest.NewScheduleHandler(repos.Class, repos.Academic, infra.DB)

	restHandler := rest.NewRouter(
		services.Student, services.Mutation, services.Registration,
		services.Attendance, services.Assessment, services.Finance,
		services.Correspondence, services.Class, services.Enrollment,
		services.Stats, services.Teacher, services.Kurikulum,
		services.Schedule, services.Jurusan, services.Preference,
		services.Auth, services.Public, services.User, services.School,
		scheduleHandler, services.Disposition, services.Signature,
		infra.MinIO, services.TeacherHistory, services.PtkLeave,
		services.PtkArchive, services.PtkDashboard, services.BukuInduk,
		services.SKPembagian, services.Statistik, services.SKP, services.BukuKontrol,
		services.Logistic, services.KIR, services.LetterType,
		services.Procurement, services.GoodsReceipt, services.Document, services.Disposal,
		services.Borrowing, services.Maintenance, services.SarprasReport,
	)
	restServer := server.NewRESTServer(
		server.RESTServerConfig{Port: getEnvAsInt("REST_PORT", 8081)},
		restHandler,
	)

	return &AppServers{GRPC: grpcServer, HTTP: httpServer, REST: restServer}
}
