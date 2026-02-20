# Moneytor 專案評估報告

## 一、明顯錯誤與拼寫

### 1. 變數拼寫錯誤（`api/transactiontypes.go`）
- **位置**：約第 16、22 行
- **問題**：`transcationTypes` 應為 `transactionTypes`（少一個 s、多一個 c）
- **影響**：可讀性與命名一致性

### 2. 測試案例名稱拼寫錯誤（`api/accounts_test.go`）
- **位置**：約第 48 行
- **問題**：`IlleagalCurrnecyID` 應為 `IllegalCurrencyID`（Illegal、Currency 拼錯）
- **影響**：測試報告與搜尋時較難理解

### 3. 未使用的 import（`main.go`）
- **位置**：約第 14 行
- **問題**：`_ "github.com/jackc/pgx/v5/stdlib"` 僅用於 `database/sql` 驅動，專案使用 `pgxpool` 直接連線，未使用 `database/sql`
- **影響**：多餘依賴，可移除以保持乾淨

---

## 二、API 與錯誤處理

### 4. CreateAccount 未處理外鍵錯誤（`api/accounts.go`）
- **問題**：無效的 `currencyID`（不存在的貨幣）會讓 DB 回傳 FK 違規，目前一律回傳 `500 Internal Server Error`
- **建議**：與 `api/categories.go`、`api/entries.go` 一致，偵測 `db.ErrorCode(err) == db.ForeignKeyViolation` 時回傳 `422 Unprocessable Entity`，並可回傳較友善的訊息（例如「貨幣不存在」）

### 5. CreateEntry 中「分類不存在」的 HTTP 狀態碼（`api/entries.go`）
- **問題**：`GetCategory` 回傳 `sql.ErrNoRows`（分類不存在）時，目前回傳 `400 Bad Request`
- **建議**：資源不存在較符合 REST 慣例的是 `404 Not Found`，可將此情況改為 404

### 6. 錯誤訊息直接回傳給客戶端（`api/server.go` 的 `errResponse`）
- **問題**：`errResponse(err)` 會把 `err.Error()` 直接放在 JSON 的 `error` 欄位，可能包含 SQL、路徑等內部資訊
- **建議**：
  - 對外：回傳通用、安全的訊息（例如「處理失敗」）
  - 詳細錯誤僅寫入 log（zerolog），並可依 `config.Env` 在 DEV 時才回傳詳細內容

---

## 三、業務邏輯與資料完整性

### 7. 帳戶餘額可能為負（`database/sqlc/tx_create_entry.go` + migrations）
- **問題**：`UpdateAccountBalance` 為 `balance + amount`，支出時 `amount` 為負數，若餘額不足會產生負餘額；migrations 中沒有 `CHECK (balance >= 0)`
- **建議**：
  - 若業務不允許負餘額：在 `CreateEntryTx` 內先查詢帳戶餘額，若更新後會小於 0 則回傳業務錯誤（例如 422）；或是在 DB 加 `CHECK (balance >= 0)`
  - 若允許透支：在文件或註解中說明

### 8. 建立 Entry 時未驗證金額正負（`api/entries.go`）
- **問題**：`createEntryRequest.Amount` 僅 `binding:"required"`，可傳 0 或負數；`ResolverEntryAmount` 只做正負號轉換，不擋 0
- **建議**：若業務上收入/支出金額應為正數，可加上 `binding:"gt=0"` 或自訂 validator，並在錯誤訊息說明「金額須大於 0」

---

## 四、設定與部署

### 9. 設定檔路徑依賴工作目錄（`main.go`）
- **問題**：`utils.LoadConfig("config.json")` 使用相對路徑，執行時若工作目錄不是專案根目錄（例如從其他目錄執行二進位），會找不到設定檔
- **建議**：
  - 支援環境變數覆寫路徑（例如 `CONFIG_PATH`）；或
  - 從執行檔所在目錄、或固定順序搜尋（例如當前目錄、專案根目錄）

### 10. 伺服器未實作優雅關閉（`api/server.go`）
- **問題**：`Start` 使用 `gin.Run(address)`，收到 SIGTERM/SIGINT 時會直接結束，不等待進行中請求完成
- **建議**：使用 `http.Server` + `ListenAndServe`，並在 `main` 中監聽 OS signal，收到後呼叫 `server.Shutdown(ctx)`，設定合理 timeout（例如 10 秒）

---

## 五、API 設計與一致性

### 11. 缺少「列出帳戶」API
- **現狀**：有 `ListCategories`、`ListCurrencies`、`ListEntries`，但沒有 `ListAccounts`；僅能建立帳戶，無法查詢帳戶列表
- **影響**：前端若要顯示帳戶清單、或讓使用者在建立 Entry 時選擇帳戶，會缺少對應 API
- **建議**：在 `database/queries/accounts.sql` 新增 `ListAccounts`（可支援分頁），並在 `api/server.go` 註冊 `GET /api/v1/accounts`

---

## 六、資料庫與錯誤處理細節

### 12. `db.ErrForeignKeyViolation` 的用途（`database/sqlc/errors.go`）
- **問題**：`ErrForeignKeyViolation` 是只設了 `Code` 的 `*pgconn.PgError`，其他欄位為零值。目前 API 是用 `ErrorCode(err) == db.ForeignKeyViolation` 判斷，這樣沒問題；但若有人寫 `errors.Is(err, db.ErrForeignKeyViolation)` 會不如預期（真實 DB 錯誤是另一個 `*pgconn.PgError` 實例）
- **建議**：
  - 若只打算用 `ErrorCode()` 判斷，可考慮移除 `ErrForeignKeyViolation` 變數，或在註解中說明「僅供測試或比對 Code，不要用 errors.Is」
  - 測試中需回傳 FK 錯誤時，可繼續用自建 `pgconn.PgError{Code: db.ForeignKeyViolation}` 或現有變數，並在文件註明用法

---

## 七、測試

### 13. 測試迴圈中的 `defer ctrl.Finish()`（`api/entries_test.go`、`api/accounts_test.go`）
- **現狀**：在 `for _, testCase := range testCases` 內使用 `defer ctrl.Finish()`，所有 defer 在函式結束時才執行
- **影響**：行為正確（每個 case 的 ctrl 都會被 Finish），但可讀性較差，且若未來在迴圈內加其他資源清理，容易混淆
- **建議**：改為在迴圈內直接 `defer ctrl.Finish()` 搭配子測試 `t.Run(testCase.name, ...)`，或在每輪結尾明確呼叫 `ctrl.Finish()`，意圖較清楚

---

## 八、小結

| 類型               | 數量 | 說明                                   |
|--------------------|------|----------------------------------------|
| 拼寫/明顯錯誤      | 3    | 變數名、測試名、未使用 import          |
| 錯誤處理/狀態碼    | 3    | Account FK、Entry 404、錯誤訊息外洩    |
| 業務/資料完整性   | 2    | 負餘額、金額 ≤ 0                       |
| 部署/運維          | 2    | 設定路徑、優雅關閉                     |
| API 設計           | 1    | 缺少 ListAccounts                      |
| 程式庫用法         | 1    | ErrForeignKeyViolation 語意            |
| 測試風格           | 1    | defer 在迴圈內                         |

整體架構（Gin + sqlc + pgx + 分層）清楚，測試也有使用 mock。優先建議先處理：**拼寫與未使用 import**、**CreateAccount 的 FK 處理**、**錯誤訊息不要直接回傳給客戶端**、以及是否要**禁止負餘額與零金額**；其餘可依產品需求與時程逐步調整。
