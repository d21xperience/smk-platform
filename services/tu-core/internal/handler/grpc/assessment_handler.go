package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"sekolah-platform/services/tu-core/internal/domain/assessment/valueobjects"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
	pb "sekolah-platform/services/tu-core/proto/assessment/v1"
)

type AssessmentGRPCHandler struct {
	pb.UnimplementedAssessmentServiceServer
	service *service.AssessmentService
}

func NewAssessmentGRPCHandler(svc *service.AssessmentService) *AssessmentGRPCHandler {
	return &AssessmentGRPCHandler{service: svc}
}

func (h *AssessmentGRPCHandler) CreateAssessment(ctx context.Context, req *pb.CreateAssessmentRequest) (*pb.CreateAssessmentResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing operational context")
	}

	assessmentType, _ := valueobjects.NewAssessmentType(req.GetAssessmentType())

	serviceReq := service.CreateAssessmentServiceRequest{
		StudentID:      req.GetStudentId(),
		SubjectID:      req.GetSubjectId(),
		AssessmentType: assessmentType,
		Score:          req.GetScore(),
		MaxScore:       req.GetMaxScore(),
		AssessmentDate: req.GetAssessmentDate().AsTime(),
		Semester:       int(req.GetSemester()),
		Notes:          req.GetNotes(),
	}

	a, err := h.service.CreateAssessment(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.CreateAssessmentResponse{
		Assessment: mapAssessmentResponseToProto(a),
	}, nil
}

func (h *AssessmentGRPCHandler) UpdateAssessment(ctx context.Context, req *pb.UpdateAssessmentRequest) (*pb.UpdateAssessmentResponse, error) {
	serviceReq := service.UpdateAssessmentServiceRequest{
		Score:    req.GetScore(),
		MaxScore: req.GetMaxScore(),
		Notes:    req.GetNotes(),
	}

	a, err := h.service.UpdateAssessment(ctx, req.GetAssessmentId(), serviceReq)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.UpdateAssessmentResponse{
		Assessment: mapAssessmentResponseToProto(a),
	}, nil
}

func (h *AssessmentGRPCHandler) GetAssessment(ctx context.Context, req *pb.GetAssessmentRequest) (*pb.GetAssessmentResponse, error) {
	a, err := h.service.GetAssessment(ctx, req.GetAssessmentId())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.GetAssessmentResponse{Assessment: mapAssessmentResponseToProto(a)}, nil
}

func (h *AssessmentGRPCHandler) GetStudentAssessments(ctx context.Context, req *pb.GetStudentAssessmentsRequest) (*pb.GetStudentAssessmentsResponse, error) {
	result, err := h.service.GetStudentAssessments(ctx, req.GetStudentId(), int(req.GetSemester()), req.GetSubjectId())
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbAssessments []*pb.Assessment
	for _, a := range result.Assessments {
		pbAssessments = append(pbAssessments, mapAssessmentResponseToProto(&a))
	}

	var pbAverages []*pb.SubjectAverage
	for _, avg := range result.Averages {
		pbAverages = append(pbAverages, mapSubjectAverageToProto(&avg))
	}

	return &pb.GetStudentAssessmentsResponse{
		Assessments: pbAssessments,
		Averages:    pbAverages,
	}, nil
}

func (h *AssessmentGRPCHandler) GetSubjectAssessments(ctx context.Context, req *pb.GetSubjectAssessmentsRequest) (*pb.GetSubjectAssessmentsResponse, error) {
	assessments, err := h.service.GetSubjectAssessments(ctx, req.GetSubjectId(), req.GetClassId(), int(req.GetSemester()))
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbAssessments []*pb.Assessment
	for _, a := range assessments {
		pbAssessments = append(pbAssessments, mapAssessmentResponseToProto(a))
	}

	return &pb.GetSubjectAssessmentsResponse{
		Assessments: pbAssessments,
	}, nil
}

func (h *AssessmentGRPCHandler) ListAssessments(ctx context.Context, req *pb.ListAssessmentsRequest) (*pb.ListAssessmentsResponse, error) {
	assessments, err := h.service.ListAssessments(ctx, req.GetStudentId(), req.GetSubjectId(), req.GetAssessmentType(), int(req.GetSemester()))
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbAssessments []*pb.Assessment
	for _, a := range assessments {
		pbAssessments = append(pbAssessments, mapAssessmentResponseToProto(a))
	}

	return &pb.ListAssessmentsResponse{
		Assessments: pbAssessments,
		Total:       int32(len(pbAssessments)),
	}, nil
}

func (h *AssessmentGRPCHandler) GetStudentAverage(ctx context.Context, req *pb.GetStudentAverageRequest) (*pb.GetStudentAverageResponse, error) {
	result, err := h.service.GetStudentAverage(ctx, req.GetStudentId(), int(req.GetSemester()))
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbAverages []*pb.SubjectAverage
	for _, avg := range result.Averages {
		pbAverages = append(pbAverages, mapSubjectAverageToProto(&avg))
	}

	return &pb.GetStudentAverageResponse{
		Averages:       pbAverages,
		OverallAverage: result.OverallAverage,
	}, nil
}

// Mapping helpers
func mapAssessmentResponseToProto(a *service.AssessmentResponse) *pb.Assessment {
	if a == nil {
		return nil
	}
	return &pb.Assessment{
		Id:               a.ID,
		SchoolId:         a.SchoolID,
		AcademicPeriodId: a.AcademicPeriodID,
		StudentId:        a.StudentID,
		SubjectId:        a.SubjectID,
		AssessmentType:   a.AssessmentType,
		Score:            a.Score,
		MaxScore:         a.MaxScore,
		Percentage:       a.Percentage,
		AssessmentDate:   timestamppb.New(a.AssessmentDate),
		Semester:         int32(a.Semester),
		Notes:            a.Notes,
		CreatedBy:        a.CreatedBy,
		CreatedAt:        timestamppb.New(a.CreatedAt),
		UpdatedAt:        timestamppb.New(a.UpdatedAt),
	}
}

func mapSubjectAverageToProto(avg *service.SubjectAverageResponse) *pb.SubjectAverage {
	if avg == nil {
		return nil
	}
	return &pb.SubjectAverage{
		SubjectId:        avg.SubjectID,
		SubjectName:      avg.SubjectName,
		Average:          avg.Average,
		TotalAssessments: int32(avg.TotalAssessments),
		Grade:            avg.Grade,
	}
}

// Ensure time package is used
var _ = time.Now
