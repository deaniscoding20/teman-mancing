# 🎣 TEMAN MANCING
## Fishing Community Platform

**Tagline: Cari Spot. Cari Teman. Gas Mancing!**

---

# 1. PROJECT OVERVIEW

## 1.1 Tentang Project

Teman Mancing adalah platform komunitas digital untuk para pemancing di Indonesia.

Platform ini dibuat untuk menghubungkan pemancing agar dapat:
- Menemukan lokasi atau spot mancing baru.
- Berbagi informasi spot mancing.
- Berinteraksi dan berdiskusi dengan pemancing lain.
- Membuat agenda atau event mancing bersama.
- Mencari teman mancing.
- Melakukan komunikasi melalui personal chat dan grup chat event.

Teman Mancing bukan hanya direktori lokasi memancing, tetapi juga merupakan social community platform khusus pemancing.

## 1.2 Product Vision

Menjadi platform komunitas memancing yang membantu pemancing Indonesia menemukan tempat memancing, mendapatkan informasi yang bermanfaat, dan membangun koneksi dengan sesama pemancing.

## 1.3 Target Users

1. Pemancing pemula.
2. Pemancing berpengalaman.
3. Komunitas mancing.
4. Pemancing yang ingin mencari teman memancing.
5. Pemancing yang ingin mengeksplorasi spot baru.

## 1.4 Core Features

Aplikasi memiliki fitur utama:

1. Authentication & User Profile.
2. Fishing Spots.
3. Fishing Events (Janjian Mancing).
4. Chat & Messaging.
5. Notifications.

---

# 2. IMPORTANT INSTRUCTIONS FOR AI AGENT

Dokumen ini merupakan sumber utama requirement proyek.

AI Agent yang mengerjakan project ini harus mengikuti instruksi berikut.

### 2.1 General Rules

- Pahami seluruh README sebelum mulai melakukan perubahan kode.
- Gunakan teknologi yang sudah ditentukan.
- Jangan mengganti teknologi utama tanpa alasan teknis yang kuat.
- Jangan membuat aplikasi hanya berupa UI mockup.
- Semua fitur utama harus terhubung dengan backend dan database.
- Gunakan clean code dan modular architecture.
- Hindari membuat satu file dengan terlalu banyak tanggung jawab.
- Jangan membuat business logic di dalam HTTP handlers.
- Jangan mengakses database langsung dari frontend.
- Jangan menyimpan password dalam bentuk plain text.
- Jangan menyimpan gambar dalam database.
- Jangan membuat tombol yang tidak memiliki fungsi.
- Jangan menggunakan dummy data sebagai pengganti backend yang sebenarnya.
- Gunakan dummy data hanya untuk development, preview, atau database seeding.
- Jangan menghapus fitur yang sudah selesai ketika mengembangkan fitur lain.

### 2.2 Instructions for Hermes Agent + Gemini

Bertindak sebagai Senior Full-Stack Engineer dan Software Architect.

Saat mengerjakan project:

1. Analisis requirement sebelum coding.
2. Periksa struktur repository sebelum membuat file baru.
3. Gunakan teknologi yang telah ditentukan dalam README.
4. Buat implementation plan untuk pekerjaan yang kompleks.
5. Implementasikan fitur secara bertahap.
6. Periksa integrasi frontend dan backend.
7. Jalankan testing setelah implementasi.
8. Perbaiki error sebelum melanjutkan ke fitur berikutnya.
9. Dokumentasikan perubahan arsitektur yang penting.
10. Jangan mengklaim fitur selesai jika belum benar-benar diimplementasikan dan diuji.

Jika terdapat informasi teknis yang belum ditentukan:
- Gunakan best practice.
- Pilih solusi yang sederhana dan maintainable.
- Hindari overengineering.
- Dokumentasikan keputusan yang dibuat.
- Jangan menanyakan hal yang sudah dijelaskan dalam README.

### 2.3 Development Workflow

Setiap kali mengerjakan task:

1. Baca requirement.
2. Inspect existing code.
3. Identifikasi file yang perlu dibuat atau diubah.
4. Implementasikan solusi.
5. Jalankan formatter dan linting.
6. Jalankan unit test atau integration test yang relevan.
7. Periksa error.
8. Update dokumentasi jika diperlukan.
9. Laporkan hasil pekerjaan secara ringkas.

---

# 3. TECHNOLOGY STACK

Gunakan teknologi berikut sebagai standar utama.

## 3.1 Backend

**Language: Golang**

Backend menggunakan Golang sebagai bahasa utama.

Technology:
- Go.
- Gin Framework.
- GORM.
- PostgreSQL Driver.
- go-playground/validator.
- golang-jwt/jwt jika menggunakan JWT.
- Uber Zap atau Go slog untuk logging.
- golang-migrate untuk database migration.
- coder/websocket untuk WebSocket.

Tanggung jawab backend:
- Authentication.
- Authorization.
- Business logic.
- Database management.
- REST API.
- File upload.
- Realtime chat.
- Notification management.

## 3.2 Frontend

**Framework: Next.js**

Frontend bebas dikembangkan selama menggunakan teknologi modern dan kompatibel dengan backend.

Rekomendasi:

- Next.js App Router.
- React.
- TypeScript.
- Tailwind CSS.
- Shadcn/UI.
- TanStack Query.
- Zustand jika diperlukan.
- React Hook Form.
- Zod.
- Leaflet.
- Socket.IO Client atau native WebSocket.

