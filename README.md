
[TOC]

# 介紹
這是為了練習golang而建立的side project，模擬一個遊戲支付錢包服務。

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
| Key | description|
| -------- | -------- |
|host|MySQL的主機位置|
|port|MySQL的Port|
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

| Method | URL                           | Describe                                       |
| ------ | ----------------------------- | ---------------------------------------------- |
| POST   | /swclient/test/setaccount     | 建立新玩家                                     |
| GET    | /player/check/:account        | 檢查玩家帳號是否存在                           |
| GET    | /transaction/balance/:account | 取得玩家錢包餘額                               |
| GET    | /transaction/record/:mtcode   | 查詢交易紀錄                                   |
| POST   | /transaction/game/bets        | 批次下注                                       |
| POST   | /transaction/game/wins        | 批次派彩                                       |
| POST   | /transaction/game/refunds     | 批次下注退款                                   |
| POST   | /transaction/game/cancel      | 用來取消注單refund狀態的注單，使其可以正常派彩 |
| POST   | /transaction/game/amends      | 批次多人修改派彩結果                           |
| POST   | /transaction/game/amend      | 單人修改派彩結果                           |


---


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
    "currency": "Coin"
  },
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2023-01-20T01:14:48.9266001-04:00"
  }
}
```

## Record
### 查詢交易紀錄
<font color="red">**每個行為的Record範例，請參考Record Sample** </font> 

**Request**
    
URL: `/transaction/record/:mtcode`

Method: `GET`

Headers :<br/>

`wtoken:config裡設定的token`<br/>

**Path Variables** 

| 參數名稱   | 型別| 必填  | 描述| 
| ----------- | -------- | -------- | -----------|
| mtcode     	  | string   | 必填 | 交易代碼 <font color="red"><br/>※請輸入欲查詢的mtcode </font><br/>

**Sample Request**

```bash=
curl --location --request GET '{wurl}/transaction/record/{mtcode}' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--header 'wtoken: {wtoken}'
```

**Response Format**
| 參數 | 型別 | 敘述 |
| --- | :-: | --- |
| _id         |string | 交易紀錄編號 |
| action          |string | 該交易的動作 如 bets、wins、 payoff |
| target.account    |string | 使用者帳號 |
| status.createtime    |string | 交易開始時間 |
| status.endtime    |string | 交易結束時間 |
| status.status    |string | 交易狀態  
| before    |number | 交易前餘額<br/>|
| balance    |number | 交易後餘額<br/> |
| currency    |string | 幣別<br/>|
| event    |array | 事件的陣列，會回傳多個mtcode <font color="red">全部的mtcode皆是紀錄在該陣列中</font>  |
| event.mtcode    |string | 交易代碼，為唯一不重複的值 |
| event.amount    |number | 該筆交易的金額<br/><font color="red"></font> |
| event.eventtime    |string | 我方發送時間|
| status.code   |string | 狀態編碼|
| status.message   |string | 狀態訊息  |
| status.datetime   |string | 回傳時間  |

**Sample Record Success**

```json=

