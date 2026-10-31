#!/usr/bin/env bash
set -uo pipefail
cd /c/C/projects/cbt-quasar/sekolah-platform

echo "===== 1. GIT STATUS ====="
git status --short

echo ""
echo "===== 2. BUILD TU-CORE ====="
(cd services/tu-core && go build ./... 2>&1)

echo ""
echo "===== 3. BUILD EDGE-BFF ====="
(cd services/edge-bff && go build ./... 2>&1)

echo ""
echo "===== 4. GREP EDGE-BFF REFS ====="
grep -rn 'pbPreference\|NewPreferenceServiceClient\|NewContextHandler\|contextHandler' services/edge-bff/ 2>&1 || echo "(no matches)"

echo ""
echo "===== 5. MIGRATION STATUS ====="
migrate -path services/tu-core/migrations \
  -database "postgres://-smk_dev:-secret123@localhost:5432/-smk_platform?sslmode=disable" \
  version 2>&1

echo ""
echo "===== 6. SEMESTER TABLE ====="
docker exec smk-postgres psql -U -smk_dev -d -smk_platform -c "SELECT * FROM semester;" 2>&1

echo ""
echo "===== 7. PREFERENCE TABLE ====="
docker exec smk-postgres psql -U -smk_dev -d -smk_platform -c "SELECT * FROM user_period_preference;" 2>&1

echo ""
echo "===== 8. CURL PERIODS ====="
curl -sS -i http://localhost:8081/api/v1/context/periods \
  -H "x-school-id: school-abc" -H "x-user-id: user-123" 2>&1 | head -30

echo ""
echo "===== 9. CURL PREFERENCE GET ====="
curl -sS -i http://localhost:8081/api/v1/users/me/preferences \
  -H "x-school-id: school-abc" -H "x-user-id: user-123" 2>&1 | head -20

echo ""
echo "===== 10. CURL PREFERENCE PUT ====="
curl -sS -i -X PUT http://localhost:8081/api/v1/users/me/preferences \
  -H "Content-Type: application/json" \
  -H "x-school-id: school-abc" -H "x-user-id: user-123" \
  -d '{"semester_id":"20252"}' 2>&1 | head -20

echo ""
echo "===== DONE ====="
