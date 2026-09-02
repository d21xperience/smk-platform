#!/bin/bash
# scripts/test-registration-workflow.sh
# Workflow test untuk Domain Registration (PPDB)

HOST="localhost:50051"

COMMON_HEADERS=(
  -plaintext
  -H "x-school-id: SCHOOL_001"
  -H "x-academic-year: 2024/2025"
  -H "x-semester: 1"
  -H "x-academic-period: 2024-2025-1"
  -H "x-user-id: 1"
  -H "x-user-role: admin"
  -H "x-request-id: registration-test-$(date +%s)"
)

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}  REGISTRATION SERVICE WORKFLOW TEST${NC}"
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

# STEP 1: Create Registration (Draft)
print_step "1/6" "CREATE REGISTRATION (DRAFT)"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"calon_nisn\": \"0077777777\",
  \"calon_nama\": \"Dewi Lestari\",
  \"calon_jenis_kelamin\": \"FEMALE\",
  \"calon_tempat_lahir\": \"Bandung\",
  \"calon_tanggal_lahir\": \"2008-07-22T00:00:00Z\",
  \"calon_alamat\": \"Jl. Sudirman No. 75\",
  \"calon_telepon\": \"081277777777\",
  \"calon_email\": \"dewi@example.com\",
  \"nama_ayah\": \"Budi Lestari\",
  \"nama_ibu\": \"Siti Aminah\",
  \"telepon_ortu\": \"081288888888\",
  \"pekerjaan_ayah\": \"PNS\",
  \"pekerjaan_ibu\": \"Guru\",
  \"asal_sekolah\": \"SMP Negeri 1 Bandung\",
  \"jurusan_dipilih\": \"Rekayasa Perangkat Lunak\",
  \"alasan_memilih\": \"Tertarik dengan programming\"
}" $HOST registration.v1.RegistrationService/CreateRegistration 2>&1)

if echo "$result" | grep -q '"id"'; then
  print_success "Registration created"
  REGISTRATION_ID=$(echo "$result" | grep '"id"' | head -1 | sed 's/.*"id": "\([^"]*\)".*/\1/')
  echo -e "    🆔 Registration ID: ${BLUE}$REGISTRATION_ID${NC}"
else
  print_error "Failed to create registration"
  echo "$result"
  exit 1
fi

# STEP 2: Submit Registration
print_step "2/6" "SUBMIT REGISTRATION"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"registration_id\": \"$REGISTRATION_ID\"
}" $HOST registration.v1.RegistrationService/SubmitRegistration 2>&1)

if echo "$result" | grep -q '"status": "SUBMITTED"'; then
  print_success "Registration submitted"
else
  print_error "Failed to submit registration"
  echo "$result"
  exit 1
fi

# STEP 3: Verify Registration
print_step "3/6" "VERIFY REGISTRATION"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"registration_id\": \"$REGISTRATION_ID\",
  \"catatan\": \"Berkas lengkap, data valid\"
}" $HOST registration.v1.RegistrationService/VerifyRegistration 2>&1)

if echo "$result" | grep -q '"status": "VERIFIED"'; then
  print_success "Registration verified"
else
  print_error "Failed to verify registration"
  echo "$result"
  exit 1
fi

# STEP 4: Approve Registration (Side-effect: Create Student)
print_step "4/6" "APPROVE REGISTRATION"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"registration_id\": \"$REGISTRATION_ID\",
  \"catatan\": \"Diterima sebagai siswa baru\"
}" $HOST registration.v1.RegistrationService/ApproveRegistration 2>&1)

if echo "$result" | grep -q '"status": "APPROVED"'; then
  print_success "Registration approved"
else
  print_error "Failed to approve registration"
  echo "$result"
  exit 1
fi

# STEP 5: Verify Side-Effect (Student Created)
print_step "5/6" "VERIFY SIDE-EFFECT (Student should be created)"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"search\": \"Dewi Lestari\",
  \"limit\": 10,
  \"offset\": 0
}" $HOST student.v1.StudentService/ListStudents 2>&1)

if echo "$result" | grep -q '"firstName": "Dewi"'; then
  print_success "Student created from registration"
  echo "$result" | grep -E '"(studentId|firstName|nisn)"' | head -3 | sed 's/^/    /'
else
  print_error "Student not found (side-effect failed)"
  echo "$result" | head -10
  exit 1
fi

# STEP 6: List Registrations
print_step "6/6" "LIST REGISTRATIONS"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"status\": \"APPROVED\",
  \"search\": \"\",
  \"limit\": 10,
  \"offset\": 0
}" $HOST registration.v1.RegistrationService/ListRegistrations 2>&1)

if echo "$result" | grep -q "$REGISTRATION_ID"; then
  print_success "Registration found in list"
  echo "$result" | grep -E '"(id|calonNama|status)"' | head -3 | sed 's/^/    /'
else
  print_error "Registration not found in list"
  echo "$result" | head -10
fi

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${GREEN}  🎉 REGISTRATION WORKFLOW TEST COMPLETE${NC}"
echo -e "${BLUE}============================================${NC}"