{
  "data": {
    "_id": "59672a547aa48000019260cf",
    "action": "bet",
    "target": {
      "account": "Eason"
    },
    "status": {
      "createtime": "2017-07-13T04:07:48.644-04:00",
      "endtime": "2017-07-13T04:07:48.673-04:00",
      "status": "success",
      "message": "success"
    },
    "before": 8164082.95,
    "balance": 8164072.95,
    "currency": "Coin",
    "event": [
      {
        "mtcode": "testbet1123456:TEST",
        "amount": 10,
        "eventtime": "2022-10-05T05:08:41-04:00"
      }
    ]
  },
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2017-07-13T04:08:02-04:00"
  }
}
```

**Sample Record not found**

```json=
{
  "data": null,
  "status": {
    "code": "1014",
    "message": "record not found",
    "datetime": "2017-07-13T04:08:02-04:00"
  }
}
```

## Batch Bets
### 批次下注

:::warning
※單一玩家，批次裡可以不同round和mtcode。每一批次最多200筆 </br>
※只有全部成功或全部失敗，沒有『部份』成功或失敗的情境 
:::


**Request**
    
URL: `/transaction/game/bets`

Method: `POST`

Headers :<br/>

`Content-Type:application/json`<br/>
`wtoken:config裡設定的token`<br/>



**Request Parameter** 
| 參數名稱   | 型別     | 必填  | 描述|
|-| :-:|:-:|:-|
| account     	  | string   | 必填 | 使用者帳號<br/>|
| session     	  | string | 選填 | SessionID|
| gamehall     	  | string   | 必填 | 遊戲廠商代號 |
| gamecode     	  | string   | 必填 | 遊戲代號|
| data       	  | text array   | 必填 | 事件資料列表用JSON包起來<br/>長度不限，但每筆win最大長度為158字元<br/>mtcode:交易代碼<br/>amount: 金額(實際從玩家錢包扣除的金額)<br/>eventtime:交易當下時間 <br/> roundid:局號 |
| createTime     	  | string   | 必填 |成單時間 <br/> 此時間可與注單的createtime對應<br/> |
| genre     	  | string   | 選填 | 遊戲項目<br/>

 **Sample Request :**

```bash=
curl --location --request POST '{Wurl}/transaction/game/bets' \
--header 'Content-Type: application/json' \
--header 'wtoken: {Wtoken}' \
--data-raw '{
  "account": "{account1}",
  "gamehall": "{gamehall}",
  "gamecode": "{gamecode}",
  "session": "{session}",
  "genre": "{genre}" ,
  "data": [
    {
      "mtcode": "{mtcode1}",
      "amount": 10,
      "roundid": "{round1}",
      "eventtime": "2020-05-20T05:20:00-04:00"
    },
    {
      "mtcode": "{mtcode2}",
      "amount": 20,
      "roundid": "{round2}",
      "eventtime": "2020-05-20T05:20:00-04:00"
    },
    {
      "mtcode": "{mtcode3}",
      "amount": 30,
      "roundid": "{round3}",
      "eventtime": "2020-05-20T05:20:00-04:00"
    }
   
  ],
  "createTime": "2020-05-20T05:20:00-04:00"
}'
```
**Response Format** 

| 參數名稱 | 型別 | 描述 |
| --- |:-:| --- |
| balance      |number | 執行動作後的餘額<br/> |
| currency     |string | 幣別<br/>|
| status.code   |string | 狀態編碼 |
| status.message   |string | 狀態訊息  |
| status.datetime   |string | 回傳時間  |

**Sample Response** 
```json=
{
  "data": {
    "balance": 600250,
    "currency": "Coin"
  },
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2023-01-19T22:15:21.6076001-05:00"
  }
}

