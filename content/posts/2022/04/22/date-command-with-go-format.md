---
title: "สร้าง date command ที่ format ด้วย layout แบบของ Go time package"
date: 2022-04-22T22:32:10+07:00
draft: false
---

ใน Unix คำสั่ง date จะใช้ format แบบ Y, m, d อะไรแบบนี้ซึ่งเอาจริงๆพอเขียน Go บ่อยว่าภาษาอื่น ทำให้จำพวกนี้ไม่ค่อยได้ แต่จำวิธี format date ของ Go ได้มากกว่า วันนี้เลยลองเขียน command สำหรับแสดง date ขึ้นมาใหม่เองโดยใช้ format layout แบบของ Go แทน

<!--more-->

ใน Unix มีคำสั่ง date ไว้แสดงวันที่ปัจจุบัน ซึ่งเวลาสั่ง date เฉยๆจะได้แบบนี้

```
$ date
Fri Apr 22 22:34:02 +07 2022
```

เราสามารถใส่ format ได้โดยใช้ option เป็น string ที่มี + ข้างหน้า เช่น

```
% date "+%Y/%m/%d"
2022/04/22
```

จะเห็นว่าใช้ format string ตามรูปแบบของ Unix 

ที่นี้เราจะสร้างโปรแกรมขึ้นมาใหม่แทน date ให้ใช้ format layout แบบของ Go [time package](https://pkg.go.dev/time#Time.Format) แทน

โค้ดสั้นๆ ง่ายๆแค่รับ flag แล้วเอาไปใส่ method Format แบบนี้

```go
package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	format := flag.String("f", time.RFC3339, "date time format layout")
	flag.Parse()

	fmt.Println(time.Now().Format(*format))
}
```

เวลาใช้งานก็สั่ง ชื่อ command ได้เลยหรือใส่ option -f ตามด้วย format แบบ Go ได้แบบนี้

```
$ ./dt
2022-04-22T22:42:11+07:00

$ ./dt -f "2006/01/02"
2022/04/22
```