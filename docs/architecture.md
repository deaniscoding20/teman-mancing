# Teman Mancing — Architecture Design Document

## 1. System Overview

**Teman Mancing** adalah platform komunitas digital berbasis web untuk para pemancing di Indonesia dengan tagline *"Cari Spot. Cari Teman. Gas Mancing!"*. 

Aplikasi ini mengadopsi pola arsitektur **Modular Monolith** dengan pendekatan **Clean Architecture** (Onion / Hexagonal Architecture) pada backend, serta **Component-Driven Modular Architecture** pada frontend Next.js App Router.

---

## 2. High-Level Architecture Diagram

```text
                               +-----------------------------+
                               |        CLIENT BROWSER       |
                               |    (Desktop / Mobile Web)   |
                               +--------------+--------------+
                                              |
                          HTTPS / REST API    |   WSS (WebSocket)
                                              v
                               +-----------------------------+
                               |     NEXT.JS APP (FRONTEND)  |
                               |  SSR / CSR / Tailwind CSS   |
                               +--------------+--------------+
                                              |
                                              | Reverse Proxy / REST & WSS
                                              v
+-----------------------------------------------------------------------------------------+
|                                GO BACKEND (MODULAR MONOLITH)                            |
|                                                                                         |
|  +-----------------------------------------------------------------------------------+  |
|  | DELIVERY LAYER (HTTP Handlers, Gin Engine, Route Groups, DTOs, WebSocket Handler) |  |
|  +-----------------------------------------+-----------------------------------------+  |
|                                            |                                            |
|  +-----------------------------------------v-----------------------------------------+  |
|  | APPLICATION LAYER (Use Cases, Workflows, Application Services, DTO Converters)    |  |
|  +-----------------------------------------+-----------------------------------------+  |
|                                            |                                            |
|  +-----------------------------------------v-----------------------------------------+  |
|  | DOMAIN LAYER (Entities, Domain Interfaces, Value Objects, Business Rules)         |  |
|  +-----------------------------------------^-----------------------------------------+  |
|                                            |                                            |
|  +-----------------------------------------+-----------------------------------------+  |
|  | INFRASTRUCTURE LAYER (PostgreSQL GORM Repos, MinIO S3 Client, Zap Logger, Caches) |  |
|  +-----------------------------------------------------------------------------------+  |
+------------------------------------+------------------------------------+---------------+
                                     |                                    |
                                     v                                    v
                       +---------------------------+        +---------------------------+
                       |    POSTGRESQL DATABASE    |        |    MINIO OBJECT STORAGE   |
                       |    - Relational Storage   |        |    - avatars                  |
                       |    - UUID Primary Keys    |        |    - fishing-spots            |
                       |    - ACID Transactions    |        |    - event-covers             |
                       |    - Persistent Chat      |        |    - S3 Compatible            |
                       +---------------------------+        +---------------------------+
```

---

## 3. Backend Architecture Layers

Backend Golang membagi tanggung jawab ke dalam 4 layer utama yang tidak boleh saling melanggar aturan dependensi (*Dependency Rule*: dependensi hanya boleh mengarah ke dalam / layer domain):

### 3.1 Domain Layer (`internal/domain/`)
- **Tanggung Jawab**: Mendefinisikan model domain, value objects, domain errors, dan repository interface yang murni.
- **Batasan**: Tidak boleh mengimpor package eksternal seperti Gin, GORM, database driver, atau library storage.
- **Modul**:
  - `user`: Akun pengguna, profil, password hashing rule, domain validation.
  - `fishing_spot`: Spot mancing, koordinat geolokasi, jenis ikan, teknik mancing, rekomendasi umpan.
  - `comment`: Komentar pada spot dan aturan kepemilikan komentar.
  - `event`: Agenda janjian mancing, kuota peserta, biaya (parkir, sewa perahu), visibility (public vs private).
  - `participant`: Status pendaftaran (pending, approved, rejected, cancelled), membership role.
  - `chat`: Percakapan (direct vs group), pesan, participant membership.
  - `notification`: Notifikasi sistem, unread state, event triggers.

### 3.2 Application Layer (`internal/application/`)
- **Tanggung Jawab**: Mengatur alur logika bisnis (Use Cases), orkestrasional database transactions, DTO input/output, dan interaksi antar domain.
- **Contoh Use Cases**:
  - `RegisterUserUseCase`, `LoginUserUseCase`
  - `CreateFishingSpotUseCase`, `ListFishingSpotsUseCase`, `DeleteSpotUseCase`
  - `CreateEventUseCase`, `JoinPublicEventUseCase` (ACID transaction dengan row-lock & quota check), `ApproveJoinRequestUseCase`
  - `SendMessageUseCase`, `GetConversationHistoryUseCase`
  - `UploadMediaUseCase`

### 3.3 Infrastructure Layer (`internal/infrastructure/`)
- **Tanggung Jawab**: Implementasi teknis konkrit dari interface yang didefinisikan pada layer domain/aplikasi.
- **Komponen**:
  - `database`: PostgreSQL connection pool menggunakan GORM, database migrations menggunakan `golang-migrate`.
  - `storage`: Client MinIO/S3 untuk upload, generate presigned URL, dan sanitasi file.
  - `websocket`: Connection registry, client session management, room broadcasting, dan heartbeat loop.
  - `logger`: Structured logging menggunakan Uber Zap.
  - `cache`: Abstraksi memory / Redis (optional).

