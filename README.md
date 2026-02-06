# Technical Test - Golang (Approval Workflow)

## Cara Menjalankan Aplikasi

### 1) Menjalankan secara lokal
1. Pastikan Go sudah terpasang (lihat versi pada go.mod).
2. Salin konfigurasi:
   - Copy `.env.example` menjadi `.env`.
3. Sesuaikan nilai di `.env` sesuai koneksi database.
4. Jalankan aplikasi:
   - `go run src/main.go`
5. Aplikasi berjalan pada port `APP_PORT` (default: 3000).

### 2) Menjalankan dengan Docker
1. Pastikan Docker & Docker Compose aktif.
2. Jalankan:
   - `docker compose up --build`
3. Aplikasi berjalan pada port yang diatur di `.env` (contoh: 8000).

## Konfigurasi Database
Konfigurasi dibaca dari `.env` atau environment variables. Berikut variable yang digunakan:

- `DB_HOST` (default: `localhost`)
- `DB_PORT` (default: `3306`)
- `DB_USER` (default: `root`)
- `DB_PASSWORD` (default: ``)
- `DB_NAME` (default: `test`)
- `DB_MIGRATE` (default: `false`)

Contoh `.env`:

```
APP_ENV=dev
APP_PORT=3000
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=test
DB_MIGRATE=false
JWT_SECRET=your-secret-key
JWT_ACCESS_EXP_MINUTES=15
JWT_REFRESH_EXP_DAYS=7
JWT_RESET_PASSWORD_EXP_MINUTES=30
JWT_VERIFY_EMAIL_EXP_MINUTES=60
```

## Daftar Endpoint API
Semua endpoint berada di prefix `/v1`.

### Auth (Public)
- `POST /v1/auth/register`
- `POST /v1/auth/login`

### Protected (JWT Required)
Gunakan header:
- `Authorization: Bearer <token>`

#### Workflows
- `POST /v1/workflows`
- `GET /v1/workflows`
- `GET /v1/workflows/:workflowId`

#### Steps
- `POST /v1/workflows/:workflowId/steps`
- `GET /v1/workflows/:workflowId/steps`
- `GET /v1/workflows/:workflowId/steps/:stepId`
- `PUT /v1/workflows/:workflowId/steps/:stepId`
- `DELETE /v1/workflows/:workflowId/steps/:stepId`

#### Requests
- `POST /v1/requests`
- `GET /v1/requests/:requestId`
- `POST /v1/requests/:requestId/approve`
- `POST /v1/requests/:requestId/reject`

## Swagger API Documentation

### Mengakses Dokumentasi Swagger

Setelah aplikasi berjalan, buka browser dan akses:

```
http://localhost:8000/swagger
```

Swagger JSON specification tersedia di:

```
http://localhost:8000/swagger.json
```

### Fitur Swagger UI

Swagger UI menyediakan:

1. **Dokumentasi Lengkap**: Setiap endpoint disertai deskripsi, parameter, request/response body, dan contoh response.
2. **Try It Out**: Fitur untuk menguji API langsung dari browser tanpa tools tambahan.
3. **Schema Definition**: Visualisasi struktur request dan response dalam bentuk yang mudah dipahami.
4. **Authentication**: Support untuk JWT Bearer token authentication. Masukkan token di interface Authorize button.

### Menggunakan Swagger untuk Testing

1. Buka Swagger UI di `http://localhost:8000/swagger`
2. Untuk protected endpoints, klik tombol **"Authorize"** dan masukkan JWT token:
   ```
   Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
   ```
3. Setelah authorized, semua protected endpoints dapat diakses
4. Klik endpoint yang ingin ditest, lalu klik **"Try it out"**
5. Masukkan request parameters atau body
6. Klik **"Execute"** untuk mengirim request
7. Response akan ditampilkan di bawah beserta status code dan response body

### Contoh Response Format

**Login Success Response:**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFsZmluZm9yd29ya0BnbWFpbC5jb20iLCJleHAiOjE3NzA1MjI5MzAsImlhdCI6MTc3MDM0MjkzMCwic3ViIjoxfQ.HxyTFhapS1M_zWH-2BnT7wyyZCUeUWDMr2oxT_U2584",
    "user": {
      "id": 1,
      "name": "alfin",
      "email": "alfinforwork@gmail.com"
    }
  }
}
```

**Error Response:**
```json
{
  "status": "error",
  "message": "invalid email or password",
  "data": null
}
```

### Regenerasi Dokumentasi Swagger

Jika menambah atau mengubah endpoint, regenerasi dokumentasi dengan command:

```bash
cd src
swag init -g main.go -o ../docs
```

Dokumentasi akan di-update secara otomatis di `docs/swagger.json` dan `docs/swagger.yaml`.

## Penjelasan Design Decision
- **Layered architecture**: pemisahan `controller`, `service`, `model`, dan `router` untuk memudahkan maintenance dan testing.
- **Service layer**: logika bisnis approval ditempatkan di `service` agar reusable dan mudah diuji.
- **GORM**: ORM untuk mempercepat pengelolaan database dan migrasi model.
- **Fiber**: framework HTTP ringan dengan performa baik.
- **JWT**: autentikasi sederhana untuk protected endpoints.
- **Environment-based config**: konfigurasi via `.env` / env vars supaya mudah di-deploy.

## Concurrency (Approve Endpoint)
- **Implementasi**: approval dijalankan di dalam database transaction dengan row-level lock (`SELECT ... FOR UPDATE`) pada data request.
- **Tujuan**: mencegah double approval dan race condition saat beberapa request approve terjadi bersamaan.
- **Catatan**: seluruh pembacaan dan update status dilakukan dalam satu transaksi agar konsisten.

## Asumsi atau Trade-off (Flow API)
- **Create Request**: selalu membuat request pada `CurrentStep = 1` dan status awal `PENDING`. Jika akumulasi `amount` sudah memenuhi `min_amount` sampai step berjalan, request dapat langsung naik level atau menjadi `APPROVED` jika tidak ada step berikutnya.
- **Approve Request**: hanya bisa dilakukan ketika status `PENDING`. Untuk approval type `API`, approval hanya terjadi jika `amount` >= `min_amount` terakumulasi sampai step berjalan; jika tidak memenuhi, status tetap `PENDING`.
- **Reject Request**: ketika di-reject, status berubah menjadi `REJECTED` dan tidak bisa di-approve kembali.
- **Approval sekali**: request yang sudah `APPROVED`/`REJECTED` akan ditolak untuk approval berikutnya.
- **Validasi utama**: mengikuti rule yang disyaratkan (workflow name wajib, step level unik per workflow, amount > 0).

## Asumsi atau Trade-off (Teknis)
- **In-memory test DB**: unit test menggunakan SQLite in-memory, lebih cepat namun berbeda dari MySQL production.
- **JWT sederhana**: tidak ada refresh token rotation atau blacklist.
- **Single service instance**: setup service dibuat langsung dari DB di router, belum memakai dependency injection container.
