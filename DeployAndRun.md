# 🚀 Panduan Deploy & Run: Teman Mancing

Dokumen ini ditujukan bagi AI Agent atau DevOps Engineer yang akan melakukan **clone repository** dan **deployment platform Teman Mancing pada VPS (Virtual Private Server)** berbasis Linux (Ubuntu / Debian).

---

## 1. Arsitektur & Port Default

Platform Teman Mancing menggunakan arsitektur **Modular Monolith**:
* **Frontend**: Next.js 14 (App Router, Tailwind CSS, TypeScript) ➔ **Port 3003** (atau 3080)
* **Backend**: Golang 1.23 (Gin Framework, Clean Architecture, WebSocket) ➔ **Port 8080**
* **Database**: PostgreSQL 16 ➔ **Port 5432**
* **Object Storage**: MinIO (S3-Compatible Storage) ➔ **Port 9000** (API) & **Port 9001** (Console Web)

---

## 2. Prasyarat Sistem pada VPS

Pastikan VPS telah terinstal tools berikut:
* **OS**: Ubuntu 22.04 / 24.04 LTS atau Debian 12
* **Git**: `git --version`
* **Docker**: Versi 24.0+ (`docker --version`)
* **Docker Compose**: Plugin compose v2 (`docker compose version`)
* **Port yang harus dibuka (Firewall / Security Group)**:
  * `3003` (Akses Web Frontend)
  * `8080` (Akses REST API & WebSocket)
  * `9001` (Akses MinIO Web Console - Opsional)
  * `80` & `443` (Jika menggunakan reverse proxy Nginx / Caddy untuk domain HTTPS)

---

## 3. Langkah Cepat Deployment (Metode Docker Compose)

Metode ini adalah **metode standar dan direkomendasikan** karena seluruh service (PostgreSQL, MinIO, auto bucket provisioning, Golang backend, dan Next.js frontend) berjalan secara terisolasi dan otomatis terhubung melalui Docker network internal.

### Langkah 1: Clone Repository
```bash
git clone https://github.com/deaniscoding20/teman-mancing.git
cd teman-mancing
```

### Langkah 2: Setup File Environment (.env)
Salin template environment di direktori docker:
```bash
cp infrastructure/docker/.env.example infrastructure/docker/.env
```

Bila diperlukan, sesuaikan nilai pada `infrastructure/docker/.env`:
```env
POSTGRES_USER=teman_mancing_user
POSTGRES_PASSWORD=ganti_password_rahasia_anda_disini
POSTGRES_DB=teman_mancing
POSTGRES_PORT=5432

MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=ganti_minio_password_disini
MINIO_API_PORT=9000
MINIO_CONSOLE_PORT=9001

MINIO_BUCKET_AVATARS=avatars
MINIO_BUCKET_SPOTS=fishing-spots
MINIO_BUCKET_EVENTS=event-covers

BACKEND_PORT=8080
FRONTEND_PORT=3003

JWT_SECRET=buat_kunci_jwt_panjang_minimal_32_karakter_acak
```

Pastikan juga menyalin file `.env` untuk backend dan frontend:
```bash
# Backend .env
cp backend/.env.example backend/.env

# Frontend .env
cp frontend/.env.example frontend/.env
```

> **Catatan untuk Frontend Production:**
> Pada `frontend/.env`, jika Anda mengakses dari IP publik VPS atau domain:
> ```env
> NEXT_PUBLIC_API_URL=http://<IP_VPS_ANDA>:8080/api/v1
> NEXT_PUBLIC_WS_URL=ws://<IP_VPS_ANDA>:8080/ws
> ```

### Langkah 3: Build & Jalankan Container
Jalankan seluruh service dalam background mode:
```bash
cd infrastructure/docker
docker compose up -d --build
```

### Langkah 4: Verifikasi Status Container
Periksa apakah seluruh service berjalan normal (`healthy` atau `running`):
```bash
docker compose ps
```

Container yang harus aktif:
* `teman_mancing_postgres` (Up, healthy)
* `teman_mancing_minio` (Up, healthy)
* `teman_mancing_minio_init` (Exited 0 / completed successfully)
* `teman_mancing_backend` (Up, healthy)
* `teman_mancing_frontend` (Up, running)

---

## 4. Eksekusi Migrasi Database (PostgreSQL)

Skema database PostgreSQL dapat dijalankan langsung ke container postgres:

```bash
# Masuk ke root repo
cd teman-mancing

# Jalankan skema awal (15 tabel + ekstensi pgcrypto UUID)
docker exec -i teman_mancing_postgres psql -U teman_mancing_user -d teman_mancing < backend/migrations/000001_init_schema.up.sql

# Jalankan seeding data master (15 jenis ikan & 9 teknik mancing Indonesia)
docker exec -i teman_mancing_postgres psql -U teman_mancing_user -d teman_mancing < backend/migrations/000002_seed_initial_data.up.sql
```

Untuk memverifikasi tabel yang telah terbentuk:
```bash
docker exec -it teman_mancing_postgres psql -U teman_mancing_user -d teman_mancing -c "\dt"
```

---

## 5. Uji Coba & Health Check

Lakukan verifikasi bahwa API dan Web telah aktif:

### 1. Test Health Check Endpoint
```bash
curl -s http://localhost:8080/api/v1/health | jq .
```
Response yang diharapkan:
```json
{
  "success": true,
  "message": "System health status",
  "data": {
    "status": "running",
    "postgres": "healthy",
    "minio": "healthy",
    "time": "2026-11-02T..."
  }
}
```

### 2. Test Frontend
```bash
curl -I http://localhost:3003/
```
Response yang diharapkan: `HTTP/1.1 200 OK`.

