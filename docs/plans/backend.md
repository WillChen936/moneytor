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
