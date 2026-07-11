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

### Check Phone (Login Flow)

**POST** `/auth/check-phone`

Request:
```json
{ "phone_number": "081234567890" }
```

Response `200 OK` — sudah punya PIN:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "nomor ditemukan",
  "data": { "has_pin": true }
}
```

Response `200 OK` — belum punya PIN:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "nomor ditemukan, silakan buat PIN terlebih dahulu",
  "data": {
    "has_pin": false,
    "session_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

Response `404 Not Found`:
```json
{
  "status": { "code": 404, "isSuccess": false },
  "message": "nomor HP tidak terdaftar",
  "data": null
}
```

Response `400 Bad Request` (belum verifikasi OTP):
```json
{
  "status": { "code": 400, "isSuccess": false },
  "message": "akun belum diverifikasi, silakan daftar terlebih dahulu",
  "data": null
}
```

---

### Set PIN

**POST** `/auth/set-pin`

Header:
```
X-Session-Token: <session_token dari check-phone>
```

Request:
```json
{ "pin": "123456" }
```

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "PIN berhasil dibuat",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "message": "PIN berhasil dibuat"
  }
}
```

> Setelah set PIN, simpan `token` ini sebagai `Authorization: Bearer <token>`.

Response `400 Bad Request` (PIN bukan 6 digit angka):
```json
{
  "status": { "code": 400, "isSuccess": false },
  "message": "PIN harus 6 digit angka",
  "data": null
}
```

Response `401 Unauthorized` (session expired):
```json
{
  "status": { "code": 401, "isSuccess": false },
  "message": "session tidak valid atau sudah kedaluwarsa",
  "data": null
}
```

---

### Login

**POST** `/auth/login`

Request:
```json
{
  "phone_number": "081234567890",
  "pin": "123456"
}
```

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "message": "Login berhasil"
  }
}
```

Response `404 Not Found`:
```json
{
  "status": { "code": 404, "isSuccess": false },
  "message": "nomor HP tidak terdaftar",
  "data": null
}
```

Response `401 Unauthorized` (PIN salah):
```json
{
  "status": { "code": 401, "isSuccess": false },
  "message": "PIN tidak valid",
  "data": null
}
```

---

### Logout

**POST** `/auth/logout`

Header:
```
Authorization: Bearer <token>
```

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "Logout berhasil",
  "data": { "message": "Logout berhasil" }
}
```

> Token yang sudah dipakai untuk logout langsung diblacklist di server. Request berikutnya ke endpoint mana pun yang butuh `Authorization: Bearer <token>` dengan token yang sama akan ditolak (401), meski token tsb belum kedaluwarsa. Frontend wajib menghapus token dari storage lokal setelah logout berhasil.

Response `401 Unauthorized` (token tidak ada / tidak valid):
```json
{
  "status": { "code": 401, "isSuccess": false },
  "message": "token diperlukan",
  "data": null
}
```

Response `401 Unauthorized` (token sudah pernah di-logout sebelumnya):
```json
{
  "status": { "code": 401, "isSuccess": false },
  "message": "token sudah tidak berlaku, silakan login kembali",
  "data": null
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

## Income Passport

Semua endpoint passport memerlukan header:
```
Authorization: Bearer <token dari Step 2>
```

---

### GET Income Passport

**GET** `/passport`

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "income passport berhasil dimuat",
  "data": {
    "eligibility": {
      "is_eligible": true,
      "days_of_data": 63,
      "min_required": 30,
      "entries_verified": 63
    },
    "active_passport": {
      "income_passport_id": "uuid",
      "passport_number": "SKT-2026-RH-A3F7B2C1",
      "emi_value": 6850000,
      "period_type": "3_bulan",
      "period_label": "Apr–Jun 2026",
      "risk_level": "RENDAH",
      "issued_at": "2026-06-29"
    }
  }
}
```

> `active_passport` = `null` jika user belum pernah menerbitkan passport.

---

### Preview Passport (Konfirmasi Penerbitan)

**GET** `/passport/preview`

Query Params:

