---
title: "Pointer of any type T"
date: 2022-11-26T10:33:38+07:00
draft: false
---

ไม่ใช่ทุก type จะมี literal value ที่ใช้ operator `&` เพื่อเอา pointer ออกมาได้

มีแค่ struct, map, array, slice ที่มี literal value แล้วใช้ `&` เพื่อเอา pointer ออกมาได้ทันที

ทีนี้ถ้าเป็นค่าอื่นๆเช่น string, int, bool หรือค่าที่เรียกจากฟังก์ชันเช่น time.Date, time.Now เราจะใช้ `&` ไม่ได้เช่น

<!--more-->

```go
// แบบนี้ compile ไม่ผ่านนะ
fmt.Println(&1)
fmt.Println(&"abc")
fmt.Println(&true)
fmt.Println(&time.Now())
```

ทางแก้ก่อนหน้านี้คือเราต้องสร้างตัวแปรมาเก็บก่อนเช่น

```go
now := time.Now()
fmt.Println(&now)

num := 10
fmt.Println(&num)
```

นอกจากนั้นยังมีอีกท่าคือสร้างฟังก์ชัน ที่รับ value เข้าไปแล้วส่ง pointer ของ value นั้นออกมาเช่น

```go
func pointerTime(t time.Time) *time.Time { return &t }
func pointerString(s string) *string { return &s }
func pointerInt(n int) *int { return &n }

func main() {
        fmt.Println(pointerTime(time.Now()))
        fmt.Println(pointerString("abc"))
        fmt.Println(pointerInt(10))
}
```

จะเห็นว่า function เหมือนกันหมดเลย ต่างกันแค่ type ของ parameter ซึ่ง pattern แบบนี้เราใช้ generic ที่มีมาตั้งแต่ Go 1.18 ช่วยยุบเหลือฟังก์ชันเดียวได้แบบนี้

```go
package main

import (
	"fmt"
	"time"
)

func pointerOf[T any](v T) *T { return &v }

func main() {
	fmt.Println(pointerOf(time.Now()))
	fmt.Println(pointerOf("abc"))
	fmt.Println(pointerOf(10))
}
```

เท่านี้เราก็ไม่ต้องสร้างตัวแปรใหม่แล้ว
