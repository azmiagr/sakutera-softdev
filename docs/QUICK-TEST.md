# SAKUTERA ML — Quick Test Guide (Live)

> Base URL live: `https://irzadzulhika29-sakutera-ml.hf.space`
> Token: `3Rl1B5AsCW0uzfcSEHdUOGXxthoLZ9MJ`

---

## Daftar Endpoint

| Method | Endpoint | Auth | Fungsi |
|--------|----------|------|-------|
| GET | `/v1/health` | ❌ | Health check |
| POST | `/v1/forecast` | ✅ | EMI, tren, proyeksi, risiko defisit |
| POST | `/v1/passport-summary` | ✅ | Ringkasan passport 3/6/12 bulan |
| GET | `/docs` | ❌ | Swagger UI |

---

## 1. Health Check

```bash
curl https://irzadzulhika29-sakutera-ml.hf.space/v1/health
```

**Response:**
```json
{"status":"ok","model":"prophet","version":"1.0.0"}
```

---

## 2. Forecast — Estimasi Dashboard

### Request

```bash
curl -X POST https://irzadzulhika29-sakutera-ml.hf.space/v1/forecast \
  -H "Content-Type: application/json" \
  -H "X-Service-Token: 3Rl1B5AsCW0uzfcSEHdUOGXxthoLZ9MJ" \
  -d '{
    "user_id": "usr_demo",
    "as_of_date": "2026-06-29",
    "transactions": [
      { "date": "2026-06-01", "amount": 280000, "source": "gojek", "category": "ride" },
      { "date": "2026-06-02", "amount": 175000, "source": "grab", "category": "food" },
      { "date": "2026-06-03", "amount": 310000, "source": "gojek", "category": "ride" },
      { "date": "2026-06-04", "amount": 250000, "source": "gojek", "category": "ride" },
      { "date": "2026-06-05", "amount": 220000, "source": "grab", "category": "food" },
      { "date": "2026-06-06", "amount": 340000, "source": "gojek", "category": "ride" },
      { "date": "2026-06-07", "amount": 365000, "source": "gojek", "category": "ride" },
      { "date": "2026-06-08", "amount": 295000, "source": "gojek", "category": "ride" },
      { "date": "2026-06-09", "amount": 210000, "source": "grab", "category": "food" },
      { "date": "2026-06-10", "amount": 330000, "source": "gojek", "category": "ride" }
    ],
    "options": {
      "forecast_horizon_days": 30,
      "include_explanation": true
    }
  }'
```

### Response `200 OK`

```json
{
  "user_id": "usr_test",
  "as_of_date": "2026-06-29",
  "data_sufficiency": {
    "days_of_data": 29,
    "transaction_count": 5,
    "is_sufficient": false,
    "min_required_days": 30
  },
  "emi": {
    "value": 1300250,
    "currency": "IDR",
    "confidence": "low"
  },
  "trend": {
    "direction": "stable",
    "change_pct": 0.0,
    "vs_period": "previous_month"
  },
  "month_to_date": {
    "total_earned": 1235000,
    "days_recorded": 5
  },
  "forecast_remaining_month": {
    "projected_additional": 43750,
    "projected_month_total": 1278750,
    "horizon_days": 30
  },
  "deficit_risk": {
    "level": "RENDAH",
    "score": 0.109,
    "projected_deficit_date": null,
    "estimated_shortfall": 0,
    "days_ahead": null,
    "should_notify": false
  },
  "explanation": "Estimasi gaji bulananmu Rp1.300.250, relatif stabil dibanding bulan lalu. Estimasi akan semakin akurat seiring kamu rutin mencatat.",
  "model": {
    "name": "prophet",
    "version": "1.0.0"
  },
  "generated_at": "2026-06-30T04:51:23.631862Z"
}
```

### Field Response

| Field | Tipe | Deskripsi |
|-------|------|-----------|
| `emi.value` | int | Estimasi gaji setara bulanan (30 hari) |
| `emi.confidence` | string | `low` (<30 hari) / `medium` / `high` (≥90 hari) |
| `trend.direction` | string | `up` / `down` / `stable` vs bulan lalu |
| `trend.change_pct` | float | Persentase perubahan |
| `month_to_date.total_earned` | int | Total penghasilan bulan berjalan |
| `forecast_remaining_month.projected_month_total` | int | Proyeksi total bulan ini |
| `deficit_risk.level` | string | `RENDAH` / `SEDANG` / `TINGGI` |
| `deficit_risk.score` | float | 0–1 (0.00–0.33 RENDAH, 0.34–0.66 SEDANG, 0.67–1.0 TINGGI) |
| `deficit_risk.should_notify` | bool | `true` = BE harus kirim notifikasi |
| `explanation` | string | Kalimat Bahasa Indonesia siap tampil |

---

## 3. Passport Summary — Income Passport

### Request