Frontend bertanggung jawab untuk:
- User interface.
- User interaction.
- Form validation.
- API integration.
- Realtime communication.
- Responsive layout.

## 3.3 Database

**PostgreSQL**

Gunakan PostgreSQL sebagai database utama.

ORM:
- GORM.

Migration:
- golang-migrate.

Database harus memiliki relational structure yang baik.

Gunakan:
- UUID sebagai primary key.
- Foreign keys.
- Database indexes.
- Unique constraints.
- Transactions.
- TIMESTAMPTZ untuk tanggal dan waktu.
- Soft delete jika diperlukan.

## 3.4 Object Storage

**MinIO**

Gunakan MinIO untuk menyimpan seluruh file gambar.

Contoh file:
- User avatar.
- Fishing spot photos.
- Event cover.
- Fishing spot gallery.

Gunakan S3-compatible API.

Backend Golang bertanggung jawab atas upload dan penghapusan file.

Database hanya menyimpan object key, bukan binary image.

## 3.5 Realtime

Gunakan WebSocket untuk:

- Personal chat.
- Event group chat.
- Realtime notifications.
- Message updates.

PostgreSQL digunakan untuk menyimpan seluruh pesan secara persistent.

Redis dapat ditambahkan untuk Pub/Sub ketika backend sudah berjalan pada beberapa instance.

## 3.6 Maps

Gunakan:
- Leaflet.
- OpenStreetMap.

Peta digunakan untuk menampilkan lokasi fishing spot dan lokasi event.

## 3.7 Infrastructure

Gunakan:
- Docker.
- Docker Compose.
- Nginx (optional).
- Git.
- GitHub Actions (optional).

---

# 4. SOFTWARE ARCHITECTURE

## 4.1 Architecture Pattern

Gunakan:

**Modular Monolith + Clean Architecture**

Jangan menggunakan Microservices untuk tahap awal.

Alasan:
- Lebih mudah dikembangkan.
- Deployment lebih sederhana.
- Business logic tetap terpisah.
- Tidak membutuhkan komunikasi antar service yang kompleks.
- Mudah di-maintain.
- Bisa dikembangkan menjadi Microservices jika dibutuhkan.

## 4.2 Architecture Layers

Backend harus memiliki 4 layer utama.

### Domain Layer

Berisi:
- Entities.
- Business rules.
- Domain interfaces.
- Domain validations.

Tidak boleh bergantung pada Gin, GORM, maupun infrastructure.

### Application Layer

Berisi:
- Use cases.
- Application services.
- DTO.
- Repository interfaces.
- Business workflows.

Contoh:
- RegisterUser.
- LoginUser.
- CreateSpot.
- JoinEvent.
- ApproveJoinRequest.
- SendMessage.

### Infrastructure Layer

Berisi implementasi:
- PostgreSQL repository.
- MinIO storage.
- WebSocket infrastructure.
- JWT atau session service.
- Redis jika digunakan.
- Logging.

### Delivery Layer

Berisi:
- HTTP handlers.
- REST API routes.
- Request validation.
- Middleware.
- WebSocket handlers.

HTTP handler hanya bertugas menerima request, memvalidasi format request, memanggil application layer, dan mengembalikan response.

Jangan menempatkan business logic kompleks di HTTP handlers.

## 4.3 High-Level Architecture

```text
                   USERS
                     |
                     v
               NEXT.JS APP
                     |
              HTTPS / REST API
                     |
                     v
               GO BACKEND
                     |
        +------------+------------+
        |            |            |
        v            v            v
      AUTH        BUSINESS     WEBSOCKET
      MODULE       MODULE       MODULE
        |            |            |
        +------------+------------+
                     |
           +---------+---------+
           |                   |
           v                   v
       POSTGRESQL             MINIO
        DATABASE          OBJECT STORAGE
           |
           v
       PERSISTENT DATA

       REDIS (OPTIONAL)
       Cache / PubSub
```

Seluruh modul backend berjalan dalam satu codebase Golang.

---

# 5. PROJECT STRUCTURE

Gunakan monorepo dengan pemisahan frontend dan backend.

```text
teman-mancing/
│
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── user/
│   │   │   ├── fishing_spot/
│   │   │   ├── comment/
│   │   │   ├── event/
│   │   │   ├── participant/
│   │   │   ├── chat/
│   │   │   └── notification/
│   │   │
│   │   ├── application/
│   │   │   ├── auth/
│   │   │   ├── user/
│   │   │   ├── fishing_spot/
│   │   │   ├── comment/
│   │   │   ├── event/
│   │   │   ├── participant/
│   │   │   ├── chat/
│   │   │   └── notification/
│   │   │
│   │   ├── infrastructure/
│   │   │   ├── database/
│   │   │   ├── storage/
│   │   │   ├── websocket/
│   │   │   ├── cache/
│   │   │   └── logger/
│   │   │
│   │   ├── delivery/
│   │   │   ├── http/
│   │   │   │   ├── handlers/
│   │   │   │   ├── middleware/
│   │   │   │   ├── routes/
│   │   │   │   └── dto/
│   │   │   │
│   │   │   └── websocket/
│   │   │
│   │   └── config/
│   │
│   ├── migrations/
│   ├── tests/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── .env.example
│
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── app/
│   │   ├── components/
│   │   ├── features/
│   │   │   ├── auth/
│   │   │   ├── spots/
│   │   │   ├── events/
│   │   │   ├── chat/
│   │   │   ├── notifications/
│   │   │   └── profile/
│   │   ├── hooks/
│   │   ├── lib/
│   │   ├── services/
│   │   ├── stores/
│   │   ├── types/
│   │   └── utils/
│   │
│   ├── package.json
│   ├── Dockerfile
│   └── .env.example
│
├── infrastructure/
│   ├── docker/
│   │   ├── docker-compose.yml
│   │   └── nginx/
│   │
│   └── scripts/
│
├── docs/
│   ├── architecture.md
│   ├── database.md
│   └── api.md
│
├── .gitignore
└── README.md
```

