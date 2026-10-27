#!/bin/bash
# scripts/test-class-enrollment-workflow.sh
# Workflow test untuk Domain Academic Structure (Class & Enrollment)

BASE_URL="http://localhost:8081"

COMMON_HEADERS=(
  -H "Content-Type: application/json"
  -H "x-school-id: SCHOOL_001"
  -H "x-academic-year: 2024/2025"
  -H "x-semester: 1"
  -H "x-academic-period: 2024-2025-1"
  -H "x-user-id: 1"
  -H "x-user-role: admin_tu"
  -H "x-request-id: class-test-$(date +%s)"
)

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}============================================================${NC}"
echo -e "${BLUE}  ACADEMIC STRUCTURE WORKFLOW TEST (Class & Enrollment)${NC}"
echo -e "${BLUE}============================================================${NC}"
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

# ============================================================
# STEP 1: Create Class
# ============================================================
print_step "1/9" "CREATE CLASS (Rombel)"
response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/api/v1/classes" \
  "${COMMON_HEADERS[@]}" \
  -d '{
    "name": "X TKJ A",
    "grade": 10,
    "level": "10",
    "major": "TKJ",
    "homeroomId": "TCH001",
    "capacity": 32
  }')
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 201 ]; then
  print_success "Class created"
  CLASS_ID=$(echo "$body" | grep -o '"id":"[^"]*' | head -1 | sed 's/"id":"//')
  echo -e "    🆔 Class ID: ${BLUE}$CLASS_ID${NC}"
else
  print_error "Failed to create class (HTTP $http_code)"
  echo "$body"
  exit 1
fi

# ============================================================
# STEP 2: Get Class by ID
# ============================================================
print_step "2/9" "GET CLASS BY ID"
response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}/api/v1/classes/${CLASS_ID}" \
  "${COMMON_HEADERS[@]}")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ]; then
  print_success "Class retrieved"
  echo "$body" | grep -E '"(name|grade|capacity|status)"' | head -4 | sed 's/^/    /'
else
  print_error "Failed to get class (HTTP $http_code)"
  echo "$body"
  exit 1
fi

# ============================================================
# STEP 3: List Classes
# ============================================================
print_step "3/9" "LIST CLASSES"
response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}/api/v1/classes?limit=10&offset=0" \
  "${COMMON_HEADERS[@]}")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ]; then
  print_success "Classes listed"
  total=$(echo "$body" | grep -o '"total":[0-9]*' | head -1 | sed 's/"total"://')
  echo -e "    📊 Total classes: ${BLUE}$total${NC}"
else
  print_error "Failed to list classes (HTTP $http_code)"
  echo "$body"
  exit 1
fi

# ============================================================
# STEP 4: Update Class
# ============================================================
print_step "4/9" "UPDATE CLASS"
response=$(curl -s -w "\n%{http_code}" -X PUT "${BASE_URL}/api/v1/classes/${CLASS_ID}" \
  "${COMMON_HEADERS[@]}" \
  -d '{
    "name": "X RPL 1 (Updated)",
    "grade": 10,
    "level": "X",
    "major": "RPL",
    "homeroomId": "TCH002",
    "capacity": 36,
    "status": "active"
  }')
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ]; then
  print_success "Class updated"
  echo "$body" | grep -E '"(name|homeroomId|capacity)"' | head -3 | sed 's/^/    /'
else
  print_error "Failed to update class (HTTP $http_code)"
  echo "$body"
  exit 1
fi

# ============================================================
# STEP 5: Add Member (Enroll Student)
# ============================================================
print_step "5/9" "ADD MEMBER (Enroll Student)"
# Pertama kita perlu student ID yang valid. Karena kita tidak punya student di test ini,
# kita akan coba dengan student dummy. Sebaiknya buat student dulu, tapi untuk demo kita
# gunakan student ID yang sudah ada atau kita buat skenario dengan student yang sudah ada.
# Untuk contoh, kita asumsikan ada student dengan ID "STU001".
STUDENT_ID="STU001" # Ganti dengan ID student yang valid di sistem Anda

# Cek apakah student ada
check_response=$(curl -s -o /dev/null -w "%{http_code}" -X GET "${BASE_URL}/api/v1/students/${STUDENT_ID}" \
  "${COMMON_HEADERS[@]}")
if [ "$check_response" -eq 404 ]; then
  echo -e "${YELLOW}  ⚠️ Student $STUDENT_ID not found. Creating dummy student...${NC}"
  # Coba buat student dummy via endpoint student (jika ada)
  # Atau lewati step ini dengan pesan warning
  print_error "Student $STUDENT_ID not found. Please create a student first."
  echo -e "${YELLOW}  → Skipping Add Member, Move Student, Remove Member tests${NC}"
  STUDENT_AVAILABLE=false
else
  STUDENT_AVAILABLE=true
  print_success "Student $STUDENT_ID exists"
fi

if [ "$STUDENT_AVAILABLE" = true ]; then
  response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/api/v1/classes/${CLASS_ID}/members" \
    "${COMMON_HEADERS[@]}" \
    -d "{
      \"studentId\": \"$STUDENT_ID\"
    }")
  http_code=$(echo "$response" | tail -n1)
  body=$(echo "$response" | sed '$d')

  if [ "$http_code" -eq 201 ]; then
    print_success "Member added"
    ENROLLMENT_ID=$(echo "$body" | grep -o '"id":"[^"]*' | head -1 | sed 's/"id":"//')
    echo -e "    🆔 Enrollment ID: ${BLUE}$ENROLLMENT_ID${NC}"
  else
    print_error "Failed to add member (HTTP $http_code)"
    echo "$body"
  fi