```bash
curl -X POST https://irzadzulhika29-sakutera-ml.hf.space/v1/passport-summary \
  -H "Content-Type: application/json" \
  -H "X-Service-Token: 3Rl1B5AsCW0uzfcSEHdUOGXxthoLZ9MJ" \
  -d '{
    "user_id": "usr_demo",
    "as_of_date": "2026-06-29",
    "transactions": [
      { "date": "2026-06-01", "amount": 280000, "source": "gojek", "category": "ride" },
      { "date": "2026-06-02", "amount": 175000, "source": "grab", "category": "food" },
      { "date": "2026-06-03", "amount": 310000, "source": "gojek", "category": "ride" },
      { "date": "2026-06-04", "amount": 250000, "source": "gojek", "category": "ride" },
      { "date": "2026-06-05", "amount": 220000, "source": "grab", "category": "food" }
    ]
  }'
```

### Response `200 OK`

```json
{
  "user_id": "usr_demo",
  "as_of_date": "2026-06-29",
  "emi_history": {
    "last_3_months": 1235000,
    "last_6_months": 617500,
    "last_12_months": 308750
  },
  "stability": {
    "score": 1.0,
    "label": "STABIL",
    "coefficient_of_variation": 0.0
  },
  "risk": {
    "level": "RENDAH",
    "score": 0.0
  },
  "period_covered": {
    "from": "2026-06-01",
    "to": "2026-06-29",
    "active_days": 5
  },
  "model": {
    "name": "prophet",
    "version": "1.0.0"
  },
  "generated_at": "2026-06-30T01:36:42.997766Z"
}
```

### Field Response

| Field | Tipe | Deskripsi |
|-------|------|-----------|
| `emi_history.last_3_months` | int | EMI rata-rata 3 bulan terakhir |
| `emi_history.last_6_months` | int | EMI rata-rata 6 bulan terakhir |
| `emi_history.last_12_months` | int | EMI rata-rata 12 bulan terakhir |
| `stability.label` | string | `STABIL` (CV≤0.25) / `CUKUP STABIL` / `FLUKTUATIF` |
| `stability.score` | float | 0–1, makin tinggi makin stabil |
| `risk.level` | string | `RENDAH` / `SEDANG` / `TINGGI` |
| `period_covered.from` | date | Tanggal transaksi terawal |
| `period_covered.active_days` | int | Jumlah hari ada transaksi |

---

## 4. Swagger UI (Dokumentasi Interaktif)

Buka di browser:
```
https://irzadzulhika29-sakutera-ml.hf.space/docs
```

Bisa test langsung lewat tombol "Try it out" di browser.

---

## 5. Error Responses

Semua error format konsisten:
```json
{ "error": { "code": "...", "message": "...", "details": {} } }
```

### Contoh: Token Salah (401)

```bash
curl -X POST https://irzadzulhika29-sakutera-ml.hf.space/v1/forecast \
  -H "Content-Type: application/json" \
  -H "X-Service-Token: wrong-token" \
  -d '{"user_id":"u","as_of_date":"2026-06-29","transactions":[{"date":"2026-06-01","amount":100000,"source":"gojek"}]}'
```

```json
{"error":{"code":"INVALID_SERVICE_TOKEN","message":"Token service tidak valid.","details":{"header":"X-Service-Token"}}}
```

### Contoh: Transaksi Kosong (422)

```bash
curl -X POST https://irzadzulhika29-sakutera-ml.hf.space/v1/forecast \
  -H "Content-Type: application/json" \
  -H "X-Service-Token: 3Rl1B5AsCW0uzfcSEHdUOGXxthoLZ9MJ" \
  -d '{"user_id":"u","as_of_date":"2026-06-29","transactions":[]}'
```

```json
{"error":{"code":"INSUFFICIENT_DATA","message":"Minimal 1 transaksi diperlukan untuk menghitung estimasi.","details":{"transaction_count":0}}}
```

### Daftar Error Lengkap

| HTTP | `error.code` | Penyebab |
|------|-------------|---------|
| 400 | `VALIDATION_ERROR` | Field hilang / tipe salah / `amount` ≤ 0 |
| 400 | `FUTURE_DATE` | Transaksi bertanggal setelah `as_of_date` |
| 401 | `INVALID_SERVICE_TOKEN` | Header `X-Service-Token` salah/tidak ada |
| 422 | `INSUFFICIENT_DATA` | Array `transactions` kosong |
| 500 | `MODEL_ERROR` | Kegagalan internal model |

---

## Quick Reference

```
Base URL:  https://irzadzulhika29-sakutera-ml.hf.space
Token:     3Rl1B5AsCW0uzfcSEHdUOGXxthoLZ9MJ
Header:    X-Service-Token: <token>
Currency:  IDR (integer rupiah penuh)
Tanggal:   YYYY-MM-DD
```

### Minimal Request Body (forecast)

```json
{
  "user_id": "string",
  "as_of_date": "YYYY-MM-DD",
  "transactions": [
    { "date": "YYYY-MM-DD", "amount": 250000, "source": "gojek" }
  ]
}
```

| Field | Wajib | Aturan |
|-------|-------|-------|
| `user_id` | ✅ | String tidak kosong |
| `as_of_date` | ✅ | Tanggal "hari ini", format `YYYY-MM-DD` |
| `transactions` | ✅ | Array, minimal 1 item |
| `transactions[].date` | ✅ | Tidak boleh setelah `as_of_date` |
| `transactions[].amount` | ✅ | Integer > 0 |
| `transactions[].source` | ✅ | String tidak kosong |
| `transactions[].category` | ❌ | Opsional |
| `options` | ❌ | Default: horizon 30, explanation true |
```
