---
title: "[Go] ลบบาง query param ออกจาก URL ด้วย package net/url"
date: 2023-03-13T17:54:36+07:00
draft: false
---

วันนี้มีโจทย์ให้ต้องลบบาง query param ออกจาก URL ที่ได้รับมา ซึ่งจะใช้ package `net/url` ช่วยในการ parse แล้วลบบางค่าออกได้

<!--more-->

เช่นเรามี URL แบบนี้

```txt
https://www.example.com?foo=hello&bar=world
```

แล้วเราอยากจะเปลี่ยนให้เหลือแค่

```txt
https://www.example.com?bar=world
```

เราจะใช้ package `net/url` เขามาช่วย เริ่มแรกให้เรา parse URL string ให้เป็น value ของ URL ก่อนแบบนี้

```go
	urlStr := "https://www.example.com?foo=hello&bar=world"
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}
```

จากจากนั้นก็เรียก method `Query()` เพื่อให้ได้ query params ในรูปแบบของ `url.Values` type จากนั้นก็เรียก method `Del` ลบ query param ตามชื่อที่ต้องการลบออก

```go
	qParams := u.Query()
	qParams.Del("foo")
```

สุดท้ายก็เรียก method `Encode()` เพื่อเซตค่ากลับไปที่ field `RawQuery` ของ type URL

```go
	u.RawQuery = qParams.Encode()
```

ตัวอย่างโค้ดสุดท้ายหน้าตาแบบนี้

```go
package main

import (
	"fmt"
	"log"
	"net/url"
)

func main() {
	urlStr := "https://www.example.com?foo=hello&bar=world"
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}

	qParams := u.Query()
	qParams.Del("foo")

	u.RawQuery = qParams.Encode()

	fmt.Println(u)

        // output:
        // https://www.example.com?bar=world
}
```
