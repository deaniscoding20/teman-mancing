# 📋 Pull Request Changes Info

## 📌 PR Title
`feat(ci-cd): implement GitHub Actions CI/CD pipeline and automated VPS deployment`

---

## 🌿 Branch Info
* **Head Branch**: `feature/ci-cd-pipeline`
* **Base Branch**: `main`

---

## 🎯 Summary of Changes
Pull Request ini mengimplementasikan pipeline **Continuous Integration (CI)** dan **Continuous Deployment (CD)** lengkap menggunakan GitHub Actions, skrip automasi deployment VPS, serta standarisasi alur kerja Pull Request dengan dokumentasi `changesinfo.md`.

### 1. Continuous Integration (`.github/workflows/ci.yml`)
* **Trigger**: 
  * `pull_request` ke branch `main` (memvalidasi setiap PR sebelum di-merge).
  * `push` ke branch `main`.
* **Jobs**:
  * `backend-test`: Menjalankan dependency download, full unit & integration testing Golang (`go test -v -race ./tests/...`), dan verifikasi kompilasi binary statis (`go build`).
  * `frontend-test`: Menjalankan `npm ci`, typecheck TypeScript (`npx tsc --noEmit`), dan kompilasi Next.js production build (`npm run build`).
  * `docker-validate`: Memvalidasi sintaks `docker-compose.yml` agar tidak terjadi kesalahan konfigurasi container.

### 2. Continuous Deployment (`.github/workflows/cd.yml`)
* **Trigger**: `push` ke branch `main` (otomatis terpicu ketika suatu Pull Request di-merge ke `main`).
* **Jobs**:
  * `verify`: Memverifikasi ulang seluruh tes backend dan frontend.
  * `deploy-vps`: Melakukan SSH ke VPS secara otomatis menggunakan GitHub Repository Secrets (`VPS_HOST`, `VPS_SSH_KEY`, `VPS_USERNAME`, `VPS_PORT`).
  * Mengeksekusi `git fetch origin main`, `git reset --hard origin/main`, dan menjalankan skrip `deploy.sh`.
  * Terdapat fallback anggun (*graceful fallback*) jika secrets belum dikonfigurasi pada GitHub repository settings.

### 3. Automated VPS Deployment Script (`infrastructure/scripts/deploy.sh`)
* Mengotomatisasi siklus hidup deployment di server VPS:
  1. Verifikasi ketersediaan file `.env`.
  2. Eksekusi `docker compose up -d --build`.
  3. Menunggu container database PostgreSQL dan MinIO hingga berstatus `healthy`.
  4. Menjalankan database migrations (`000001_init_schema.up.sql` dan `000002_seed_initial_data.up.sql`).
  5. Menjalankan health check curl pada backend (`:8080/api/v1/health`) dan frontend (`:3080/`).

### 4. Deployment Documentation (`DeployAndRun.md`)
* Menambahkan bagian panduan konfigurasi GitHub Repository Secrets untuk mengaktifkan auto-deploy.
* Menjelaskan aturan alur kerja PR dan pembuatan file `changesinfo.md` untuk kolaborasi multi-agent.

---

## 📂 Files Changed
| File | Status | Keterangan |
|---|---|---|
| `.github/workflows/ci.yml` | **Baru** | GitHub Actions workflow untuk CI testing (Go + Next.js + Docker) |
| `.github/workflows/cd.yml` | **Baru** | GitHub Actions workflow untuk auto-deploy ke VPS saat merge ke `main` |
| `infrastructure/scripts/deploy.sh` | **Baru** | Shell script deployment otomatis dengan healthcheck & migrasi DB |
| `DeployAndRun.md` | **Diperbarui** | Panduan setup GitHub Secrets CI/CD dan alur kerja PR |
| `changesinfo.md` | **Baru** | Dokumen rincian perubahan untuk proses review dan merge |

---

## ✅ Verification & Test Results
Seluruh pengujian lokal telah dijalankan dan lulus 100%:
* **Backend Tests**: `go test -v ./tests/...` ➔ **PASS** (100% lulus pada semua modul).
* **Backend Build**: `go build -ldflags="-w -s" -o bin/api ./cmd/api` ➔ **SUCCESS**.
* **Frontend Typecheck**: `npx tsc --noEmit` ➔ **0 errors (PASS)**.
* **Frontend Build**: `npm run build` ➔ **13 rute App terkompilasi optimal (PASS)**.
* **Docker Compose Syntax**: `docker compose config` ➔ **Valid**.

---

## 🔍 Checklist untuk Reviewer / Agent Lain
Sebelum melakukan merge ke branch `main`, pastikan hal berikut:
- [x] Pipeline CI workflow mencakup pengujian backend Golang dan frontend Next.js.
- [x] Script `deploy.sh` memiliki izin eksekusi (`chmod +x`).
- [x] File `.env.example` tidak membocorkan kredensial produksi.
- [x] Dokumentasi `DeployAndRun.md` memuat panduan GitHub Secrets dengan jelas.
- [x] Tidak ada file binary atau `node_modules` yang ter-commit.

---

## 🚀 Cara Merge Pull Request Ini
Bagi Agent lain atau reviewer:
1. Periksa tab **Files changed** pada Pull Request ini di GitHub.
2. Pastikan workflow check **CI Pipeline** berstatus hijau (semua checks lulus).
3. Klik tombol **Squash and merge** atau **Merge pull request**.
4. Setelah ter-merge ke `main`, workflow **CD Auto-Deploy** akan otomatis berjalan.
