# Candle

Ref: Japanese Candlestick Charting Techniques https://book.douban.com/subject/2124790/

test

# Candlestick charting data

## refs
TODO:  
https://blog.csdn.net/Eumenides_max/article/details/144694349

## TODOs
https://g.co/gemini/share/4d933ed3ad57

## charting
项目地址：go-echarts GitHub：https://github.com/go-echarts/go-echarts  
文档：go-echarts Handbook：https://go-echarts.github.io/go-echarts/#/  
示例：go-echarts Examples：https://github.com/go-echarts/examples  



## Real time

东方财富网的API接口  
```shell
curl "http://push2.eastmoney.com/api/qt/stock/get?secid=1.600519&fields=f43,f57,f58,f59,f60,f61"

# 实时
curl https://tsanghi.com/api/fin/stock/XSHE/realtime?token=demo&ticker=300059
# 历史
curl https://tsanghi.com/api/fin/stock/XSHE/daily\?token\=demo\&order\=2\&ticker\=300059
```

参数说明：
```shell
secid：股票代码，1.600519 表示上海证券交易所的贵州茅台（600519）。如果是深圳证券交易所的股票，secid 前缀为 0.，例如平安银行（000001）的代码为 0.1。
fields：指定需要返回的字段，例如：
f43：最新价
f57：涨跌幅
f58：涨跌额
f59：开盘价
f60：最高价
f61：最低价
```
示例输出：
```json
{
  "data": {
    "f43": "1800.00",
    "f57": "0.50",
    "f58": "9.00",
    "f59": "1790.00",
    "f60": "1805.00",
    "f61": "1795.00"
  }
}
```
