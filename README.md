# Nuxt + Go Project

โปรเจกต์ Full-stack พัฒนาด้วย **Nuxt 4** (Frontend) และ **Go / Gin** (Backend)

---

## 🛠 สิ่งที่ต้องติดตั้ง (Prerequisites)

1. **[Bun](https://bun.sh/)** (JavaScript Runtime & Package Manager สำหรับ Frontend)
   ```bash
   curl -fsSL https://bun.sh/install | bash
   ```
2. **[Go](https://go.dev/)** (แนะนำเวอร์ชัน 1.22 ขึ้นไป สำหรับ Backend)
3. **[PostgreSQL](https://www.postgresql.org/)** (ฐานข้อมูลสำหรับ Backend)

---

## ⚙️ การตั้งค่า Environment Variables (`.env`)

สร้างไฟล์ `.env` ไว้ที่โฟลเดอร์ Root ของโปรเจกต์ (`/Nuxt+Go/.env`):

```env
# Database Connection (PostgreSQL)
DATABASE_URL=postgres://username:password@localhost:5432/dbname

# Cloudflare R2 Storage (สำหรับเก็บไฟล์/รูปภาพ)
R2_ENDPOINT=your_cloudflare_account_id
R2_ACCESS_KEY_ID=your_r2_access_key_id
R2_SECRET_ACCESS_KEY=your_r2_secret_access_key
R2_BUCKET_NAME=your_bucket_name
R2_PUBLIC_DOMAIN=https://your-public-domain.com

# AI API Keys
OPENAI_API_KEY=your_openai_api_key
GEMINI_API_KEY=your_gemini_api_key
```

---

## 🚀 วิธีการรันโปรเจกต์ (Getting Started)

### 1. Frontend (Nuxt 4)

เปิด Terminal ที่ 1 แล้วเข้าไปที่โฟลเดอร์ `frontend`:

```bash
cd frontend

# ติดตั้ง dependencies ด้วย Bun
bun install

# เริ่มรันโหมด Development
bun run dev
```

> Frontend จะเริ่มทำงานที่: [http://localhost:3000](http://localhost:3000)

---

### 2. Backend (Go + Gin)

เปิด Terminal ที่ 2 แล้วเข้าไปที่โฟลเดอร์ `backend`:

```bash
cd backend

# ดาวน์โหลด Go modules ที่จำเป็น
go mod tidy

# เริ่มรัน Server
go run main.go
```

> Backend API จะเริ่มทำงานที่: [http://localhost:8080](http://localhost:8080)

---

## 📁 โครงสร้างโปรเจกต์ (Project Structure)

```text
├── frontend/           # Nuxt 4 (Vue 3, TailwindCSS, Bun)
│   ├── app/            # Vue Pages & Components
│   ├── server/         # Nitro Server Routes
│   ├── nuxt.config.ts  # Nuxt Configuration
│   └── package.json    # Frontend Dependencies
│
├── backend/            # Go Backend (Gin Framework)
│   ├── main.go         # Entry point & API Routes
│   ├── go.mod          # Go Dependencies
│   └── go.sum
│
├── .env.example        # ตัวอย่าง Environment Variables
├── .gitignore          # ระบุไฟล์ที่ไม่ต้องการนำขึ้น Git
└── README.md           # คู่มือโปรเจกต์
```