Jangan membuat semua logic dalam satu package.

Setiap module harus memiliki tanggung jawab yang jelas.

---

# 6. AUTHENTICATION & USER PROFILE

## 6.1 Register

Pengguna dapat membuat akun dengan:
- Nama lengkap.
- Username.
- Email.
- Password.
- Foto profil (optional).

Rules:
- Email harus unik.
- Username harus unik.
- Password wajib di-hash.
- Validasi input.
- Berikan feedback ketika registrasi berhasil atau gagal.

## 6.2 Login

Pengguna dapat login menggunakan email dan password.

Fitur:
- Login.
- Logout.
- Remember session.
- Forgot password.
- Session management.

Gunakan password hashing Argon2id atau bcrypt.

Gunakan secure cookie atau mekanisme token yang aman.

Jangan menyimpan password atau refresh token dalam bentuk plain text.

## 6.3 User Profile

Setiap pengguna memiliki:
- Foto profil.
- Nama.
- Username.
- Bio.
- Lokasi atau domisili.
- Jumlah spot yang dibagikan.
- Jumlah event yang dibuat.
- Jumlah event yang diikuti.

Pengguna dapat mengedit profil miliknya.

---

# 7. FISHING SPOTS

Buat section bernama **Spot Mancing**.

Section ini menjadi direktori komunitas untuk berbagi dan menemukan lokasi memancing.

## 7.1 Spot Listing

Tampilkan daftar spot dalam bentuk card.

Informasi:
- Foto utama.
- Nama spot.
- Lokasi.
- Jenis ikan.
- Teknik memancing.
- Nama pembuat.
- Jumlah komentar.

Fitur:
- Search berdasarkan nama.
- Search berdasarkan lokasi.
- Filter jenis ikan.
- Filter teknik memancing.
- Filter wilayah.
- Sorting spot terbaru.
- Pagination.

## 7.2 Spot Detail

Tampilkan:

**Informasi Spot**
- Nama spot.
- Deskripsi.
- Galeri foto.
- Alamat.
- Latitude.
- Longitude.
- Jenis ikan.
- Teknik memancing.
- Umpan rekomendasi.
- Nama pembuat.
- Tanggal dibuat.

**Map**
- Tampilkan lokasi pada peta.
- Gunakan Leaflet + OpenStreetMap.
- Tampilkan marker lokasi.
- Sediakan link petunjuk arah.

**Comments**
- Pengguna dapat memberikan komentar.
- Pengguna dapat melihat komentar lain.
- Pengguna hanya dapat menghapus komentarnya sendiri.
- Tampilkan nama dan avatar pengguna.
- Tampilkan waktu komentar.

## 7.3 Create Spot

Pengguna login dapat membuat spot.

Form:
- Nama spot.
- Deskripsi.
- Upload foto.
- Lokasi melalui map.
- Alamat.
- Jenis ikan (multiple selection).
- Teknik memancing (multiple selection).
- Umpan rekomendasi.

Fitur:
- Create.
- Edit.
- Delete.
- Upload multiple photos.

Pengguna hanya boleh mengedit dan menghapus spot miliknya sendiri.

---

# 8. FISHING EVENTS / JANJIAN MANCING

Buat section bernama **Janjian Mancing**.

Fitur ini digunakan untuk membuat kegiatan memancing bersama.

## 8.1 Event Listing

Tampilkan card event dengan informasi:

- Cover event.
- Nama event.
- Lokasi tujuan.
- Tanggal.
- Waktu.
- Meeting point.
- Tingkat kesulitan.
- Jumlah peserta.
- Nama pembuat.
- Public atau Private.
- Status event.

Filter:
- Lokasi.
- Tanggal.
- Tingkat kesulitan.
- Public / Private.
- Status pendaftaran.

## 8.2 Create Event

Form pembuatan event.

### Informasi Utama

- Nama event.
- Deskripsi.
- Cover event.
- Lokasi tujuan spot.
- Tanggal pelaksanaan.
- Jam mulai.
- Jam selesai (optional).
- Tempat janjian.
- Tingkat kesulitan.

### Informasi Biaya

- Apakah ada biaya parkir?
- Estimasi biaya parkir.
- Apakah membutuhkan penyewaan perahu?
- Estimasi biaya sewa perahu.
- Biaya tambahan lainnya.

### Pengaturan Event

Jenis event:

**Public Event**
- Semua pengguna dapat melihat event.
- Pengguna dapat langsung join.
- Jika berhasil join, otomatis masuk grup chat event.