| Param | Wajib | Keterangan |
|-------|-------|-----------|
| `period` | ✅ | `3_bulan` \| `6_bulan` \| `12_bulan` |

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "preview passport berhasil dimuat",
  "data": {
    "period_type": "3_bulan",
    "period_label": "Apr–Jun",
    "period_start": "2026-04-01",
    "period_end": "2026-06-30",
    "emi_value": 6850000,
    "stability_label": "STABIL",
    "trend_direction": "up",
    "trend_change_pct": 12.0,
    "risk_level": "RENDAH",
    "risk_score": 0.15,
    "total_entries": 63
  }
}
```

> `emi_value` = total penghasilan dalam periode dibagi jumlah bulan.  
> `stability_label`: `STABIL` (trend up) · `CUKUP STABIL` (stable) · `FLUKTUATIF` (down).

Response `400 Bad Request` (forecast belum ada):
```json
{
  "status": { "code": 400, "isSuccess": false },
  "message": "data forecast belum tersedia, catat minimal 30 hari transaksi terlebih dahulu",
  "data": null
}
```

---

### Terbitkan Passport

**POST** `/passport`

Request:
```json
{
  "period": "3_bulan"
}
```

| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|-----------|
| `period` | string | ✅ | `3_bulan` \| `6_bulan` \| `12_bulan` |

Response `201 Created`:
```json
{
  "status": { "code": 201, "isSuccess": true },
  "message": "income passport berhasil diterbitkan",
  "data": {
    "income_passport_id": "uuid",
    "passport_number": "SKT-2026-RH-A3F7B2C1",
    "emi_value": 6850000,
    "period_type": "3_bulan",
    "period_label": "Apr–Jun 2026",
    "risk_level": "RENDAH",
    "issued_at": "2026-06-29"
  }
}
```

> `passport_number` format: `SKT-{TAHUN}-{2 HURUF ACAK}-{8 KARAKTER HASH TERAKHIR}`.

Response `400 Bad Request` (period tidak valid):
```json
{
  "status": { "code": 400, "isSuccess": false },
  "message": "request tidak valid",
  "data": "Key: 'IssuePassportRequest.Period' Error:Field validation for 'Period' failed on the 'oneof' tag"
}
```

---

## Dashboard

### GET Dashboard

**GET** `/dashboard`

Header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "dashboard berhasil dimuat",
  "data": {
    "user_full_name": "Rudi",
    "today_income": 285000,
    "valid_chain_count": 63,
    "forecast": {
      "emi_value": 6850000,
      "confidence": "high",
      "trend_direction": "up",
      "trend_change_pct": 12.0,
      "month_to_date_income": 285000,
      "forecast_total": 6850000,
      "risk_level": "RENDAH",
      "risk_score": 0.15,
      "is_data_sufficient": true,
      "days_of_data": 63,
      "transaction_count": 63,
      "model_name": "prophet"
    },
    "recent_transactions": [
      {
        "transaction_id": "uuid",
        "source_name": "Gojek",
        "source_provider": "GoPay",
        "amount": 285000,
        "transaction_date": "2026-06-29",
        "category": "Trip sore",
        "hash_prefix": "a3E7b2c1"
      }
    ],
    "trend_data": [
      { "date": "2026-06-01", "amount": 350000 },
      { "date": "2026-06-02", "amount": 0 },
      { "date": "2026-06-03", "amount": 420000 }
    ]
  }
}
```

> `forecast` akan `null` jika user belum pernah mencatat transaksi sama sekali.

---

## Buku Kas Digital

### GET Ledger

**GET** `/ledger`

Header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Query Params (opsional):

| Param | Default | Keterangan |
|-------|---------|-----------|
| `period` | `all` | Filter periode: `all` \| `bulan_ini` \| `3_bulan` |
| `source_id` | — | UUID sumber penghasilan (dari `/transactions/sources`) |

