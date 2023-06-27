---
title: สรุป นิยามของ Software Engineering จาก Titus Winters และ Russ Cox
date: 2023-06-27T07:13:00+07:00
draft: false
---

นิยามของ Software Engineering โดย Titus Winters บอกไว้ว่า

"Software Engineering is Programming integrated over time"

ส่วนบทความของ Russ cox "What is Software Engineering?" ซึ่งก็ยกนิยามของ Software Engineering ของ Titus Winters มาขยายความต่อไปอีกว่า

"Software engineering is what happens to programming when you add time and other programmers"

<!--more-->

Russ บอกว่า Programming ก็คือการเขียนโปรแกรมให้มันทำงานได้ เรามีปัญหาที่ต้องแก้ เราเขียนโค้ด แล้วก็รันมัน แล้วก็ได้คำตอบของปัญหา นั่นคือ Programming

Programming โดยตัวมันเองก็ยากอยู่แล้ว ทีนี้ถ้าเราต้องทำให้ Code ที่เราเขียน ยังคงทำงานต่อไปได้เรื่อยๆ มีการเปลี่ยนแปลงเรื่อยๆเมื่อ เวลา ผ่านไป ต้องทำงานร่วมกันกับโปรแกรมเมอร์หลายๆคนล่ะ

เราอาจจะเริ่มนึกถึง version control system เช่น Git เพื่อเทรคความเปลี่ยนแปลงที่เกิดขึ้น, นึกถึง unit test เพื่อช่วยเช็คว่าโปรแกรมไม่มีบัคเกิดขึ้นเมื่อเวลาผ่านไปหรือมีการแก้ไข, เราต้องนึกถึงเรื่องการจัดการ module การใช้งาน design patterns เพื่อแบ่งโค้ดของเป็นส่วนๆแล้วทำให้ทีมทำงานแยกจากกันได้ง่ายขึ้น, เราจะต้องหาเครื่องมือมาช่วยป้องกัน bug เพื่อหา bug ได้แต่เนิ่น

จะเห็นว่านี่ละคืองานของ Software Engineering ที่มีอะไรมากกว่าเขียน Program ให้ทำงานได้ อีกมากมาย มีอะไรให้ศึกษาต่ออีกเยอะเพื่อทำให้โปรแกรมเรายังคงทำงานได้ เปลี่ยนแปลงได้ เมื่อเวลาผ่านไป

ลองอ่านของ Russ cox เต็มๆได้ต่อที่นี่ https://research.swtch.com/vgo-eng
และฟังคลิปของ Titus Winters ได้ที่นี่ https://www.youtube.com/watch?v=tISy7EJQPzI&t=8m17s
