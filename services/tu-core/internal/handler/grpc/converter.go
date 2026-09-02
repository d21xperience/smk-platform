package grpc

import (
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	domainmodels "sekolah-platform/services/tu-core/internal/domain/student/models"
	pb "sekolah-platform/services/tu-core/proto/student/v1"
)

// StudentToProto mengkonversi domain Student ke proto Student
func StudentToProto(s *domainmodels.Student) *pb.Student {
	if s == nil {
		return nil
	}

	enrollments := make([]*pb.Enrollment, 0, len(s.Enrollments()))
	for _, e := range s.Enrollments() {
		enrollments = append(enrollments, &pb.Enrollment{
			EnrollmentId:   e.EnrollmentID(),
			SchoolId:       e.SchoolID(),
			PeriodId:       e.PeriodID(),
			ClassId:        e.ClassID(),
			EnrollmentDate: timestamppb.New(e.EnrollmentDate()),
			Status:         e.Status(),
		})
	}

	return &pb.Student{
		StudentId: s.StudentID(),
		SchoolId:  s.SchoolID(),
		Nisn:      s.NISN().Value(),
		Nis:       s.NIS().Value(),
		FullName: &pb.FullName{
			FirstName:  s.FullName().FirstName(),
			MiddleName: s.FullName().MiddleName(),
			LastName:   s.FullName().LastName(),
		},
		BirthDate: timestamppb.New(s.BirthDate()),
		Gender:    genderToProto(s.Gender().Code()),
		Address: &pb.Address{
			Street:     s.Address().Street,
			RtRw:       s.Address().RtRw,
			Village:    s.Address().Village,
			District:   s.Address().District,
			City:       s.Address().City,
			PostalCode: s.Address().PostalCode,
		},
		ContactInfo: &pb.ContactInfo{
			Phone: s.ContactInfo().Phone,
			Email: s.ContactInfo().Email,
		},
		GuardianInfo: &pb.GuardianInfo{
			Name:       s.GuardianInfo().Name,
			Relation:   s.GuardianInfo().Relation,
			Phone:      s.GuardianInfo().Phone,
			Occupation: s.GuardianInfo().Occupation,
		},
		Status:      statusToProto(s.Status().Code()),
		Enrollments: enrollments,
		CreatedAt:   timestamppb.New(s.CreatedAt()),
		UpdatedAt:   timestamppb.New(s.UpdatedAt()),
	}
}

// ProtoToGenderString mengkonversi proto Gender ke string code
func ProtoToGenderString(g pb.Gender) string {
	switch g {
	case pb.Gender_GENDER_MALE:
		return "MALE"
	case pb.Gender_GENDER_FEMALE:
		return "FEMALE"
	default:
		return ""
	}
}

// genderToProto mengkonversi gender code ke proto enum
func genderToProto(code string) pb.Gender {
	switch code {
	case "MALE":
		return pb.Gender_GENDER_MALE
	case "FEMALE":
		return pb.Gender_GENDER_FEMALE
	default:
		return pb.Gender_GENDER_UNSPECIFIED
	}
}

// statusToProto mengkonversi domain status code (English) ke proto enum (Indonesian)
func statusToProto(code string) pb.StudentStatus {
	switch strings.ToUpper(code) {
	case "ACTIVE":
		return pb.StudentStatus_STUDENT_STATUS_ACTIVE
	case "TRANSFERRED":
		return pb.StudentStatus_STUDENT_STATUS_TRANSFERRED
	case "GRADUATED":
		return pb.StudentStatus_STUDENT_STATUS_GRADUATED
	case "DROPPED", "DROPOUT":
		return pb.StudentStatus_STUDENT_STATUS_DROPPED
	case "ALUMNI":
		return pb.StudentStatus_STUDENT_STATUS_ALUMNI
	default:
		return pb.StudentStatus_STUDENT_STATUS_UNSPECIFIED
	}
}

// mapDomainStatusToProto adalah alias untuk statusToProto (untuk konsistensi naming di handler)
func mapDomainStatusToProto(domainStatus string) pb.StudentStatus {
	return statusToProto(domainStatus)
}

// mapFrontendStatusToDomain mengkonversi status Indonesian dari frontend ke domain English
func mapFrontendStatusToDomain(frontendStatus string) string {
	switch strings.ToUpper(frontendStatus) {
	case "AKTIF":
		return "ACTIVE"
	case "MUTASI_KELUAR":
		return "TRANSFERRED"
	case "LULUS":
		return "GRADUATED"
	case "DO":
		return "DROPPED"
	case "ALUMNI":
		return "ALUMNI"
	default:
		return ""
	}
}

// TimestampToTime mengkonversi proto Timestamp ke time.Time
func TimestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
