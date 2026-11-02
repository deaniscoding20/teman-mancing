# Teman Mancing — REST API Specifications

## 1. API Overview & Conventions

- **Base URL**: `/api/v1`
- **Content Type**: `application/json; charset=utf-8`
- **Authentication**: Bearer Token via Header `Authorization: Bearer <jwt_token>` atau Secure HttpOnly Cookie.
- **WebSocket Endpoint**: `/ws`

---

## 2. Standard JSON Response Envelope

Seluruh endpoint REST API mematuhi format response standar konsisten berikut:

### 2.1 Success Response

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 20,
    "total_items": 150,
    "total_pages": 8
  }
}
```

*Keterangan:* `meta` bersifat opsional dan wajib disertakan pada endpoint yang memiliki pagination.

### 2.2 Error Response

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "message": "email is already registered"
    }
  ]
}
```

### 2.3 Standard HTTP Status Codes
- `200 OK`: Permintaan berhasil diproses dan mengembalikan payload data.
- `201 Created`: Entitas baru berhasil dibuat (e.g. register, create spot, create event).
- `400 Bad Request`: Payload request tidak valid, parameter query salah, atau validasi form gagal.
- `401 Unauthorized`: Token JWT tidak ada, expired, atau tanda tangan tidak valid.
- `403 Forbidden`: Pengguna tidak memiliki hak akses terhadap resource ini (bukan pemilik spot/event).
- `404 Not Found`: Resource yang diminta tidak ditemukan di database.
- `409 Conflict`: Terjadi konflik data unik (e.g. email sudah terdaftar, membership event sudah ada).
- `422 Unprocessable Entity`: Logika bisnis tidak dapat diproses (e.g. kuota event penuh).
- `500 Internal Server Error`: Kesalahan sistem internal yang tidak tertangani.

---

## 3. Endpoints Catalog

### 3.1 Health & System

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/v1/health` | Public | Status kesehatan API, koneksi PostgreSQL, dan MinIO |
| `GET` | `/api/v1/version` | Public | Informasi versi build dan metadata backend |

---

### 3.2 Authentication (`/api/v1/auth`)

| Method | Path | Auth | Request Body | Description |
|---|---|---|---|---|
| `POST` | `/auth/register` | Public | `{ "name", "username", "email", "password" }` | Registrasi user baru |
| `POST` | `/auth/login` | Public | `{ "email", "password" }` | Login & dapatkan token |
| `POST` | `/auth/logout` | User | - | Logout & revoke session |
| `POST` | `/auth/refresh` | User | - | Refresh access token |
| `GET` | `/auth/me` | User | - | Dapatkan profil user saat ini |

---

### 3.3 Users (`/api/v1/users`)

| Method | Path | Auth | Request Body | Description |
|---|---|---|---|---|
| `GET` | `/users/:username` | Public | - | Profil publik pemancing berdasarkan username |
| `PATCH` | `/users/me` | User | `{ "name", "bio", "location" }` | Update profil sendiri |
| `POST` | `/users/me/avatar` | User | `multipart/form-data (file)` | Upload foto avatar profil ke MinIO |

---

### 3.4 Fishing Spots (`/api/v1/spots`)

| Method | Path | Auth | Request / Query | Description |
|---|---|---|---|---|
| `GET` | `/spots` | Public | `page, limit, q, fish_type, technique, sort` | Daftar spot dengan card & filter |
| `POST` | `/spots` | User | `{ "name", "description", "address", "latitude", "longitude", "fish_type_ids", "technique_ids", "baits" }` | Buat spot mancing baru |
| `GET` | `/spots/:id` | Public | - | Detail lengkap spot, peta, ulasan |
| `PATCH` | `/spots/:id` | Owner | `{ "name", "description", ... }` | Update spot (hanya pemilik) |
| `DELETE` | `/spots/:id` | Owner | - | Hapus spot (hanya pemilik) |
| `POST` | `/spots/:id/photos` | Owner | `multipart/form-data (photos)` | Upload multiple foto galeri ke MinIO |
| `DELETE` | `/spots/:id/photos/:photoId` | Owner | - | Hapus foto dari galeri MinIO |

---

### 3.5 Spot Comments (`/api/v1/spots/:id/comments`)

| Method | Path | Auth | Request Body | Description |
|---|---|---|---|---|
| `GET` | `/spots/:id/comments` | Public | `page, limit` | Daftar komentar pada spot |
| `POST` | `/spots/:id/comments` | User | `{ "comment": "..." }` | Tambah komentar baru |
| `DELETE` | `/comments/:id` | Owner | - | Hapus komentar sendiri |

---

### 3.6 Fishing Events / Janjian Mancing (`/api/v1/events`)

| Method | Path | Auth | Request / Query | Description |
|---|---|---|---|---|
| `GET` | `/events` | Public | `page, limit, q, date, difficulty, visibility, status` | Listing agenda janjian mancing |
| `POST` | `/events` | User | `{ "name", "description", "destination_name", "latitude", "longitude", "event_date", "start_time", "visibility", "max_participants", ... }` | Buat agenda event mancing |
| `GET` | `/events/:id` | Public | - | Detail lengkap event & daftar peserta |
| `PATCH` | `/events/:id` | Creator | `{ "name", ... }` | Update detail event |
| `DELETE` | `/events/:id` | Creator | - | Batalkan / hapus event |
| `POST` | `/events/:id/join` | User | - | Join public event langsung / kirim request join |
| `POST` | `/events/:id/leave` | Participant | - | Keluar dari kepesertaan event |
| `GET` | `/events/:id/requests` | Creator | - | Daftar permohonan join (private event) |
| `POST` | `/events/:id/requests/:userId/approve` | Creator | - | Setujui permohonan peserta |
| `POST` | `/events/:id/requests/:userId/reject` | Creator | - | Tolak permohonan peserta |

---

### 3.7 Chat & Messaging (`/api/v1/conversations`)

| Method | Path | Auth | Request Body | Description |
|---|---|---|---|---|
| `GET` | `/conversations` | User | - | Daftar percakapan aktif (DM & group event) |
| `POST` | `/conversations` | User | `{ "type": "direct", "participant_id": "uuid" }` | Inisialisasi DM dengan pemancing lain |
| `GET` | `/conversations/:id/messages` | Member | `page, limit, before` | Riwayat pesan chat |
| `POST` | `/conversations/:id/messages` | Member | `{ "message": "..." }` | Kirim pesan teks (fallback REST) |
| `GET` | `/ws` | User | `token=<jwt>` | WebSocket connection untuk realtime chat |

---

### 3.8 Notifications (`/api/v1/notifications`)

| Method | Path | Auth | Request Body | Description |
|---|---|---|---|---|
| `GET` | `/notifications` | User | `page, limit, unread_only` | Daftar notifikasi user |
| `PATCH` | `/notifications/:id/read` | Recipient | - | Tandai notifikasi sudah dibaca |
| `POST` | `/notifications/read-all` | Recipient | - | Tandai seluruh notifikasi telah dibaca |
