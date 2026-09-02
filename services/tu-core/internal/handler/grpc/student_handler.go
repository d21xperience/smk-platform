package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	platformcontext "sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/service"
	pb "sekolah-platform/services/tu-core/proto/student/v1"
)

type StudentGRPCHandler struct {
	pb.UnimplementedStudentServiceServer
	service *service.StudentService
}

func NewStudentGRPCHandler(svc *service.StudentService) *StudentGRPCHandler {
	return &StudentGRPCHandler{service: svc}
}

// extractOpCtxFromMetadata mengambil Operational Context dari gRPC Headers.
func extractOpCtxFromMetadata(ctx context.Context) (*platformcontext.OperationalContext, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "metadata tidak ditemukan")
	}

	getVal := func(key string) string {
		if values := md.Get(key); len(values) > 0 {
			return values[0]
		}
		return ""
	}

	schoolID := getVal("x-school-id")
	academicYear := getVal("x-academic-year")
	semesterStr := getVal("x-semester")
	userID := getVal("x-user-id")
	role := getVal("x-user-role")

	semester, _ := strconv.Atoi(semesterStr)

	opCtx, err := platformcontext.New(schoolID, academicYear, semester, userID, role)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid operational context: %v", err)
	}

	return opCtx, nil
}

// RegisterStudent menangani pembuatan data siswa
func (h *StudentGRPCHandler) RegisterStudent(ctx context.Context, req *pb.RegisterStudentRequest) (*pb.RegisterStudentResponse, error) {
	opCtx, err := extractOpCtxFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	// Gabungkan nama
	fullName := strings.TrimSpace(fmt.Sprintf("%s %s %s",
		req.GetFullName().GetFirstName(),
		req.GetFullName().GetMiddleName(),
		req.GetFullName().GetLastName()))

	// Gabungkan address
	addr := req.GetAddress()
	addressStr := strings.TrimSpace(fmt.Sprintf("%s, RT/RW%s, %s, %s, %s, %s",
		addr.GetStreet(), addr.GetRtRw(), addr.GetVillage(),
		addr.GetDistrict(), addr.GetCity(), addr.GetPostalCode()))

	// PERBAIKAN: Konversi Gender enum ke string code yang benar
	// "GENDER_MALE" → "MALE", "GENDER_FEMALE" → "FEMALE"
	genderCode := protoGenderToCode(req.GetGender())

	serviceReq := service.CreateStudentRequest{
		NIS:         req.GetNis(),
		NISN:        req.GetNisn(),
		Name:        fullName,
		Gender:      genderCode, // ← PERBAIKAN: gunakan "MALE"/"FEMALE", bukan "GENDER_MALE"
		BirthDate:   req.GetBirthDate().AsTime(),
		BirthPlace:  "", // Bisa ditambahkan jika ada di proto
		Address:     addressStr,
		Phone:       req.GetContactInfo().GetPhone(),
		Email:       req.GetContactInfo().GetEmail(),
		ParentName:  req.GetGuardianInfo().GetName(),
		ParentPhone: req.GetGuardianInfo().GetPhone(),
		ParentEmail: req.GetGuardianInfo().GetOccupation(),
		ClassID:     "",
	}

	userIDInt, _ := strconv.Atoi(opCtx.UserID)
	requestID := "req-" + opCtx.UserID

	log.Printf("[RegisterStudent] 📝 Creating student: NIS=%s, NISN=%s, Name=%s, School=%s, Period=%s",
		serviceReq.NIS, serviceReq.NISN, serviceReq.Name, opCtx.SchoolID, opCtx.AcademicPeriodID)

	result, err := h.service.CreateStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, userIDInt, opCtx.Role, requestID, serviceReq)
	if err != nil {
		log.Printf("[RegisterStudent] ❌ Error: %v", err)
		return nil, mapServiceError(err)
	}

	log.Printf("[RegisterStudent] ✅ Student created: ID=%s", result.ID)

	return &pb.RegisterStudentResponse{
		Student: mapToProtoStudent(result),
	}, nil
}

// EnrollStudent menangani pendaftaran siswa ke kelas tertentu
func (h *StudentGRPCHandler) EnrollStudent(ctx context.Context, req *pb.EnrollStudentRequest) (*pb.EnrollStudentResponse, error) {
	opCtx, err := extractOpCtxFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	userIDInt, _ := strconv.Atoi(opCtx.UserID)
	requestID := "req-" + opCtx.UserID

	err = h.service.EnrollStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, userIDInt, opCtx.Role, requestID, req.GetStudentId(), req.GetEnrollmentId(), req.GetClassId())
	if err != nil {
		return nil, mapServiceError(err)
	}

	student, err := h.service.GetStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.GetStudentId())
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.EnrollStudentResponse{Student: mapToProtoStudent(student)}, nil
}

