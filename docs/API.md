# Sakutera API Documentation

**Base URL:** `http://localhost:8080/api/v1`  
**Format:** JSON  
**Encoding:** UTF-8

---

## Onboarding Flow

```
Step 1 → POST /auth/register         (kirim nomor HP + nama)
Step 2 → POST /auth/verify-otp       (verifikasi OTP via WhatsApp)
Step 3 → POST /onboarding/work-platform    (pilih pekerjaan)
Step 4 → POST /onboarding/income-source   (pilih sumber penghasilan)
```

Setelah Step 2 berhasil, simpan `token` dari response — dipakai sebagai Bearer token untuk Step 3 & 4.

---

## Auth

### Step 1 — Register

**POST** `/auth/register`

Request:
```json
{
  "phone_number": "081234567890",
  "full_name": "Budi Santoso"
}
```

Response `200 OK`:
```json
{
  "status": {
    "code": 200,
    "isSuccess": true
  },
  "message": "Kode OTP telah dikirim via WhatsApp",
  "data": {
    "session_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "phone_masked": "0812****7890",
    "message": "Kode OTP telah dikirim via WhatsApp"
  }
}
```

Response `409 Conflict` (nomor sudah terdaftar):
```json
{
  "status": {
    "code": 409,
    "isSuccess": false
  },
  "message": "nomor HP sudah terdaftar",
  "data": null
}
```

Response `400 Bad Request` (field kosong):
```json
{
  "status": {
    "code": 400,
    "isSuccess": false
  },
  "message": "request tidak valid",
  "data": "Key: 'RegisterRequest.PhoneNumber' Error:Field validation for 'PhoneNumber' failed on the 'required' tag"
}
```

---

### Step 2 — Verify OTP

**POST** `/auth/verify-otp`

Header:
```
X-Session-Token: <session_token dari Step 1>
```

Request:
```json
{
  "code": "123456"
}
```

Response `200 OK`:
```json
{
  "status": {
    "code": 200,
    "isSuccess": true
  },
  "message": "Akun berhasil diverifikasi",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "message": "Akun berhasil diverifikasi"
  }
}
```

> **Simpan `token` ini** — digunakan sebagai `Authorization: Bearer <token>` untuk semua endpoint selanjutnya.

Response `400 Bad Request` (OTP salah):
```json
{
  "status": {
    "code": 400,
    "isSuccess": false
  },
  "message": "kode OTP tidak valid",
  "data": null
}
```

Response `401 Unauthorized` (session expired):
```json
{
  "status": {
    "code": 401,
    "isSuccess": false
  },
  "message": "session tidak valid atau sudah kedaluwarsa",
  "data": null
}
```

---

## Onboarding

Semua endpoint onboarding memerlukan header:
```
Authorization: Bearer <token dari Step 2>
```

---

### Step 3a — List Pekerjaan

**GET** `/onboarding/work-categories`

Header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Response `200 OK`:
```json
{
  "status": {
    "code": 200,
    "isSuccess": true
  },
  "message": "data pekerjaan berhasil diambil",
  "data": {
    "categories": [
      {
        "work_category_id": "550e8400-e29b-41d4-a716-446655440001",
        "name": "Pengemudi Ojol",
        "platforms": [
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440010", "name": "Gojek" },
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440011", "name": "Grab" },
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440012", "name": "Maxim" },
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440013", "name": "InDriver" }
        ]
      },
      {
        "work_category_id": "550e8400-e29b-41d4-a716-446655440002",
        "name": "Kurir Online",
        "platforms": [
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440020", "name": "ShopeeFood" },
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440021", "name": "GoSend" },
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440022", "name": "J&T" },
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440023", "name": "JNE" }
        ]
      },
      {
        "work_category_id": "550e8400-e29b-41d4-a716-446655440003",
        "name": "Freelancer",
        "platforms": [
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440030", "name": "Fiverr" },
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440031", "name": "Upwork" },
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440032", "name": "Freelancer.com" }
        ]
      },
      {
        "work_category_id": "550e8400-e29b-41d4-a716-446655440004",
        "name": "Pedagang UMKM",
        "platforms": [
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440040", "name": "Tokopedia" },
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440041", "name": "Shopee" },
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440042", "name": "Warung Offline" }
        ]
      },
      {
        "work_category_id": "550e8400-e29b-41d4-a716-446655440005",
        "name": "Pekerjaan Lainnya",
        "platforms": [
          { "work_platform_id": "550e8400-e29b-41d4-a716-446655440050", "name": "Lainnya" }
        ]
      }
    ]
  }
}
```