**Private Event**
- Event dapat ditampilkan di listing sesuai pengaturan visibility.
- Pengguna harus mengirim request untuk bergabung.
- Creator harus approve request.
- Jika disetujui, pengguna otomatis menjadi anggota grup chat.

Pengaturan tambahan:
- Maksimal peserta (optional).
- Registration deadline (optional).
- Status event.

## 8.3 Event Detail

Tampilkan:
- Nama.
- Cover.
- Deskripsi.
- Lokasi.
- Peta.
- Tanggal dan waktu.
- Meeting point.
- Tingkat kesulitan.
- Biaya parkir.
- Biaya sewa perahu.
- Biaya tambahan.
- Nama creator.
- Daftar peserta.
- Jumlah peserta.
- Status.
- Tombol join atau request join.

## 8.4 Event Participation

Gunakan status:

- Pending.
- Approved.
- Rejected.
- Cancelled.

Untuk Public Event:
- Pengguna dapat langsung bergabung.
- Validasi kapasitas.
- Hindari duplicate membership.

Untuk Private Event:
- Pengguna mengirim request.
- Creator mendapatkan notifikasi.
- Creator dapat approve atau reject.
- Setelah approve, pengguna menjadi peserta.

## 8.5 Event Management

Creator dapat:
- Edit event.
- Cancel event.
- Melihat peserta.
- Melihat join requests.
- Approve request.
- Reject request.
- Mengeluarkan peserta jika diperlukan.

Peserta dapat meninggalkan event.

Creator tidak boleh mengelola event milik pengguna lain.

---

# 9. CHAT & MESSAGING

Buat section bernama **Chat**.

Ada dua jenis percakapan.

## 9.1 Personal Chat

Pengguna dapat chat dengan pengguna lain.

Fitur:
- Membuat percakapan personal.
- Mengirim pesan teks.
- Melihat riwayat pesan.
- Menampilkan avatar.
- Menampilkan nama.
- Menampilkan timestamp.
- Unread indicator.
- Read receipt (optional).
- Realtime messaging.

## 9.2 Event Group Chat

Setiap event memiliki grup chat tersendiri.

Rules:
- Creator otomatis menjadi anggota grup.
- Pengguna yang join public event otomatis menjadi anggota grup.
- Pengguna yang di-approve pada private event otomatis menjadi anggota grup.
- Hanya anggota event yang dapat membaca dan mengirim pesan.
- Pengguna yang meninggalkan event kehilangan akses sesuai kebijakan aplikasi.
- Riwayat pesan tetap disimpan.

Informasi grup:
- Nama grup menggunakan nama event.
- Foto grup menggunakan cover event.
- Daftar anggota.
- Jumlah anggota.

## 9.3 Chat Interface

**Sidebar**
- Personal conversations.
- Event group conversations.
- Avatar.
- Nama.
- Last message.
- Last message time.
- Unread count.

**Chat Window**
- Header.
- Conversation information.
- Message history.
- Message input.
- Send button.

## 9.4 Realtime Requirements

Gunakan WebSocket.

Alur:

1. User membuka chat.
2. Frontend membuat WebSocket connection.
3. Backend melakukan authentication.
4. Backend memvalidasi conversation membership.
5. User bergabung ke conversation room.
6. Pesan disimpan ke PostgreSQL.
7. Backend broadcast pesan ke anggota yang berhak menerima.

Implementasikan:
- Connection management.
- Heartbeat.
- Reconnection.
- Authorization.
- Persistent message storage.

Jika menggunakan beberapa instance backend, tambahkan Redis Pub/Sub.

---

# 10. NOTIFICATIONS

Buat sistem notifikasi.

Jenis notifikasi:

- Komentar baru pada spot milik pengguna.
- Join request event.
- Join request approved.
- Join request rejected.
- Event updated.
- Event cancelled.
- Pesan baru.
- Notifikasi aktivitas event.

Fitur:
- Notification bell.
- Unread count.
- Notification list.
- Mark as read.
- Mark all as read.
- Realtime notification jika memungkinkan.

---

# 11. DATABASE DESIGN

Gunakan PostgreSQL.

Seluruh tabel utama menggunakan UUID sebagai primary key.

Gunakan TIMESTAMPTZ untuk timestamp.

## 11.1 users

| Column | Type |
|---|---|
| id | UUID |
| name | VARCHAR |
| username | VARCHAR UNIQUE |
| email | VARCHAR UNIQUE |
| password_hash | TEXT |
| avatar_key | TEXT |
| bio | TEXT |
| location | VARCHAR |
| created_at | TIMESTAMPTZ |
| updated_at | TIMESTAMPTZ |

## 11.2 fishing_spots

| Column | Type |
|---|---|
| id | UUID |
| user_id | UUID FK |
| name | VARCHAR |
| description | TEXT |
| address | TEXT |
| latitude | DECIMAL |
| longitude | DECIMAL |
| created_at | TIMESTAMPTZ |
| updated_at | TIMESTAMPTZ |

## 11.3 fishing_spot_photos

| Column | Type |
|---|---|
| id | UUID |
| spot_id | UUID FK |
| object_key | TEXT |
| is_primary | BOOLEAN |
| created_at | TIMESTAMPTZ |

## 11.4 fish_types

| Column | Type |
|---|---|
| id | UUID |
| name | VARCHAR UNIQUE |

## 11.5 fishing_spot_fish_types

