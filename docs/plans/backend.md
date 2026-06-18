# Backend Plan

## TODO

### 補上 zerolog 請求 log middleware

**背景：** commit `d580b42`（2026-05-31）把 `gin.Default()` 改成 `gin.New()` + `gin.Recovery()`，
順手拿掉了 Gin 內建的 `gin.Logger()`（因為它的格式跟 `main.go` 已設定的 zerolog 不一致）。
原本意圖是改用 zerolog 寫自訂的請求 log middleware，但這部分一直沒補上，
導致目前後端收到任何請求都不會印 log，無法得知誰打了哪支 API。

**要做：**
- 在 `api/` 寫一個 gin middleware，記錄每個請求的 method、路徑、狀態碼、耗時
- log 走 zerolog（與 `main.go` 的設定一致），不要用 `gin.Logger()`
- 在 `server.go` 用 `router.Use(...)` 掛上去（放在 `gin.Recovery()` 附近）

---

## 上 Production 前要改進的事

以下是學習階段為了簡化而暫時妥協、正式上線前需要補強的項目。

### 認證安全性強化

**現況（學習階段簡化版）：**
- 登入時 access token 和 refresh token 都放在 response body（JSON）回傳
- 前端把 access token 存在 `localStorage`

**問題：**
- `localStorage` 可被 XSS 讀取，token 有被竊取風險（前端側問題，但需後端配合修）

**要做：**
- 登入 / refresh endpoint 改用 `Set-Cookie` 下發 **refresh token**，並加上 `HttpOnly`、`Secure`、`SameSite` 屬性（JS 讀不到，防 XSS 竊取）
- access token 維持短命，由前端負責過期後呼叫 refresh
- 對應前端改動見 frontend plan 的「認證安全性強化」

### CORS 設定

**背景：** 本 session 為了讓後端維持純 API，移除了原本硬編在 `server.go` 的 `corsMiddleware`。
目前 dev 環境靠前端 Vite proxy 避開 CORS。

**正式版：**
- 若前後端**跨來源**部署（不同網域 / port）→ 後端需重新加上 CORS middleware，並用設定檔管理允許來源清單（不要硬編）
- 若**同源**部署（後端服務靜態檔，或前後端在同一反向代理後面）→ 不需要 CORS