> UUIDs di atas adalah contoh ilustrasi. UUID asli akan berbeda sesuai data yang di-seed ke database.

Response `401 Unauthorized` (token tidak ada / tidak valid):
```json
{
  "status": {
    "code": 401,
    "isSuccess": false
  },
  "message": "token diperlukan",
  "data": null
}
```

---

### Step 3b — Pilih Pekerjaan

**POST** `/onboarding/work-platform`

Header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

Request (gunakan `work_platform_id` dari response GET di atas):
```json
{
  "work_platform_id": "550e8400-e29b-41d4-a716-446655440010"
}
```

Response `200 OK`:
```json
{
  "status": {
    "code": 200,
    "isSuccess": true
  },
  "message": "pekerjaan berhasil disimpan",
  "data": {
    "message": "pekerjaan berhasil disimpan"
  }
}
```

Response `400 Bad Request` (UUID format salah):
```json
{
  "status": {
    "code": 400,
    "isSuccess": false
  },
  "message": "work_platform_id tidak valid",
  "data": null
}
```

Response `404 Not Found` (UUID tidak ada di DB):
```json
{
  "status": {
    "code": 404,
    "isSuccess": false
  },
  "message": "pekerjaan tidak ditemukan",
  "data": null
}
```

---

### Step 4a — List Sumber Penghasilan

**GET** `/onboarding/income-sources`

Header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Response `200 OK`:
```json
{
  "status": {
    "code": 200,
    "isSuccess": true
  },
  "message": "sumber penghasilan berhasil diambil",
  "data": {
    "sources": [
      {
        "type": "manual",
        "label": "Input Manual",
        "description": "Catat sendiri dalam kurang dari 10 detik. Tersedia di semua kondisi, tanpa koneksi platform.",
        "is_available": true
      },
      {
        "type": "ewallet",
        "label": "Hubungkan GoPay / OVO",
        "description": "Transaksi masuk otomatis. Kami hanya membaca data masuk — tidak bisa transfer atau tarik saldo.",
        "is_available": false
      }
    ]
  }
}
```

---

### Step 4b — Pilih Sumber Penghasilan

**POST** `/onboarding/income-source`

Header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

Request:
```json
{
  "income_source_type": "manual"
}
```

Response `200 OK`:
```json
{
  "status": {
    "code": 200,
    "isSuccess": true
  },
  "message": "sumber penghasilan berhasil disimpan",
  "data": {
    "message": "sumber penghasilan berhasil disimpan"
  }
}
```

Response `400 Bad Request` (tipe tidak valid — hanya `manual` yang diizinkan saat ini):
```json
{
  "status": {
    "code": 400,
    "isSuccess": false
  },
  "message": "request tidak valid",
  "data": "Key: 'SelectIncomeSourceRequest.IncomeSourceType' Error:Field validation for 'IncomeSourceType' failed on the 'oneof' tag"
}
```

---

## Error Reference

| HTTP Code | Situasi |
|-----------|---------|
| 400 | Request body tidak valid / field wajib kosong / nilai tidak sesuai |
| 401 | Token tidak ada, expired, atau tidak valid |
| 404 | Resource tidak ditemukan (misal: work_platform_id tidak ada) |
| 409 | Konflik data (misal: nomor HP sudah terdaftar aktif) |
| 500 | Internal server error |

---

## Testing dengan curl

### Full onboarding flow:

```bash
# Step 1 — Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "081234567890", "full_name": "Budi Santoso"}'

# Step 2 — Verify OTP (ganti SESSION_TOKEN dan OTP_CODE)
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -H "X-Session-Token: SESSION_TOKEN" \
  -d '{"code": "123456"}'

# Step 3a — List pekerjaan (ganti JWT_TOKEN)
curl -X GET http://localhost:8080/api/v1/onboarding/work-categories \
  -H "Authorization: Bearer JWT_TOKEN"

# Step 3b — Pilih pekerjaan (ganti JWT_TOKEN dan PLATFORM_ID)
curl -X POST http://localhost:8080/api/v1/onboarding/work-platform \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer JWT_TOKEN" \
  -d '{"work_platform_id": "PLATFORM_ID"}'

# Step 4a — List sumber penghasilan
curl -X GET http://localhost:8080/api/v1/onboarding/income-sources \
  -H "Authorization: Bearer JWT_TOKEN"

# Step 4b — Pilih sumber penghasilan
curl -X POST http://localhost:8080/api/v1/onboarding/income-source \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer JWT_TOKEN" \
  -d '{"income_source_type": "manual"}'
```
