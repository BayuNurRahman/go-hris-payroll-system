# Proyek Ujian Akhir: Rancang Bangun REST API HRIS & Payroll System

## 📋 Tujuan Proyek
Bangun backend REST API independen untuk sistem HRIS dan Payroll dengan arsitektur Clean Architecture, keamanan RBAC, proteksi IDOR, serta transaksi database ACID untuk konsistensi data payroll.

## 🏗️ Struktur Folder (Clean Architecture)
hris-payroll/
├── config/
│   └── database.go           # Koneksi PostgreSQL & konfigurasi GORM
├── domain/                   # Entity bisnis murni dan interface kontrak
│   ├── auth.go               # TokenBlacklist, DTO login/register, interface auth
│   ├── employee.go           # Employee, DepartmentBudget, interface employee
│   ├── payroll.go            # Payroll, DTO payroll, interface payroll
│   └── leave.go              # LeaveRequest, DTO leave request, interface leave
├── internal/
│   ├── auth/
│   │   ├── delivery/
│   │   │   └── auth_handler.go      # Handler login, register, logout
│   │   └── usecase/
│   │       └── auth_usecase.go      # Validasi email, hashing bcrypt, JWT, blacklist token
│   ├── employee/
│   │   ├── repository/
│   │   │   └── employee_repository.go
│   │   ├── usecase/
│   │   │   └── employee_usecase.go   # Validasi domain @company.co.id
│   │   └── delivery/
│   │       └── employee_handler.go
│   ├── payroll/
│   │   ├── repository/
│   │   │   └── payroll_repository.go # Transaksi ACID + FOR UPDATE
│   │   ├── usecase/
│   │   │   └── payroll_usecase.go    # Hitung gaji total dan validasi anggaran
│   │   └── delivery/
│   │       └── payroll_handler.go    # POST /api/v1/payroll/process
│   └── leave/
│       ├── repository/
│       │   └── leave_repository.go
│       ├── usecase/
│       │   └── leave_usecase.go      # Query defensif untuk mitigasi IDOR
│       └── delivery/
│           └── leave_handler.go      # GET/POST /api/v1/leaves
├── middleware/
│   ├── auth_middleware.go    # Verifikasi JWT & blacklist token
│   └── rbac_middleware.go    # Middleware RequireRole (HRD/EMPLOYEE)
├── internal/bootstrap/
│   └── seed.go               # Seed data awal department dan user HRD
├── .env
├── main.go                   # Routing utama dan dependency injection
└── plan.md                   # Catatan proyek dan checklist

## ✅ Fitur yang Harus Diimplementasi

### 1. Account Management & Authentication
- HRD dapat mendaftarkan karyawan dengan field email, password, dan role.
- Role hanya boleh HRD atau EMPLOYEE.
- Email harus menggunakan domain @company.co.id.
- Password dienkripsi menggunakan bcrypt.
- Login menghasilkan JWT.
- Logout menambahkan token ke blacklist database.
- Endpoint internal dilindungi JWT dan middleware RBAC.

### 2. Payroll Management
- HRD dapat memproses payroll batch via POST /api/v1/payroll/process.
- Payload berisi array object dengan EmployeeID dan Bonus.
- Proses memakai transaksi database (db.Transaction).
- Mengambil base salary dari employee.
- Total Paid = Base Salary + Bonus.
- Budget departemen dikurangi sesuai kebutuhan.
- Pakai row-level locking FOR UPDATE pada budget departemen.
- Jika budget tidak mencukupi, rollback transaksi dan kirim error terstruktur.

### 3. Leave Management & IDOR Protection
- EMPLOYEE dapat mengajukan cuti via POST /api/v1/leaves.
- EMPLOYEE dapat melihat data cutinya sendiri via GET /api/v1/leaves/:id.
- HRD dapat melihat semua data cuti.
- Employee tidak boleh melihat/ubah leave milik orang lain.

### 4. Code Quality & Clean Architecture
- Pemisahan layer domain, repository, usecase, delivery.
- DTO dipakai untuk request/response.
- Business logic tidak bercampur dengan framework Gin.
- Struktur kode bersih dan modular.

### 🔄 Update Implementasi Terbaru
- Seed data sekarang memastikan department budget dibuat lebih dulu sebelum akun HRD dibuat, agar relasi foreign key tetap aman.
- Department menggunakan ID eksplisit (1 untuk Engineering, 2 untuk HR) untuk menjaga konsistensi data awal.
- Konfigurasi database saat startup mengacu pada file .env, termasuk port PostgreSQL yang digunakan.

## ✅ Checklist Progress
- [x] Struktur folder dibuat
- [x] Konfigurasi database dan .env
- [x] Auth flow (register/login/logout)
- [x] Middleware JWT dan RBAC
- [x] Payroll flow dengan transaksi
- [x] Leave flow dengan proteksi IDOR
- [x] Seed data awal
- [x] Uji integrasi endpoint dengan PostgreSQL
- [x] Dokumentasi API dan contoh request