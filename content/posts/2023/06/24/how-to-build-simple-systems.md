---
title: สรุป How to build simple systems จากการพูดของ Rich Hickey เรื่อง Simple Made Easy
date: 2023-06-24T09:42:00+07:00
draft: false
---

ได้กลับไปฟังและอ่าน Notes จากการพูดของ Rich Hickey เรื่อง [Simple Made Easy](https://www.infoq.com/presentations/Simple-Made-Easy/) อีกครั้ง
ชอบในส่วนของท้ายๆที่พูดถึงวิธีการสร้างระบบที่ simple ขอสรุปตามที่ตัวเองเข้าใจไว้สักหน่อยดังนี้

<!--more-->

## Build simple systems/component

### Abstraction for Simplicity

- อย่างแรกคือจะทำ component ที่ simple องค์ประกอบ หรือ subcomponent ที่เราเลือกมาใช้ก็ต้อง simple ก่อน
- Rich ให้นิยาม Abstract เอาไว้ว่า drawn (something) away ไม่ใช่ hiding complexity คือเป็นการจัดองค์ประกอบนะ ไม่ใช่กวาดเอาไปซ่อน
- ออกแบบโดยถามคำถามว่า What, Who, When, Where, Why and How
- นอกจากนั้น ให้ประโยคนี้ช่วย "I don't know; I don't want to know" เพื่อเช็คว่าสิ่งที่เราทำ มันจำเป็นต้องให้ component นี้รู้ในรายละเอียดตรงนี้มั้ย จะได้เป็นจุดเริ่มต้นในการ Abstract มันออกไป เพื่อให้เกิดความเรียบง่ายขึ้น

### What

- ตอบให้ได้ว่า component นี้มันจะทำอะไร ทำหน้าที่อะไร
- จัด abstraction โดยแบ่งเป็นกลุ่มของฟังก์ชันเล็กๆที่หน้าที่เกี่ยวข้องกัน กำหนด specification ของมัน จะเห็นว่ากระบวนการตรงนี้เรายังไม่สนใจ How ว่ามันทำงานยังไง แต่เราสนใจว่ามันทำอะไรได้บ้าง
- ใช้งานกลไก polymorphism คือสนใจแค่ set of functions ซึ่งแต่ละภาษาก็มีกลไกหรือชื่อเรียกต่างกัน แต่ก็แนวคิดเดียวกัน ไม่ว่าจะ interface, protocol หรือ typeclass
- พวก input, output ของแต่ละ function ที่เราออกแบบให้ใช้เป็น value ปกติ หรือเป็นพวก abstraction component อย่าง interface, protocol หรือ typeclass เพื่อที่จะได้ไม่ต้องรู้ implementation detail โดยตรง และยังเป็น polymorphism เปลี่ยน concret value ได้ขอแค่มัน implements ตาม specification ก็พอ

### Who

- Who คือการจัดการว่า สิ่งที่เราจะทำ (what) นั้นทำกับข้อมูลอะไร
- Entities ก็คือตัว component ที่ implements detail ของ abstraction ที่เรากำหนดขึ้นด้วย
- แน่นอนว่าเราแบบ component ย่อยๆ แล้วเรามาประกอบเป็น component ที่ใหญ่ขึ้น
- อย่าให้ parent จัดการ subcomponent หรือรู้รายละเอียดของ subcomponent ให้แค่ inject subcomponent เข้ามาประกอบกันพอแล้ว

### How

- คือจังหวะที่เราสนใจรายละเอียดในการ implement แล้ว
- เป็นจังหวะที่เราทำการร้อยเรียงเชื่อมต่อ component ต่างๆ แต่ก็พยายามเชื่อม abstraction ของ component ผ่าน polymorphism แทนการจัดการ component ตรงๆ

### When, Where

- ลดการที่สอง component ที่ต้องติดต่อกัน แล้วต้องรู้ว่า component อื่นๆนั้นอยู่ไหน และจะทำอะไรเมื่อไหร่ออกไปโดยใช้ queue เข้ามาช่วย พอ component เชื่อมต่อกับ queue โดยตรง ก็ไม่ต้องไปจัดการ component อื่นๆโดยตรง คุยกับ queue พอให้ queue จัดการให้ว่าจะมีค่ามาจากไหน เมื่อไหร่

### Why

- Why part คือส่วนของ Policy และ Rule ของโปรแกรม คือส่วนขอพวก logic ต่างๆของโปรแกรม ส่วนนี้จะมีความซับซ้อนกับพวก control flow ต่างๆเช่นพวก if, else, loop อะไรแบบนี้
- ส่วนนี้เป็นไปได้พยายามทำให้เป็น declarative คือแยกพวก step จริงๆออกจากส่วนที่ declare policy and rules

### Information is Simple

- Information is Simple ตามนี้เลย อยากทำให้ information มัน complex โดยการเอามันไปซ่อนใน class ซ่อนใน object ที่มี method จัดการ

### Simplicity Made Easy

สุดท้าย Rich Hickey สรุปไว้ว่า

- Choose simple constructs over complexity-generating constructs - It's the artifacts, not the authoring (ใช้ของที่ simple สร้างของที่ simple ถ้าเจออะไรที่ complex ก็เปลี่ยนซะ)
- Create abstractions with simplicity as a basis
- Simplify the problem space before you start (ลองมองปัญหาใหม่ ปรับให้มันเป็นของง่ายๆเพื่อที่จะได้แก้ง่ายๆตั้งแต่เริ่มต้น)
- Simplicity often means making more things, not fewer (ทำให้เรียบง่ายอาจจะทำให้ต้องสร้างหลายๆสิ่งก็ได้ Simplify ไม่ใช่ทำให้ทำน้อยสิ่งเสมอไป)

สุดท้ายยังไงก็อย่าลืมไปฟังต้นฉบับนะครับ อันนี้สรุปจากที่ผมฟัง ผมอาจจะฟังผิด แปลผิด หรือตีความผิดอยู่ก็เป็นได้
