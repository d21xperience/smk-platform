#!/bin/bash

BASE_URL="http://localhost:8081/api/v1"
SCHOOL_ID="SCHOOL_001"
ACADEMIC_PERIOD_ID="2024-2025-1"
USER_ID="1"
USER_ROLE="admin"

# Test create student
curl -X POST "$BASE_URL/students" \
  -H "Content-Type: application/json" \
  -H "x-school-id: $SCHOOL_ID" \
  -H "x-academic-period-id: $ACADEMIC_PERIOD_ID" \
  -H "x-user-id: $USER_ID" \
  -H "x-user-role: $USER_ROLE" \
  -d '{
    "nis": "2024001",
    "nisn": "0012345678",
    "name": "Test Student",
    "gender": "MALE",
    "birthDate": "2008-01-01T00:00:00Z",
    "birthPlace": "Bandung",
    "address": "Jl. Test No. 1",
    "phone": "081234567890",
    "email": "test@example.com",
    "parentName": "Parent Name",
    "parentPhone": "081111111111",
    "parentEmail": "parent@example.com",
    "classId": "cls-001"
  }'

# Test list students
curl -X GET "$BASE_URL/students?limit=10&offset=0" \
  -H "x-school-id: $SCHOOL_ID" \
  -H "x-academic-period-id: $ACADEMIC_PERIOD_ID" \
  -H "x-user-id: $USER_ID" \
  -H "x-user-role: $USER_ROLE"
