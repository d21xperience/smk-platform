#!/bin/bash
# scripts/test-attendance-workflow.sh
# Workflow test untuk Domain Attendance

HOST="localhost:50051"

COMMON_HEADERS=(
  -plaintext
  -H "x-school-id: SCHOOL_001"
  -H "x-academic-year: 2024/2025"
  -H "x-semester: 1"
  -H "x-academic-period: 2024-2025-1"
  -H "x-user-id: 1"
  -H "x-user-role: admin"
  -H "x-request-id: attendance-test-$(date +%s)"
)

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}  ATTENDANCE SERVICE WORKFLOW TEST${NC}"
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
print_step "1/5" "CREATE STUDENT (Prerequisite)"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"nis\": \"2024777\",
  \"nisn\": \"0077777777\",
  \"full_name\": {\"first_name\": \"Andi\", \"last_name\": \"Pratama\"},
  \"birth_date\": \"2008-05-10T00:00:00Z\",
  \"gender\": \"GENDER_MALE\",
  \"address\": {\"street\": \"Jl. Merdeka No. 1\"},
  \"contact_info\": {\"phone\": \"081234567890\"},
  \"guardian_info\": {\"name\": \"Budi\", \"phone\": \"081111111111\"}
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

# STEP 2: Create Attendance Session
print_step "2/5" "CREATE ATTENDANCE SESSION"
SESSION_DATE=$(date -u +"%Y-%m-%dT00:00:00Z")
START_TIME=$(date -u +"%Y-%m-%dT07:00:00Z")
END_TIME=$(date -u +"%Y-%m-%dT08:00:00Z")

result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"class_id\": \"cls-x-rpl-1\",
  \"subject_id\": \"sub-mtk\",
  \"teacher_id\": \"teacher-001\",
  \"session_date\": \"$SESSION_DATE\",
  \"start_time\": \"$START_TIME\",
  \"end_time\": \"$END_TIME\",
  \"notes\": \"Sesi pagi\"
}" $HOST attendance.v1.AttendanceService/CreateSession 2>&1)

if echo "$result" | grep -q '"id"'; then
  print_success "Session created"
  SESSION_ID=$(echo "$result" | grep '"id"' | head -1 | sed 's/.*"id": "\([^"]*\)".*/\1/')
  echo -e "    🆔 Session ID: ${BLUE}$SESSION_ID${NC}"
else
  print_error "Failed to create session"
  echo "$result"
  exit 1
fi

# STEP 3: Submit Attendance
print_step "3/5" "SUBMIT ATTENDANCE"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"session_id\": \"$SESSION_ID\",
  \"records\": [
    {
      \"student_id\": \"$STUDENT_ID\",
      \"status\": \"PRESENT\",
      \"note\": \"Hadir tepat waktu\"
    }
  ]
}" $HOST attendance.v1.AttendanceService/SubmitAttendance 2>&1)

if echo "$result" | grep -q '"status": "SUBMITTED"'; then
  print_success "Attendance submitted"
else
  print_error "Failed to submit attendance"
  echo "$result"
  exit 1
fi

# STEP 4: Get Student Attendance
print_step "4/5" "GET STUDENT ATTENDANCE"
# ✅ PERBAIKAN: Gunakan %-m untuk menghilangkan leading zero (misal: 08 menjadi 8) agar valid di JSON
CURRENT_MONTH=$(date +%-m)
CURRENT_YEAR=$(date +%Y)

result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"student_id\": \"$STUDENT_ID\",
  \"month\": $CURRENT_MONTH,
  \"year\": $CURRENT_YEAR
}" $HOST attendance.v1.AttendanceService/GetStudentAttendance 2>&1)

if echo "$result" | grep -q '"presentCount": 1'; then
  print_success "Student attendance retrieved"
  echo "$result" | grep -E '"(presentCount|absentCount|attendancePercentage)"' | sed 's/^/    /'
else
  print_error "Failed to get student attendance"
  echo "$result" | head -10
fi

# STEP 5: List Sessions
print_step "5/5" "LIST SESSIONS"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"class_id\": \"cls-x-rpl-1\",
  \"status\": \"SUBMITTED\",
  \"limit\": 10,
  \"offset\": 0
}" $HOST attendance.v1.AttendanceService/ListSessions 2>&1)

if echo "$result" | grep -q "$SESSION_ID"; then
  print_success "Session found in list"
  echo "$result" | grep -E '"(id|classId|status)"' | head -3 | sed 's/^/    /'
else
  print_error "Session not found in list"
  echo "$result" | head -10
fi

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${GREEN}  🎉 ATTENDANCE WORKFLOW TEST COMPLETE${NC}"
echo -e "${BLUE}============================================${NC}"
