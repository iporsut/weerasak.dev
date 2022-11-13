---
title: "Go Frameworks"
date: 2022-11-12T19:10:51+07:00
draft: false
---

Framework ของเราไม่เท่ากัน เวลาคนที่คุ้นเคยกับ web framework ของ stack ที่ตัวเองใช้อยู่ แล้วมาหัดใช้ [Go](https://go.dev) ก็มักจะถามหา Framework ว่า Go นิยมใช้ตัวไหนกันบ้าง แต่สุดท้ายก็มักจะไม่เจอ framework แบบที่เคยใช้ หรือแบบที่คิดไว้ ก็ผิดหวัง หรือรู้สึกใช้งานยากกันไป

<!--more-->

ขอยกตัวอย่าง web frameworks ดังๆของ 2 ภาษาแล้วกัน คือ [Ruby on Rails](https://rubyonrails.org) ของภาษา Ruby และ [Laravel](https://laravel.com) ของ PHP

ซึ่ง Laravel ก็ได้รับ inspiration มาจาก Rails นั่นละ คือเป็น MVC (Model/View/Controller) web framework ซึ่งตัว famework เองนั้น pack library จำเป็นในการจัดการ model แล้วจัดเก็บ model ลง Database ส่วนใหญ่คือเตรียม ORM library มาให้เองแล้ว
ส่วน view ก็มี template engine มาให้เอง และ controler ก็เตรียม class มาให้เป็น template เพื่อ implements request handler พร้อมวิธี config routing ให้เรียบร้อยแล้ว

ทีนี้ Go frameworks ที่เป็นลักษณะแบบ Rails หรือ Laravel มีหรือเปล่า? จริงๆคือมี เช่น [Beego](https://github.com/beego/beego) และ [Go Buffalo](https://gobuffalo.io) แต่ framework ลักษณะนี้ ไม่เป็นที่นิยมมากนักสำหรับคนใช้งาน Go


แล้ว Framework แบบไหนล่ะที่เป็นที่นิยมสำหรับคนใช้งาน Go?

รูปแบบ Framework ที่นิยมในการใช้งาน Go หรือจริงๆจะเรียกว่าเป็น Library ก็ได้นั้น คือการไม่ล็อคสเปค หรือบังคับว่าส่วนของ Model, View หรือ Controller จะต้องใช้ library ไหน คนใช้งานสามารถเลือก libray เองได้แล้วเอามาประกอบกันให้เว็บเราทำงานได้เหมือนกัน จริงๆคือไม่ได้ยึดว่าต้องเป็น MVC pattern ด้วยซ้ำ

ตัวอย่างเช่น ถ้าเราทำ API endpoints ง่ายๆ ไม่ซับซ้อน เผลอๆมี 1 endpoints เท่านั้น เราไม่จำเป็นต้องใช้ 3rd party library ด้วยซ้ำ เราสามารถใช้ standard package `net/http` ได้เลย ส่วนถ้าต้องมีการจัดเก็บข้อมูล ต้องต่อฐานข้อมูล ก็มี package `database/sql` ให้ใช้อยู่แล้ว แค่ต้องหา package driver สำหรับ DBMS ที่เราเลือกใช้ได้เอง

ถ้าเราอยากได้ frameworks ที่จัดการ HTTP ที่ซับซ้อนยิ่งขึ้น เช่นต้องการ map path params เราก็มีทางเลือกทั้ง gorilla mux library หรือจะใช้ [Gin](https://github.com/gin-gonic/gin), [Echo](https://echo.labstack.com) ก็ได้ที่นอกจากจัดการ HTTP แล้วยังมี library ช่วยจัดการ request validation การจัดการ errors มีระบบ middleware ให้เราแยก concern ต่างๆที่ใช้ในหลายๆ endpoint ออกมาเมื่อ reuse ได้

ส่วนการจัดการ database เราก็มีทางเลือกนอกจาก `database/sql` คือ library ในแนวที่เป็น ORM อย่าง [GORM](https://gorm.io) หรือ [ent](https://entgo.io) framework ช่วยได้เช่นกัน

หรือถ้าฐานข้อมูลเราไม่ได้เป็นแบบ SQL เราก็แค่ใช้ sdk library ของ ฐานข้อมูลนั้นๆร่วมกันกับ framework อย่าง Gin, Echo ได้เลยเช่นเลือกใช้ Firestore หรือ MongoDB เป็นต้น

สรุปก็คือ สำหรับ Go frameworks แบบ Rails, Laravel นั้นมี แค่ไม่เป็นที่นิยมเท่า frameworks แต่แยกการทำงานเฉพาะทางเช่น จัดการ HTTP handlers เท่านั้น แล้วให้ทางเลือกคนใช้สามารถประกอบของอื่นๆให้มาทำงานร่วมกันได้เช่นเลือก database library มาใช้เอง

ส่วนตัวข้อดีของการใช้งาน Frameworks อย่าง Rails, Laravel ก็คือเลือกมาให้แล้ว ทำให้เริ่มได้ง่ายและเร็ว แต่ข้อเสียก็คือถ้าความต้องการของเรา ไม่ตรงไปตามสิ่งที่เตรียมมาให้ จะเริ่มทำให้เราต้องใช้ท่ายาก ทำให้โค้ดเริ่มจะผิดเพี้ยนไปจาก pattern ที่ framework เตรียมมาให้ สุดท้ายกลับดูแลยากกว่าเดิม

ส่วนแบบที่คนใช้ Go นิยมใช้ แน่นอนว่าเริ่มยากกว่าถ้าเพิ่งใช้มาใช้งาน เพราะก็ไม่รู้จะเลือกอะไรมาประกอบกัน อันนี้ก็ต้องลองค้นดูหรือถามดูว่าคนส่วนใหญ่ใช้อะไร สุดท้ายก็ต้องลองทำจริงๆแล้วดูแล้วว่าเราถนัดแบบไหน และแบบไหนที่ทำให้เราดูแลระบบเราได้ง่ายสุด
