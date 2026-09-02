#!/bin/bash
# scripts/test-assessment-workflow.sh
# Workflow test untuk Domain Assessment

HOST="localhost:50051"

# Gunakan timestamp agar data unik, tapi pastikan panjangnya valid
TIMESTAMP=$(date +%s)

# ✅ PERBAIKAN:
# TIMESTAMP biasanya 10 digit (misal: 1725090000)
# NISN wajib maksimal 10 digit. Kita ambil 8 digit terakhir, lalu tambah "00" di depan.
TEST_NISN="00${TIMESTAMP: -8}"

# NIS wajib maksimal 20 digit. Kita ambil 6 digit terakhir, lalu tambah "2025" di depan.
TEST_NIS="2025${TIMESTAMP: -6}"

COMMON_HEADERS=(
  -plaintext
  -H "x-school-id: SCHOOL_001"
  -H "x-academic-year: 2024/2025"
  -H "x-semester: 1"
  -H "x-academic-period: 2024-2025-1"
  -H "x-user-id: 1"
  -H "x-user-role: admin"
  -H "x-request-id: assessment-test-${TIMESTAMP}"
)

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}  ASSESSMENT SERVICE WORKFLOW TEST${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""

print_step() {
  echo ""
  echo -e "${YELLOW}▶ [$1] $2${NC}"
}

print_success() {
  echo -e "${GREEN}  ✅ $1${NC}"
}

print_error() {
  echo -e "${RED}  ❌ $1${NC}"
}

# STEP 1: Create Student (Prerequisite)
print_step "1/6" "CREATE STUDENT (Prerequisite) - NIS: $TEST_NIS"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"nis\": \"$TEST_NIS\",
  \"nisn\": \"$TEST_NISN\",
  \"full_name\": {\"first_name\": \"Rina\", \"last_name\": \"Wulandari\"},
  \"birth_date\": \"2008-03-15T00:00:00Z\",
  \"gender\": \"GENDER_FEMALE\",
  \"address\": {\"street\": \"Jl. Sudirman No. 50\"},
  \"contact_info\": {\"phone\": \"081234567890\"},
  \"guardian_info\": {\"name\": \"Siti\", \"phone\": \"081111111111\"}
}" $HOST student.v1.StudentService/RegisterStudent 2>&1)

if echo "$result" | grep -q '"studentId"'; then
  print_success "Student created"
  STUDENT_ID=$(echo "$result" | grep '"studentId"' | head -1 | sed 's/.*"studentId": "\([^"]*\)".*/\1/')
  echo -e "    🆔 Student ID: ${BLUE}$STUDENT_ID${NC}"
else
  print_error "Failed to create student"
  echo "$result"
  exit 1
fi

# STEP 2: Create Assessment (DAILY)
print_step "2/6" "CREATE ASSESSMENT (DAILY - Matematika)"
ASSESSMENT_DATE=$(date -u +"%Y-%m-%dT00:00:00Z")

result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"student_id\": \"$STUDENT_ID\",
  \"subject_id\": \"sub-mtk\",
  \"assessment_type\": \"DAILY\",
  \"score\": 85,
  \"max_score\": 100,
  \"assessment_date\": \"$ASSESSMENT_DATE\",
  \"semester\": 1,
  \"notes\": \"Ulangan harian bab 1\"
}" $HOST assessment.v1.AssessmentService/CreateAssessment 2>&1)

if echo "$result" | grep -q '"id"'; then
  print_success "Assessment created"
  ASSESSMENT_ID=$(echo "$result" | grep '"id"' | head -1 | sed 's/.*"id": "\([^"]*\)".*/\1/')
  echo -e "    🆔 Assessment ID: ${BLUE}$ASSESSMENT_ID${NC}"
  echo "$result" | grep -E '"(score|percentage)"' | sed 's/^/    /'
else
  print_error "Failed to create assessment"
  echo "$result"
  exit 1
fi

# STEP 3: Create Another Assessment (MIDTERM)
print_step "3/6" "CREATE ASSESSMENT (MIDTERM - Matematika)"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"student_id\": \"$STUDENT_ID\",
  \"subject_id\": \"sub-mtk\",
  \"assessment_type\": \"MIDTERM\",
  \"score\": 78,
  \"max_score\": 100,
  \"assessment_date\": \"$ASSESSMENT_DATE\",
  \"semester\": 1,
  \"notes\": \"Ujian tengah semester\"
}" $HOST assessment.v1.AssessmentService/CreateAssessment 2>&1)

if echo "$result" | grep -q '"id"'; then
  print_success "Midterm assessment created"
else
  print_error "Failed to create midterm assessment"
  echo "$result"
  exit 1
fi

# STEP 4: Get Student Assessments
print_step "4/6" "GET STUDENT ASSESSMENTS (Semester 1)"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"student_id\": \"$STUDENT_ID\",
  \"semester\": 1,
  \"subject_id\": \"\"
}" $HOST assessment.v1.AssessmentService/GetStudentAssessments 2>&1)

# ✅ PERBAIKAN: Cari "assessmentType
# STEP 5: Get Student Average
print_step "5/6" "GET STUDENT AVERAGE (Semester 1)"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"student_id\": \"$STUDENT_ID\",
  \"semester\": 1
}" $HOST assessment.v1.AssessmentService/GetStudentAverage 2>&1)

if echo "$result" | grep -q '"overallAverage"'; then
  print_success "Student average retrieved"
  echo "$result" | grep -E '"(overallAverage|average|grade)"' | head -5 | sed 's/^/    /'
else
  print_error "Failed to get student average"
  echo "$result" | head -10
fi

# STEP 6: Update Assessment (Koreksi Nilai)
print_step "6/6" "UPDATE ASSESSMENT (Koreksi Nilai)"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"assessment_id\": \"$ASSESSMENT_ID\",
  \"score\": 90,
  \"max_score\": 100,
  \"notes\": \"Ulangan harian bab 1 (koreksi)\"
}" $HOST assessment.v1.AssessmentService/UpdateAssessment 2>&1)

if echo "$result" | grep -q '"score": 90'; then
  print_success "Assessment updated successfully"
  echo "$result" | grep -E '"(score|percentage)"' | sed 's/^/    /'
else
  print_error "Failed to update assessment"
  echo "$result"
fi

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${GREEN}  🎉 ASSESSMENT WORKFLOW TEST COMPLETE${NC}"
echo -e "${BLUE}============================================${NC}"
