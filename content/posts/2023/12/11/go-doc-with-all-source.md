---
title: "go doc เพื่อแสดง source ทั้งหมดของ package"
date: 2023-12-11T16:05:00+07:00
draft: false
---

```
go doc --all --src -u
```

ใช้คำสั่งนี้เพื่อ generate go doc และทำให้มันรวมโค้ดของทั้ง package ทุกไฟล์  ทั้ง exported / unexported ออกไปสู่ standard out  ให้เราเห็นทีเดียวได้เลย
 
เช่นใน package เรามีไฟล์ a.go, b.go, c.go
 
สั่ง `go doc --all --src -u` ใน folder ของ package มันจะรวมโค้ดของ a.go, b.go, c.go ให้เราเห็นทีเดียวได้เลย
