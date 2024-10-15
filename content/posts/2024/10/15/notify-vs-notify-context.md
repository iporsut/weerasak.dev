---
title: "signal.Notify vs singal.NotifyContext"
date: 2024-10-15T13:48:20+07:00
draft: false
---

ใน package signal จะมี 2 ฟังก์ชันที่ใช้ในการรับค่า interrupt จาก OS คือ Notify และ NotifyContext สองฟังก์ชันนี้ต่างกันยังไงบ้าง สรุปได้ดังนี้

<!--more-->

เวลาทำ gracefully shutdown ที่ต้องการรับ Interrupt signal จาก OS ใน Go จะมีทางเลือกที่เราใช้บ่อย ๆ คือใช้ ฟังก์ชัน signal.Notify และ signal.NotifyContext จากแพคเกจ signal

จุดที่ต่างกันของสองฟังก์ชันนี้คือ Notify จะให้เราส่ง channel เข้าไปเพื่อ register ว่าเวลามี interrupt signal ที่เราสนใจเข้ามา จะให้มันส่งมาเก็บไว้ใน channel

เราสามารถดึงค่า (receiving) ออกมาจาก channel แล้วเช็คได้ว่ามันเป็น interrupt signal อะไร

https://pkg.go.dev/os/signal#Notify

ส่วน NotifyContext นั้นจะให้เราส่ง parent context และ list ของ signal ที่เราสนใจเข้าไป แล้วมันจะสร้าง context ใหม่มาให้เราที่จะถูก cancel เมื่อมี signal ตามที่เราลิสต์ไว้ถูกส่งจาก OS มาที่ process ของโปรแกรมเรา

https://pkg.go.dev/os/signal#NotifyContext

เราสามารถเช็คว่ามี signal มาแล้วหรือยังได้เองจากการเช็คว่ามีค่าใน channel ctx.Done() แล้วหรือยัง

จะเห็นว่าเราเช็คว่าเกิดการ interrupt ได้เหมือนกันทั้งคู่ จุดต่างคือ ที่ทางในการเช็ค signal อันแรก จาก channel ที่เราสร้างขึ้นเอง อีกอันจาก ctx.Done() ของ context ที่ NotifyContext สร้างมาให้

เพราะฉะนั้น ถ้าเราสนใจแค่ว่ามี interrupt เกิดขึ้นแล้ว ไม่สนใจว่าเป็น signal แบบไหน การใช้ NotifyContext ก็ดูจะเรียบง่าย และเราสามารถส่ง context ที่ได้ต่อไปยัง function อื่น ๆ ที่ต้องการ context ต่อไปได้เลยเพื่อเป็นช่องทาง cancel process การทำงานเมื่อมี interrupt เกิดขึ้น

แต่ถ้าเมื่อไหร่เราต้องการลง detail ว่าเกิด interrupt แบบไหนกันแน่เกิดขึ้น การใช้ Notify และมี channel รอรับก็จะดีกว่าเพราะเราดึงค่าออกมาดูได้