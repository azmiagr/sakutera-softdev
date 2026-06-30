# BACKEND PRODUCT REQUIREMENTS DOCUMENT (PRD)

**Produk:** Sakutera — Platform Identitas Ekonomi Terverifikasi
**Versi Dokumen:** 1.1.0 (Backend & Security Focus)
**Status:** DRAFT

---

## 1. Ringkasan Eksekutif & Tujuan Backend

Sakutera adalah infrastruktur identitas ekonomi digital untuk pekerja gig dan informal. Dari sisi _backend_, sistem ini bertugas mengelola tiga mekanisme inti secara aman dan berkinerja tinggi:

1. **Verifiable Ledger:** Pemrosesan buku kas digital berbasis rantai hash (_hash-chained_) yang tidak dapat dimanipulasi.
2. **Smart Salary Estimator:** Orkestrasi data _time-series_ dari _ledger_ ke layanan _Machine Learning_ untuk menghasilkan estimasi gaji.
3. **Income Passport:** Pembuatan dokumen identitas ekonomi berbasis PDF yang dilengkapi tanda tangan digital dan manajemen akses (_consent-based_).

---

## 2. Arsitektur Teknis & Tech Stack

Sistem menggunakan arsitektur _microservices_ dengan pemisahan tanggung jawab yang jelas untuk mendukung _horizontal scaling_ dan keamanan infrastruktur.

- **API Gateway & Reverse Proxy:** **Caddy** via `Caddyfile`. Dipilih karena otomatisasi _SSL termination_ (HTTPS) terintegrasi, konfigurasi yang lebih ringkas, dan performa yang sangat baik untuk meneruskan _traffic_ ke layanan _backend_.
- **Application Layer (Core Backend):** **Golang + Gin Framework**. Dipilih karena performa tinggi dan konkurensi efisien (via _goroutines_), krusial untuk kalkulasi _hash-chain_ tanpa latensi tinggi dan pencegahan _race conditions_.
- **AI / ML Microservice:** **Python + FastAPI + Prophet/LSTM** untuk _forecasting time-series_ penghasilan.
- **Database Relasional:** **PostgreSQL**. Dipilih untuk menggantikan MariaDB karena dukungan tipe data tingkat lanjut (seperti JSONB untuk log dinamis), performa konkurensi yang superior, dan keandalan tinggi dalam menjaga integritas data transaksional yang kompleks.
- **Cache & Rate Limiting:** **Redis** (Versi 7.0+) untuk _session management_ dan _caching_ hasil komputasi ML.
- **Message Queue:** **RabbitMQ / Redis Streams** untuk pemrosesan asinkron (misal: _webhook e-wallet_ diteruskan ke prosesor Golang).
- **Object Storage:** **S3-compatible storage** untuk menyimpan file PDF _Income Passport_ dan foto bukti transaksi.
- **Infrastruktur & Deployment:** **Docker** untuk kontainerisasi layanan, memastikan portabilitas dan konsistensi _environment_ antara tahap _development_ dan _production_.

---

## 3. Persyaratan Fungsional Backend (Services Breakdown)

### 3.1. Auth & Account Service

- **OTP & Sesi:** Sistem harus mengelola pendaftaran via nomor telepon (+62) dengan verifikasi OTP yang berlaku 5 menit (Maksimal 3 percobaan/10 menit).
- **Autentikasi Lanjutan:** Mendukung _login_ dengan PIN 6 digit dan pembaruan token sesi (JWT) secara otomatis (_sliding session_) dengan _timeout_ 30 menit inaktivitas.
- **Right to Erasure:** _Endpoint_ khusus untuk penghapusan permanen akun dan data privasi pengguna, kecuali data yang diwajibkan regulasi.

### 3.2. Ledger & Transaction Service (Core Integrity)

