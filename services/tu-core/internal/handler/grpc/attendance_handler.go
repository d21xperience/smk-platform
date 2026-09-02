package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
	pb "sekolah-platform/services/tu-core/proto/attendance/v1"
)

type AttendanceGRPCHandler struct {
	pb.UnimplementedAttendanceServiceServer
	service *service.AttendanceService
}

func NewAttendanceGRPCHandler(svc *service.AttendanceService) *AttendanceGRPCHandler {
	return &AttendanceGRPCHandler{service: svc}
}

func (h *AttendanceGRPCHandler) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing operational context")
	}

	serviceReq := service.CreateSessionRequest{
		ClassID:     req.GetClassId(),
		SubjectID:   req.GetSubjectId(),
		TeacherID:   req.GetTeacherId(),
		SessionDate: req.GetSessionDate().AsTime(),
		StartTime:   req.GetStartTime().AsTime(),
		EndTime:     req.GetEndTime().AsTime(),
		Notes:       req.GetNotes(),
	}

	session, err := h.service.CreateSession(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.CreateSessionResponse{
		Session: mapSessionResponseToProto(session),
	}, nil
}

func (h *AttendanceGRPCHandler) SubmitAttendance(ctx context.Context, req *pb.SubmitAttendanceRequest) (*pb.SubmitAttendanceResponse, error) {
	// 1. Map proto records ke service.AttendanceRecordInput
	var records []service.AttendanceRecordInput
	for _, r := range req.GetRecords() {
		records = append(records, service.AttendanceRecordInput{
			StudentID: r.GetStudentId(),
			Status:    r.GetStatus(),
			Note:      r.GetNote(),
		})
	}

	// 2. Gunakan SubmitSessionAttendanceRequest (yang sudah didefinisikan di attendance_service.go)
	serviceReq := service.SubmitSessionAttendanceRequest{
		SessionID: req.GetSessionId(),
		Records:   records,
	}

	// 3. Panggil service
	session, err := h.service.SubmitAttendance(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.SubmitAttendanceResponse{
		Session: mapSessionResponseToProto(session),
	}, nil
}

func (h *AttendanceGRPCHandler) GetSession(ctx context.Context, req *pb.GetSessionRequest) (*pb.GetSessionResponse, error) {
	session, err := h.service.GetSession(ctx, req.GetSessionId())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.GetSessionResponse{Session: mapSessionResponseToProto(session)}, nil
}

func (h *AttendanceGRPCHandler) GetStudentAttendance(ctx context.Context, req *pb.GetStudentAttendanceRequest) (*pb.GetStudentAttendanceResponse, error) {
	result, err := h.service.GetStudentAttendance(ctx, req.GetStudentId(), int(req.GetMonth()), int(req.GetYear()))
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbRecords []*pb.AttendanceRecord
	for _, r := range result.Records {
		pbRecords = append(pbRecords, mapRecordResponseToProto(&r))
	}

	return &pb.GetStudentAttendanceResponse{
		Records: pbRecords,
		Summary: mapSummaryResponseToProto(&result.Summary),
	}, nil
}

func (h *AttendanceGRPCHandler) ListSessions(ctx context.Context, req *pb.ListSessionsRequest) (*pb.ListSessionsResponse, error) {
	var startDate, endDate time.Time
	if req.GetStartDate() != nil {
		startDate = req.GetStartDate().AsTime()
	}
	if req.GetEndDate() != nil {
		endDate = req.GetEndDate().AsTime()
	}

	sessions, err := h.service.ListSessions(ctx, req.GetClassId(), req.GetStatus(), startDate, endDate)
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbSessions []*pb.AttendanceSession
	for _, s := range sessions {
		pbSessions = append(pbSessions, mapSessionResponseToProto(s))
	}

	return &pb.ListSessionsResponse{
		Sessions: pbSessions,
		Total:    int32(len(pbSessions)),
	}, nil
}

func (h *AttendanceGRPCHandler) GetAttendanceSummary(ctx context.Context, req *pb.GetAttendanceSummaryRequest) (*pb.GetAttendanceSummaryResponse, error) {
	summaries, err := h.service.GetAttendanceSummary(ctx, req.GetClassId(), int(req.GetMonth()), int(req.GetYear()))
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbSummaries []*pb.StudentAttendanceSummary
	for _, s := range summaries {
		pbSummaries = append(pbSummaries, &pb.StudentAttendanceSummary{
			StudentId:   s.StudentID,
			StudentName: s.StudentName,
			Summary:     mapSummaryResponseToProto(&s.Summary),
		})
	}

	return &pb.GetAttendanceSummaryResponse{
		Students: pbSummaries,
	}, nil
}

// Mapping helpers
func mapSessionResponseToProto(s *service.SessionResponse) *pb.AttendanceSession {
	if s == nil {
		return nil
	}

	var pbRecords []*pb.AttendanceRecord
	for _, r := range s.Records {
		pbRecords = append(pbRecords, mapRecordResponseToProto(&r))
	}

	return &pb.AttendanceSession{
		Id:               s.ID,
		SchoolId:         s.SchoolID,
		AcademicPeriodId: s.AcademicPeriodID,
		ClassId:          s.ClassID,
		SubjectId:        s.SubjectID,
		TeacherId:        s.TeacherID,
		SessionDate:      timestamppb.New(s.SessionDate),
		StartTime:        timestamppb.New(s.StartTime),
		EndTime:          timestamppb.New(s.EndTime),
		Status:           s.Status,
		Notes:            s.Notes,
		Records:          pbRecords,
		CreatedAt:        timestamppb.New(s.CreatedAt),
		UpdatedAt:        timestamppb.New(s.UpdatedAt),
	}
}

func mapRecordResponseToProto(r *service.RecordResponse) *pb.AttendanceRecord {
	if r == nil {
		return nil
	}
	return &pb.AttendanceRecord{
		Id:        r.ID,
		SessionId: r.SessionID,
		StudentId: r.StudentID,
		Status:    r.Status,
		Note:      r.Note,
		CreatedAt: timestamppb.New(r.CreatedAt),
		UpdatedAt: timestamppb.New(r.UpdatedAt),
	}
}

func mapSummaryResponseToProto(s *service.AttendanceSummaryResponse) *pb.AttendanceSummary {
	if s == nil {
		return nil
	}
	return &pb.AttendanceSummary{
		TotalSessions:        int32(s.TotalSessions),
		PresentCount:         int32(s.PresentCount),
		AbsentCount:          int32(s.AbsentCount),
		LateCount:            int32(s.LateCount),
		SickCount:            int32(s.SickCount),
		PermissionCount:      int32(s.PermissionCount),
		AttendancePercentage: s.AttendancePercentage,
	}
}
