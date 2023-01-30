# 介紹
這是我為了練習golang而建立的side project，模擬遊戲支付錢包服務。

可以透過API來：建立玩家、確認玩家是否存在、下注、派彩、取消下注、修改派彩金額。

# 使用到的資料庫

- MySQL：玩家的帳號、餘額都存在MySQL

- MongoDB：每一筆交易紀錄，都存在MongoDB

# 使用到的框架

- Gin
- Gorm
