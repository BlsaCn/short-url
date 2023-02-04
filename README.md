<h1>BlsaCn/short-url</h1>
<p>

</p>

## 短地址转换(个人学习使用)

## 安装(依赖redis服务)

```
go get github.com/BlsaCn/short-url
```

使用：
==

    1、创建短链接：http://127.0.0.1:8000/api/shorten
    Method：POST
    Params：{"url":"https://www.baidu.com", "expiration_min":10}

    2、短链接信息：http://127.0.0.1:8000/api/info?shortLink=3
    Method：GET

    3、重定向到长链接：http://127.0.0.1:8000/api/3
    Method：GET
