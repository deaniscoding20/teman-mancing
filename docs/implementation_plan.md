# Teman Mancing — Implementation Plan

Dokumen ini mendefinisikan roadmap implementasi terstruktur untuk platform **Teman Mancing** berdasarkan README.md. Proyek dikerjakan secara bertahap (phased rollout) dengan verifikasi ketat pada setiap fase sebelum melangkah ke fase berikutnya.

---

## Roadmap Overview

| Phase | Title | Scope Utama | Status |
|---|---|---|---|
| **Phase 1** | **Project Foundation** | Scaffolding monorepo, Docker Compose, DB migrations, Backend Clean Arch foundation, Frontend Next.js foundation, CI/testing base | **In Progress (Current Milestone)** |
| **Phase 2** | **Authentication & User Profile** | Register, Login, Logout, Session/JWT, Password Hashing, Profile management, MinIO Avatar upload | Pending |
| **Phase 3** | **Fishing Spots** | Spot directory, Leaflet map, CRUD spots, Multi-photo MinIO upload, Fish types & Techniques, Comments | Pending |
| **Phase 4** | **Fishing Events / Janjian Mancing** | Event listing, Create event, Cost breakdowns, Public/Private workflows, Capacity race condition locks, Approvals | Pending |
| **Phase 5** | **Realtime Chat & Messaging** | Direct messaging, Event group chat, Coder/WebSocket, PostgreSQL message persistence, Unread badge | Pending |
| **Phase 6** | **Notifications** | Event & spot activity notifications, in-app bell, realtime pushes | Pending |
| **Phase 7** | **Finalization & Polishing** | Responsive UI audit, End-to-end testing, security hardening, performance tuning | Pending |

---

## Detailed Milestone: Phase 1 — Project Foundation

### Objective
Membangun fondasi arsitektural yang kokoh, modular, dan teruji untuk backend (Golang) dan frontend (Next.js), serta orkestrasi container menggunakan Docker Compose untuk PostgreSQL dan MinIO.

### Task Breakdown

#### 1. Repository & Infrastructure Setup
- [x] Inisialisasi struktur monorepo (`backend/`, `frontend/`, `infrastructure/`, `docs/`).
- [x] Konfigurasi `.gitignore` komprehensif untuk Go binaries, Next.js artifacts, environment secrets, dan persistent volumes.
- [x] Buat `infrastructure/docker/docker-compose.yml` yang mencakup:
  - `postgres` (PostgreSQL 16 Alpine, persistent volume, healthcheck via `pg_isready`).
  - `minio` (S3 object storage, console UI port 9001, API port 9000, persistent volume).
  - `minio-create-buckets` (Init container menggunakan `minio/mc` untuk provisioning bucket `avatars`, `fishing-spots`, dan `event-covers`).
  - `backend` (Golang Gin container, multi-stage build, healthcheck).
  - `frontend` (Next.js container, multi-stage build).
- [x] Sediakan environment templates:
  - Root / infra `.env.example`
  - `backend/.env.example`
  - `frontend/.env.example`

#### 2. Database Schema & Migration Foundation
- [x] Buat SQL migration files menggunakan konvensi `golang-migrate`:
  - `backend/migrations/000001_init_schema.up.sql` (15 tabel utama, UUID extension, Foreign Keys, Unique constraints, indexes).
  - `backend/migrations/000001_init_schema.down.sql` (Reversible drop statements).
  - `backend/migrations/000002_seed_initial_data.up.sql` (Master data jenis ikan & teknik mancing di Indonesia).
  - `backend/migrations/000002_seed_initial_data.down.sql`.
- [x] Dokumentasikan skema relasional di `docs/database.md`.

#### 3. Backend Architecture Scaffolding (Golang)
- [x] Inisialisasi `go.mod` (Go 1.23, module `teman-mancing/backend`).
- [x] Terapkan 4 Layer Clean Architecture:
  - **Domain**: Entitas domain murni (`user`, `fishing_spot`, `event`, `chat`, `notification`).
  - **Application**: Abstraksi repository interface dan standard application error definitions.
  - **Infrastructure**:
    - `logger`: Structured logging berbasis Zap.
    - `database`: GORM PostgreSQL connection pool dengan ping & connection tuning.
    - `storage`: MinIO S3 client dengan auto bucket checking.
  - **Delivery**:
    - `http/dto`: Standard JSON response envelop (`success`, `message`, `data`, `meta`, `errors`).
    - `http/middleware`: Request logger, panic recovery, CORS handler.
    - `http/handlers`: Health check handler (`/api/v1/health`) dan version handler (`/api/v1/version`).
    - `http/routes`: Gin engine router setup.
- [x] Implementasikan entry point `cmd/api/main.go` dengan graceful shutdown.
- [x] Buat backend unit tests untuk memverifikasi middleware, response DTOs, domain models, dan config loader.

#### 4. Frontend Foundation (Next.js + TypeScript)
- [x] Inisialisasi konfigurasi Next.js App Router:
  - `package.json` dengan dependensi React, Next.js, Lucide icons, Tailwind CSS, TanStack Query.
  - `tsconfig.json` dengan path alias `@/*` mengarah ke `./src/*`.
  - `tailwind.config.ts` dengan warna identitas Teman Mancing:
    - Primary: `#123B50`
    - Secondary: `#35A7B5`
    - Accent: `#7BAE7F`
    - Background: `#F5F7F7`
    - Surface: `#FFFFFF`
    - Text: `#172B36`
  - Font konfigurasi: Plus Jakarta Sans / Inter.
- [x] Folder struktur modular sesuai README Section 5:
  - `src/features/{auth,spots,events,chat,notifications,profile}`
  - `src/components/ui/`, `src/hooks/`, `src/lib/`, `src/services/`, `src/stores/`, `src/types/`, `src/utils/`.
- [x] Buat layout utama dengan Navbar responsif, brand title "Teman Mancing", dan footer informatif.
- [x] Implementasikan API client `src/lib/api-client.ts` untuk komunikasi dengan backend.
- [x] Validasi frontend build dengan `npm run build` / type checking.

---

## Phase 2 Preview: Authentication & User Profile
Langkah berikutnya setelah fondasi tervalidasi:
1. Domain entity methods untuk password hashing dan validasi format username & email.
2. User repository implementation dengan PostgreSQL GORM.
3. Use case `RegisterUser`, `LoginUser`, `GetMe`, `UpdateProfile`.
4. JWT token generation & verification service (Access Token + Refresh Token).
5. Delivery HTTP handlers & route grouping `/api/v1/auth` dan `/api/v1/users`.
6. Frontend Auth Forms (Login & Register dengan Zod & React Hook Form).
