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
