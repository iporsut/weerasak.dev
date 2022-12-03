---
title: "วิธีเทสโค้ดที่ใช้ Time.Now"
date: 2022-12-03T07:49:43+07:00
draft: false
---

โค้ดที่เราใช้ `time.Now()` ใน package ทำให้โค้ดเราไม่สามารถเขียนเทสแล้วเช็คได้ว่า ค่าที่ได้ตอนเทสคือะไร ทางเลือกนึงง่ายๆที่ใช้อยู่ตอนนี้เพื่อให้ทดสอบโค้ดที่ใช้ `time.Now()` ได้คือใช้ตัวแปรเก็บฟังก์ชัน time.Now เอาไว้ซะ เราจะได้เปลี่ยนได้ตอนเขียนเทส

<!--more-->

วิธีก็ง่ายๆถ้าเรามีโค้ดที่ใช้ `time.Now()` ใน package เช่น

```go
package customer

type Customer struct {
        UpdatedAt time.Time
}

func (c *Customer) Save() {
        c.UpdatedAt = time.Now()
}
```

เวลาเขียนเทสเราจะเช็คไม่ได้ว่า `c.UpdatedAt` มีค่าเป็นอะไร

```go
package customer

func TestSaveCustomer(t *testing.T) {
        var c Customer
        c.Save()

        if c.UpdatedAt.Equal(__??__) {
                t.Error("UpdatedAt is not correct")
        }
}
```

ทางแก้ก็อย่างที่บอก ใช้ตัวแปรใน scope ของ pacakge เก็บฟังก์ชันเอาไว้

```go
package customer

var timeNow = time.Now // ไม่ใส่ () นะเพราะเราจะเก็บตัวฟังก์ชัน ไม่ได้เรียกฟังก์ชันตรงนี้

type Customer struct {
        UpdatedAt time.Time
}

func (c *Customer) Save() {
        c.UpdatedAt = timeNow() // เรียกฟังก์ชันผ่านตัวแปรที่เราเก็บฟังก์ชันเอาไว้แทน
}
```

ทีนี้ในเทสเราก็เปลี่ยน `timeNow` เป็นฟังก์ชันที่ตอบกลับ time.Time ที่เราต้องการแทนเช่น

```go
package customer

func TestSaveCustomer(t *testing.T) {
        now := time.Now() // เก็บค่า time.Now() ตอนเทสเอาไว้ไปเช็คได้อีกที
        timeNow = func() time.Time { // เปลี่ยน timeNow เป็นฟังก์ชันใหม่ที่ตอบกลับเป็นค่า now ที่เราเก็บเอาไว้ก่อน
                return now
        }
        defer func() { timeNow = time.Now } // เปลี่ยนกลับเมื่อเทสฟังก์ชันทำงานจบ

        var c Customer
        c.Save()

        if !c.UpdatedAt.Equal(now) { // เอามาไว้เช็คตรงนี้ได้แล้ว
                t.Error("UpdatedAt is not correct")
        }
}
```

เท่านี้เราก็กำหนด time.Now() ตอนเขียนเทสได้แล้ว
