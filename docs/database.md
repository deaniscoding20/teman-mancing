# Teman Mancing — Database Schema Design & Specifications

## 1. Overview
Database utama platform **Teman Mancing** menggunakan **PostgreSQL 16**.
Seluruh tabel utama menggunakan:
- **UUIDv4** sebagai Primary Key untuk skalabilitas, keamanan ID (mencegah enum attacks), dan desentralisasi pembuatan ID.
- **TIMESTAMPTZ** (`timestamp with time zone`) untuk konsistensi zona waktu seluruh wilayah Indonesia (WIB, WITA, WIT).
- **Foreign Keys** dengan relasi integritas referensial yang ketat (`ON DELETE CASCADE` atau `ON DELETE SET NULL`).
- **Composite Unique Constraints** untuk mencegah duplicate joins dan memberships.
- **B-Tree Indexes** pada foreign keys dan kolom pencarian untuk performa query optimal.

---

## 2. Entity Relationship Diagram (ERD)

```text
       +--------------------+
       |       users        |
       +---------+----------+
                 | 1
                 |
                 +-----------------------------------------+-----------------------------------------+
                 | N (creator)                             | N (author)                              | N (sender)
                 v                                         v                                         v
       +---------+----------+                    +---------+----------+                    +---------+----------+
       |       events       |<--+                |  fishing_spots     |                    |      messages      |
       +---------+----------+   |                +---------+----------+                    +---------+----------+
                 | 1            |                          | 1                                       ^
                 |              |                          |                                         | N
                 +--------------|------------+             +---------------+---------------+         |
                 |              |            |             |               |               |         |
                 | 1:1          |            v 1:N         v 1:N           v 1:N           v 1:N     |
                 v              |       +----+------+ +----+-----+   +-----+------+  +-----+-----+   |
       +---------+----------+   |       |   event   | |  photos  |   |  comments  |  |   baits   |   |
       |   conversations    |   |       |participants |    +     |   +-----+------+  +-----------+   |
       +---------+----------+   |       +-----------+ |          |         |                         |
                 | 1            |                     v          v         v                         |
                 |              |                  (MinIO)    users       users                      |
                 +--------------+                                                                    |
                 | N                                                                                 |
                 v                                                                                   |
       +---------+----------+                                                                        |
       |   conversation_    |                                                                        |
       |    participants    |                                                                        |
       +--------------------+                                                                        |
                 |                                                                                   |
                 +-----------------------------------------------------------------------------------+
```

---

## 3. Detailed Table Specifications

