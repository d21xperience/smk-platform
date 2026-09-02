package service

import (
	"context"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/assessment/engine"
	"sekolah-platform/services/tu-core/internal/domain/assessment/models"
	"sekolah-platform/services/tu-core/internal/domain/assessment/valueobjects"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/google/uuid"
)

type AssessmentService struct {
	repo        repository.AssessmentRepository
	studentRepo repository.StudentRepository
	assEngine   *engine.AssessmentEngine
	publisher   EventPublisher
}

func NewAssessmentService(
	repo repository.AssessmentRepository,
	studentRepo repository.StudentRepository,
	assEngine *engine.AssessmentEngine,
	publisher EventPublisher,
) *AssessmentService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &AssessmentService{
		repo:        repo,
		studentRepo: studentRepo,
		assEngine:   assEngine,
		publisher:   publisher,
	}
}

type CreateAssessmentServiceRequest struct {
	StudentID      string
	SubjectID      string
	AssessmentType valueobjects.AssessmentType
	Score          float64
	MaxScore       float64
	AssessmentDate time.Time
	Semester       int
	Notes          string
}

func (s *AssessmentService) CreateAssessment(ctx context.Context, req CreateAssessmentServiceRequest) (*AssessmentResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	data := models.AssessmentData{
		StudentID:      req.StudentID,
		SubjectID:      req.SubjectID,
		AssessmentType: req.AssessmentType,
		Score:          req.Score,
		MaxScore:       req.MaxScore,
		AssessmentDate: req.AssessmentDate,
		Semester:       req.Semester,
		Notes:          req.Notes,
	}

	result := s.assEngine.CreateAssessment(data, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("CreateAssessment", "Assessment", "", "invalid assessment data", result.Error)
	}

	if err := s.repo.Create(ctx, result.Assessment); err != nil {
		return nil, NewServiceError("CreateAssessment", "Assessment", result.Assessment.ID, "failed to save to database", err)
	}

	event := DomainEvent{
		EventType:     "AssessmentCreated",
		AggregateType: "Assessment",
		AggregateID:   result.Assessment.ID,
		Payload: map[string]interface{}{
			"assessmentId":   result.Assessment.ID,
			"studentId":      result.Assessment.StudentID,
			"subjectId":      result.Assessment.SubjectID,
			"assessmentType": result.Assessment.AssessmentType.String(),
			"score":          result.Assessment.Score,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toAssessmentResponse(result.Assessment), nil
}

type UpdateAssessmentServiceRequest struct {
	Score    float64
	MaxScore float64
	Notes    string
}

func (s *AssessmentService) UpdateAssessment(ctx context.Context, assessmentID string, req UpdateAssessmentServiceRequest) (*AssessmentResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	a, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, assessmentID)
	if err != nil {
		return nil, NewServiceError("UpdateAssessment", "Assessment", assessmentID, "failed to get assessment", err)
	}

	result := s.assEngine.UpdateAssessment(a, req.Score, req.MaxScore, req.Notes)
	if result.Error != nil {
		return nil, NewServiceError("UpdateAssessment", "Assessment", assessmentID, "failed to update", result.Error)
	}

	if err := s.repo.Update(ctx, result.Assessment); err != nil {
		return nil, NewServiceError("UpdateAssessment", "Assessment", assessmentID, "failed to save to database", err)
	}

	event := DomainEvent{
		EventType:     "AssessmentUpdated",
		AggregateType: "Assessment",
		AggregateID:   result.Assessment.ID,
		Payload: map[string]interface{}{
			"assessmentId": result.Assessment.ID,
			"studentId":    result.Assessment.StudentID,
			"subjectId":    result.Assessment.SubjectID,
			"newScore":     result.Assessment.Score,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toAssessmentResponse(result.Assessment), nil
}

func (s *AssessmentService) GetAssessment(ctx context.Context, assessmentID string) (*AssessmentResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	a, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, assessmentID)
	if err != nil {
		return nil, NewServiceError("GetAssessment", "Assessment", assessmentID, "failed to get assessment", err)
	}

	return toAssessmentResponse(a), nil
}

func (s *AssessmentService) GetStudentAssessments(ctx context.Context, studentID string, semester int, subjectID string) (*StudentAssessmentsResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	assessments, err := s.repo.GetByStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID, semester, subjectID)
	if err != nil {
		return nil, NewServiceError("GetStudentAssessments", "Assessment", studentID, "failed to get assessments", err)
	}

	// Calculate averages per subject
	averages := calculateSubjectAverages(assessments)

	var responses []AssessmentResponse
	for _, a := range assessments {
		responses = append(responses, *toAssessmentResponse(&a))
	}

	return &StudentAssessmentsResponse{
		Assessments: responses,
		Averages:    averages,
	}, nil
}

func (s *AssessmentService) GetSubjectAssessments(ctx context.Context, subjectID, classID string, semester int) ([]*AssessmentResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	assessments, err := s.repo.GetBySubject(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, subjectID, classID, semester)
	if err != nil {
		return nil, NewServiceError("GetSubjectAssessments", "Assessment", subjectID, "failed to get assessments", err)
	}

	var responses []*AssessmentResponse
	for _, a := range assessments {
		responses = append(responses, toAssessmentResponse(&a))
	}

	return responses, nil
}

func (s *AssessmentService) ListAssessments(ctx context.Context, studentID, subjectID, assessmentType string, semester int) ([]*AssessmentResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter := repository.AssessmentFilter{
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		StudentID:        studentID,
		SubjectID:        subjectID,
		AssessmentType:   assessmentType,
		Semester:         semester,
		Limit:            100,
		Offset:           0,
	}

	assessments, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, NewServiceError("ListAssessments", "Assessment", "", "failed to list assessments", err)
	}

	var responses []*AssessmentResponse
	for _, a := range assessments {
		responses = append(responses, toAssessmentResponse(&a))
	}

	return responses, nil
}