```

<a href="##Record" title="Transfer-Wallet-GameBoy"> **Sample Record - Batch bets**</a>

```json=
{
  "data": {
    "_id": "59672a547aa48000019260cf",
    "action": "bets",
    "target": {
      "account": "{account1}"
    },
    "status": {
      "createtime": "2017-07-13T04:07:48.644-04:00",
      "endtime": "2017-07-13T04:07:48.673-04:00",
      "status": "success",
      "message": "success"
    },
    "before": 600310,
    "balance": 600250,
    "currency": "Coin",
    "event": [
      {
        "mtcode": "{mtcode1}",
        "amount": 10,
        "eventtime": "2020-05-20T05:20:00-04:00"
      },
      {
        "mtcode": "{mtcode2}",
        "amount": 20,
        "eventtime": "2020-05-20T05:20:00-04:00"
      },
      {
        "mtcode": "{mtcode3}",
        "amount": 30,
        "eventtime": "2020-05-20T05:20:00-04:00"
      }
    ]
  },
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2023-01-13T04:08:02-04:00"
  }
}
```

##  Wins
### <font color="red">**多人**</font>多派彩。*如該局是輸的，wins送amount = 0*
:::warning
※多人多批次，每一批次(ucode)最多200筆
※有『部份』成功或失敗的情境。請參考reponse範例
:::

**Request**

URL: `/transaction/game/wins` 
 
Method: `POST`

Headers :<br/>

`Content-Type:application/json`<br/>
`wtoken:config裡設定的token`<br/>

**Request Parameter**

| 參數名稱   | 型別     | 必填  | 描述|
|-| :-:|:-:|:-|
| account | string  |必填| 玩家帳號<br/>※最大長度為36字元 | 
| event       	  | text array    | 必填 | 事件資料列表用JSON包起來<br/>長度不限，但每筆win最大長度為158字元<br/>mtcode:交易代碼<br/>amount:金額(實際給玩家錢包的金額)<br/>roundid:注單號<br/>eventtime:交易當下時間<br/>gamecode:遊戲編碼<br/>gamehall:遊戲廠商代號 <br/> validbet:有效投注| 
| eventtime | string  |必填| 該批交易的時間 | 
| ucode     	  | string   | 必填 | 交易批次代碼 ，回傳派彩結果時需帶此資訊  |

**Sample Request :**
```bash=
curl --location --request POST '{wurl}/transaction/game/wins' \
--header 'Content-Type: application/json' \
--header 'wtoken: {wtoken}' \
--data-raw '{
  "list": [
    {
      "account": "{account1}",
      "eventtime": "2022-04-27T04:33:26-04:00",
      "ucode": "{ucode1}",
      "event": [
        {
          "mtcode": "{mtcode1}",
          "amount": 10,
          "validbet":10,
          "roundid": "{round1}",
          "eventtime": "2022-04-27T04:33:26-04:00",
          "gamecode": "{gamecode}",
          "gamehall": "{gamehall}"
        },
        {
          "mtcode": "{mtcode2}",
          "amount": 20,
          "validbet":20,
          "roundid": "{round2}",
          "eventtime": "2022-04-27T04:33:26-04:00",
          "gamecode": "{gamecode}",
          "gamehall": "{gamehall}"
        },
        {
          "mtcode": "{mtcode3}",
          "amount": 30,
          "validbet":30,
          "roundid": "{round3}",
          "eventtime": "2022-04-27T04:33:26-04:00",
          "gamecode": "{gamecode}",
          "gamehall": "{gamehall}"
        }
      ]
    },
    {
      "account": "{account2}",
      "eventtime": "2022-04-27T04:33:26-04:00",
      "ucode": "{ucode2}",
      "event": [
        {
          "mtcode": "{mtcode4}",
          "amount": 10,
          "validbet":10,
          "roundid": "{round4}",
          "eventtime": "2022-04-27T04:33:26-04:00",
          "gamecode": "{gamecode}",
          "gamehall": "{gamehall}"
        },
        {
          "mtcode": "{mtcode5}",
          "amount": 20,
          "validbet":20,
          "roundid": "{round5}",
          "eventtime": "2022-04-27T04:33:26-04:00",
          "gamecode": "{gamecode}",
          "gamehall": "{gamehall}"
        },
        {
          "mtcode": "{mtcode6}",
          "amount": 30,
          "validbet":30,
          "roundid": "{round6}",
          "eventtime": "2022-04-27T04:33:26-04:00",
          "gamecode": "{gamecode}",
          "gamehall": "{gamehall}"
        }
      ]
    }
  ]
}'
```

**Response Format**

| 參數名稱 | 型別 | 描述 |
| :- |:-:| :- |
| data.success.account  <br/>(data.failed.account )     |string | 玩家帳號 |
|  data.success.ucode  <br/>(data.failed.ucode )    |string | 交易批碼，回傳派彩結果時需附代此資訊<br/> 成功 Success 失敗 Failed  |
|  data.success.balance      |number | 執行後的餘額<br/> |
|  data.success.currency     |string | 幣別<br/>|
| data.failed.code      |string | 錯誤代碼 |
| data.failed.message      |string | 錯誤訊息 |
| status.code   |string | 狀態編碼 |
| status.message   |string | 狀態訊息  |
| status.datetime   |string | 回傳時間  |

**Sample Response**
```json=
{
  "data":{
      "success":[
        {
            "account":"{account1}" 
            "balance":600280 ,
            "currency":"Coin",
            "ucode": "{ucode1}"
        }
      ],
      "failed":[
        {
            "account":"{account2}" ,
            "code":"transaction:1003" ,
            "message":"xxxxxx" ,
            "ucode": "{ucode2}"
        }
      ],
  },
  "status": {
      "code": "0",
      "message": "Success",
      "datetime": "2022-01-27T04:33:26-04:00"
  }
}
```

## Refunds
### 批次押注退款

:::warning
※批次退款為單一玩家。每一批次最多200筆 </br>
※只有全部成功或全部失敗，沒有『部份』成功或失敗的情境 
:::
**Request**
    
URL: `/transaction/game/refunds`

Method: `POST`

Headers :<br/>


`Content-Type:application/json`<br/>
`wtoken:config裡設定的token`<br/>


**Request Parameter** 
| 參數名稱   | 型別     | 必填  | 描述|
|-| :-:|:-:|:-|
| mtcode     	  | string array   | 必填 | 交易代碼  <br> <font color="red">※使用的mtcode為欲被退還的mtcode </font> |
    
**Sample Request** 


```bash=
curl --location --request POST '{Wurl}/transaction/game/refunds' 
--header 'Content-Type: application/json' 
--header 'wtoken:{Wtoken}
--data-raw '{
    "mtcode": [
        "{{mtcode1}}",
        "{{mtcode2}}"
    ]
}'
```


**Response Format** 

| 參數名稱 | 型別 | 描述 |
| --- |:-:| --- |
| balance      |number | refund後的餘額<br/><font color="red"></font> |
| currency     |string | 幣別|
| status.code   |string | 狀態編碼|
| status.message   |string | 狀態訊息  |
| status.datetime   |string | 回傳時間  |

**Sample Response** 
```json=
{
  "data": {
    "balance": 600280,
    "currency": "Coin"
  },
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2017-01-20T01:14:48.9266001-04:00"
  }
}

