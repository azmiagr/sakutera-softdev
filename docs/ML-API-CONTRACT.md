# KONTRAK API & PRD — ML FORECASTING SERVICE

## SAKUTERA — Smart Salary Estimator

| Atribut Dokumen | Detail |
|---|---|
| Nomor Dokumen | SKT-ML-API-2026-001 |
| Versi | 1.0.0 (Baseline) |
| Status | DRAFT — Untuk Implementasi |
| Komponen | ML Forecasting Service (Python + FastAPI + Prophet) |
| Konsumen Utama | Backend Service (Golang/Gin) |
| Referensi | PRD SKT-PRD-2026-001 — Modul F-04, F-05 |

---

## 1. RINGKASAN & PRINSIP DESAIN

ML Forecasting Service adalah **microservice stateless** yang bertugas mengubah riwayat
transaksi penghasilan menjadi 3 output analitik: **EMI** (estimasi gaji setara bulanan),
**proyeksi penghasilan**, dan **skor risiko defisit arus kas**.

### 1.1 Prinsip Wajib

| # | Prinsip | Konsekuensi |
|---|---|---|
| P1 | **Stateless** — service tidak punya database sendiri | Semua data riwayat DIKIRIM oleh BE di setiap request. ML tidak menyimpan apa pun |
| P2 | **Internal-only** — tidak boleh diakses dari FE/publik | Hanya BE (Golang) yang boleh memanggil. FE → BE → ML, tidak pernah FE → ML |
| P3 | **ML tidak pernah menghambat user** | Endpoint berat dipanggil async oleh BE; hasil di-cache BE di Redis |
| P4 | **Idempotent** — input sama → output sama | Tidak ada efek samping. Bisa di-retry & di-scale horizontal dengan aman |
| P5 | **BE pegang otoritas data & auth** | ML percaya data yang dikirim BE sudah tervalidasi & terverifikasi (hash-chain) |

### 1.2 Posisi dalam Alur Sistem

```
┌──────┐  1. minta dashboard   ┌─────────────┐  3. kirim riwayat   ┌──────────────┐
│  FE  │ ───────────────────►  │   BE (Go)   │ ─────────────────►  │ ML (FastAPI) │
│ Next │ ◄───────────────────  │  + MariaDB  │ ◄─────────────────  │   Prophet    │
└──────┘  6. response JSON     └─────────────┘  4. hasil analitik  └──────────────┘
                                 2. ambil riwayat        5. (opsional) cache ke Redis
                                    dari DB
```

> **Catatan:** FE tidak pernah berkomunikasi langsung dengan ML. Data riwayat selalu
> berasal dari MariaDB, dibaca oleh BE, lalu dikirim ke ML sebagai payload request.

---

## 2. INFORMASI DASAR

| Aspek | Nilai |
|---|---|
| Base URL (internal) | `http://ml-service:8000` (di dalam Docker network) |
| Versioning | Prefix path `/v1` |
| Format | JSON (request & response), `Content-Type: application/json` |
| Encoding | UTF-8 |
| Mata uang | IDR (integer, tanpa desimal — rupiah penuh) |
| Tanggal | ISO 8601 `YYYY-MM-DD` (tanggal), `YYYY-MM-DDTHH:mm:ssZ` (timestamp, UTC) |
| Autentikasi | Header `X-Service-Token: <token>` — shared secret antar service internal |
| Dokumentasi live | `GET /docs` (Swagger UI bawaan FastAPI) |

---

## 3. MODEL DATA BERSAMA

### 3.1 Object `Transaction` (dikirim BE ke ML)

Representasi satu entri ledger penghasilan. BE memetakan baris MariaDB ke object ini.

```json
{
  "date": "2026-06-15",
  "amount": 250000,
  "source": "gojek",
  "category": "ride"
}
```

| Field | Tipe | Wajib | Deskripsi |
|---|---|---|---|
| `date` | string (date) | ✅ | Tanggal transaksi. Tidak boleh di masa depan |
| `amount` | integer | ✅ | Nominal penghasilan dalam rupiah penuh. Harus > 0 |
| `source` | string | ✅ | Sumber penghasilan (mis. `gojek`, `grab`, `manual`, `fiverr`) |
| `category` | string | ❌ | Kategori opsional (mis. `ride`, `food`, `project`). Bantu akurasi model |

