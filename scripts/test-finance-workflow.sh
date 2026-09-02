#!/bin/bash
# scripts/test-finance-workflow.sh
# Workflow test untuk Domain Finance

HOST="localhost:50051"

# Gunakan timestamp agar data unik
TIMESTAMP=$(date +%s)
TEST_NIS="2026${TIMESTAMP: -6}"
TEST_NISN="00${TIMESTAMP: -8}"

COMMON_HEADERS=(
  -plaintext
  -H "x-school-id: SCHOOL_001"
  -H "x-academic-year: 2024/2025"
  -H "x-semester: 1"
  -H "x-academic-period: 2024-2025-1"
  -H "x-user-id: 1"
  -H "x-user-role: admin"
  -H "x-request-id: finance-test-${TIMESTAMP}"
)

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}  FINANCE SERVICE WORKFLOW TEST${NC}"
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
print_step "1/7" "CREATE STUDENT (Prerequisite) - NIS: $TEST_NIS"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"nis\": \"$TEST_NIS\",
  \"nisn\": \"$TEST_NISN\",
  \"full_name\": {\"first_name\": \"Dewi\", \"middle_name\": \"\", \"last_name\": \"Sartika\"},
  \"birth_date\": \"2008-06-20T00:00:00Z\",
  \"gender\": \"GENDER_FEMALE\",
  \"address\": {\"street\": \"Jl. Gatot Subroto No. 25\", \"rt_rw\": \"02/05\", \"village\": \"Cicadas\", \"district\": \"Cidadap\", \"city\": \"Bandung\", \"postal_code\": \"40132\"},
  \"contact_info\": {\"phone\": \"081234567890\", \"email\": \"dewi@example.com\"},
  \"guardian_info\": {\"name\": \"Hendra\", \"relation\": \"father\", \"phone\": \"081111111111\", \"occupation\": \"PNS\"}
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

# STEP 2: Create Invoice (SPP Januari)
print_step "2/7" "CREATE INVOICE (SPP Januari 2025)"
DUE_DATE=$(date -u +"%Y-01-15T00:00:00Z")

result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"student_id\": \"$STUDENT_ID\",
  \"component_type\": \"SPP\",
  \"month\": 1,
  \"year\": 2025,
  \"amount\": 250000,
  \"due_date\": \"$DUE_DATE\",
  \"notes\": \"SPP bulan Januari\"
}" $HOST finance.v1.FinanceService/CreateInvoice 2>&1)

if echo "$result" | grep -q '"id"'; then
  print_success "Invoice created"
  INVOICE_ID=$(echo "$result" | grep '"id"' | head -1 | sed 's/.*"id": "\([^"]*\)".*/\1/')
  echo -e "    🆔 Invoice ID: ${BLUE}$INVOICE_ID${NC}"
  echo "$result" | grep -E '"(amount|status)"' | head -2 | sed 's/^/    /'
else
  print_error "Failed to create invoice"
  echo "$result"
  exit 1
fi

# STEP 3: Create Invoice (SPP Februari)
print_step "3/7" "CREATE INVOICE (SPP Februari 2025)"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"student_id\": \"$STUDENT_ID\",
  \"component_type\": \"SPP\",
  \"month\": 2,
  \"year\": 2025,
  \"amount\": 250000,
  \"due_date\": \"$DUE_DATE\",
  \"notes\": \"SPP bulan Februari\"
}" $HOST finance.v1.FinanceService/CreateInvoice 2>&1)

if echo "$result" | grep -q '"id"'; then
  print_success "Second invoice created"
  INVOICE_ID_2=$(echo "$result" | grep '"id"' | head -1 | sed 's/.*"id": "\([^"]*\)".*/\1/')
else
  print_error "Failed to create second invoice"
  echo "$result"
  exit 1
fi

# STEP 4: Create Payment (Partial - Rp 150.000)
print_step "4/7" "CREATE PAYMENT (Partial - Rp 150.000 via TRANSFER)"
PAYMENT_DATE=$(date -u +"%Y-%m-%dT10:00:00Z")

result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"invoice_id\": \"$INVOICE_ID\",
  \"amount\": 150000,
  \"payment_method\": \"TRANSFER\",
  \"payment_date\": \"$PAYMENT_DATE\",
  \"reference_number\": \"TRF-$(date +%s)\",
  \"notes\": \"Pembayaran sebagian via BCA\"
}" $HOST finance.v1.FinanceService/CreatePayment 2>&1)

if echo "$result" | grep -q '"status": "PARTIAL"'; then
  print_success "Partial payment created - Status: PARTIAL"
  echo "$result" | grep -E '"(paidAmount|status)"' | head -2 | sed 's/^/    /'
else
  print_error "Failed to create partial payment"
  echo "$result"
  exit 1
fi

# STEP 5: Create Payment (Remaining - Rp 100.000)
print_step "5/7" "CREATE PAYMENT (Remaining - Rp 100.000 via CASH)"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"invoice_id\": \"$INVOICE_ID\",
  \"amount\": 100000,
  \"payment_method\": \"CASH\",
  \"payment_date\": \"$PAYMENT_DATE\",
  \"reference_number\": \"CASH-$(date +%s)\",
  \"notes\": \"Pelunasan via kas\"
}" $HOST finance.v1.FinanceService/CreatePayment 2>&1)

if echo "$result" | grep -q '"status": "PAID"'; then
  print_success "Full payment - Status: PAID"
  echo "$result" | grep -E '"(paidAmount|status|paidAt)"' | head -3 | sed 's/^/    /'
else
  print_error "Failed to complete payment"
  echo "$result"
  exit 1
fi

# STEP 6: Get Student Outstanding
print_step "6/7" "GET STUDENT OUTSTANDING"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"student_id\": \"$STUDENT_ID\"
}" $HOST finance.v1.FinanceService/GetStudentOutstanding 2>&1)

if echo "$result" | grep -q '"totalOutstanding": 250000'; then
  print_success "Student outstanding retrieved"
  echo "$result" | grep -E '"(totalOutstanding|unpaidInvoices)"' | head -2 | sed 's/^/    /'
else
  print_error "Failed to get student outstanding"
  echo "$result" | head -10
fi

# STEP 7: Get Payment Summary
print_step "7/7" "GET PAYMENT SUMMARY (Januari 2025)"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"month\": 1,
  \"year\": 2025,
  \"component_type\": \"SPP\"
}" $HOST finance.v1.FinanceService/GetPaymentSummary 2>&1)

if echo "$result" | grep -q '"totalPaid"'; then
  print_success "Payment summary retrieved"
  echo "$result" | grep -E '"(totalInvoices|paidInvoices|totalPaid|paymentPercentage)"' | sed 's/^/    /'
else
  print_error "Failed to get payment summary"
  echo "$result" | head -10
fi

echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${GREEN}  🎉 FINANCE WORKFLOW TEST COMPLETE${NC}"
echo -e "${BLUE}============================================${NC}"
