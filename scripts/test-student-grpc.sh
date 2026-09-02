#!/bin/bash
# scripts/test-student-grpc-simple.sh

HOST="localhost:50051"
STUDENT_ID="550e8400-e29b-41d4-a716-446655440000"

case "$1" in
  register)
    echo "📝 Registering student..."
    grpcurl -plaintext \
      -H "x-school-id: SCHOOL_001" \
      -H "x-academic-year: 2024/2025" \
      -H "x-semester: 1" \
      -H "x-academic-period: 2024-2025-1" \
      -H "x-user-id: 1" \
      -H "x-user-role: admin" \
      -H "x-request-id: req-$(date +%s)" \
      -d "{
        \"student_id\": \"$STUDENT_ID\",
        \"nis\": \"2024001\",
        \"nisn\": \"0012345678\",
        \"full_name\": {\"first_name\": \"Ahmad\", \"middle_name\": \"\", \"last_name\": \"Fauzi\"},
        \"birth_date\": \"2007-05-15T00:00:00Z\",
        \"gender\": \"GENDER_MALE\",
        \"address\": {\"street\": \"Jl. Merdeka No. 10\", \"rt_rw\": \"01/02\", \"village\": \"Sukajadi\", \"district\": \"Coblong\", \"city\": \"Bandung\", \"postal_code\": \"40132\"},
        \"contact_info\": {\"phone\": \"081234567890\", \"email\": \"ahmad@example.com\"},
        \"guardian_info\": {\"name\": \"Budi Hartono\", \"relation\": \"father\", \"phone\": \"08111222333\", \"occupation\": \"Wiraswasta\"}
      }" \
      $HOST student.v1.StudentService/RegisterStudent
    ;;

  get)
    echo "🔍 Getting student..."
    grpcurl -plaintext \
      -H "x-school-id: SCHOOL_001" \
      -H "x-academic-year: 2024/2025" \
      -H "x-semester: 1" \
      -H "x-academic-period: 2024-2025-1" \
      -H "x-user-id: 1" \
      -H "x-user-role: admin" \
      -d "{\"student_id\": \"$STUDENT_ID\"}" \
      $HOST student.v1.StudentService/GetStudent
    ;;

  list)
    echo "📋 Listing students..."
    grpcurl -plaintext \
      -H "x-school-id: SCHOOL_001" \
      -H "x-academic-year: 2024/2025" \
      -H "x-semester: 1" \
      -H "x-academic-period: 2024-2025-1" \
      -H "x-user-id: 1" \
      -H "x-user-role: admin" \
      -d '{"search": "", "status": "", "limit": 10, "offset": 0}' \
      $HOST student.v1.StudentService/ListStudents
    ;;

  *)
    echo "Usage: $0 {register|get|list}"
    exit 1
    ;;
esac