> ML **tidak** menerima/memerlukan field hash, `tx_id`, atau data sensitif pengguna.
> Cukup tanggal + nominal + sumber. Privasi by design.

---

## 4. ENDPOINT

### 4.1 `GET /v1/health` — Health Check

Untuk liveness/readiness probe (Docker, orchestrator).

**Response `200 OK`**
```json
{ "status": "ok", "model": "prophet", "version": "1.0.0" }
```

---

### 4.2 `POST /v1/forecast` — Estimasi Dashboard (Utama)

Endpoint inti untuk **F-04.01, F-04.04, F-04.05, F-04.06**. BE memanggil ini setelah ada
transaksi baru (async) atau saat perlu menyegarkan dashboard.

#### Request Body
```json
{
  "user_id": "usr_8f3a21",
  "as_of_date": "2026-06-29",
  "transactions": [
    { "date": "2026-06-01", "amount": 280000, "source": "gojek", "category": "ride" },
    { "date": "2026-06-02", "amount": 175000, "source": "grab", "category": "food" },
    { "date": "2026-06-03", "amount": 310000, "source": "gojek", "category": "ride" }
  ],
  "options": {
    "forecast_horizon_days": 30,
    "include_explanation": true
  }
}
```

| Field | Tipe | Wajib | Default | Deskripsi |
|---|---|---|---|---|
| `user_id` | string | ✅ | — | ID pengguna (untuk logging/tracing, bukan untuk lookup data) |
| `as_of_date` | string (date) | ✅ | — | Tanggal "hari ini" perspektif perhitungan. ML proyeksi dari titik ini |
| `transactions` | array\<Transaction\> | ✅ | — | SELURUH riwayat yang relevan. BE yang menentukan rentang (mis. 90–365 hari terakhir) |
| `options.forecast_horizon_days` | integer | ❌ | 30 | Berapa hari ke depan diproyeksikan |
| `options.include_explanation` | boolean | ❌ | true | Sertakan penjelasan bahasa awam (F-04.07) |

#### Response `200 OK`
```json
{
  "user_id": "usr_8f3a21",
  "as_of_date": "2026-06-29",
  "data_sufficiency": {
    "days_of_data": 45,
    "transaction_count": 132,
    "is_sufficient": true,
    "min_required_days": 30
  },
  "emi": {
    "value": 6850000,
    "currency": "IDR",
    "confidence": "medium"
  },
  "trend": {
    "direction": "up",
    "change_pct": 15.2,
    "vs_period": "previous_month"
  },
  "month_to_date": {
    "total_earned": 5120000,
    "days_recorded": 28
  },
  "forecast_remaining_month": {
    "projected_additional": 1730000,
    "projected_month_total": 6850000,
    "horizon_days": 30
  },
  "deficit_risk": {
    "level": "SEDANG",
    "score": 0.62,
    "projected_deficit_date": "2026-07-08",
    "estimated_shortfall": 450000,
    "days_ahead": 9,
    "should_notify": true
  },
  "explanation": "Estimasi gaji bulananmu Rp6.850.000, naik 15% dari bulan lalu. Berdasarkan polamu, sekitar 8 Juli penghasilan diperkirakan menipis ~Rp450.000. Pertimbangkan menambah jam kerja minggu depan.",
  "model": { "name": "prophet", "version": "1.0.0" },
  "generated_at": "2026-06-29T14:30:00Z"
}
```

#### Penjelasan Field Response Penting

| Field | Makna | Requirement |
|---|---|---|
| `data_sufficiency.is_sufficient` | `false` jika data < 30 hari → EMI ditandai low confidence (lihat §6) | R-04 |
| `emi.value` | **Equivalent Monthly Income** — estimasi gaji setara bulanan | F-04.01 |
| `emi.confidence` | `low` / `medium` / `high` — naik seiring banyaknya data | F-04.07 |
| `trend.direction` | `up` / `down` / `stable` — vs periode pembanding | F-04.03 |
| `forecast_remaining_month.projected_month_total` | Proyeksi total bulan berjalan | F-04.04 |
| `deficit_risk.level` | `RENDAH` / `SEDANG` / `TINGGI` | F-04.05 |
| `deficit_risk.estimated_shortfall` | Estimasi besar defisit (rupiah) | F-04.05 |
| `deficit_risk.should_notify` | `true` jika BE harus kirim notifikasi (level ≥ SEDANG & ≥7 hari sebelum) | F-04.06 |
| `explanation` | Kalimat bahasa awam Indonesia, siap tampil ke user | F-04.07 |