Filter dapat dikombinasikan. Contoh: `?period=bulan_ini&source_id=UUID`

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "buku kas berhasil diambil",
  "data": {
    "summary": {
      "total_entries": 63,
      "chain_valid": true
    },
    "transactions": [
      {
        "transaction_id": "uuid",
        "source_name": "Gojek",
        "source_provider": "GoPay",
        "amount": 285000,
        "transaction_date": "2026-06-29",
        "transaction_time": "14:22",
        "is_verified": true,
        "hash_prefix": "a3f7b2c1"
      }
    ],
    "total": 4
  }
}
```

> `chain_valid` selalu `true` — semua entri di ledger berstatus `success` (dijaga di level CreateTransaction).  
> `transaction_time` format `"HH:MM"` (diambil dari kolom `transaction_date` yang bertipe `datetime`).  
> `total_entries` = jumlah transaksi setelah filter diterapkan.

Response `400 Bad Request` (period tidak valid):
```json
{
  "status": { "code": 400, "isSuccess": false },
  "message": "period tidak valid, gunakan: all, bulan_ini, atau 3_bulan",
  "data": null
}
```

Response `400 Bad Request` (source_id bukan UUID):
```json
{
  "status": { "code": 400, "isSuccess": false },
  "message": "source_id tidak valid",
  "data": null
}
```

---

## Transaksi (Catat Penghasilan)

Semua endpoint transaksi memerlukan header:
```
Authorization: Bearer <token dari Step 2>
```

---

### List Sumber Penghasilan (Dropdown)

**GET** `/transactions/sources`

Header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Query Params:

| Param | Wajib | Keterangan |
|-------|-------|-----------|
| `provider` | ❌ | Filter berdasarkan jenis sumber penghasilan / payment provider, contoh: `GoPay`, `OVO`, `ShopeePay`, `Manual`. Jika kosong, semua sumber aktif dikembalikan. |

Contoh: `GET /transactions/sources?provider=GoPay`

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "sumber penghasilan berhasil diambil",
  "data": {
    "sources": [
      { "transaction_source_id": "uuid-1", "name": "Gojek",      "provider": "GoPay"     },
      { "transaction_source_id": "uuid-2", "name": "Grab",       "provider": "OVO"       },
      { "transaction_source_id": "uuid-3", "name": "ShopeeFood", "provider": "ShopeePay" },
      { "transaction_source_id": "uuid-4", "name": "GoSend",     "provider": "GoPay"     },
      { "transaction_source_id": "uuid-5", "name": "Tokopedia",  "provider": "Manual"    },
      { "transaction_source_id": "uuid-6", "name": "Lainnya",    "provider": "Manual"    }
    ]
  }
}
```

Response saat `?provider=GoPay`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "sumber penghasilan berhasil diambil",
  "data": {
    "sources": [
      { "transaction_source_id": "uuid-1", "name": "Gojek",  "provider": "GoPay" },
      { "transaction_source_id": "uuid-4", "name": "GoSend", "provider": "GoPay" }
    ]
  }
}
```

> Gunakan `transaction_source_id` dari response ini untuk isi field `transaction_source_id` saat Catat Penghasilan. Frontend menampilkan sebagai `"Gojek · GoPay"` (Name · Provider).

---

### Catat Penghasilan

**POST** `/transactions`

Header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

Request:
```json
{
  "amount": 285000,
  "transaction_source_id": "uuid-dari-sources",
  "transaction_date": "2026-06-29",
  "description": "Trip sore"
}
```

| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|-----------|
| `amount` | number | ✅ | Nominal dalam IDR, harus > 0 |
| `transaction_source_id` | string (UUID) | ✅ | Dari endpoint GET /transactions/sources |
| `transaction_date` | string | ✅ | Format `YYYY-MM-DD` |
| `description` | string | ❌ | Catatan opsional |

Response `201 Created`:
```json
{
  "status": { "code": 201, "isSuccess": true },
  "message": "penghasilan berhasil dicatat ke ledger",
  "data": {
    "transaction_id": "uuid",
    "amount": 285000,
    "current_hash": "a3e7b2c1d4f5e6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1",
    "status": "success",
    "message": "penghasilan berhasil dicatat ke ledger"
  }
}
```

> `current_hash` adalah SHA-256 dari data transaksi yang di-chain dengan hash transaksi sebelumnya. Ini yang ditampilkan di UI sebagai 8 karakter pertama (e.g., `a3E7b2c1`).

Response `400 Bad Request` (amount tidak valid):
```json
{
  "status": { "code": 400, "isSuccess": false },
  "message": "request tidak valid",
  "data": "Key: 'CreateTransactionRequest.Amount' Error:Field validation for 'Amount' failed on the 'gt' tag"
}
```

Response `400 Bad Request` (format tanggal salah):
```json
{
  "status": { "code": 400, "isSuccess": false },
  "message": "format transaction_date tidak valid, gunakan YYYY-MM-DD",
  "data": null
}
```

Response `404 Not Found` (source tidak ditemukan):
```json
{
  "status": { "code": 404, "isSuccess": false },
  "message": "sumber penghasilan tidak ditemukan",
  "data": null
}
```

---

### List Transaksi

**GET** `/transactions?limit=20`

Header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Query Params (opsional):

| Param | Default | Keterangan |
|-------|---------|-----------|
| `limit` | 20 | Jumlah transaksi yang dikembalikan |

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "transaksi berhasil diambil",
  "data": {
    "transactions": [
      {
        "transaction_id": "uuid",
        "source_name": "Gojek",
        "source_provider": "GoPay",
        "amount": 285000,
        "transaction_date": "2026-06-29",
        "category": "Trip sore",
        "hash_prefix": "a3E7b2c1"
      },
      {
        "transaction_id": "uuid",
        "source_name": "Grab",
        "source_provider": "OVO",
        "amount": 195000,
        "transaction_date": "2026-06-27",
        "category": "",
        "hash_prefix": "c1a9e3e7"
      }
    ],
    "total": 2
  }
}
```

