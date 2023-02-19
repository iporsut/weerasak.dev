---
title: "กระจายการทำงานหลายๆ งาน ผ่านหลายๆ goroutine ง่ายๆ ด้วย package errgroup"
date: 2023-02-19T08:01:48+07:00
draft: false
---

บางการทำงานเช่นโหลดข้อมูลจากไฟล์หลายๆ record แล้วเขียนลง database เราไม่จำเป็นต้องให้มันทำงานทีละ 1 record เสร็จก่อนแล้วค่อยไปทำเขียน record ถัดไปก็ได้ เราสามารถแยกการทำงานให้ทำในหลายๆ goroutine แล้วใช้ความสามารถของเครื่อง (CPU) ให้เต็มที่ถ้ามีหลายๆ CPU หลายๆ core ซึ่ง [package errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) ช่วยเราทำแบบนี้ได้พร้อมจัดการ error และ cancel context ได้อีกด้วย

<!--more-->

ตัวอย่างเช่นเรามีโค้ดที่เขียนข้อมูล users ลง database แบบนี้

```go
for _, user := range users {
        err := db.Insert(ctx, uesr)
        if err != nil {
                return err
        }
}
```

การ Insert ก็จะทำทีละ user จบก่อนค่อยไป Insert user ถัดไป

เราสามารถเขา [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) มาช่วยได้ โค้ดจะเป็นแบบนี้

```go
grp, grpCtx := errgroup.WithContext(ctx) // (1)

for _, user := range users {
        user := user
        grp.Go(func() err { // (2)
                return db.Insert(grpCtx, user)
        })
}

if err := grp.Wait(); err != nil { // (3)
        return err
}
```

โดยมีวิธีใช้งานแบบนี้

1. เราสร้าง errgroup value ขึ้นมาก่อนโดยใช้ func `WithContext` ซึ่งจะรับ context เข้าไปด้วยแล้วสร้าง Group value กับ Group context กลับมาให้เรา
2. เราเอา grp ไปเรียกใช้ method `Go` ซึ่งจะรับค่าเป็น anonymous function ที่ต้อง return error กลับออกมา ซึ่งฟังก์ชันที่เราส่งไปให้ `Go` เนี่ยแหละคือสิ่งที่เราอยากให้มันทำงานแบบ concurrent เช่นตามตัวอย่างคือ inesrt user เข้าไป, ในกรณีนี้เราต้องเปลี่ยน context ที่เคยส่งในตอนเรียก Insert เป็น group context (grpCtx) แทนด้วย เวลาเกิดมีข้อผิดพลาดในการทำงานสัก 1 ครั้งจะได้ยกเลิกการทำงานของการ insert ที่เหลือไปด้วยเลย
3. เรียก method `Wait()` เพื่อรอให้การทำงานของทุกๆ goroutine ที่ method `Go` สร้างไว้ทำงานเสร็จ แล้วก็จะส่งค่า error กลับมา ถ้าไม่มี error ก็ได้ค่า nil ถ้ามีก็ได้ค่า error ค่าแรกที่เกิดขึ้นนั่นเอง

เท่านี้เราก็ได้การทำงานแบบ concurrent ง่ายๆ ที่สามารถจัดการ context และ error ให้เราได้ด้วย ได้แล้ว
