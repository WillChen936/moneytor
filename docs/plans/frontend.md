# Frontend Plan

## 技術棧

- **框架：** React + Vite + TypeScript
- **UI 元件庫：** shadcn/ui（基於 Tailwind CSS）
- **路由：** react-router-dom
- **HTTP Client：** axios
- **位置：** `/frontend`（monorepo，與後端同一個 repo）

---

## 開發階段

### Phase 1 — 環境建置
**目標：** 能在瀏覽器看到一個 React 頁面跑起來

- 用 Vite 建立 React + TypeScript 專案到 `/frontend`
- 確認開發伺服器可以啟動（`npm run dev`）
- 清理 Vite 預設樣板，留下乾淨的起點

---

### Phase 2 — 加入路由
**目標：** 能在不同 URL 之間切換頁面（但頁面內容暫時是空白）

- 安裝 `react-router-dom`
- 建立以下幾個空白頁面元件：`LoginPage`、`DashboardPage`、`EntriesPage`
- 設定路由：`/login` → LoginPage，`/` → DashboardPage，`/entries` → EntriesPage

---

### Phase 3 — 加入 UI 元件庫
**目標：** 有基本的視覺樣式，不用自己寫 CSS

- 安裝 `shadcn/ui`（基於 Tailwind CSS）
- 設定 Tailwind
- 在 LoginPage 用 shadcn 元件排出一個登入表單的外觀（暫時不串 API）

---

### Phase 4 — 串接登入 API
**目標：** 填入帳號密碼後，能真的打到後端 API 拿到 token

- 安裝 `axios`，建立 API client（設定 base URL）
- 實作登入表單的 `onSubmit`，呼叫 `POST /api/v1/users/login`
- 把拿到的 `access_token` 存到 `localStorage`
- 登入成功後跳轉到 DashboardPage

---

### Phase 5 — 全域認證狀態
**目標：** App 記得你有沒有登入，沒登入就擋在 LoginPage

- 用 React `Context` 建立 `AuthContext`，存放 token 和登入狀態
- 建立 `ProtectedRoute` 元件：未登入就 redirect 到 `/login`
- 把 DashboardPage 和 EntriesPage 包進 ProtectedRoute

---

### Phase 6 — Dashboard 頁面
**目標：** 登入後能看到所有帳戶和餘額

- 打 `GET /api/v1/accounts` 拿帳戶列表
- 用卡片方式顯示每個帳戶的名稱和餘額
- 加入「新增帳戶」按鈕和簡單的 Modal 表單（串接 `POST /api/v1/accounts`）

---

### Phase 7 — 收支記錄頁面
**目標：** 能看到收支明細，並新增一筆記錄

- 打 `GET /api/v1/entries` 顯示記錄列表（日期、分類、金額、帳戶）
- 新增記錄的表單（帳戶、分類、金額、類型選擇）
- 串接 `POST /api/v1/entries`

---

### Phase 8 — 轉帳功能
**目標：** 能在兩個帳戶之間轉帳

- 新增轉帳表單（來源帳戶、目標帳戶、金額）
- 串接 `POST /api/v1/transfers`

---

### Phase 9 — 導覽列與整體佈局
**目標：** App 有一致的外觀和導覽

- 建立 `Layout` 元件，包含側邊欄或頂部導覽
- 套用到所有頁面
- 加入登出功能（清除 token，跳回 login）

---

## 備註

- 每個 Phase 只專注一個概念，完成並確認能跑起來再進下一個
- 後端缺少的 API（Update / Delete / 統計）等前端做到需要時再補

---

## 上 Production 前要改進的事

以下是學習階段為了簡化而暫時妥協、正式上線前需要補強的項目。

### 認證安全性強化

**現況（學習階段簡化版）：**
- access token 直接存在 `localStorage`
- 只用 access token，沒有實作 token 過期後自動 refresh
- token 欄位用 `data.accessToken`（後端回傳在 response body）

**問題：**
- `localStorage` 可被 XSS 讀取，token 有被竊取風險
- access token 過期後沒有 refresh 機制，使用者會突然被登出

**要做：**
- refresh token 改由後端用 HttpOnly Cookie 下發（JS 讀不到，較安全），需後端配合（見 backend plan）
- access token 維持短命；可考慮放記憶體而非 localStorage
- 在 `src/lib/api.ts` 的 interceptor 加上「access token 過期 → 自動用 refresh token 換新 token → 重送原請求」的邏輯
- 後端已同時回傳 access / refresh token（`loginResponse`），基礎已具備

### 開發環境的 CORS 取巧

**現況：** dev 用 Vite proxy（`vite.config.ts` 的 `server.proxy`）讓前端同源打 `/api`，避開 CORS。

**正式版：** 若前後端跨來源部署（不同網域/port），後端需正式設定 CORS 允許清單（即本 session 移除的 `corsMiddleware`）；若同源部署（後端服務靜態檔或反向代理）則不需要。