```


<a href="##Record" title="Transfer-Wallet-GameBoy"> ***Sample Record - Batch bets Refunds***</a>

```json=

{
  "data": {
    "_id": "59672a547aa48000019260cf",
    "action": "bets",
    "target": {
      "account": "{account1}"
    },
    "status": {
      "createtime": "2017-07-13T04:07:48.644-04:00",
      "endtime": "2017-07-13T04:07:48.673-04:00",
      "status": "success",
      "message": "success"
    },
    "before": 600250,
    "balance": 600280,
    "currency": "Coin",
    "event": [
      {
        "mtcode": "{mtcode1}",
        "amount": 10,
        "eventtime": "2020-05-20T05:20:00-04:00",
        "status": "refund"
      },
      {
        "mtcode": "{mtcode2}",
        "amount": 20,
        "eventtime": "2020-05-20T05:20:00-04:00"
        "status": "refund",
      },
      {
        "mtcode": "{mtcode3}",
        "amount": 30,
        "eventtime": "2020-05-20T05:20:00-04:00"
      }
    ]
  },
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2023-07-13T04:08:02-04:00"
  }
}
```

## Cancel
### 取消refund狀態，可進行派彩
:::warning
※批次取消為單一玩家。每一批次最多200筆 </br>
※只有全部成功或全部失敗，沒有『部份』成功或失敗的情境 
:::
**Request**
    
URL: `/transaction/game/cancel `

Method: `POST`

Headers :<br/>


`Content-Type:application/json`<br/>
`wtoken:config裡設定的token`<br/>


**Request Parameter** 
| 參數名稱   | 型別     | 必填  | 描述|
|-| :-:|:-:|:-|
| mtcode     	  | string array   | 必填 | 交易代碼  <br> <font color="red">※使用的mtcode為取消退還的mtcode</font> <br/> |

**Sample Request :**
```bash=
curl --location --request POST '{Wurl}/transaction/game/cancel' 
--header 'Content-Type: application/json' 
--header 'wtoken:{Wtoken}
--data-raw '{
    "mtcode": [
        "{{mtcode1}}"
    ]
}'
```

**Response Format** 

| 參數名稱 | 型別 | 描述 |
| :- |:-:| :- |
| balance      |number | cancel後的餘額<br/> |
| currency     |string | 幣別|
| status.code   |string | 狀態編碼 |
| status.message   |string | 狀態訊息  |
| status.datetime   |string | 回傳時間  |

**Sample Response** 
```json=
{
  "data": {
    "balance": 600270,
    "currency": "Coin"
  },
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2017-01-20T01:14:48.9266001-05:00"
  }
}

