#!/bin/bash
# scripts/test-mutation-workflow.sh
# Workflow test untuk Domain Mutation

HOST="localhost:50051"

COMMON_HEADERS=(
  -plaintext
  -H "x-school-id: SCHOOL_001"
  -H "x-academic-year: 2024/2025"
  -H "x-semester: 1"
  -H "x-academic-period: 2024-2025-1"
  -H "x-user-id: 1"
  -H "x-user-role: admin"
  -H "x-request-id: mutation-test-$(date +%s)"
)

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}  MUTATION SERVICE WORKFLOW TEST${NC}"
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

# STEP 1: Buat Siswa Baru (Prasyarat untuk Mutasi Keluar)
print_step "1/4" "CREATE STUDENT (Prerequisite)"
STUDENT_NIS="2024888"
STUDENT_NISN="0088888888"

result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"nis\": \"$STUDENT_NIS\",
  \"nisn\": \"$STUDENT_NISN\",
  \"full_name\": {\"first_name\": \"Siti\", \"middle_name\": \"\", \"last_name\": \"Nurhaliza\"},
  \"birth_date\": \"2008-08-15T00:00:00Z\",
  \"gender\": \"GENDER_FEMALE\",
  \"address\": {\"street\": \"Jl. Merdeka No. 50\", \"rt_rw\": \"02/03\", \"village\": \"Cicendo\", \"district\": \"Bandung\", \"city\": \"Bandung\", \"postal_code\": \"40171\"},
  \"contact_info\": {\"phone\": \"081299999999\", \"email\": \"siti@example.com\"},
  \"guardian_info\": {\"name\": \"Bapak Siti\", \"relation\": \"father\", \"phone\": \"081288888888\", \"occupation\": \"PNS\"}
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

# STEP 2: Ajukan Mutasi Keluar
print_step "2/4" "CREATE MUTATION (KELUAR)"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"student_id\": \"$STUDENT_ID\",
  \"mutation_type\": \"KELUAR\",
  \"tujuan_sekolah\": \"SMK Negeri 1 Jakarta\",
  \"alasan\": \"Mengikuti orang tua pindah tugas\"
}" $HOST mutation.v1.MutationService/CreateMutation 2>&1)

if echo "$result" | grep -q '"id"'; then
  print_success "Mutation created"
  MUTATION_ID=$(echo "$result" | grep '"id"' | head -1 | sed 's/.*"id": "\([^"]*\)".*/\1/')
  echo -e "    🆔 Mutation ID: ${BLUE}$MUTATION_ID${NC}"
else
  print_error "Failed to create mutation"
  echo "$result"
  exit 1
fi

# STEP 3: Approve Mutation
print_step "3/4" "APPROVE MUTATION"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"mutation_id\": \"$MUTATION_ID\",
  \"catatan\": \"Berkas lengkap, disetujui\"
}" $HOST mutation.v1.MutationService/ApproveMutation 2>&1)

if echo "$result" | grep -q '"status": "DISETUJUI"'; then
  print_success "Mutation approved successfully"
else
  print_error "Failed to approve mutation"
  echo "$result"
  exit 1
fi

# STEP 4: Verifikasi Side-Effect (Status Siswa Berubah)
print_step "4/4" "VERIFY SIDE-EFFECT (Student Status should be TRANSFERRED)"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{\"student_id\": \"$STUDENT_ID\"}" \
  $HOST student.v1.StudentService/GetStudent 2>&1)

if echo "$result" | grep -q '"status": "STUDENT_STATUS_TRANSFERRED"'; then
  print_success "Side-effect verified! Student status changed to TRANSFERRED"
  echo "$result" | grep '"status"' | sed 's/^/    /'
else
  print_error "Side-effect failed! Student status did not change"
  echo "$result" | grep '"status"' | sed 's/^/    /'
  exit 1
fi

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${GREEN}  🎉 MUTATION WORKFLOW TEST COMPLETE${NC}"
echo -e "${BLUE}============================================${NC}"
