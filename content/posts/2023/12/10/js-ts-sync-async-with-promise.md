---
title: "JavaScript/TypeScript Sync an Async with Promise"
date: 2023-12-10T06:38:00+07:00
draft: false
---

เราใช้ Promise และ async/await syntax ช่วยจัดการร้อยเรียงโค้ดที่จะทำงานหลังจากมีการเรียกใช้ asynchronous function ต่างๆ
แต่ส่วนใหญ่เรามักจะ await หรือ then กันทันทีตรงจุดที่เราเรียกใช้งาน async function แต่จริงๆแล้วเราสามารถเรียกแบบ ไม่มี await แล้วเอา Promise object ไป await หรือ then ทีหลังได้

<!--more-->

ตัวอย่างโค้ดที่เราเขียนกันบ่อยๆ คือ

```ts
async function asyncFunction() {
  // do multiple asynchronous tasks
  return "result";
}

async function main() {
  const result = await asyncFunction();
  console.log(result);
}
```

แต่บางครั้งตอนเราเรียก async ที่จุดนึง แต่ไม่ได้ต้องการ await response ตรงนั้นทันที ต้องการไป await ที่อื่น เราก็เรียกฟังก์ชันนั้นแบบปกตินี่แหละ ให้มันตอบกลับ Promise Object กลับมาให้เรา แล้วค่อยเอา Promise Object ไป await หรือจะ then ก็ได้

สถานการณ์ที่ผมเจอคือมี external library ให้เตรียม callback function ให้มัน โดยเขียนไว้ใน useEffect ของ React แล้วมันจะตอบกลับ unsubscribe function กลับมาให้ ซึ่งเราก็ต้องส่ง unsubscribe ให้กับ useEffect เพื่อเรียกตอนมันจะ re-render ครั้งถัดไป ทีนี้ก่อนจะเรียก subscribe function โดยใส่ callback เรามี asynchronous process นิดนึงก่อนถึงจะเตรียม callback function ได้

สุดท้ายเลยได้ใช้ Promise object ช่วย sync ให้กับ useEffect รอการทำงาน async ที่เตรียม callback ทำงานจบก่อนแล้วจะได้ unsubscribe function กลับมาให้เรียกในตอนท้ายหลังจากนั้น

```ts
useEffect(() => {
  const promise = asyncFunction().then((result) => {
    return externalLibrary.subscribe((data) => {
      // do something with data and result
    });
  });
  return () => {
    promise.then(() => {
      unsubscribe();
    });
  };
}, []);
```

เท่านี้เราก็ทำให้ useEffect ทำงานได้ถูกต้องตามที่ต้องการแล้ว คือเรียก unsubscribe() ตอน re-render ครั้งถัดไป