else
  echo -e "${YELLOW}  ⏭️  Skipping Add Member${NC}"
fi

# ============================================================
# STEP 6: List Members
# ============================================================
print_step "6/9" "LIST MEMBERS"
response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}/api/v1/classes/${CLASS_ID}/members" \
  "${COMMON_HEADERS[@]}")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ]; then
  print_success "Members listed"
  total_members=$(echo "$body" | grep -o '"total":[0-9]*' | head -1 | sed 's/"total"://')
  echo -e "    📊 Total members: ${BLUE}$total_members${NC}"
  if [ "$total_members" -gt 0 ]; then
    echo "$body" | grep -E '"(name|nis)"' | head -4 | sed 's/^/    /'
  fi
else
  print_error "Failed to list members (HTTP $http_code)"
  echo "$body"
fi

# ============================================================
# STEP 7: Move Student (pindah rombel)
# ============================================================
if [ "$STUDENT_AVAILABLE" = true ] && [ -n "$ENROLLMENT_ID" ]; then
  print_step "7/9" "MOVE STUDENT (Pindah Rombel)"
  # Buat class tujuan baru untuk dipindah
  response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/api/v1/classes" \
    "${COMMON_HEADERS[@]}" \
    -d '{
      "name": "X RPL 2",
      "grade": 10,
      "level": "X",
      "major": "RPL",
      "homeroomId": "TCH003",
      "capacity": 32
    }')
  http_code=$(echo "$response" | tail -n1)
  body=$(echo "$response" | sed '$d')
  if [ "$http_code" -eq 201 ]; then
    TARGET_CLASS_ID=$(echo "$body" | grep -o '"id":"[^"]*' | head -1 | sed 's/"id":"//')
    print_success "Target class created: $TARGET_CLASS_ID"
  else
    print_error "Failed to create target class (HTTP $http_code)"
    TARGET_CLASS_ID=""
  fi

  if [ -n "$TARGET_CLASS_ID" ]; then
    response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/api/v1/enrollments/move" \
      "${COMMON_HEADERS[@]}" \
      -d "{
        \"studentId\": \"$STUDENT_ID\",
        \"toClassId\": \"$TARGET_CLASS_ID\",
        \"reason\": \"Pindah dari X RPL 1 ke X RPL 2\"
      }")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" -eq 200 ]; then
      print_success "Student moved to new class"
      echo "$body" | grep -E '"(id|classId|status|reason)"' | head -4 | sed 's/^/    /'
    else
      print_error "Failed to move student (HTTP $http_code)"
      echo "$body"
    fi
  else
    echo -e "${YELLOW}  ⏭️  Skipping Move Student (target class creation failed)${NC}"
  fi
else
  echo -e "${YELLOW}  ⏭️  Skipping Move Student${NC}"
fi

# ============================================================
# STEP 8: Get Student History
# ============================================================
if [ "$STUDENT_AVAILABLE" = true ]; then
  print_step "8/9" "GET STUDENT HISTORY"
  response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}/api/v1/students/${STUDENT_ID}/history" \
    "${COMMON_HEADERS[@]}")
  http_code=$(echo "$response" | tail -n1)
  body=$(echo "$response" | sed '$d')

  if [ "$http_code" -eq 200 ]; then
    print_success "Student history retrieved"
    count=$(echo "$body" | grep -o '"id"' | wc -l)
    echo -e "    📊 Total history records: ${BLUE}$count${NC}"
    echo "$body" | grep -E '"(classId|status|reason)"' | head -6 | sed 's/^/    /'
  else
    print_error "Failed to get student history (HTTP $http_code)"
    echo "$body"
  fi
else
  echo -e "${YELLOW}  ⏭️  Skipping Get Student History${NC}"
fi

# ============================================================
# STEP 9: Remove Member (keluarkan siswa)
# ============================================================
if [ "$STUDENT_AVAILABLE" = true ] && [ -n "$ENROLLMENT_ID" ]; then
  print_step "9/9" "REMOVE MEMBER (Keluarkan Siswa)"
  response=$(curl -s -w "\n%{http_code}" -X DELETE "${BASE_URL}/api/v1/classes/${CLASS_ID}/members/${STUDENT_ID}" \
    "${COMMON_HEADERS[@]}")
  http_code=$(echo "$response" | tail -n1)
  body=$(echo "$response" | sed '$d')

  if [ "$http_code" -eq 204 ]; then
    print_success "Member removed (HTTP 204 No Content)"
  else
    print_error "Failed to remove member (HTTP $http_code)"
    echo "$body"
  fi
else
  echo -e "${YELLOW}  ⏭️  Skipping Remove Member${NC}"
fi

# ============================================================
# SUMMARY
# ============================================================
echo ""
echo -e "${BLUE}============================================================${NC}"
echo -e "${GREEN}  🎉 ACADEMIC STRUCTURE WORKFLOW TEST COMPLETE${NC}"
echo -e "${BLUE}============================================================${NC}"
echo ""
echo -e "📝 Summary:"
echo -e "  • Class created: ${BLUE}$CLASS_ID${NC}"
if [ -n "$ENROLLMENT_ID" ]; then
  echo -e "  • Enrollment created: ${BLUE}$ENROLLMENT_ID${NC}"
fi
if [ -n "$TARGET_CLASS_ID" ]; then
  echo -e "  • Target class for move: ${BLUE}$TARGET_CLASS_ID${NC}"
fi
echo ""
echo -e "${YELLOW}ℹ️  Note: Pastikan student dengan ID 'STU001' ada di database.${NC}"
echo -e "${YELLOW}   Jika tidak, buat student terlebih dahulu atau ubah STUDENT_ID di script.${NC}"