> **Pembagian tugas notifikasi:** ML hanya menentukan `should_notify` (logika risiko).
> Pengiriman push/SMS aktual dilakukan **Notification Service di BE**, bukan ML.

---

### 4.3 `POST /v1/passport-summary` — Ringkasan Income Passport

Untuk **F-05.02**. Dipanggil BE (synchronous) saat user menerbitkan Income Passport.
BE butuh angka EMI 3/6/12 bulan + tren stabilitas untuk dimasukkan ke PDF.

#### Request Body
```json
{
  "user_id": "usr_8f3a21",
  "as_of_date": "2026-06-29",
  "transactions": [ "...riwayat minimal 12 bulan jika tersedia..." ]
}
```

#### Response `200 OK`
```json
{
  "user_id": "usr_8f3a21",
  "as_of_date": "2026-06-29",
  "emi_history": {
    "last_3_months": 6850000,
    "last_6_months": 6420000,
    "last_12_months": 5980000
  },
  "stability": {
    "score": 0.74,
    "label": "STABIL",
    "coefficient_of_variation": 0.26
  },
  "risk": {
    "level": "RENDAH",
    "score": 0.31
  },
  "period_covered": {
    "from": "2025-07-01",
    "to": "2026-06-29",
    "active_days": 312
  },
  "model": { "name": "prophet", "version": "1.0.0" },
  "generated_at": "2026-06-29T14:35:00Z"
}
```

| Field | Makna |
|---|---|
| `emi_history.last_3/6/12_months` | EMI dihitung pada 3 jendela waktu (F-05.02) |
| `stability.label` | `STABIL` / `CUKUP STABIL` / `FLUKTUATIF` — dari koefisien variasi penghasilan |
| `stability.score` | 0–1, makin tinggi makin stabil |
| `risk.level` | Skor risiko untuk dicantumkan di passport |

> ML hanya menyediakan **angka**. Pembuatan PDF, digital signature, dan kode unik
> dokumen adalah tugas **Passport Service (BE)** — di luar scope ML.

---

## 5. DEFINISI PERHITUNGAN

### 5.1 EMI (Equivalent Monthly Income)

EMI = estimasi penghasilan untuk **satu bulan penuh (30 hari)** berdasarkan pola historis,
bukan sekadar penjumlahan. Prophet memodelkan tren + musiman (weekend, hari besar) lalu
diagregasi ke setara bulanan. Untuk jendela waktu (3/6/12 bln), EMI = rata-rata bulanan
penghasilan dalam jendela tersebut, dihaluskan model.

### 5.2 Skor Risiko Defisit

1. Prophet memproyeksikan penghasilan harian `forecast_horizon_days` ke depan.
2. ML membandingkan proyeksi kumulatif dengan baseline kebutuhan (turunan dari rata-rata
   historis pengguna).
3. Skor 0–1 dipetakan ke level:

| Skor | Level | `should_notify` |
|---|---|---|
| 0.00 – 0.33 | RENDAH | false |
| 0.34 – 0.66 | SEDANG | true (jika proyeksi defisit ≥7 hari ke depan) |
| 0.67 – 1.00 | TINGGI | true |

### 5.3 Cold-Start (data < 30 hari) — Mitigasi R-04

Jika `days_of_data < 30`:
- `data_sufficiency.is_sufficient = false`
- `emi.confidence = "low"`
- ML pakai estimasi statistik sederhana (rata-rata bergerak), **bukan** Prophet penuh.
- `explanation` menyertakan disclaimer: *"Estimasi akan semakin akurat seiring kamu rutin mencatat."*

---

## 6. PENANGANAN ERROR

Semua error memakai format konsisten:

```json
{
  "error": {
    "code": "INSUFFICIENT_DATA",
    "message": "Minimal 1 transaksi diperlukan untuk menghitung estimasi.",
    "details": { "transaction_count": 0 }
  }
}
```

| HTTP Status | `error.code` | Kapan terjadi |
|---|---|---|
| `400 Bad Request` | `VALIDATION_ERROR` | Field wajib hilang, tipe data salah, `amount` ≤ 0 |
| `400 Bad Request` | `FUTURE_DATE` | Ada transaksi bertanggal di masa depan |
| `401 Unauthorized` | `INVALID_SERVICE_TOKEN` | `X-Service-Token` salah/tidak ada |
| `422 Unprocessable` | `INSUFFICIENT_DATA` | `transactions` kosong (tidak ada apa pun untuk dihitung) |
| `500 Internal` | `MODEL_ERROR` | Kegagalan internal model (BE sebaiknya retry / fallback) |

> Catatan: data < 30 hari **bukan** error — tetap `200 OK` dengan `is_sufficient: false`
> dan confidence rendah. Error `INSUFFICIENT_DATA` hanya untuk array benar-benar kosong.

---

## 7. PERSYARATAN NON-FUNGSIONAL

| ID | Persyaratan | Target | Ref PRD |
|---|---|---|---|
| ML-NF-01 | Latensi `/v1/forecast` (P95, data ~90 hari) | < 1,5 detik | NF-02 |
| ML-NF-02 | Latensi `/v1/passport-summary` (P95) | < 5 detik | NF-03 |
| ML-NF-03 | Service stateless — mendukung horizontal scaling | tanpa refactor | NF-14 |
| ML-NF-04 | Seluruh teks `explanation` & label dalam Bahasa Indonesia | 100% | NF-16 |
| ML-NF-05 | Test coverage logika model & endpoint | > 80% | NF-17 |
| ML-NF-06 | Dokumentasi API otomatis (OpenAPI 3.0 via FastAPI) | tersedia di `/docs` | NF-18 |

---

## 8. CONTOH ALUR END-TO-END

Skenario: Rudi (ojol) membuka dashboard.

```
1. FE  → GET /api/dashboard            (FE minta ke BE, bawa JWT user)
2. BE  : validasi JWT, query MariaDB → ambil 90 transaksi terakhir Rudi
3. BE  → POST http://ml-service:8000/v1/forecast
         Header: X-Service-Token: <secret>
         Body  : { user_id, as_of_date, transactions:[...90 entri...] }
4. ML  : Prophet fit → hitung EMI, tren, risiko defisit
5. ML  → 200 OK { emi, trend, deficit_risk, explanation, ... }
6. BE  : (opsional) cache hasil ke Redis; gabung dengan data lain
7. BE  → 200 OK ke FE { emi, grafik, peringatan, ... }
8. FE  : render dashboard — "Gaji bulananmu Rp6.850.000 ↑15%"
```

Yang perlu diingat: langkah 3–5 adalah satu-satunya titik ML terlibat. Sumber data tetap
MariaDB (langkah 2). Saat development, langkah 2–3 disimulasikan oleh **dummy data
generator** yang menghasilkan array `transactions` berbentuk sama persis.

---

## 9. RINGKASAN PEMBAGIAN TANGGUNG JAWAB

| Tugas | Pemilik |
|---|---|
| Simpan/baca riwayat transaksi (MariaDB) | BE (Go) |
| Hash-chain, validasi ledger, auth user (JWT) | BE (Go) |
| Kirim array transaksi ke ML | BE (Go) |
| Hitung EMI, tren, skor risiko defisit | **ML (FastAPI)** |
| Tentukan `should_notify` (logika risiko) | **ML (FastAPI)** |
| Kirim notifikasi push/SMS aktual | BE (Notification Service) |
| Generate PDF, digital signature, kode unik Passport | BE (Passport Service) |
| Cache hasil ML ke Redis | BE (Go) |
| Render UI / grafik | FE (Next.js) |

---

— AKHIR DOKUMEN —
Kontrak API ML Forecasting Service v1.0.0 | SAKUTERA | Tim Engineers SABI SAMA KITA