| Column | Type |
|---|---|
| spot_id | UUID FK |
| fish_type_id | UUID FK |

## 11.6 fishing_techniques

| Column | Type |
|---|---|
| id | UUID |
| name | VARCHAR UNIQUE |

## 11.7 spot_techniques

| Column | Type |
|---|---|
| spot_id | UUID FK |
| technique_id | UUID FK |

## 11.8 recommended_baits

| Column | Type |
|---|---|
| id | UUID |
| spot_id | UUID FK |
| name | VARCHAR |
| description | TEXT |

## 11.9 spot_comments

| Column | Type |
|---|---|
| id | UUID |
| spot_id | UUID FK |
| user_id | UUID FK |
| comment | TEXT |
| created_at | TIMESTAMPTZ |
| updated_at | TIMESTAMPTZ |

## 11.10 events

| Column | Type |
|---|---|
| id | UUID |
| creator_id | UUID FK |
| name | VARCHAR |
| description | TEXT |
| cover_image_key | TEXT |
| fishing_spot_id | UUID nullable |
| destination_name | VARCHAR |
| destination_address | TEXT |
| latitude | DECIMAL |
| longitude | DECIMAL |
| event_date | DATE |
| start_time | TIME |
| end_time | TIME nullable |
| meeting_point | TEXT |
| difficulty | VARCHAR |
| has_parking_fee | BOOLEAN |
| parking_fee | DECIMAL nullable |
| requires_boat | BOOLEAN |
| boat_rental_fee | DECIMAL nullable |
| additional_cost | TEXT |
| visibility | VARCHAR |
| max_participants | INTEGER nullable |
| registration_deadline | TIMESTAMPTZ nullable |
| status | VARCHAR |
| created_at | TIMESTAMPTZ |
| updated_at | TIMESTAMPTZ |

## 11.11 event_participants

| Column | Type |
|---|---|
| id | UUID |
| event_id | UUID FK |
| user_id | UUID FK |
| role | VARCHAR |
| status | VARCHAR |
| joined_at | TIMESTAMPTZ |

Tambahkan unique constraint untuk event_id dan user_id agar pengguna tidak memiliki membership ganda.

## 11.12 conversations

| Column | Type |
|---|---|
| id | UUID |
| type | VARCHAR |
| event_id | UUID nullable |
| created_at | TIMESTAMPTZ |

Type:
- direct.
- group.

Event ID harus unique untuk grup chat event.

## 11.13 conversation_participants

| Column | Type |
|---|---|
| id | UUID |
| conversation_id | UUID FK |
| user_id | UUID FK |
| joined_at | TIMESTAMPTZ |

## 11.14 messages

| Column | Type |
|---|---|
| id | UUID |
| conversation_id | UUID FK |
| sender_id | UUID FK |
| message | TEXT |
| created_at | TIMESTAMPTZ |
| updated_at | TIMESTAMPTZ |

Tambahkan tabel message_reads jika membutuhkan read receipt per pengguna.

## 11.15 notifications

| Column | Type |
|---|---|
| id | UUID |
| user_id | UUID FK |
| type | VARCHAR |
| title | VARCHAR |
| message | TEXT |
| reference_id | UUID nullable |
| is_read | BOOLEAN |
| created_at | TIMESTAMPTZ |

## 11.16 Additional Tables

Jika diperlukan, buat tabel:

- sessions.
- refresh_tokens.
- message_reads.
- event_invitations.
- user_follows.
- fishing_activity_logs.

Tabel tambahan tersebut bukan prioritas MVP kecuali diperlukan oleh implementasi inti.

---

# 12. IMPORTANT BUSINESS RULES

Bagian ini sangat penting.

AI Agent harus memastikan semua business rules diterapkan pada backend, bukan hanya frontend.

## 12.1 Event Join Transaction

Ketika pengguna berhasil bergabung ke event:

1. Validasi pengguna sudah login.
2. Validasi event tersedia.
3. Validasi event belum dibatalkan.
4. Validasi registration deadline.
5. Validasi kapasitas.
6. Validasi pengguna belum menjadi peserta.
7. Buat event_participants.
8. Tambahkan pengguna ke conversation_participants.
9. Commit transaction.

Seluruh proses harus menggunakan database transaction.

Jika salah satu operasi gagal, rollback seluruh transaction.

Untuk private event, proses membership baru dilakukan setelah creator menyetujui request.

## 12.2 Prevent Race Condition

Jika terdapat satu slot tersisa dan dua pengguna mencoba join secara bersamaan:

- Gunakan database transaction.
- Gunakan row-level locking jika diperlukan.
- Validasi kapasitas di dalam transaction.
- Pastikan jumlah peserta tidak melebihi batas.

## 12.3 Chat Authorization

Pengguna hanya boleh mengakses conversation jika memiliki membership yang valid.

Jangan mempercayai conversation ID yang dikirim frontend.

Validasi akses pada setiap permintaan REST API maupun WebSocket.

## 12.4 Ownership

- Pengguna hanya dapat mengedit spot miliknya.
- Pengguna hanya dapat menghapus komentar miliknya.
- Creator hanya dapat mengelola event miliknya.
- Pengguna hanya dapat mengelola profil sendiri.

---

# 13. REST API DESIGN

Base URL:

```text
/api/v1
```

Gunakan JSON untuk request dan response.