---

## Error Reference

| HTTP Code | Situasi |
|-----------|---------|
| 400 | Request body tidak valid / field wajib kosong / nilai tidak sesuai |
| 401 | Token tidak ada, expired, tidak valid, atau sudah di-logout (blacklisted) |
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

---

### Logout:

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer JWT_TOKEN"
```

---

### Dashboard & Transaksi:

```bash
# Dashboard
curl -X GET http://localhost:8080/api/v1/dashboard \
  -H "Authorization: Bearer JWT_TOKEN"

# List sumber penghasilan (untuk dropdown Catat Penghasilan)
curl -X GET http://localhost:8080/api/v1/transactions/sources \
  -H "Authorization: Bearer JWT_TOKEN"

# Catat Penghasilan (ganti SOURCE_ID dengan UUID dari /transactions/sources)
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer JWT_TOKEN" \
  -d '{
    "amount": 285000,
    "transaction_source_id": "SOURCE_ID",
    "transaction_date": "2026-06-29",
    "description": "Trip sore"
  }'

# List transaksi
curl -X GET http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer JWT_TOKEN"

# List transaksi dengan limit custom
curl -X GET "http://localhost:8080/api/v1/transactions?limit=10" \
  -H "Authorization: Bearer JWT_TOKEN"
```

---

### Buku Kas Digital:

```bash
# Semua transaksi
curl -X GET http://localhost:8080/api/v1/ledger \
  -H "Authorization: Bearer JWT_TOKEN"

# Filter bulan ini
curl -X GET "http://localhost:8080/api/v1/ledger?period=bulan_ini" \
  -H "Authorization: Bearer JWT_TOKEN"

# Filter 3 bulan terakhir
curl -X GET "http://localhost:8080/api/v1/ledger?period=3_bulan" \
  -H "Authorization: Bearer JWT_TOKEN"

# Filter by source (ganti SOURCE_ID dengan UUID dari /transactions/sources)
curl -X GET "http://localhost:8080/api/v1/ledger?source_id=SOURCE_ID" \
  -H "Authorization: Bearer JWT_TOKEN"

# Kombinasi filter
curl -X GET "http://localhost:8080/api/v1/ledger?period=bulan_ini&source_id=SOURCE_ID" \
  -H "Authorization: Bearer JWT_TOKEN"
```

---

### Income Passport:

```bash
# Cek eligibilitas & passport aktif
curl -X GET http://localhost:8080/api/v1/passport \
  -H "Authorization: Bearer JWT_TOKEN"

# Preview sebelum terbitkan (pilih periode)
curl -X GET "http://localhost:8080/api/v1/passport/preview?period=3_bulan" \
  -H "Authorization: Bearer JWT_TOKEN"

curl -X GET "http://localhost:8080/api/v1/passport/preview?period=6_bulan" \
  -H "Authorization: Bearer JWT_TOKEN"

curl -X GET "http://localhost:8080/api/v1/passport/preview?period=12_bulan" \
  -H "Authorization: Bearer JWT_TOKEN"

# Terbitkan passport
curl -X POST http://localhost:8080/api/v1/passport \
  -H "Authorization: Bearer JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"period": "3_bulan"}'

# Lihat daftar akses
curl -X GET http://localhost:8080/api/v1/passport/access \
  -H "Authorization: Bearer JWT_TOKEN"

# Berikan akses ke organisasi
curl -X POST http://localhost:8080/api/v1/passport/access \
  -H "Authorization: Bearer JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"organization_id": "ORG_UUID", "data_scope": ["emi","stability","risk"], "expires_in_days": 30}'

