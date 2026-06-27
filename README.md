# HRIS & Payroll API

Backend REST API sederhana untuk sistem HRIS dan Payroll berbasis Go, Gin, GORM, dan PostgreSQL.

## Fitur Utama
- Autentikasi JWT
- Registrasi karyawan oleh HRD
- Role-based access control (HRD/EMPLOYEE)
- Payroll batch processing dengan transaksi database ACID
- Leave request dengan proteksi IDOR
- Seed data awal untuk HRD dan department

## Persyaratan
- Go 1.22+
- PostgreSQL
- Git

## Instalasi
1. Pastikan PostgreSQL sudah aktif.
2. Buat database sesuai nilai di file .env.
   ```sql
   CREATE DATABASE hrisdb;
   ```
3. Sesuaikan konfigurasi di .env:
   ```env
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=secret45
   DB_NAME=hrisdb
   DB_PORT=5432
   JWT_SECRET=your-secret-key
   SEED_PASSWORD=Admin123!
   PORT=8080
   ```
4. Jalankan perintah berikut:
   ```bash
   go mod tidy
   go run .
   ```
5. Saat server berjalan, aplikasi akan otomatis:
   - migrate tabel,
   - seed department,
   - seed user HRD awal.

### Catatan Implementasi Terbaru
- Seed data sekarang membuat department budget terlebih dahulu sebelum membuat akun HRD, sehingga relasi foreign key tetap valid.
- Department dibuat dengan ID eksplisit: 1 untuk Engineering dan 2 untuk HR.
- Konfigurasi database sepenuhnya mengikuti nilai .env, termasuk host, port, username, password, dan nama database.

## Kredensial Seed HRD
- Email: hrd@company.co.id
- Password: nilai SEED_PASSWORD pada .env

## Endpoint API
Semua endpoint yang membutuhkan autentikasi menggunakan header:
```http
Authorization: Bearer <token>
```

### 1. Auth
#### Login
- Method: POST
- Path: /api/v1/auth/login
- Body:
```json
{
  "email": "hrd@company.co.id",
  "password": "Admin123!"
}
```

#### Register (HRD only)
- Method: POST
- Path: /api/v1/auth/register
- Header: Bearer token HRD
- Body:
```json
{
  "name": "Jane Doe",
  "email": "jane@company.co.id",
  "password": "Password123!",
  "role": "EMPLOYEE",
  "base_salary": 8000000,
  "department_id": 1
}
```

#### Logout
- Method: POST
- Path: /api/v1/auth/logout
- Header: Bearer token aktif

### 2. Payroll
#### Process Payroll (HRD only)
- Method: POST
- Path: /api/v1/payroll/process
- Header: Bearer token HRD
- Body:
```json
[
  {
    "employee_id": 1,
    "bonus": 500000
  }
]
```

### 3. Leave
#### Create Leave Request
- Method: POST
- Path: /api/v1/leaves
- Header: Bearer token EMPLOYEE
- Body:
```json
{
  "start_date": "2026-07-01",
  "end_date": "2026-07-03",
  "reason": "Sick leave"
}
```

#### Get Leave Detail
- Method: GET
- Path: /api/v1/leaves/:id
- Header: Bearer token EMPLOYEE atau HRD

## Catatan Keamanan
- Role HRD dapat mengakses fitur administratif.
- Role EMPLOYEE hanya dapat mengakses data cuti miliknya sendiri.
- Token logout akan diblacklist agar tidak bisa dipakai lagi.
- Payroll memakai transaksi database dan row-level locking untuk mencegah inkonsistensi data.

## Testing
Jalankan:
```bash
go test ./...
```
