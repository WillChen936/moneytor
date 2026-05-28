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
