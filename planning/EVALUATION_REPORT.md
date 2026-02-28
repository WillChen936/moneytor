# Moneytor 專案評估報告

**審視日期：** 2025-02-20  
**範圍：** 整體專案架構、API、資料庫、測試、設定與部署、程式風格與安全。

---

## 一、專案概覽

| 項目 | 說明 |
|------|------|
| 專案名稱 | Moneytor |
| 語言／框架 | Go 1.25、Gin、pgx、sqlc |
| 結構 | `main.go` → `api`（HTTP handlers）→ `database/sqlc`（Store/Queries）→ PostgreSQL |
| 測試 | `go test ./...` 通過（api、database/sqlc 有單元／整合測試，使用 mock） |

整體分層清楚：進入點 → API 層 → Store 介面 → sqlc 產生之查詢。依賴注入（Store）利於測試與替換實作。

---

## 二、已改善或已決策項目（與先前報告對照）

- **main.go 未使用 import**：已無多餘 import。
- **CreateAccount FK 處理**：currency_id 不存在時回傳 422，並可依 `ErrorCode(err) == db.ForeignKeyViolation` 判斷。
- **CreateEntry 分類不存在**：先 `GetCategory`，不存在時回傳 404。
- **函式命名 createRandomTransactionType**：已修正。
- **測試案例名 IllegalCurrencyID**：已採用。
- **ListAccounts API（GET /api/v1/accounts）**：已實作，分頁參數 `pageId`、`pageSize`。
- **API 參數命名**：對外已統一 **camelCase**（query 與 JSON body）。
- **CreateEntry 金額驗證**：`amount` 已加 `binding:"required,gt=0"`，不允許 ≤ 0。
- **帳戶餘額可為負**：已決策允許（如信用卡透支），不在 DB 加 `CHECK (balance >= 0)`，理由見 `CreateEntryTx` 註解。

---

## 三、尚待改善項目

### 3.1 錯誤訊息直接回傳給客戶端（`api/server.go` 的 `errResponse`）

- **問題**：`errResponse(err)` 將 `err.Error()` 直接放入 JSON 的 `error` 欄位，可能洩漏 SQL、路徑、內部實作等資訊。
- **建議**：
  - 對外：回傳通用、安全訊息（例如「處理失敗」或依情境「參數錯誤」）。
  - 詳細錯誤僅寫入 log（zerolog），並可依 `config.Env == "DEV"` 在開發環境才回傳詳細內容。

---

### 3.2 設定檔路徑依賴工作目錄（`main.go`）

- **問題**：`utils.LoadConfig("config.json")` 使用相對路徑，若執行時工作目錄不是專案根目錄（例如從其他目錄執行二進位），會找不到設定檔。
- **建議**：
  - 支援環境變數覆寫路徑（例如 `CONFIG_PATH`）；或
  - 依固定順序搜尋（當前目錄、執行檔所在目錄、專案根目錄等），並在文件說明。

---

### 3.3 伺服器未實作優雅關閉（`api/server.go`）

- **問題**：`Start` 使用 `gin.Run(address)`，收到 SIGTERM/SIGINT 時程式直接結束，不等待進行中請求完成。
- **建議**：改用 `http.Server` + `ListenAndServe`，在 `main` 中監聽 OS signal，收到後呼叫 `server.Shutdown(ctx)`，並設定合理 timeout（例如 10 秒）。

---

### 3.4 `ErrForeignKeyViolation` 的用途（`database/sqlc/errors.go`）

- **問題**：`ErrForeignKeyViolation` 為只設了 `Code` 的 `*pgconn.PgError`，其他欄位為零值。目前 API 以 `ErrorCode(err) == db.ForeignKeyViolation` 判斷，行為正確；但若有人使用 `errors.Is(err, db.ErrForeignKeyViolation)` 會不如預期（真實 DB 錯誤為另一實例）。
- **建議**：
  - 若僅以 `ErrorCode()` 判斷，可在註解中說明「僅供比對 Code，勿用 `errors.Is`」；或考慮移除該變數，改為常數與輔助函式。
  - 測試中需模擬 FK 錯誤時，可繼續使用自建 `pgconn.PgError{Code: db.ForeignKeyViolation}` 或現有變數，並在套件註解註明用法。

---

### 3.5 函式命名（`api/entries.go`）

- **位置**：約第 32、99 行。
- **問題**：`ResolverEntryAmount` 應為 `ResolveEntryAmount`（動詞 Resolve，而非 Resolver）。
- **影響**：與常見「動詞 + 名詞」命名一致，較易閱讀與搜尋。

---

### 3.6 測試迴圈中的 `defer ctrl.Finish()`（`api/*_test.go`）

- **現狀**：在 `for _, testCase := range testCases` 內使用 `defer ctrl.Finish()`，所有 defer 在函式結束時才執行，行為正確。
- **建議**：改為在迴圈內使用子測試 `t.Run(testCase.name, ...)` 並在子測試內 `defer ctrl.Finish()`，或在每輪結尾明確呼叫 `ctrl.Finish()`，意圖較清楚，也利於並行與篩選執行。

---

## 四、其他觀察與建議

### 4.1 文件與設定

- **README**：目前僅專案名稱，建議補充：專案簡介、如何執行（含 `make server`、DB 與 migrate）、主要 API 概覽、測試指令。
- **設定檔**：`config.json` 已在 `.gitignore`，建議新增 `config.example.json`（不含敏感值），方便新成員與部署參考。

### 4.2 API 設計

- 路由與參數風格一致（camelCase），健康檢查 `GET /api/v1/health` 存在。
- 若未來要對外提供文件，可考慮 OpenAPI/Swagger 或至少一份簡要 API 列表（路徑、方法、主要參數）。

### 4.3 資料庫與 sqlc

- migrations 與 sqlc 設定分離清楚，`sqlc.yaml` 使用 pgx、型別覆寫合理。
- `CreateEntryTx` 註解已說明負餘額為刻意允許，有助於後續維護。

---

## 五、小結

| 類型 | 數量 | 說明 |
|------|------|------|
| 錯誤處理／安全 | 1 | 錯誤訊息外洩（errResponse） |
| 部署／運維 | 2 | 設定路徑、優雅關閉 |
| 程式庫／語意 | 1 | ErrForeignKeyViolation 用法說明 |
| 命名／風格 | 1 | ResolveEntryAmount |
| 測試風格 | 1 | 迴圈內 defer ctrl.Finish |
| 文件／範例 | 2 | README、config 範例（建議） |

**優先建議**：先處理 **錯誤訊息不要直接回傳給客戶端**，再依需求排程處理 **設定路徑** 與 **優雅關閉**；其餘可依產品時程與可讀性需求逐步調整。整體架構與測試狀況良好，適合作為後續功能擴充的基礎。
