---
title: "[JavaScript] parse date time แบบกำหนด timezone ด้วย dayjs library"
date: 2023-02-28T17:14:51+07:00
draft: false
---

[Dayjs](https://day.js.org) เป็น JavaScript library เพื่อจัดการกับ Date Time สามารถ parse, validates, manipulates, and displays date time ให้เราได้ ตอนนี้ก็ใช้ library นี้อยู่ แต่ว่าพอจะ parse date time ที่ต้องการระบุ timezone เช่น Asia/Bangkok สำหรับภาษาไทย จะมีท่าที่ต้องทำเพิ่มนิดหน่อย

<!--more-->

หลังจากติดตั้ง dayjs แล้วไม่ว่าจะผ่าน npm หรือ yarn แล้วเวลาจะใช้งานก็ต้อง import มาก่อน แบบนี้ (ผมใช้ import syntax นะ ใครจะใช้ require ก็ได้ไม่มีปัญหา)

```ts
import dayjs from "dayjs";
```

แต่ว่าถ้าเราจะใช้ฟังก์ชันที่สามารถระบุ timezone ได้เราต้อง import module เพิ่มเติมแบบนี้

```ts
import timezone from "dayjs/plugin/timezone";
```

นอกจากนั้นถ้าเราต้องการ parse date ด้วย format ที่เราต้องการด้วย เราต้อง import เพิ่มอีก module แบบนี้

```ts
import customParseFormat from "dayjs/plugin/customParseFormat";
```

จะเห็นว่า dayjs ใช่ plugin pattern ในการเพิ่มความสามารถให้กับ library แต่ว่าแค่ import ยังไม่พอ เราต้องเรียก function extend เพื่อเพิ่มความสามารถด้วยแบบนี้

```ts
dayjs.extend(timezone);
dayjs.extend(customParseFormat);
```

หลังจากนั้นเราถึงจะเรียกฟังก์ชันในการ parse date with custom format ด้วย timezone ที่ต้องการได้แบบนี้

```ts
const dateString = "28/02/2023, 17:23";
const dateTime = dayjs.tz(dateString, "DD/MM/YYYY, HH:mm", "Asia/Bangkok");
```

เท่านี้ก็ใช้งานได้แล้ว

ref: https://day.js.org/docs/en/timezone/parsing-in-zone