```


<a href="##Record" title="Transfer-Wallet-GameBoy"> ***Sample Record - Cancel***</a>

```json=

{
  "data": {
    "_id": "59672a547aa48000019260cf",
    "action": "bets",
    "target": {
      "account": "{account1}"
    },
    "status": {
      "createtime": "2017-07-13T04:07:48.644-04:00",
      "endtime": "2017-07-13T04:07:48.673-04:00",
      "status": "success",
      "message": "success"
    },
    "before": 600280,
    "balance": 600270,
    "currency": "Coin",
    "event": [
      {
        "mtcode": "{mtcode1}",
        "amount": 10,
        "eventtime": "2020-05-20T05:20:00-04:00",
        "status": "cancel"
      },
      {
        "mtcode": "{mtcode2}",
        "amount": 20,
        "eventtime": "2020-05-20T05:20:00-04:00"
        "status": "refund",
      },
      {
        "mtcode": "{mtcode3}",
        "amount": 30,
        "eventtime": "2020-05-20T05:20:00-04:00"
      }
    ]
  },
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2023-01-13T04:08:02-04:00"
  }
}
```
## Amends
### 批次多人修改派果
:::warning
※多人多批次，每一批次(ucode)最多200筆 </br>
※有『部份』成功或失敗的情境。請參考reponse範例
:::

**Request**

URL: `/transaction/game/amends `

Method: `POST`

Headers :<br/>

`Content-Type:application/json`<br/>
`wtoken:config裡設定的token`<br/>

**Request Parameter**
| 參數名稱 | 型別 | 必填 | 描述|
|-| :-:|:-:|:-|
| account     	  | string   | 必填 | 使用者帳號|
| event       	  | text array | 必填 | 事件資料列表用JSON包起來<br/>長度不限，但每筆win最大長度為158字元<br/>mtcode:交易代碼<br/>amount:金額<br/>action:行為 分別為 credit 補款 ,debit 扣款<br/>roundid:注單號<br/>gamecode:遊戲編碼 <br/> validbet:有效投注 |
| eventtime     	  | string   | 必填 |交易當下時間<br/> 此時間可與注單的createtime對應 <br/>※最大長度為35字元|
| amount     	  | number    | 必填 | 該批交易金額(實際從玩家錢包"扣除"或"增加"的金額)<br/> |
| action     	  | string   | 必填 | 批次執行的錢包行為<br/>**credit 補款 /debit 扣款**    |
| ucode     	  | string   | 必填 | 交易批碼，回傳派彩結果時需附帶此資訊<br/>  |
    
**Sample Request :**
```bash=
curl --location --request POST '{wurl}/transaction/game/amends' \
--header 'Content-Type: application/json' \
--header 'wtoken: {wtoken}' \
--data-raw '{
  "list": [
    {
      "account": "{account1}",
      "eventtime": "2019-12-13T08:39:08.906Z",
      "amount": 20,
      "action": "credit",
      "ucode": "{ucode1}",
      "event": [
        {
          "mtcode": "{mtcode1}",
          "amount": 100,
          "validbet": 55,
          "action": "debit",
          "roundid": "{round1}",
          "eventtime": "2019-12-13T08:39:08.906Z",
          "gamecode": "{gamecode}"
        },
        {
          "mtcode": "{mtcode2}",
          "amount": 120,
          "validbet": 20,
          "action": "credit",
          "roundid": "{round2}",
          "eventtime": "2019-12-13T08:39:08.906Z",
          "gamecode": "{gamecode}"
        }
      ]
    },
    {
      "account": "{account2}",
      "eventtime": "2019-12-13T08:39:08.906Z",
      "amount": 20,
      "action": "debit",
      "ucode": "{ucode2}",
      "event": [
        {
          "mtcode": "{mtcode3}",
          "amount": 100,
          "validbet": 55,
          "action": "debit",
          "roundid": "{round3}",
          "eventtime": "2019-12-13T08:39:08.906Z",
          "gamecode": "{gamecode}"
        },
        {
          "mtcode": "{mtcode4}",
          "amount": 80,
          "validbet": 20,
          "action": "credit",
          "roundid": "{round4}",
          "eventtime": "2019-12-13T08:39:08.906Z",
          "gamecode": "{gamecode}"
        }
      ]
    }
  ]
}'
```

**Response Format** 

| 參數名稱 | 型別 | 描述 |
| :- |:-:| :- |
| failed.code      |string | 錯誤代碼 |
| failed.message      |string | 錯誤訊息 |
| success.account  <br/>(failed.ucode )     |string | 帳號 |
|  success.ucode  <br/>(failed.ucode )    |string | 交易批碼，回傳派彩結果時需附帶此資訊  |
|success.before|number|執行動作前的餘額<br/>|
|  success.balance      |number | 執行動作後的餘額
| status.code   |string | 狀態編碼 |
| status.message   |string | 狀態訊息  |
| status.datetime   |string | 回傳時間  |

**Sample Response :**
```json=
{
    "data": {
        "failed": [
            {
                "account": "{account1}",
                "code": "1006:wallet",
                "message": "Player not found.",
                "ucode": "{ucode1}"
            }
        ],
        "success": [
            {
                "account": "{account2}",
                "currency": "Coin",
                "before": 8970720,
                "balance": 8970700,
                "ucode": "{ucode2}"
            }
        ]
    },
    "status": {
        "code": "0",
        "message": "Success",
        "datetime": "2019-08-23T05:13:22-04:00"
    }
}