# Cabut akses
curl -X PATCH http://localhost:8080/api/v1/passport/access/CONSENT_UUID/revoke \
  -H "Authorization: Bearer JWT_TOKEN"

# Riwayat akses
curl -X GET http://localhost:8080/api/v1/passport/access/logs \
  -H "Authorization: Bearer JWT_TOKEN"
```

---

## Kelola Akses Data

### Get Consents (Screen 12)

**GET** `/passport/access`

Header: `Authorization: Bearer <token>`

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "daftar akses berhasil dimuat",
  "data": {
    "consents": [
      {
        "consent_id": "uuid",
        "organization_name": "Koperasi Sejahtera Jawa",
        "organization_type": "fintech",
        "granted_at": "29 Jun 2026",
        "data_scope": ["EMI", "Tren Stabilitas", "Risiko"],
        "expires_at": "29 Jul 2026",
        "days_remaining": 30,
        "status": "active",
        "status_label": "AKTIF",
        "purpose": ""
      }
    ]
  }
}
```

> `status_label`: `"AKTIF"` (aktif), `"X HR LAGI"` (≤3 hari sebelum kedaluwarsa), `"DICABUT"` (dicabut), `"KEDALUWARSA"` (expired).

---

### Grant Access (Bagikan ke Pihak Baru)

**POST** `/passport/access`

Header: `Authorization: Bearer <token>`

Request:
```json
{
  "organization_id": "uuid-organisasi",
  "data_scope": ["emi", "stability", "risk"],
  "expires_in_days": 30,
  "purpose": "Penilaian kredit motor"
}
```

> `data_scope` valid values: `"emi"`, `"stability"`, `"risk"`, `"full"`. `expires_in_days` = 0 berarti tidak ada batas waktu.

Response `201 Created`:
```json
{
  "status": { "code": 201, "isSuccess": true },
  "message": "akses berhasil diberikan",
  "data": null
}
```

Response `400 Bad Request` (passport belum diterbitkan):
```json
{
  "status": { "code": 400, "isSuccess": false },
  "message": "income passport belum diterbitkan",
  "data": null
}
```

Response `409 Conflict` (akses sudah ada):
```json
{
  "status": { "code": 409, "isSuccess": false },
  "message": "akses untuk organisasi ini sudah aktif",
  "data": null
}
```

---

### Revoke Access (Cabut Akses)

**PATCH** `/passport/access/:consent_id/revoke`

Header: `Authorization: Bearer <token>`

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "akses berhasil dicabut",
  "data": null
}
```

Response `404 Not Found`:
```json
{
  "status": { "code": 404, "isSuccess": false },
  "message": "akses tidak ditemukan",
  "data": null
}
```

---

## Riwayat Akses Data

### Get Access Logs (Screen 13)

**GET** `/passport/access/logs?filter=`

Header: `Authorization: Bearer <token>`

Query params:
- `filter` (opsional): `""` (semua) atau `"income_passport"`

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "riwayat akses berhasil dimuat",
  "data": {
    "logs": [
      {
        "access_log_id": "uuid",
        "organization_name": "Koperasi Sejahtera Jawa",
        "organization_type": "fintech",
        "accessed_at": "29 Jun 2026 · 14:32 WIB",
        "data_scope": ["EMI", "Tren Stabilitas", "Risiko"],
        "consent_status": "active",
        "status_label": "VALID",
        "note": "Kamu sudah diberitahu"
      },
      {
        "access_log_id": "uuid",
        "organization_name": "Bank XYZ Leasing",
        "organization_type": "bank",
        "accessed_at": "05 Jun 2026 · 11:20 WIB",
        "data_scope": ["Akses penuh"],
        "consent_status": "revoked",
        "status_label": "DICABUT",
        "note": "Akses dicabut oleh kamu · 8 Jun"
      }
    ]
  }
}
```

> `status_label`: `"VALID"` jika consent masih aktif, `"DICABUT"` jika sudah dicabut.

---

## Daftar Organisasi

### Get Organizations

**GET** `/organizations`

Header: `Authorization: Bearer <token>`

Response `200 OK`:
```json
{
  "status": { "code": 200, "isSuccess": true },
  "message": "daftar organisasi berhasil dimuat",
  "data": [
    { "organization_id": "uuid", "name": "Koperasi Sejahtera Jawa", "type": "fintech" },
    { "organization_id": "uuid", "name": "PT BPR Sentosa Digital", "type": "bank" }
  ]
}
```