### 3. Akses via Browser Pengguna:
* **Web Utama**: `http://<IP_VPS>:3003` (atau domain)
* **API Version**: `http://<IP_VPS>:8080/api/v1/version`
* **MinIO Console**: `http://<IP_VPS>:9001` (Login menggunakan `MINIO_ROOT_USER` dan `MINIO_ROOT_PASSWORD`)

---

## 6. Opsi Deployment Tanpa Docker (Bare-Metal / Systemd)

Jika VPS Anda tidak menggunakan Docker dan ingin menjalankan aplikasi secara native:

### A. Persiapan Runtime
1. **Golang 1.23+**:
   ```bash
   wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   ```
2. **Node.js 20+ & NPM**:
   ```bash
   curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
   sudo apt-get install -y nodejs
   ```

### B. Menjalankan Backend (Golang)
```bash
cd backend
cp .env.example .env
go build -ldflags="-w -s" -o bin/api ./cmd/api
./bin/api
```
*(Gunakan `systemd` atau `supervisord` untuk menjaga proses tetap aktif di background).*

### C. Menjalankan Frontend (Next.js)
```bash
cd frontend
cp .env.example .env
npm install
npm run build
PORT=3080 npm start
```
*(Gunakan `pm2` untuk process manager Next.js: `pm2 start npm --name "frontend" -- run start -- -p 3080`).*

---

## 7. Setup Nginx Reverse Proxy & SSL (Domain Production)

Jika Anda ingin mengarahkan domain (contoh: `temanmancing.id`) dengan SSL gratis Let's Encrypt:

### Konfigurasi Nginx (`/etc/nginx/sites-available/teman-mancing`)
```nginx
server {
    listen 80;
    server_name temanmancing.id www.temanmancing.id;

    # Frontend Next.js (Port 3080)
    location / {
        proxy_pass http://127.0.0.1:3080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # Backend REST API (Port 8080)
    location /api/ {
        proxy_pass http://127.0.0.1:8080/api/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Realtime WebSocket (/ws)
    location /ws {
        proxy_pass http://127.0.0.1:8080/ws;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
    }
}
```

Aktifkan site dan pasang SSL:
```bash
sudo ln -s /etc/nginx/sites-available/teman-mancing /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
sudo certbot --nginx -d temanmancing.id -d www.temanmancing.id
```

---

## 8. Troubleshooting Umum

| Gejala / Error | Penyebab | Solusi |
|---|---|---|
| `Cannot connect to the Docker daemon` | Service docker belum berjalan di VPS | Jalankan `sudo systemctl start docker && sudo systemctl enable docker` |
| `port is already allocated: 3080` atau `8080` | Port bentrok dengan service lain | Cek proses yang menduduki port: `sudo lsof -i :3080` atau ganti port di `.env` |
| `MinIO bucket not found` | Bucket belum dibuat saat startup | Container `minio-create-buckets` otomatis membuat bucket `avatars`, `fishing-spots`, dan `event-covers`. Jalankan ulang: `docker compose up minio-create-buckets` |
| `Permission denied (publickey)` saat git clone | SSH Key belum terdaftar di GitHub | Gunakan HTTPS clone: `git clone https://github.com/deaniscoding20/teman-mancing.git` |
| Halaman web tidak bisa diakses dari luar | Firewall VPS memblokir port | Buka port via UFW: `sudo ufw allow 3080/tcp && sudo ufw allow 8080/tcp` |

---

## 10. CI/CD Pipeline & GitHub Actions (Auto-Deploy)

Repository ini telah dilengkapi dengan pipeline otomatis menggunakan GitHub Actions:

### Alur Kerja (Workflow)
1. **Pull Request Validation (CI)**:
   * Setiap Pull Request yang diajukan ke branch `main` akan secara otomatis menjalankan:
     * Pengujian Unit & Integrasi Backend Golang (`go test -v -race ./tests/...`)
     * Pengecekan Typecheck TypeScript (`npx tsc --noEmit`)
     * Verifikasi Build Next.js (`npm run build`)
     * Validasi sintaks `docker-compose.yml`
   * **Aturan PR**: Setiap kali membuat Pull Request baru, wajib menyertakan file `changesinfo.md` di root proyek agar agent lain / reviewer dapat memahami perubahan dan menyetujui merge.
2. **Auto-Deploy to VPS (CD)**:
   * Ketika Pull Request di-merge ke branch `main`, GitHub Actions akan otomatis melakukan deploy ke VPS melalui SSH dan mengeksekusi skrip `infrastructure/scripts/deploy.sh`.

### Setup GitHub Repository Secrets
Untuk mengaktifkan auto-deploy ke VPS, tambahkan secrets berikut di repository GitHub Anda (**Settings ➔ Secrets and variables ➔ Actions ➔ New repository secret**):

| Secret Name | Deskripsi | Contoh |
|---|---|---|
| `VPS_HOST` | IP Publik atau Hostname VPS Anda | `123.45.67.89` |
| `VPS_USERNAME` | User SSH pada VPS | `root` atau `ubuntu` |
| `VPS_SSH_KEY` | Private SSH Key (ed25519 / RSA) untuk akses ke VPS | `-----BEGIN OPENSSH PRIVATE KEY----- ...` |
| `VPS_PORT` | Port SSH VPS (Opsional, default 22) | `22` |
| `VPS_DEPLOY_PATH` | Lokasi folder project di VPS (Opsional) | `/home/ubuntu/apps/teman-mancing` |

---

## 11. Menjalankan Unit Testing di VPS
Untuk memverifikasi integritas kode backend pada server VPS:
```bash
cd backend
go test -v ./tests/...
```
Semua test case (Auth, Spots, Events, Race-condition locks, Chat WebSocket, Notifications) harus berstatus **`PASS`**.
