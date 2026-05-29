# Backend Plan

## 開發順序

### 階段 1：User + 認證

**設計決策**
- Token 方式：Bearer Token（Header），適合前後端分離架構
- Refresh Token：實作，存於 `sessions` 表（DB 優先，之後遷移至 Redis）
- Categories：per-user，建立 User 時用 transaction 插入預設分類

**實作步驟**
1. 設計 ERD（`docs/erd.md`）
2. DB Migration：新增 `users`、`sessions` 表；`accounts`/`categories` 加 `user_id` FK
3. 更新 sqlc queries：新增 users/sessions query，現有 query 補 `user_id` 條件
4. `make sqlc` 重新生成程式碼
5. 實作 `POST /users`（註冊，bcrypt hash 密碼）
6. 實作 `POST /users/login`（登入，回傳 access token + refresh token）
7. 實作 `POST /users/refresh`（換發 access token）
8. 實作 JWT middleware，受保護路由套上 middleware
9. 更新現有 handlers，從 context 取 `user_id` 帶入 query
10. `make mockdb` 重新生成 mock，補齊測試

### 階段 2：完善功能與修復問題
- 修掉下方 Bug 清單
- 完成未完成功能（Transfer）

---

# Backend Issues

## Bug

### 1. `createAccount` — `balance` 不允許為 0
- **位置**: `api/accounts.go:12`
- **問題**: binding tag 用 `gt=0`，導致餘額為 0 的帳戶無法建立
- **修法**: 改為 `gte=0`

### 2. `listAccounts` / `listCategories` / `listEntries` — `pageSize` 最小值限制不合理
- **位置**: `api/accounts.go:43`, `api/categories.go:41`, `api/entries.go:59`
- **問題**: `min=5` 強制每次至少取 5 筆，語意上不正確
- **修法**: 改為 `min=1`

### 3. `ResolverEntryAmount` 不應 export
- **位置**: `api/entries.go:100`
- **問題**: 僅在 `api` package 內部使用，不應對外暴露
- **修法**: 改為小寫 `resolverEntryAmount`，同步更新測試

## 未完成功能

### 4. Transfer 類型未完整實作
- **位置**: `api/entries.go:107`, `database/sqlc/tx_create_entry.go:34`
- **問題**: `TransactionTypeTransfer = 3` 已定義，但 `resolverEntryAmount` 碰到它直接回傳 error；`tx_create_entry.go` 的 `ToAccountID` 永遠是 `Valid: false`，跨帳戶轉帳邏輯尚未實作
- **修法**: 實作 Transfer 的金額邏輯（from_account 扣款、to_account 入款）及對應的 transaction
