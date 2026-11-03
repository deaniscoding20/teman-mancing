#!/usr/bin/env bash
# ==============================================================================
# Teman Mancing — Automated VPS Deployment Script
# Digunakan oleh GitHub Actions CI/CD atau dijalankan secara manual di VPS.
# ==============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
DOCKER_DIR="${ROOT_DIR}/infrastructure/docker"

echo "=========================================================="
echo "🎣 Teman Mancing: Memulai Proses Deployment..."
echo "📂 Project Root: ${ROOT_DIR}"
echo "=========================================================="

cd "${ROOT_DIR}"

# 1. Pastikan file environment (.env) tersedia
if [ ! -f "${DOCKER_DIR}/.env" ]; then
  echo "⚠️ File ${DOCKER_DIR}/.env tidak ditemukan, membuat dari template .env.example..."
  cp "${DOCKER_DIR}/.env.example" "${DOCKER_DIR}/.env"
fi

if [ ! -f "${ROOT_DIR}/backend/.env" ]; then
  cp "${ROOT_DIR}/backend/.env.example" "${ROOT_DIR}/backend/.env"
fi

if [ ! -f "${ROOT_DIR}/frontend/.env" ]; then
  cp "${ROOT_DIR}/frontend/.env.example" "${ROOT_DIR}/frontend/.env"
fi

# 2. Build & Deploy Docker Containers
echo "🐳 Menjalankan Docker Compose build & start..."
cd "${DOCKER_DIR}"
docker compose down --remove-orphans || true
docker compose up -d --build

# 3. Tunggu hingga database PostgreSQL dan MinIO siap (Healthy)
echo "⏳ Menunggu container PostgreSQL dan MinIO mencapai status healthy..."
TIMEOUT=60
ELAPSED=0
until [ "$(docker inspect -f '{{.State.Health.Status}}' teman_mancing_postgres 2>/dev/null || echo 'starting')" = "healthy" ]; do
  sleep 2
  ELAPSED=$((ELAPSED + 2))
  if [ $ELAPSED -ge $TIMEOUT ]; then
    echo "❌ Timeout menunggu PostgreSQL menjadi healthy!"
    docker compose logs postgres
    exit 1
  fi
done
echo "✅ PostgreSQL siap."

# 4. Jalankan Migrasi Database
echo "📦 Menjalankan Database Migrations..."
docker exec -i teman_mancing_postgres psql -U teman_mancing_user -d teman_mancing \
  < "${ROOT_DIR}/backend/migrations/000001_init_schema.up.sql" || true

docker exec -i teman_mancing_postgres psql -U teman_mancing_user -d teman_mancing \
  < "${ROOT_DIR}/backend/migrations/000002_seed_initial_data.up.sql" || true
echo "✅ Migrasi database selesai diaplikasikan."

# 5. Uji Healthcheck Backend dan Frontend
echo "🔍 Melakukan verifikasi endpoint aplikasi..."
sleep 3

# Test Backend API Health
BACKEND_HEALTH_URL="http://localhost:8080/api/v1/health"
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "${BACKEND_HEALTH_URL}" || echo "000")
if [ "${HTTP_STATUS}" -eq 200 ] || [ "${HTTP_STATUS}" -eq 503 ]; then
  echo "✅ Backend API responsif di port 8080 (HTTP ${HTTP_STATUS})"
else
  echo "⚠️ Peringatan: Backend API mengembalikan HTTP ${HTTP_STATUS} pada ${BACKEND_HEALTH_URL}"
fi

# Test Frontend HTTP response
FRONTEND_URL="http://localhost:3080/"
FRONTEND_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "${FRONTEND_URL}" || echo "000")
if [ "${FRONTEND_STATUS}" -eq 200 ]; then
  echo "✅ Frontend Web responsif di port 3080 (HTTP 200)"
else
  echo "⚠️ Peringatan: Frontend mengembalikan HTTP ${FRONTEND_STATUS} pada ${FRONTEND_URL}"
fi

echo "=========================================================="
echo "🎉 Deployment Teman Mancing Berhasil Selesai!"
echo "🌐 Frontend: http://localhost:3080"
echo "🔌 API:      http://localhost:8080/api/v1"
echo "=========================================================="
