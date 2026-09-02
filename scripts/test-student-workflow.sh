# STEP 1: Create Student (Prerequisite)
print_step "1/6" "CREATE STUDENT (Prerequisite) - NIS: $TEST_NIS"
result=$(grpcurl "${COMMON_HEADERS[@]}" -d "{
  \"nis\": \"$TEST_NIS\",
  \"nisn\": \"$TEST_NISN\",
  \"full_name\": {\"first_name\": \"Rina\", \"middle_name\": \"\", \"last_name\": \"Wulandari\"},
  \"birth_date\": \"2008-03-15T00:00:00Z\",
  \"gender\": \"GENDER_FEMALE\",
  \"address\": {\"street\": \"Jl. Sudirman No. 50\", \"rt_rw\": \"01/02\", \"village\": \"Menteng\", \"district\": \"Menteng\", \"city\": \"Jakarta Pusat\", \"postal_code\": \"10310\"},
  \"contact_info\": {\"phone\": \"081234567890\", \"email\": \"rina@example.com\"},
  \"guardian_info\": {\"name\": \"Siti Aminah\", \"relation\": \"mother\", \"phone\": \"081111111111\", \"occupation\": \"Ibu Rumah Tangga\"}
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