## Authentication

```text
POST   /auth/register
POST   /auth/login
POST   /auth/logout
POST   /auth/refresh
GET    /auth/me
```

## Users

```text
GET    /users/:username
PATCH  /users/me
POST   /users/me/avatar
```

## Fishing Spots

```text
GET    /spots
POST   /spots
GET    /spots/:id
PATCH  /spots/:id
DELETE /spots/:id
POST   /spots/:id/photos
DELETE /spots/:id/photos/:photoId
```

## Comments

```text
GET    /spots/:id/comments
POST   /spots/:id/comments
DELETE /comments/:id
```

## Events

```text
GET    /events
POST   /events
GET    /events/:id
PATCH  /events/:id
DELETE /events/:id

POST   /events/:id/join
POST   /events/:id/leave

GET    /events/:id/requests
POST   /events/:id/requests/:userId/approve
POST   /events/:id/requests/:userId/reject
```

## Chat

```text
GET    /conversations
POST   /conversations
GET    /conversations/:id/messages
POST   /conversations/:id/messages
```

WebSocket:

```text
/ws
```

## Notifications

```text
GET    /notifications
PATCH  /notifications/:id/read
POST   /notifications/read-all
```

Gunakan pagination untuk endpoint listing.

Contoh:

```text
GET /api/v1/spots?page=1&limit=20
```

## Standard Response

Success:

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": {},
  "meta": {}
}
```

Error:

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": []
}
```

Gunakan HTTP status code yang sesuai.

---

# 14. MINIO STORAGE

Gunakan bucket:

```text
avatars
fishing-spots
event-covers
```

Contoh object key:

```text
avatars/{user_id}/{uuid}.webp

fishing-spots/{spot_id}/{uuid}.webp

event-covers/{event_id}/{uuid}.webp
```

### Upload Flow

```text
Frontend
   |
   v
Go Backend
   |
   v
Validate File
   |
   v
Upload to MinIO
   |
   v
Save Object Key to PostgreSQL
   |
   v
Return File Information
```

Rules:
- Validasi MIME type.
- Validasi file size.
- Gunakan UUID.
- Jangan mempercayai file extension.
- Simpan object key ke database.
- Jangan menyimpan binary image ke PostgreSQL.
- Gunakan private bucket sebagai default.
- Gunakan presigned URL atau backend image proxy untuk akses file private.

---

# 15. AUTHENTICATION & SECURITY

Terapkan security best practices.

### Password
- Gunakan Argon2id atau bcrypt.
- Jangan simpan plain password.
- Validasi password strength.

### Session
- Gunakan secure cookies.
- HttpOnly.
- Secure pada production.
- SameSite sesuai kebutuhan.
- Implementasikan refresh token rotation jika menggunakan refresh token.

### API Security
- Authentication middleware.
- Authorization middleware.
- Rate limiting.
- Input validation.
- SQL injection protection.
- XSS protection.
- CSRF protection sesuai mekanisme authentication.
- CORS configuration.
- HTTPS pada production.

### File Upload Security
- MIME validation.
- File size validation.
- Random object key.
- Jangan izinkan executable files.
- Batasi akses storage.

### Secrets
- Jangan commit secret.
- Jangan expose MinIO credentials ke frontend.
- Jangan expose database credentials ke frontend.
- Gunakan environment variables.

---

# 16. UI/UX DESIGN GUIDELINES

## 16.1 Brand Identity

Nama: Teman Mancing.

Tagline: Cari Spot. Cari Teman. Gas Mancing!

Design direction:
- Modern.
- Clean.
- Friendly.
- Outdoor adventure.
- Community-oriented.

## 16.2 Color Palette

| Purpose | Color |
|---|---|
| Primary | #123B50 |
| Secondary | #35A7B5 |
| Accent | #7BAE7F |
| Background | #F5F7F7 |
| Surface | #FFFFFF |
| Text | #172B36 |

## 16.3 Typography

Gunakan:
- Inter.
- Plus Jakarta Sans.
- Manrope.

## 16.4 Navigation

Navbar:

1. Home.
2. Spot Mancing.
3. Janjian.
4. Chat.
5. Profile.

Tambahkan:
- Notification bell.
- User avatar.
- Login dan register untuk guest.

## 16.5 Responsive Design

Desktop:
- Navbar.
- Main content.
- Sidebar jika diperlukan.
- Grid cards.

Mobile:
- Mobile navigation.
- Bottom navigation jika sesuai.
- Responsive cards.
- Mobile-friendly forms.
- Chat layout yang nyaman.

## 16.6 UI States

Setiap halaman harus memiliki:
- Loading state.
- Empty state.
- Error state.
- Success feedback.
- Form validation feedback.

Gunakan toast notification untuk feedback pengguna.

---

# 17. APPLICATION ROUTES

| Route | Description |
|---|---|
| `/` | Landing page |
| `/login` | Login |
| `/register` | Sign up |
| `/spots` | Spot listing |
| `/spots/create` | Create spot |
| `/spots/:id` | Spot detail |
| `/spots/:id/edit` | Edit spot |
| `/events` | Event listing |
| `/events/create` | Create event |
| `/events/:id` | Event detail |
| `/events/:id/edit` | Edit event |
| `/events/:id/manage` | Manage event |
| `/chat` | Chat inbox |
| `/chat/:id` | Chat detail |
| `/notifications` | Notifications |
| `/profile/:username` | User profile |
| `/settings/profile` | Edit profile |

