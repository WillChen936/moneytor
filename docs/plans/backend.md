# Backend Plan

## 開發順序

### 階段 1：User + 認證 ✅ 完成

### 階段 2：完善功能與修復問題
- 完成未完成功能（Transfer）

---

# Backend Issues

## 未完成功能

### Transfer 類型未完整實作
- **位置**: `api/entries.go`, `database/sqlc/tx_create_entry.go`
- **問題**: `TransactionTypeTransfer = 3` 已定義，但 `resolverEntryAmount` 碰到它直接回傳 error；`tx_create_entry.go` 的 `ToAccountID` 永遠是 `Valid: false`，跨帳戶轉帳邏輯尚未實作
- **修法**: 實作 Transfer 的金額邏輯（from_account 扣款、to_account 入款）及對應的 transaction