### 3.4 Delivery Layer (`internal/delivery/`)
- **Tanggung Jawab**: Menerima request eksternal (HTTP REST dan WebSocket), memvalidasi format data (menggunakan `go-playground/validator`), memanggil Application Use Case, dan memetakan hasil ke format JSON standar.
- **Komponen**:
  - `http/handlers`: Controller per resource (AuthHandler, SpotHandler, EventHandler, ChatHandler, NotificationHandler).
  - `http/middleware`: Authentication (JWT), request logging (Zap), recovery/panic handling, CORS, rate limiting.
  - `http/routes`: Pengelompokan route `/api/v1` dengan proteksi middleware yang sesuai.
  - `http/dto`: Request/Response struct dengan tag JSON dan validation tag.
  - `websocket`: Upgrade HTTP connection ke WebSocket via `coder/websocket` dan integrasi dengan WS infrastructure.

---

## 4. Frontend Architecture

Frontend dibangun dengan Next.js App Router (TypeScript) dengan prinsip Feature-Driven Development:

```text
frontend/src/
├── app/                  # Next.js App Router (pages, layout, error, loading)
│   ├── (auth)/           # Route group auth: /login, /register
│   ├── spots/            # /spots, /spots/create, /spots/[id], /spots/[id]/edit
│   ├── events/           # /events, /events/create, /events/[id], /events/[id]/manage
│   ├── chat/             # /chat, /chat/[id]
│   ├── notifications/    # /notifications
│   ├── profile/          # /profile/[username]
│   ├── settings/         # /settings/profile
│   └── layout.tsx        # Root layout dengan Navbar, Footer, Providers
├── components/           # Reusable UI components (Shadcn/UI, Leaflet Map, Dialog, Button)
├── features/             # Business feature modules (Components, hooks, types per fitur)
│   ├── auth/
│   ├── spots/
│   ├── events/
│   ├── chat/
│   ├── notifications/
│   └── profile/
├── hooks/                # Custom React hooks (useAuth, useWebSocket, useDebounce)
├── lib/                  # Utilities, API client (fetcher), cn helper
├── services/             # API services yang membungkus endpoint backend
├── stores/               # Client-side state management (Zustand / React Context)
├── types/                # Global TypeScript definitions & API response types
└── utils/                # Helper formatting (tanggal, format Rupiah, geo math)
```

---

## 5. Critical Technical Flows

### 5.1 Event Join Transaction & Race Condition Prevention
Ketika seorang pengguna menekan tombol "Join" pada event:
1. HTTP Handler menerima request dengan token user.
2. Application Use Case memulai Database Transaction (`tx = db.Begin()`).
3. Query event dengan Row-Level Lock: `SELECT * FROM events WHERE id = ? FOR UPDATE`.
4. Validasi status: event belum cancelled, deadline belum terlewati, event berstatus `active`.
5. Hitung jumlah approved participants: `SELECT COUNT(*) FROM event_participants WHERE event_id = ? AND status = 'approved'`.
6. Jika `max_participants` terisi dan kuota penuh -> Return `ErrEventFull`, `tx.Rollback()`.
7. Jika lolos:
   - Jika Public: buat record `event_participants` dengan status `approved`. Tambahkan pengguna ke `conversation_participants` untuk grup chat event tersebut.
   - Jika Private: buat record `event_participants` dengan status `pending`. Buat notifikasi untuk creator event.
8. Commit transaction (`tx.Commit()`).

### 5.2 Realtime Chat & WebSocket Membership
1. Client melakukan handshake WebSocket ke `/ws?token=<jwt>`.
2. Backend memvalidasi token JWT saat handshake.
3. Hub mendaftarkan client ke active connections map.
4. Ketika client mengirim atau mendengarkan room conversation:
   - Backend memverifikasi bahwa `user_id` adalah anggota valid di `conversation_participants`.
   - Jika tidak terdaftar, request ditolak (Unauthorized).
5. Pesan baru disimpan terlebih dahulu ke PostgreSQL tabel `messages` sebelum di-broadcast ke anggota room yang sedang online.
6. Untuk anggota yang offline, unread status diakumulasi dan notifikasi dibuat.

### 5.3 Storage & File Upload
1. Frontend mengirim `multipart/form-data` ke Go backend.
2. Backend memvalidasi MIME type (hanya JPEG, PNG, WEBP) dan ukuran file (maksimal 5MB).
3. Backend membaca stream dan meng-upload langsung ke MinIO dengan nama file aman: `{bucket}/{entity_id}/{uuid}.webp`.
4. Backend menyimpan *object key* ke database PostgreSQL.
5. URL file disajikan via presigned URL atau API proxy publik backend sehingga MinIO credentials tidak pernah terekspos ke frontend.

---

## 6. Security Principles
- **Password Storage**: Hashing menggunakan `bcrypt` dengan cost factor 12 atau `argon2id`.
- **Session/Token**: JWT Token bertanda tangan HMAC-SHA256 dengan secret aman, payload minimal (`sub`, `username`, `exp`).
- **Authorization**: Ownership check ketat di backend pada setiap operasi update/delete (hanya pemilik spot/komentar/event yang dapat memodifikasi).
- **Injection Protection**: Parameterized queries via GORM / database/sql.
- **XSS & CORS**: Sanitasi input dan header CORS ketat yang hanya mengizinkan origin frontend yang terdaftar.