---

# 18. DOCKER & DEVELOPMENT ENVIRONMENT

Gunakan Docker Compose.

Service minimal:

```text
frontend
backend
postgres
minio
```

Redis bersifat optional untuk tahap awal.

Contoh development architecture:

```text
Frontend:  localhost:3000
Backend:   localhost:8080
Postgres:  localhost:5432
MinIO API: localhost:9000
MinIO UI:  localhost:9001
```

Ketentuan:
- Persistent volume untuk PostgreSQL.
- Persistent volume untuk MinIO.
- Health check.
- Environment variables.
- Development dan production configuration dipisahkan.
- Gunakan Docker network internal untuk komunikasi antar service.

Buat file:

```text
infrastructure/docker/docker-compose.yml
```

Pastikan seluruh service utama dapat dijalankan dengan Docker Compose.

---

# 19. ENVIRONMENT VARIABLES

Backend `.env.example`:

```env
APP_ENV=development
APP_PORT=8080

DATABASE_URL=postgres://user:password@postgres:5432/teman_mancing?sslmode=disable

JWT_SECRET=change-me

MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=change-me
MINIO_USE_SSL=false

MINIO_BUCKET_AVATARS=avatars
MINIO_BUCKET_SPOTS=fishing-spots
MINIO_BUCKET_EVENTS=event-covers

REDIS_URL=redis://redis:6379
```

