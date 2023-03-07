---
title: "[JavaScript] ความต่างระหว่างใช้ Promise.all กับ for แล้ว await ทีละอัน"
date: 2023-03-07T09:47:00+07:00
draft: false
---

ถ้าเรามีหลาย async call แล้วอยากรอให้ทำงานเสร็จหมดก่อนค่อยไปต่อ มีสองท่าที่นึกออกคือใช้ Promise.all และ for แล้ว await ทีละอัน ซึ่งมันจะมีพฤติกรรมต่างกันอยู่ ขอบันทึกเอาไว้หน่อย

<!--more-->

## ถ้าใช้ Promise.all จะเรียกแบบ async ทั้งหมดแล้วค่อยรอ

ถ้าเราใช้ await Promise.all มันจะทำให้การทำงานแบบ async เกิดขึ้นเลยแล้วก็ค่อยรอให้ทุกๆอันเสร็จ ตัวอย่างเช่น

```js
const names = ["john", "jane", "joe"];

const results = await Promise.all(names.map((name) => queryData(name)));
```

เราจะทำการ map แล้วเรียก queryData ซึ่งจะได้ลิสต์ของ Promise object ของการทำงานแบบ async ที่แต่ละ query ก็ไม่ขึ้นต่อกัน ต่างคนต่าง query ได้เลย แล้ว Promise.all ก็จับมารวมกันเพื่อให้เป็น Promise object หนึ่งอันแล้วเราก็ใช้ await เพื่อรอมัน resolve ทั้งหมดได้ array ของ result ของทุก query

## ถ้าใช้ for แล้ว await จะรอเสร็จทีละ 1 แล้วค่อยเรียกแบบ async

ทีนี้ถ้าเราไม่ได้ใช้ Promise.all แล้ววนลูปด้วย for แล้วเรียก await queryData ทีละอันแบบนี้

```js
const names = ["john", "jane", "joe"];
const results = [];
for (const name of names) {
  results.push(await queryData(name));
}
```

## ความต่าง

ถ้าดูดีๆจะเห็นว่า ถ้าเราใช้ for แล้ว await ทีละอันนั้นหมายความว่ามันจะต้องรอ await ทีละอันเสร็จก่อนแล้วไปทำอันถัดไป แต่ถ้าเราไม่จำเป็นต้องรอทีละอัน เราใช้ Promise.all ไปเลยดีกว่า

ทีนี้มีสถานการณ์ไหนต้องรอบ้าง ก็พอมีอยู่บ้างแต่คงไม่เยอะ เช่น เราต้องการเอาผลลัพธ์จาก async function ก่อนหน้าไปเป็น input ให้กับ async function ถัดไปเช่น

```js
const params = ["john", "jane", "joe"];
let result = "";

for (const p of params) {
  result = await query(result, p);
}
```

จะเห็นว่าต้องรอ result ก่อนหน้าเพื่อเป็น input ให้กับ query ถัดไป แบบนี้ก็จะใช้ Promise.all ไม่ได้แล้วต้องเลือกใช้ for แล้ว await ทีละอันแทน