func (s *AssessmentService) GetStudentAverage(ctx context.Context, studentID string, semester int) (*StudentAverageResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	assessments, err := s.repo.GetByStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID, semester, "")
	if err != nil {
		return nil, NewServiceError("GetStudentAverage", "Assessment", studentID, "failed to get assessments", err)
	}

	averages := calculateSubjectAverages(assessments)

	// Calculate overall average
	var totalAvg float64
	var count int
	for _, avg := range averages {
		totalAvg += avg.Average
		count++
	}

	overallAvg := 0.0
	if count > 0 {
		overallAvg = totalAvg / float64(count)
	}

	return &StudentAverageResponse{
		Averages:       averages,
		OverallAverage: overallAvg,
	}, nil
}

// Helper functions
func calculateSubjectAverages(assessments []models.Assessment) []SubjectAverageResponse {
	subjectMap := make(map[string]*SubjectAverageResponse)

	for _, a := range assessments {
		if _, exists := subjectMap[a.SubjectID]; !exists {
			subjectMap[a.SubjectID] = &SubjectAverageResponse{
				SubjectID: a.SubjectID,
			}
		}
		subjectMap[a.SubjectID].TotalScore += a.Percentage()
		subjectMap[a.SubjectID].TotalAssessments++
	}

	var averages []SubjectAverageResponse
	for _, avg := range subjectMap {
		avg.Average = avg.TotalScore / float64(avg.TotalAssessments)
		avg.Grade = valueobjects.CalculateGrade(avg.Average)
		averages = append(averages, *avg)
	}

	return averages
}

type StudentAssessmentsResponse struct {
	Assessments []AssessmentResponse
	Averages    []SubjectAverageResponse
}

type StudentAverageResponse struct {
	Averages       []SubjectAverageResponse
	OverallAverage float64
}

func toAssessmentResponse(a *models.Assessment) *AssessmentResponse {
	if a == nil {
		return nil
	}
	return &AssessmentResponse{
		ID:               a.ID,
		SchoolID:         a.SchoolID,
		AcademicPeriodID: a.AcademicPeriodID,
		StudentID:        a.StudentID,
		SubjectID:        a.SubjectID,
		AssessmentType:   a.AssessmentType.String(),
		Score:            a.Score,
		MaxScore:         a.MaxScore,
		Percentage:       a.Percentage(),
		AssessmentDate:   a.AssessmentDate,
		Semester:         a.Semester,
		Average:          0, // Akan dihitung di layer service jika diperlukan
		Notes:            a.Notes,
		CreatedBy:        a.CreatedBy,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
	}
}