// GraduateStudent menangani proses kelulusan
func (h *StudentGRPCHandler) GraduateStudent(ctx context.Context, req *pb.GraduateStudentRequest) (*pb.GraduateStudentResponse, error) {
	opCtx, err := extractOpCtxFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	userIDInt, _ := strconv.Atoi(opCtx.UserID)
	requestID := "req-" + opCtx.UserID

	err = h.service.GraduateStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, userIDInt, opCtx.Role, requestID, req.GetStudentId())
	if err != nil {
		return nil, mapServiceError(err)
	}

	student, _ := h.service.GetStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.GetStudentId())
	return &pb.GraduateStudentResponse{Student: mapToProtoStudent(student)}, nil
}

// TransferStudent menangani mutasi siswa
func (h *StudentGRPCHandler) TransferStudent(ctx context.Context, req *pb.TransferStudentRequest) (*pb.TransferStudentResponse, error) {
	opCtx, err := extractOpCtxFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	userIDInt, _ := strconv.Atoi(opCtx.UserID)
	requestID := "req-" + opCtx.UserID

	err = h.service.TransferStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, userIDInt, opCtx.Role, requestID, req.GetStudentId(), req.GetTargetSchool(), req.GetReason())
	if err != nil {
		return nil, mapServiceError(err)
	}

	student, _ := h.service.GetStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.GetStudentId())
	return &pb.TransferStudentResponse{Student: mapToProtoStudent(student)}, nil
}

// GetStudent mengambil detail siswa
func (h *StudentGRPCHandler) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.GetStudentResponse, error) {
	opCtx, err := extractOpCtxFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	student, err := h.service.GetStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.GetStudentId())
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.GetStudentResponse{Student: mapToProtoStudent(student)}, nil
}

// CheckNisnAvailability cek ketersediaan NISN
func (h *StudentGRPCHandler) CheckNisnAvailability(ctx context.Context, req *pb.CheckNisnRequest) (*pb.CheckAvailabilityResponse, error) {
	opCtx, err := extractOpCtxFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	isAvailable, err := h.service.CheckNISNAvailability(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.GetNisn())
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.CheckAvailabilityResponse{Available: isAvailable}, nil
}

// CheckNisAvailability cek ketersediaan NIS
func (h *StudentGRPCHandler) CheckNisAvailability(ctx context.Context, req *pb.CheckNisRequest) (*pb.CheckAvailabilityResponse, error) {
	opCtx, err := extractOpCtxFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	isAvailable, err := h.service.CheckNISAvailability(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.GetNis())
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.CheckAvailabilityResponse{Available: isAvailable}, nil
}

// === Helper Functions ===

// protoGenderToCode mengkonversi proto Gender enum ke string code yang diterima domain.
// "GENDER_MALE" → "MALE", "GENDER_FEMALE" → "FEMALE"
func protoGenderToCode(g pb.Gender) string {
	switch g {
	case pb.Gender_GENDER_MALE:
		return "MALE"
	case pb.Gender_GENDER_FEMALE:
		return "FEMALE"
	default:
		return ""
	}
}

func mapToProtoStudent(s *service.StudentResponse) *pb.Student {
	if s == nil {
		return nil
	}

	// Mapping status domain (English) ke proto enum
	statusEnum := pb.StudentStatus_STUDENT_STATUS_UNSPECIFIED
	switch strings.ToUpper(s.Status) {
	case "ACTIVE":
		statusEnum = pb.StudentStatus_STUDENT_STATUS_ACTIVE
	case "GRADUATED":
		statusEnum = pb.StudentStatus_STUDENT_STATUS_GRADUATED
	case "TRANSFERRED":
		statusEnum = pb.StudentStatus_STUDENT_STATUS_TRANSFERRED
	case "DROPPED", "DROPOUT":
		statusEnum = pb.StudentStatus_STUDENT_STATUS_DROPPED
	case "ALUMNI":
		statusEnum = pb.StudentStatus_STUDENT_STATUS_ALUMNI
	}

	// Mapping gender
	genderEnum := pb.Gender_GENDER_UNSPECIFIED
	switch strings.ToUpper(s.Gender) {
	case "MALE":
		genderEnum = pb.Gender_GENDER_MALE
	case "FEMALE":
		genderEnum = pb.Gender_GENDER_FEMALE
	}

	// Parsing nama
	names := strings.SplitN(s.Name, " ", 3)
	firstName := names[0]
	middleName := ""
	lastName := ""
	if len(names) > 1 {
		middleName = names[1]
	}
	if len(names) > 2 {
		lastName = names[2]
	}

	return &pb.Student{
		StudentId: s.ID,
		SchoolId:  s.SchoolID,
		Nisn:      s.NISN,
		Nis:       s.NIS,
		FullName: &pb.FullName{
			FirstName:  firstName,
			MiddleName: middleName,
			LastName:   lastName,
		},
		BirthDate: timestamppb.New(s.BirthDate),
		Gender:    genderEnum,
		Address: &pb.Address{
			Street:     s.Address,
			RtRw:       "",
			Village:    "",
			District:   "",
			City:       "",
			PostalCode: "",
		},
		ContactInfo: &pb.ContactInfo{
			Phone: s.Phone,
			Email: s.Email,
		},
		GuardianInfo: &pb.GuardianInfo{
			Name:       s.ParentName,
			Relation:   "Orang Tua",
			Phone:      s.ParentPhone,
			Occupation: "",
		},
		Status:    statusEnum,
		CreatedAt: timestamppb.New(s.CreatedAt),
		UpdatedAt: timestamppb.New(s.UpdatedAt),
	}
}

