---
title: "[TypeScript] lookup ค่าใน object ด้วย key ที่มี type เป็น keyof type"
date: 2023-03-18T10:19:27+07:00
draft: false
---

ถ้าเราใช้ JavaScript โดยตรงเราสามารถ access property ใน object ได้โดยใช้ dot `.` หรือใช้ indexing operator `[property_name]` แบบนี้ แต่พอเป็น TypeScript ตอนเราใช้ indexing operator จะฟ้องว่าเราใส่ผิด type เราต้อง convert string ที่เราใส่ตรง `[]` เป็น type ของ key ของ object ก่อนด้วย `keyof`

<!---more-->

ตัวอย่างโค้ดง่ายๆเช่น

```typescript
type Point = { x: number; y: number };
type P = keyof Point;

const p: Point = { x: 10, y: 20 };
const y: string = "y";
console.log(p[y]);
```

เราจะโดน compiler แจ้ง type error ว่า

```txt
Element implicitly has an 'any' type because expression of type 'string' can't be used to index type 'Point'.
  No index signature with a parameter of type 'string' was found on type 'Point'.
```

วิธีแก้คือเราต้องกำหนด type ของ `y` ให้เป็น `keyof Point` แทนแบบนี้

```typescript
type Point = { x: number; y: number };
type P = keyof Point;

const p: Point = { x: 10, y: 20 };
const y: keyof Point = "y";
console.log(p[y]);
```

แต่ถ้าเราไม่ได้เป็นคนสร้าง `y` เอง เราก็ต้อง convert type จาก string เป็น `keyof Point` เองก่อนแบบนี้

```typescript
type Point = { x: number; y: number };
type P = keyof Point;

const p: Point = { x: 10, y: 20 };
const y: string = "y";
console.log(p[y as keyof Point]);
```

ดูเรื่อง keyof เพิ่มเต็มได้ที่นี่ https://www.typescriptlang.org/docs/handbook/2/keyof-types.html
