# Backend Plan

## 開發順序

### 階段 1：User + 認證 ✅ 完成

### 階段 2：完善功能與修復問題
- 完成未完成功能（Transfer）

---

# Backend Issues

## 未完成功能

### Transfer 類型未完整實作

**位置**: `api/entries.go`, `database/sqlc/tx_create_entry.go`, `api/entries_test.go`

**問題**: `TransactionTypeTransfer = 3` 已定義，但 `resolverEntryAmount` 碰到它直接回傳 error；`tx_create_entry.go` 的 `ToAccountID` 永遠是 `Valid: false`，跨帳戶轉帳邏輯尚未實作。

#### Step 1：`database/sqlc/tx_create_entry.go`

- `CreateEntryTxParams` 加上 `ToAccountID pgtype.Int8`
- `CreateEntryTxResult` 加上 `ToAccount Account`
- `CreateEntryTx` 邏輯：
  - 建立 entry 時傳入 `arg.ToAccountID`（不再寫死 `Valid: false`）
  - 若 `arg.ToAccountID.Valid`，額外呼叫 `UpdateAccountBalance` 讓 to_account 餘額 `+amount`
  - **Deadlock 預防**：兩個帳戶永遠按 ID 大小順序更新（借鑑 simplebank `addMoney` 模式）

#### Step 2：`api/entries.go`

- `createEntryRequest` 加上 `ToAccountID int64 \`json:"toAccountId" binding:"omitempty,gt=0"\``
- `resolverEntryAmount` 新增 `TransactionTypeTransfer` case，直接回傳正數金額
- `createEntry` handler 新增 Transfer 驗證：
  - `ToAccountID == 0` → 400
  - `ToAccountID == AccountID` → 400
  - 呼叫 `GetAccount(ToAccountID, userID)` 確認帳戶存在且屬於此 user → 404 on error
- 建構 `CreateEntryTxParams` 時傳入 `ToAccountID`

#### Step 3：`api/entries_test.go`

- `InvalidTransactionTypeID` 改名為 `Transfer_MissingToAccountID`（不帶 toAccountId → 400）
- 新增 `Transfer_SameAccount`（toAccountId == accountId → 400）
- 新增 `OK_Transfer`（帶 toAccountId，mock GetAccount + CreateEntryTx → 200）
