## 1: Start with Metrics's name, e.g.:
```
axios_success_rate
```

<br/><hr/><br/>

## 2: Support 4 operators, e.g.:
- = :equal to
- != :not equal to
- =~ :match Regex
- !~ :not match Regex

<br/><hr/><br/>

## 2.1: Query with Label
```
metrics_name { label_name = lable_value }
```

<br/><hr/><br/>

## 2.2: Query with AND & Regex
```
metrics_name { label1 = lable_value1, lable2 =~ "*.domain.com" }
```

<br/><hr/><br/>

## 2.3: Query with metrics'value Comparation
```
metrics_name > 100
```

<br/><hr/><br/>

## 2.4: Query with Range (time & number)
```
metrics_name[5m]
```

<br/><hr/><br/>

## 3: Common Present functions, e.g.:
- rate() :: 取平均值 :: rate( metrics_name[5m] ) :: 計算 5 分鐘內的平均值
- sum() :: 取總值 :: sum( metrics_name ) by ?
- topk() :: 只顯示頭幾個