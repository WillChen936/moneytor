# Moneytor 專案評估報告

以下為**尚待改善**項目。已改善項目已從清單移除，包含：main.go 未使用 import、CreateAccount FK 處理、CreateEntry 分類不存在改為 404、**函式名 createRandomTransactionType**、**測試案例名 IllegalCurrencyID**、**ListAccounts API（GET /api/v1/accounts）**。

---

## 一、API 與錯誤處理

### 1. 錯誤訊息直接回傳給客戶端（`api/server.go` 的 `errResponse`）

- **問題**：`errResponse(err)` 會把 `err.Error()` 直接放在 JSON 的 `error` 欄位，可能包含 SQL、路徑等內部資訊
- **建議**：
  - 對外：回傳通用、安全的訊息（例如「處理失敗」）
  - 詳細錯誤僅寫入 log（zerolog），並可依 `config.Env` 在 DEV 時才回傳詳細內容

---

## 二、業務邏輯與資料完整性

### 2. 帳戶餘額可能為負（`database/sqlc/tx_create_entry.go` + migrations）

- **Status (decided):**Negative balance is allowed (e.g. credit card accounts may be overdrawn). The reason for not adding `CHECK (balance >= 0)` is documented in the comment on `CreateEntryTx`; migrations remain without this constraint.

### 3. 建立 Entry 時未驗證金額正負（`api/entries.go`）

- **問題**：`createEntryRequest.Amount` 僅 `binding:"required"`，可傳 0 或負數；`ResolverEntryAmount` 只做正負號轉換，不擋 0
- **建議**：若業務上收入/支出金額應為正數，可加上 `binding:"gt=0"` 或自訂 validator，並在錯誤訊息說明「金額須大於 0」

---

## 三、設定與部署

### 4. 設定檔路徑依賴工作目錄（`main.go`）

- **問題**：`utils.LoadConfig("config.json")` 使用相對路徑，執行時若工作目錄不是專案根目錄（例如從其他目錄執行二進位），會找不到設定檔
- **建議**：
  - 支援環境變數覆寫路徑（例如 `CONFIG_PATH`）；或
  - 從執行檔所在目錄、或固定順序搜尋（例如當前目錄、專案根目錄）

### 5. 伺服器未實作優雅關閉（`api/server.go`）

- **問題**：`Start` 使用 `gin.Run(address)`，收到 SIGTERM/SIGINT 時會直接結束，不等待進行中請求完成
- **建議**：使用 `http.Server` + `ListenAndServe`，並在 `main` 中監聽 OS signal，收到後呼叫 `server.Shutdown(ctx)`，設定合理 timeout（例如 10 秒）

---

## 四、API 設計與一致性

### **6. API 參數命名風格不一致**

- **現狀（已修正）：**對外 API 已統一為 **camelCase**：query 參數（如 `accountId`、`pageId`、`pageSize`）與 JSON body 皆使用 camelCase。
- **影響：**前端只需處理一種風格；文件與 SDK 無需區分 query 與 body 的命名慣例。
- **建議：**已在 API 文件中明確說明對外一律使用 camelCase。

---

## 五、資料庫與錯誤處理細節

### 7. `db.ErrForeignKeyViolation` 的用途（`database/sqlc/errors.go`）

- **問題**：`ErrForeignKeyViolation` 是只設了 `Code` 的 `*pgconn.PgError`，其他欄位為零值。目前 API 是用 `ErrorCode(err) == db.ForeignKeyViolation` 判斷，這樣沒問題；但若有人寫 `errors.Is(err, db.ErrForeignKeyViolation)` 會不如預期（真實 DB 錯誤是另一個 `*pgconn.PgError` 實例）
- **建議**：
  - 若只打算用 `ErrorCode()` 判斷，可考慮移除 `ErrForeignKeyViolation` 變數，或在註解中說明「僅供測試或比對 Code，不要用 errors.Is」
  - 測試中需回傳 FK 錯誤時，可繼續用自建 `pgconn.PgError{Code: db.ForeignKeyViolation}` 或現有變數，並在文件註明用法

---

## 六、命名與程式風格

### 8. 函式命名（`api/entries.go`）

- **位置**：約第 32、99 行
- **問題**：`ResolverEntryAmount` 應為 `ResolveEntryAmount`（動詞 Resolve，而非 Resolver）
- **影響**：與常見「動詞 + 名詞」命名一致，較易理解

---

## 七、測試

### 9. 測試迴圈中的 `defer ctrl.Finish()`（`api/entries_test.go`、`api/accounts_test.go` 等）

- **現狀**：在 `for _, testCase := range testCases` 內使用 `defer ctrl.Finish()`，所有 defer 在函式結束時才執行
- **影響**：行為正確（每個 case 的 ctrl 都會被 Finish），但可讀性較差，且若未來在迴圈內加其他資源清理，容易混淆
- **建議**：改為在迴圈內搭配子測試 `t.Run(testCase.name, ...)` 並在子測試內 `defer ctrl.Finish()`，或在每輪結尾明確呼叫 `ctrl.Finish()`，意圖較清楚

---

## 八、小結


| 類型       | 數量  | 說明                        |
| -------- | --- | ------------------------- |
| 錯誤處理/安全  | 1   | 錯誤訊息外洩                    |
| 業務/資料完整性 | 2   | 負餘額、金額 ≤ 0                |
| 部署/運維    | 2   | 設定路徑、優雅關閉                 |
| API 設計   | 1   | 參數命名一致性                   |
| 程式庫用法    | 1   | ErrForeignKeyViolation 語意 |
| 命名/風格    | 1   | ResolveEntryAmount        |
| 測試風格     | 1   | defer 在迴圈內                |


整體架構（Gin + sqlc + pgx + 分層）清楚，測試也有使用 mock。優先建議先處理：**錯誤訊息不要直接回傳給客戶端**、以及是否要**禁止負餘額與零金額**；其餘可依產品需求與時程逐步調整。