- **Webhook Listener:** Menangkap data transaksi otomatis dari platform/e-wallet terintegrasi via _webhook_ atau API resmi.
- **Validasi & Deduplikasi:** Validasi format di _backend_ (tipe data, rentang nilai, batas tanggal) dan deduplikasi menggunakan _composite key_ di PostgreSQL (platform + ID transaksi + tanggal + jumlah).
- **Kalkulasi Hash-Chain:** Sistem harus menghitung _hash_ setiap entri di _application layer_ (Golang) menggunakan rumus `SHA256(data_transaksi + hash_entri_sebelumnya)`.
- **Immutability:** Menyimpan entri _ledger_ dengan status 'terverifikasi' secara permanen; memblokir _update_ atau _delete_ pada tabel utama melalui _role-based access control_ di PostgreSQL.

### 3.3. ML Integration & Notification Service

- **Async ML Trigger:** Mengirim _event_ ke _message queue_ untuk memicu layanan _forecasting_ setiap kali ada entri _ledger_ baru.
- **Risk Engine:** Menyimpan dan menyajikan hasil perhitungan skor risiko (_Equivalent Monthly Income_/EMI).
- **Notification:** _Trigger webhook_ atau layanan notifikasi _push_ internal jika skor risiko melewati ambang batas SEDANG.

### 3.4. Income Passport & Consent Service

- **Document Generation:** Mengompilasi data _ledger_ valid (>30 hari), menghasilkan format PDF dalam waktu < 10 detik, dan menyematkan _digital signature_.
- **Consent Management:** Mengelola hierarki izin (kepada siapa, batas waktu, pencabutan izin akses) pada tingkat _database_.
- **Access Logging:** Mencatat secara permanen setiap permintaan akses verifikasi pihak ketiga ke dalam tabel log audit khusus.

---

## 4. Persyaratan Non-Fungsional (SLA & Keamanan Siber)

- **Performa & Latensi:** Waktu respons API untuk pencatatan transaksi (termasuk _hash calculation_) harus < 2 detik pada P95. _Throughput_ harus menahan >500 _concurrent users_ tanpa degradasi.
- **Keamanan Infrastruktur (CyberSecurity):** * Wajib memitigasi injeksi SQL dengan penggunaan fitur *prepared statements\* yang terintegrasi (seperti via pustaka GORM di Golang).
  - Validasi sanitasi _input_ yang ketat sebelum proses _hashing_.
  - Manajemen rahasia (_secrets management_) diletakkan di luar repositori kode (melalui _environment variables_ di _setup_ Docker).
- **Keamanan Data:** TLS 1.2+ wajib diaktifkan via Caddy untuk data _in-transit_. Enkripsi AES-256 untuk data sensitif pengguna _at-rest_.
- **Autentikasi Antar-Layanan:** Menggunakan JWT dengan rotasi kunci berkala untuk komunikasi _inter-service_.
- **Integritas Ledger:** _Backend_ harus mampu melakukan verifikasi ulang seluruh _hash-chain_ untuk 1 tahun data dalam waktu < 30 detik.
- **Kualitas Kode:** Minimal 80% _unit test coverage_ pada _backend_ dan dokumentasi API publik menggunakan spesifikasi OpenAPI 3.0.

---

## 5. Data Pipeline (Alur Pencatatan Transaksi)

_Alur sistem dari penerimaan data hingga pembaruan dashboard pengguna:_

1. **Akuisisi:** Input dari REST API (manual) atau HTTP POST _Webhook_ (e-wallet) diterima melalui _reverse proxy_ Caddy.
2. **Validasi & Deduplikasi:** Layanan Golang mengecek sanitasi _input_ dan menolak duplikasi berdasarkan _constraint_ PostgreSQL.
3. **Hashing:** Layanan Golang menghitung nilai kriptografis SHA-256.
4. **Persistensi:** Simpan struktur final ke PostgreSQL.
5. **Event Publish:** Mengirim beban kerja secara asinkron ke RabbitMQ.
6. **Processing (Python ML):** _Worker_ membaca _queue_, menghitung proyeksi EMI dengan model _time-series_.
7. **Cache Update:** Hasil komputasi di- _push_ ke Redis untuk akses _dashboard_ instan.
8. **Notifikasi:** Mesin aturan mengevaluasi _threshold_ risiko untuk _trigger_ peringatan dini.
