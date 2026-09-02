#!/bin/bash
# scripts/test-correspondence-workflow.sh

HOST="localhost:50051"
TIMESTAMP=$(date +%s)

COMMON_HEADERS=(
  -plaintext
  -H "x-school-id: SCHOOL_001"
  -H "x-academic-year: 2024/2025"
  -H "x-semester: 1"
  -H "x-academic-period: 2024-2025-1"
  -H "x-user-id: 1"
  -H "x-user-role: admin"
  -H "x-request-id: corr-test-${TIMESTAMP}"
)

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}  CORRESPONDENCE SERVICE WORKFLOW TEST${NC}"
echo -e "${BLUE}============================================${NC}"

print_step() { echo -e "\n${YELLOW}▶ [$1] $2${NC}"; }
print_success() { echo -e "${GREEN}  ✅ $1${NC}"; }
print_error() { echo -e "${RED}  ❌ $1${NC}"; }

# STEP 1: Create Outgoing Correspondence
print_step "1/5" "CREATE OUTGOING CORRESPONDENCE"
DATE_STR=$(date -u +"%Y-%m-%dT00:00:00Z")

result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"type\": \"OUTGOING\",
  \"number\": \"001/SMK/X/2024\",
  \"date\": \"$DATE_STR\",
  \"subject\": \"Undangan Rapat Koordinasi\",
  \"from\": \"Kepala SMK Negeri 1\",
  \"to\": \"Dinas Pendidikan Kota\",
  \"description\": \"Undangan untuk menghadiri rapat koordinasi tahunan\",
  \"attachment_url\": \"https://example.com/undangan.pdf\"
}" $HOST correspondence.v1.CorrespondenceService/CreateCorrespondence 2>&1)

if echo "$result" | grep -q '"id"'; then
  print_success "Correspondence created"
  CORR_ID=$(echo "$result" | grep '"id"' | head -1 | sed 's/.*"id": "\([^"]*\)".*/\1/')
  echo -e "    🆔 ID: ${BLUE}$CORR_ID${NC}"
else
  print_error "Failed to create correspondence"
  echo "$result"
  exit 1
fi

# STEP 2: Update Correspondence
print_step "2/5" "UPDATE CORRESPONDENCE"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"id\": \"$CORR_ID\",
  \"subject\": \"Undangan Rapat Koordinasi (Revisi)\",
  \"description\": \"Revisi waktu dan tempat rapat\",
  \"attachment_url\": \"https://example.com/undangan_revisi.pdf\"
}" $HOST correspondence.v1.CorrespondenceService/UpdateCorrespondence 2>&1)

if echo "$result" | grep -q '"subject": "Undangan Rapat Koordinasi (Revisi)"'; then
  print_success "Correspondence updated"
else
  print_error "Failed to update correspondence"
  echo "$result"
  exit 1
fi

# STEP 3: Archive Correspondence
print_step "3/5" "ARCHIVE CORRESPONDENCE"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"id\": \"$CORR_ID\"
}" $HOST correspondence.v1.CorrespondenceService/ArchiveCorrespondence 2>&1)

if echo "$result" | grep -q '"status": "ARCHIVED"'; then
  print_success "Correspondence archived"
else
  print_error "Failed to archive correspondence"
  echo "$result"
  exit 1
fi

# STEP 4: Get Correspondence
print_step "4/5" "GET CORRESPONDENCE"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"id\": \"$CORR_ID\"
}" $HOST correspondence.v1.CorrespondenceService/GetCorrespondence 2>&1)

if echo "$result" | grep -q '"status": "ARCHIVED"'; then
  print_success "Correspondence retrieved"
  echo "$result" | grep -E '"(type|number|status)"' | head -3 | sed 's/^/    /'
else
  print_error "Failed to get correspondence"
  echo "$result" | head -10
fi

# STEP 5: List Correspondences
print_step "5/5" "LIST CORRESPONDENCES"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"type\": \"OUTGOING\",
  \"status\": \"ARCHIVED\",
  \"limit\": 10,
  \"offset\": 0
}" $HOST correspondence.v1.CorrespondenceService/ListCorrespondences 2>&1)

if echo "$result" | grep -q "$CORR_ID"; then
  print_success "Correspondence found in list"
else
  print_error "Correspondence not found in list"
  echo "$result" | head -10
fi

echo -e "\n${BLUE}============================================${NC}"
echo -e "${GREEN}  🎉 CORRESPONDENCE WORKFLOW TEST COMPLETE${NC}"
echo -e "${BLUE}============================================${NC}"