```

<a href="##Record" title="Transfer-Wallet-GameBoy"> **Sample Record - Amends**</a>

```json=
{
  "data": {
    "_id": "59672a547aa48000019260cf",
    "action": "amends", 
    "target": {
      "account": "{account2}"
    },
    "status": {
      "createtime": "2017-07-13T04:07:48.644-04:00",
      "endtime": "2017-07-13T04:07:48.673-04:00",
      "status": "success",
      "message": "success"
    },
    "before": 8970720,
    "balance": 8970700,
    "currency": "Coin",
    "event": [
      {
        "mtcode": "{mtcode3}}",
        "amount": 100,
        "action": "debit", 
        "eventtime": "2016-01-02T15:04:05-04:00"
      },
      {
        "mtcode": "{mtcode4}",
        "amount": 80,
        "action": "credit",
        "eventtime": "2016-01-02T15:04:05-04:00"
      }
    ]
  },
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2017-07-13T04:08:02-04:00"
  }
}
```

## Amend
### <font color="red">**單人**</font>修改派果
:::warning
※單一玩家，批次裡可以不同round和mtcode。每一批次最多200筆 </br>
※只有全部成功或全部失敗，沒有『部份』成功或失敗的情境 
:::

**Request**
    
URL: `/transaction/game/amend`

Method: `POST`

Headers :<br/>

`Content-Type:application/json`<br/>
`wtoken:config裡設定的token`<br/>

**Request Parameter** 
| 參數名稱   | 型別     | 必填  | 描述|
|-| :-:|:-:|:-|
| account     	  | string   | 必填 | 使用者帳號|
| gamehall     	  | string   | 必填 | 遊戲廠商代號|
| gamecode     	  | string   | 必填 | 遊戲代號
| action     	  | string   | 必填 | 批次修改後所執行的錢包行為。(確定此批交易是補款還是扣款) 分別為 credit 補款 ,debit 扣款  |
| amount     	  | number    | 必填 | 該批交易 金額(實際從玩家錢包"扣除"或"增加"的金額)|
| createTime     	  | string   | 必填 |成單時間<br/> 此時間可與注單的createtime對應 |
| data       	  | text array    | 必填 | 事件資料列表用JSON包起來<br/>長度不限，但每筆win最大長度為158字元<br/>mtcode:交易代碼<br/>amount:差額部份<br/>action:行為 分別為 credit 補款 ,debit 扣款<br/>roundid:注單號<br/>eventtime:交易當下時間<br/> validbet:有效投注|

**Sample Request :**
```bash=
curl --location --request POST '{wurl}/transaction/game/amends' \
--header 'Content-Type: application/json' \
--header 'wtoken: {wtoken}' \
--data-raw '{
    "account": "{account1}",
    "gamehall": "xxx",
    "gamecode": "gamecode",
    "action": "credit",
    "amount": 10,
    "createTime": "2016-01-02T15:04:05-04:00",
    "data": [
        {
            "mtcode":"{mtcode1}",
            "amount":20,
            "validbet":5,
            "roundid":"{round1}",
            "eventtime":"2016-01-02T15:04:0504:00",
            "action":"credit"
        },
        {
            "mtcode":"{mtcode2}",
            "amount":10,
            "validbet":5,
            "roundid":"{round2}",
            "eventtime":"2016-01-02T15:04:0504:00",
            "action":"debit"
        }
    ]
}'
```

**Response Format**

| 參數名稱 | 型別 | 描述 |
| --- |:-:| --- |
| balance      |number | 執行動作後的餘額 |
| currency     |string | 幣別</font> |
| status.code   |string | 狀態編碼|
| status.message   |string | 狀態訊息  |
| status.datetime   |string | 回傳時間  |

**Sample Response :**
```json=
{
  "data": {
    "balance": 600240,
    "currency": "Coin"
  },
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2023-01-19T22:15:21.6076001-05:00"
  }
}

```

<a href="##Record" title="Transfer-Wallet-GameBoy"> **Sample Record - Amend**</a>

```json=
{
  "data": {
    "_id": "59672a547aa48000019260cf",
    "action": "amend", 
    "target": {
      "account": "{account1}"
    },
    "status": {
      "createtime": "2017-07-13T04:07:48.644-04:00",
      "endtime": "2017-07-13T04:07:48.673-04:00",
      "status": "success",
      "message": "success"
    },
    "before": 600230,
    "balance": 600240,
    "currency": "Coin",
    "event": [
      {
        "mtcode": "{mtcode1}",
        "amount": 20,
        "action": "credit", 
        "eventtime": "2016-01-02T15:04:05-04:00"
      },
      {
        "mtcode": "{mtcode2}",
        "amount": 10,
        "action": "debit",
        "eventtime": "2016-01-02T15:04:05-04:00"
      }
    ]
  },
  "status": {
    "code": "0",
    "message": "Success",
    "datetime": "2017-07-13T04:08:02-04:00"
  }
}
```
