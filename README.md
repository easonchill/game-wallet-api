
[TOC]

# 介紹
這是我為了練習golang而建立的side project，模擬一個遊戲支付錢包服務。

可以透過API來：建立玩家、確認玩家是否存在、下注、派彩、取消下注、修改派彩金額。

## 使用到的框架

- Gin
- Gorm

## 使用到的資料庫

- MySQL：玩家的帳號、餘額都存在MySQL

- MongoDB：每一筆交易紀錄，都存在MongoDB
  
# 前置設定
## Config
在`/config/config.yaml`裡可以設定MySQL及服務的參數

備註：因該專案為個人實驗用，放在這裡是方便修改，實務上應該避免暴露
設定，要用其他方式如 consul 等來儲存。

### MySQL
| Key | Describe |
| -------- | -------- |
|host|MySQL的主機位置|
|port|MySQL的PORT號|
|account|MySQL的使用者帳號|
|password|MySQL的登入密碼|
|DB名稱|使用的DB名稱|

### Gin
| Key | Describe |
| -------- | -------- |
|wtoken|每支API會驗證的wtoken|
|port|gin啟動時使用的port號|
|mode|gin啟動的模式|


# API Docs
## 路由

| Method | URL | Describe |
| -------- | -------- | -------- |
|POST|/swclient/test/setaccount|建立新玩家
| GET     | /player/check/:account     |  檢查玩家帳號是否存在|
|GET|/transaction/balance/:account|取得玩家錢包餘額
|GET|/transaction/record/:mtcode|查詢交易紀錄
|POST| /transaction/game/bets|批次下注
|POST| /transaction/game/wins|批次派彩
|POST| /transaction/game/refunds|批次下注退款
|POST|/transaction/game/cancel|用來取消注單refund狀態的注單，使其可以正常派彩
|POST|/transaction/game/amend|單人修改派彩結果
|POST|/transaction/game/amends|批次多人修改派彩結果

## CreatePlayer
### 建立新玩家

**Request**

URL: `/swclient/test/setaccount`

Method: `GET`

Headers :<br/>

`wtoken:config裡設定的token`<br/>


**Path Variables** 

| 參數     |  型別  | 必填 | 敘述 |
| - | :-: | :-: | :-|
| account     	  | string| 必填 |玩家帳號|
| password     	  | string| 必填 |玩家密碼|
|balance| number| 必填 | 玩家初始餘額|
| currency     	  | string| 必填 |幣別|

## CheckPlayer
### 檢查玩家帳號是否存在

**Request**

URL: `/player/check/:account`

Method: `GET`

Headers :<br/>

`wtoken:config裡設定的token`<br/>


**Path Variables** 

| 參數 | 型別 | 必填 | 敘述 |
| - | :-: | :-: | :-|
| account     	  | string| 必填 |使用者帳號 <br> |

**Sample Request**
```bash=
curl --location --request GET '{Wurl}/player/check/Testplayer' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--header 'wtoken: {Wtoken}'
```

**Response Format** 

| 參數名稱 | 參數類型 | 描述 |
| --- |:-:| --- |
| data      |bool | 玩家帳號存在時  true=是，false=否  |

**Sample Response**

***若為玩家存在時***：
```json=
{
  "data":true,
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2023-01-20T01:14:48-04:00"
  }
}
```
***若玩家不存在時***：
```json=
{
  "data":false,
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2023-01-20T01:14:48-04:00"
  }
}
```

## Balance
###  取得玩家錢包餘額

**Request**

URL: `/transaction/balance/:account`

Method: `GET`

Headers :<br/>

`wtoken:config裡設定的token`<br/>

**Path Variables** 
|參數名稱|型別|必填|描述|
|-|:-:|:-:|:-|
| account     	  | string   | 必填 | 使用者帳號<br/> |

**Query Params** 
|參數名稱|型別|必填|描述|
|-|:-:|:-:|:-|
| gamecode     	  | string   | 選填 | 遊戲代號<br/> |

**Sample Request**

```bash=
curl --location --request GET '{Wurl}/transaction/balance/Testplayer?gamecode={gamecode}' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--header 'Wtoken: {Wtoken}'
```

**Response Format** 

| 參數名稱 | 型別 | 描述 |
| --- |:-:| --- |
| balance      |number | 玩家餘額|
| currency     |string | 幣別|
| status.code   |string | 狀態編碼  |
| status.message   |string | 狀態訊息 |
| status.datetime   |string | 回傳時間  |

**Sample Response**

```json=
{
  "data": {
    "balance": 600270,
    "currency": "CNY"
  },
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2023-01-20T01:14:48.9266001-04:00"
  }
}
```