### 3.1 `users`
Menyimpan identitas akun pengguna dan profil publik.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID unik pengguna |
| `name` | VARCHAR(100) | NOT NULL | Nama lengkap |
| `username` | VARCHAR(50) | NOT NULL, UNIQUE | Username unik (@username) |
| `email` | VARCHAR(255) | NOT NULL, UNIQUE | Alamat email terdaftar |
| `password_hash` | TEXT | NOT NULL | Hash password (bcrypt / Argon2) |
| `avatar_key` | TEXT | NULLABLE | MinIO object key untuk foto profil |
| `bio` | TEXT | NULLABLE | Biografi singkat pengguna |
| `location` | VARCHAR(150) | NULLABLE | Asal domisili / kota pemancing |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu registrasi |
| `updated_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu pembaruan profil |

*Indexes:*
- `idx_users_username` on `username`
- `idx_users_email` on `email`

---

### 3.2 `fishing_spots`
Direktori spot pemancingan yang dibagikan oleh komunitas.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID unik spot |
| `user_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Pembuat spot |
| `name` | VARCHAR(200) | NOT NULL | Nama lokasi / spot |
| `description` | TEXT | NOT NULL | Informasi detail kondisi spot |
| `address` | TEXT | NOT NULL | Alamat atau petunjuk jalan |
| `latitude` | DECIMAL(10, 7) | NOT NULL | Titik koordinat latitude |
| `longitude` | DECIMAL(10, 7) | NOT NULL | Titik koordinat longitude |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| `updated_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu update |

*Indexes:*
- `idx_fishing_spots_user_id` on `user_id`
- `idx_fishing_spots_lat_lng` on `(latitude, longitude)`
- `idx_fishing_spots_created_at` on `created_at DESC`

---

### 3.3 `fishing_spot_photos`
Galeri foto untuk setiap spot mancing yang disimpan di MinIO.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID unik foto |
| `spot_id` | UUID | NOT NULL, REFERENCES fishing_spots(id) ON DELETE CASCADE | Spot terkait |
| `object_key` | TEXT | NOT NULL | MinIO object key |
| `is_primary` | BOOLEAN | NOT NULL, DEFAULT FALSE | Penanda cover utama spot |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu upload |

*Indexes:*
- `idx_fishing_spot_photos_spot_id` on `spot_id`

---

### 3.4 `fish_types` & `fishing_spot_fish_types`
Katalog jenis ikan yang tersedia pada spot mancing (Many-to-Many).

**`fish_types`**:
| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID unik jenis ikan |
| `name` | VARCHAR(100) | NOT NULL, UNIQUE | Nama ikan (e.g. Kakap Putih, Nila, Patin) |

**`fishing_spot_fish_types`**:
| Column | Type | Constraints | Description |
|---|---|---|---|
| `spot_id` | UUID | NOT NULL, REFERENCES fishing_spots(id) ON DELETE CASCADE | Relasi ke spot |
| `fish_type_id` | UUID | NOT NULL, REFERENCES fish_types(id) ON DELETE CASCADE | Relasi ke jenis ikan |
| - | - | PRIMARY KEY (`spot_id`, `fish_type_id`) | Composite PK |

---

### 3.5 `fishing_techniques` & `spot_techniques`
Katalog teknik memancing pada spot (Many-to-Many).

**`fishing_techniques`**:
| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID unik teknik |
| `name` | VARCHAR(100) | NOT NULL, UNIQUE | Nama teknik (e.g. Casting, Dasaran, Jigging) |

**`spot_techniques`**:
| Column | Type | Constraints | Description |
|---|---|---|---|
| `spot_id` | UUID | NOT NULL, REFERENCES fishing_spots(id) ON DELETE CASCADE | Relasi ke spot |
| `technique_id` | UUID | NOT NULL, REFERENCES fishing_techniques(id) ON DELETE CASCADE | Relasi ke teknik |
| - | - | PRIMARY KEY (`spot_id`, `technique_id`) | Composite PK |

---

### 3.6 `recommended_baits`
Rekomendasi umpan yang efektif pada spot tertentu.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID umpan |
| `spot_id` | UUID | NOT NULL, REFERENCES fishing_spots(id) ON DELETE CASCADE | Spot terkait |
| `name` | VARCHAR(150) | NOT NULL | Nama umpan (e.g. Udang Hidup, Cacing) |
| `description` | TEXT | NULLABLE | Catatan / racikan umpan |

*Indexes:*
- `idx_recommended_baits_spot_id` on `spot_id`

---

### 3.7 `spot_comments`
Komentar dan ulasan komunitas pada suatu spot mancing.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID komentar |
| `spot_id` | UUID | NOT NULL, REFERENCES fishing_spots(id) ON DELETE CASCADE | Spot yang dikomentari |
| `user_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Pembuat komentar |
| `comment` | TEXT | NOT NULL | Isi teks komentar |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu komentar dibuat |
| `updated_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu pembaruan |

*Indexes:*
- `idx_spot_comments_spot_id` on `spot_id`
- `idx_spot_comments_user_id` on `user_id`

---

### 3.8 `events` (Janjian Mancing)
Agenda kegiatan memancing bersama yang dibuat oleh pengguna.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID unik event |
| `creator_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Pembuat/host event |
| `name` | VARCHAR(200) | NOT NULL | Judul kegiatan |
| `description` | TEXT | NOT NULL | Deskripsi lengkap agenda kegiatan |
| `cover_image_key` | TEXT | NULLABLE | MinIO object key cover banner |
| `fishing_spot_id` | UUID | NULLABLE, REFERENCES fishing_spots(id) ON DELETE SET NULL | Spot referensi (jika ada) |
| `destination_name` | VARCHAR(200) | NOT NULL | Nama lokasi tujuan |
| `destination_address` | TEXT | NOT NULL | Alamat tujuan |
| `latitude` | DECIMAL(10, 7) | NOT NULL | Titik koordinat destinasi |
| `longitude` | DECIMAL(10, 7) | NOT NULL | Titik koordinat destinasi |
| `event_date` | DATE | NOT NULL | Tanggal pelaksanaan |
| `start_time` | TIME | NOT NULL | Waktu mulai |
| `end_time` | TIME | NULLABLE | Waktu perkiraan selesai |
| `meeting_point` | TEXT | NOT NULL | Titik kumpul / titik temu |
| `difficulty` | VARCHAR(50) | NOT NULL, CHECK (difficulty IN ('pemula', 'menengah', 'ahli')) | Tingkat kesulitan |
| `has_parking_fee` | BOOLEAN | NOT NULL, DEFAULT FALSE | Apakah ada biaya parkir |
| `parking_fee` | DECIMAL(12, 2) | NULLABLE | Estimasi biaya parkir |
| `requires_boat` | BOOLEAN | NOT NULL, DEFAULT FALSE | Apakah membutuhkan sewa kapal/perahu |
| `boat_rental_fee` | DECIMAL(12, 2) | NULLABLE | Estimasi sewa kapal |
| `additional_cost` | TEXT | NULLABLE | Rincian biaya lain (umpan bersama, makan) |
| `visibility` | VARCHAR(20) | NOT NULL, CHECK (visibility IN ('public', 'private')) | Jenis pendaftaran |
| `max_participants` | INTEGER | NULLABLE | Kuota maksimal peserta (null = tak terbatas) |
| `registration_deadline` | TIMESTAMPTZ | NULLABLE | Batas akhir pendaftaran |
| `status` | VARCHAR(30) | NOT NULL, DEFAULT 'open', CHECK (status IN ('open', 'closed', 'cancelled', 'completed')) | Status event |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| `updated_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu pembaruan |

*Indexes:*
- `idx_events_creator_id` on `creator_id`
- `idx_events_date_status` on `(event_date, status)`
- `idx_events_visibility` on `visibility`

---

### 3.9 `event_participants`
Daftar peserta dan permohonan join pada kegiatan event.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID unik kepesertaan |
| `event_id` | UUID | NOT NULL, REFERENCES events(id) ON DELETE CASCADE | Relasi event |
| `user_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Relasi user |
| `role` | VARCHAR(30) | NOT NULL, DEFAULT 'participant', CHECK (role IN ('creator', 'co-host', 'participant')) | Peran di event |
| `status` | VARCHAR(30) | NOT NULL, DEFAULT 'pending', CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')) | Status pendaftaran |
| `joined_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu bergabung / request |

*Unique Constraint & Indexes:*
- `UNIQUE (event_id, user_id)`: Mencegah pengguna mendaftar ganda ke event yang sama.
- `idx_event_participants_event_status` on `(event_id, status)`
- `idx_event_participants_user_id` on `user_id`

---

### 3.10 `conversations`
Entitas percakapan obrolan (Direct Message maupun Event Group Chat).

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID room percakapan |
| `type` | VARCHAR(20) | NOT NULL, CHECK (type IN ('direct', 'group')) | Tipe percakapan |
| `event_id` | UUID | NULLABLE, UNIQUE, REFERENCES events(id) ON DELETE CASCADE | Event ID (jika tipe group event chat) |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan room |

*Indexes:*
- `idx_conversations_event_id` on `event_id` WHERE event_id IS NOT NULL

---

### 3.11 `conversation_participants`
Daftar keanggotaan user dalam suatu room percakapan.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID keanggotaan room |
| `conversation_id` | UUID | NOT NULL, REFERENCES conversations(id) ON DELETE CASCADE | Room percakapan |
| `user_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Pengguna dalam room |
| `joined_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu bergabung |

*Unique Constraint & Indexes:*
- `UNIQUE (conversation_id, user_id)`: Mencegah record membership ganda.
- `idx_conv_participants_user` on `user_id`

---

### 3.12 `messages`
Penyimpanan riwayat pesan chat secara persisten.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID unik pesan |
| `conversation_id` | UUID | NOT NULL, REFERENCES conversations(id) ON DELETE CASCADE | Room tujuan pesan |
| `sender_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Pengirim pesan |
| `message` | TEXT | NOT NULL | Isi teks pesan |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Timestamp pengiriman |
| `updated_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Timestamp pembaruan |

*Indexes:*
- `idx_messages_conversation_created` on `(conversation_id, created_at ASC)`

---

### 3.13 `notifications`
Penyimpanan riwayat notifikasi in-app untuk pengguna.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() | ID unik notifikasi |
| `user_id` | UUID | NOT NULL, REFERENCES users(id) ON DELETE CASCADE | Penerima notifikasi |
| `type` | VARCHAR(50) | NOT NULL | Tipe (e.g. `spot_comment`, `join_request`, `request_approved`) |
| `title` | VARCHAR(255) | NOT NULL | Judul notifikasi |
| `message` | TEXT | NOT NULL | Isi ringkas pesan notifikasi |
| `reference_id` | UUID | NULLABLE | ID entitas terkait (spot_id, event_id, dsb) |
| `is_read` | BOOLEAN | NOT NULL, DEFAULT FALSE | Status terbaca |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu dikirimkan |

*Indexes:*
- `idx_notifications_user_read` on `(user_id, is_read, created_at DESC)`

---

## 4. Migration Strategy
Migrations dikelola menggunakan **golang-migrate** dalam direktori `backend/migrations/`:
- Pola penamaan file: `{version}_{title}.up.sql` dan `{version}_{title}.down.sql`.
- Migration dapat dieksekusi via CLI `golang-migrate`, script Docker compose, maupun programmatic runner di Go backend saat startup.