// mapServiceError memetakan service error ke gRPC status error.
// PERBAIKAN: Menggunakan errors.Is() untuk unwrap error yang di-wrap.
func mapServiceError(err error) error {
	if err == nil {
		return nil
	}

	// PERBAIKAN: Gunakan errors.Is() untuk match wrapped errors
	switch {
	case errors.Is(err, service.ErrNotFound):
		return status.Errorf(codes.NotFound, "data tidak ditemukan")
	case errors.Is(err, service.ErrDuplicateEntry):
		return status.Errorf(codes.AlreadyExists, "data sudah ada")
	case errors.Is(err, service.ErrInvalidInput):
		return status.Errorf(codes.InvalidArgument, "input tidak valid: %v", err)
	case errors.Is(err, service.ErrMissingOperationalContext):
		return status.Errorf(codes.InvalidArgument, "missing operational context")
	default:
		// PERBAIKAN: Log error detail untuk debugging
		log.Printf("[mapServiceError] Unhandled error: %T: %v", err, err)
		return status.Errorf(codes.Internal, "internal server error: %v", err)
	}
}

// ListStudents mengambil daftar siswa dengan filter dan pagination
func (h *StudentGRPCHandler) ListStudents(ctx context.Context, req *pb.ListStudentsRequest) (*pb.ListStudentsResponse, error) {
	opCtx, err := extractOpCtxFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	filter := service.StudentFilterRequest{
		Keyword: req.GetSearch(),
		Status:  req.GetStatus(),
		ClassID: req.GetClassId(),
		Limit:   int(req.GetLimit()),
		Offset:  int(req.GetOffset()),
	}

	result, err := h.service.ListStudents(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, filter)
	if err != nil {
		return nil, mapServiceError(err)
	}

	protoStudents := make([]*pb.Student, 0, len(result.Data))
	for _, s := range result.Data {
		protoStudents = append(protoStudents, mapToProtoStudent(&s))
	}

	return &pb.ListStudentsResponse{
		Data:    protoStudents,
		Total:   int32(result.Total),
		Limit:   int32(result.Limit),
		Offset:  int32(result.Offset),
		HasMore: result.HasMore,
	}, nil
}

// UpdateStudent memperbarui data siswa
func (h *StudentGRPCHandler) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.UpdateStudentResponse, error) {
	opCtx, err := extractOpCtxFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	updateReq := service.UpdateStudentRequest{}

	if req.GetFullName() != nil {
		updateReq.Name = strings.TrimSpace(fmt.Sprintf("%s %s %s",
			req.GetFullName().GetFirstName(),
			req.GetFullName().GetMiddleName(),
			req.GetFullName().GetLastName()))
	}
	if req.GetBirthDate() != nil {
		updateReq.BirthDate = req.GetBirthDate().AsTime()
	}
	updateReq.BirthPlace = req.GetBirthPlace()
	if req.GetGender() != pb.Gender_GENDER_UNSPECIFIED {
		updateReq.Gender = ProtoToGender(req.GetGender())
	}
	if req.GetAddress() != nil {
		addr := req.GetAddress()
		updateReq.Address = fmt.Sprintf("%s, RT/RW%s, %s, %s, %s, %s",
			addr.GetStreet(), addr.GetRtRw(), addr.GetVillage(),
			addr.GetDistrict(), addr.GetCity(), addr.GetPostalCode())
	}
	if req.GetContactInfo() != nil {
		updateReq.Phone = req.GetContactInfo().GetPhone()
		updateReq.Email = req.GetContactInfo().GetEmail()
	}
	if req.GetGuardianInfo() != nil {
		updateReq.ParentName = req.GetGuardianInfo().GetName()
		updateReq.ParentPhone = req.GetGuardianInfo().GetPhone()
		updateReq.ParentEmail = req.GetGuardianInfo().GetOccupation()
	}
	updateReq.NIS = req.GetNis()
	updateReq.NISN = req.GetNisn()
	updateReq.ClassID = req.GetClassName()

	result, err := h.service.UpdateStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.GetStudentId(), updateReq)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.UpdateStudentResponse{Student: mapToProtoStudent(result)}, nil
}

// DeleteStudent menghapus siswa (soft delete)
func (h *StudentGRPCHandler) DeleteStudent(ctx context.Context, req *pb.DeleteStudentRequest) (*pb.DeleteStudentResponse, error) {
	opCtx, err := extractOpCtxFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	err = h.service.DeleteStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.GetStudentId())
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.DeleteStudentResponse{
		Success: true,
		Message: "Siswa berhasil dihapus (soft delete)",
	}, nil
}

// ProtoToGender helper untuk konversi enum ke string
func ProtoToGender(g pb.Gender) string {
	switch g {
	case pb.Gender_GENDER_MALE:
		return "MALE"
	case pb.Gender_GENDER_FEMALE:
		return "FEMALE"
	default:
		return ""
	}
}