Frontend `.env.example`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws
```

Jangan pernah memasukkan credential production ke repository.

---

# 20. TESTING STRATEGY

Gunakan testing di seluruh layer.

## Backend

- Unit testing.
- Repository integration testing.
- API testing.
- Authentication testing.
- Authorization testing.
- WebSocket testing.
- MinIO integration testing.

## Frontend

- Component testing.
- Form validation testing.
- API integration testing.
- Responsive testing.
- End-to-end testing.

## Critical Test Cases

Pastikan skenario berikut diuji:

1. User tidak dapat join event yang sudah penuh.
2. User tidak dapat join event yang sama dua kali.
3. Private event membutuhkan approval.
4. User yang bukan anggota tidak dapat mengakses grup chat.
5. User tidak dapat mengedit spot milik orang lain.
6. User tidak dapat menghapus komentar milik orang lain.
7. File upload yang tidak valid harus ditolak.
8. Pesan tetap tersimpan setelah user disconnect.
9. Tidak terjadi overbooking ketika dua pengguna mengambil slot terakhir.
10. User yang logout tidak dapat menggunakan session yang sudah dicabut.
11. Pengguna yang meninggalkan event tidak dapat mengirim pesan ke grup event.

---

# 21. DEVELOPMENT ROADMAP

Kerjakan project secara bertahap.

Jangan mengimplementasikan seluruh fitur sekaligus tanpa testing.

## Phase 1 — Project Foundation

- [ ] Setup monorepo.
- [ ] Setup backend Golang.
- [ ] Setup frontend Next.js.
- [ ] Setup PostgreSQL.
- [ ] Setup MinIO.
- [ ] Setup Docker Compose.
- [ ] Setup database migration.
- [ ] Setup basic architecture.
- [ ] Setup logging dan error handling.

## Phase 2 — Authentication

- [ ] Register.
- [ ] Login.
- [ ] Logout.
- [ ] Session management.
- [ ] User profile.
- [ ] Avatar upload.

## Phase 3 — Fishing Spots

- [ ] Spot listing.
- [ ] Spot detail.
- [ ] Create spot.
- [ ] Edit spot.
- [ ] Delete spot.
- [ ] Multiple image upload.
- [ ] Map integration.
- [ ] Fish types.
- [ ] Fishing techniques.
- [ ] Recommended baits.
- [ ] Comment system.

## Phase 4 — Fishing Events

- [ ] Event listing.
- [ ] Create event.
- [ ] Edit event.
- [ ] Event detail.
- [ ] Public event join.
- [ ] Private event request.
- [ ] Approval system.
- [ ] Participant management.
- [ ] Event cancellation.
- [ ] Capacity validation.

## Phase 5 — Chat

- [ ] Conversation management.
- [ ] Personal chat.
- [ ] WebSocket.
- [ ] Message persistence.
- [ ] Event group chat.
- [ ] Automatic group membership.
- [ ] Read status.
- [ ] Reconnection handling.

## Phase 6 — Notifications

- [ ] Notification system.
- [ ] Event notifications.
- [ ] Comment notifications.
- [ ] Chat notifications.
- [ ] Realtime notification.
- [ ] Read/unread status.

## Phase 7 — Finalization

- [ ] Responsive optimization.
- [ ] Security review.
- [ ] Backend testing.
- [ ] Frontend testing.
- [ ] Integration testing.
- [ ] Performance optimization.
- [ ] Documentation.
- [ ] Production deployment preparation.

---

# 22. ACCEPTANCE CRITERIA

Project dianggap memenuhi MVP jika seluruh fitur utama berfungsi.

## Authentication

- [ ] Register berfungsi.
- [ ] Login berfungsi.
- [ ] Logout berfungsi.
- [ ] User dapat mengedit profil.

## Fishing Spots

- [ ] User dapat melihat spot.
- [ ] User dapat membuat spot.
- [ ] User dapat mengedit spot miliknya.
- [ ] User dapat menghapus spot miliknya.
- [ ] User dapat upload foto.
- [ ] Peta menampilkan lokasi.
- [ ] User dapat berkomentar.

## Events

- [ ] User dapat membuat event.
- [ ] User dapat melihat event.
- [ ] User dapat join public event.
- [ ] User dapat request private event.
- [ ] Creator dapat approve dan reject.
- [ ] Kapasitas peserta berfungsi.
- [ ] Creator dapat mengelola event.

## Chat

- [ ] Personal chat berfungsi.
- [ ] Event group chat berfungsi.
- [ ] Peserta otomatis masuk grup chat.
- [ ] Pesan tersimpan di database.
- [ ] Chat berjalan realtime.
- [ ] Authorization grup chat berfungsi.

## Notifications

- [ ] Notifikasi aktivitas tersedia.
- [ ] Unread count berfungsi.
- [ ] Mark as read berfungsi.

## Infrastructure

- [ ] PostgreSQL berjalan.
- [ ] MinIO berjalan.
- [ ] Backend dapat terhubung ke database.
- [ ] Backend dapat upload file ke MinIO.
- [ ] Frontend terhubung ke backend.
- [ ] Docker Compose berjalan.

---

# 23. FUTURE FEATURES (NOT MVP)

Jangan implementasikan fitur berikut sebelum seluruh MVP selesai.

Fitur pengembangan berikutnya:

### Fishing Activity / Logbook
Pengguna dapat membagikan hasil tangkapan:
- Foto ikan.
- Jenis ikan.
- Berat ikan.
- Ukuran ikan.
- Spot lokasi.
- Teknik memancing.
- Umpan yang digunakan.
- Tanggal tangkapan.

### Spot Rating
Pengguna dapat memberikan rating:
- Akses lokasi.
- Kondisi spot.
- Kualitas pengalaman.
- Fasilitas.

### Follow System
Pengguna dapat mengikuti pemancing lain.

### Event Invitation
Creator dapat mengundang pengguna tertentu ke event.

### Fishing Community
Pengguna dapat membuat komunitas berdasarkan wilayah atau minat memancing.

### Admin Dashboard
Admin dapat:
- Mengelola pengguna.
- Memoderasi spot.
- Memoderasi komentar.
- Mengelola laporan.
- Mengelola event.

---

# 24. FINAL INSTRUCTIONS

Bangun **Teman Mancing** sebagai full-stack web application.

Gunakan technology stack:

| Component | Technology |
|---|---|
| Backend | Golang |
| Framework | Gin |
| Architecture | Modular Monolith |
| Design Pattern | Clean Architecture |
| Frontend | Next.js |
| Frontend Language | TypeScript |
| UI Framework | Tailwind CSS |
| UI Components | Shadcn/UI |
| Database | PostgreSQL |
| ORM | GORM |
| Migration | golang-migrate |
| Object Storage | MinIO |
| Realtime | WebSocket |
| Cache / PubSub | Redis (optional) |
| Maps | Leaflet + OpenStreetMap |
| Container | Docker Compose |

## Priority Rules

1. Backend harus menggunakan Golang.
2. Database harus menggunakan PostgreSQL.
3. File storage harus menggunakan MinIO.
4. Frontend menggunakan Next.js sebagai rekomendasi utama.
5. Gunakan Modular Monolith, bukan Microservices.
6. Gunakan Clean Architecture.
7. Seluruh fitur harus terhubung dengan backend.
8. Terapkan authorization di backend.
9. Pastikan seluruh fitur diuji.
10. Jangan mengimplementasikan future features sebelum MVP selesai.

## How to Start

Ketika mulai mengerjakan project:

**Step 1:** Inspect repository dan baca README.

**Step 2:** Buat implementation plan.

**Step 3:** Setup project structure.

**Step 4:** Setup Docker Compose, PostgreSQL, dan MinIO.

**Step 5:** Setup backend Golang dengan Clean Architecture.

**Step 6:** Setup database migration.

**Step 7:** Implementasikan fitur sesuai roadmap.

**Step 8:** Integrasikan frontend Next.js.

**Step 9:** Testing seluruh fitur.

**Step 10:** Update dokumentasi dan laporkan status project.

Jangan melompat ke tahap berikutnya jika tahap sebelumnya masih memiliki error kritis.

Jika pekerjaan terlalu besar untuk satu sesi, simpan progress dalam dokumentasi agar dapat dilanjutkan pada sesi berikutnya.

---

# PROJECT IDENTITY

**Project Name:** Teman Mancing

**Product Category:** Fishing Community Platform

**Tagline:** Cari Spot. Cari Teman. Gas Mancing!

**Backend:** Golang

**Database:** PostgreSQL

**Object Storage:** MinIO

**Frontend:** Next.js

**Architecture:** Modular Monolith + Clean Architecture

**AI Development Agent:** Hermes Agent powered by Gemini

**Main Goal:** Membangun platform komunitas memancing yang memungkinkan pengguna berbagi spot, membuat event mancing, dan berkomunikasi dengan pemancing lainnya.

---

**END OF PROJECT REQUIREMENTS**

AI Agent: Use this document as the primary source of truth when developing Teman Mancing. Prioritize a working, secure, maintainable MVP over unnecessary complexity.
