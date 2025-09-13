
# STMADB Portal Go Backend

Backend API untuk STMADB Portal, dibangun dengan Go, Gin, dan Prisma ORM. Mendukung autentikasi JWT, manajemen pengguna, dan berbagai modul akademik.

## Fitur Utama
- Autentikasi JWT (login, refresh token, ganti password)
- Manajemen pengguna (admin, guru, siswa)
- Health check endpoint
- Dokumentasi API Swagger
- Modular: mendukung pengembangan fitur akademik lainnya

## Struktur Direktori

- `cmd/api/`         — Entry point utama API
- `cmd/seeder/`      — Seeder database
- `internal/`        — Kode aplikasi (handler, service, middleware, router, config)
- `prisma/`          — Prisma schema, migrasi, dan client Go
- `docs/`            — Dokumentasi Swagger

## Instalasi & Menjalankan

1. **Clone repo & install dependencies**
	```bash
	git clone <repo-url>
	cd stmadb-portal-go
	go mod tidy
	```

2. **Konfigurasi environment**
	- Salin `.env.example` (jika ada) atau buat `.env` sesuai kebutuhan.
	- Contoh isi `.env`:
	  ```env
	  PORT=3000
	  DATABASE_URL="mysql://root:@localhost:3306/stmadb_portal_go"
	  JWT_SECRET=rahasia-sekali-jangan-disebar
	  JWT_REFRESH_SECRET=rahasia-refresh-token
	  ```

3. **Generate Prisma Client**
	```bash
	go run github.com/steebchen/prisma-client-go generate
	```

4. **Migrasi Database**
	```bash
	npx prisma migrate dev --name init
	```

5. **Seed Database**
	```bash
	go run ./cmd/seeder/main.go
	```

6. **Jalankan API**
	```bash
	go run ./cmd/api/main.go
	```

## Dokumentasi API

Swagger UI tersedia di: [http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html)

## Contoh Endpoint

- `POST /api/v1/auth/login` — Login user
- `POST /api/v1/auth/refresh` — Refresh token
- `GET /api/v1/auth/profile` — Profil user (butuh JWT)
- `PUT /api/v1/auth/change-password` — Ganti password (butuh JWT)
- `GET /api/v1/health` — Health check

## Lisensi

MIT
