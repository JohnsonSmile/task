

## 1. 开启内网穿透服务
将tool中的sunny移到path中
```bash
sunny clientid xxxxxx
```

## 2. 开启服务
```bash
cd twitter_task
go run main.go --tag debug
```

## 3. 访问网页
内网穿透代理地址:
http://{代理地址}/index

点击登录twitter

根据提示最终auth_token和auth_secret将会被存储在数据库中。
