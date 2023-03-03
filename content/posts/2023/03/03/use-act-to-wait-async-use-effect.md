---
title: "[React] ใช้ act ครอบเวลาทดสอบ render ที่มี async useEffect"
date: 2023-03-03T07:55:57+07:00
draft: false
---

ได้โอกาสลองเขียนเทสทดสอบ React component พบว่าถ้าเรามีใช้งาน async เช่น get ข้อมูลมาแสดงใน component ผ่านทาง useEffect ซึ่งจะทำงานหลังจาก render เราต้องใช้ function act ช่วยเพื่อให้มันรอ useEffect ตรงนั้นทำงานเสร็จเรียบร้อยและ rerender ให้เรียบร้อยก่อน

<!--more-->

วิธีใช้ก็ง่ายๆคือใช้ act แล้วสร้าง async callback function ครอบการเรียก render แล้วก็ใช้ await ข้างหน้า async ด้วยเพื่อให้รอทำงานเสร็จ

```ts
        const container = document.createElement("div");

        await act(async () => {
                render(
                        <ListUser />
                ),
                {
                        container: container,
                },
        });

        expect(container).toMatchSnapshot();
```

เท่านี้การ render ก็จะทำงาน async ตรง useEffect จบก่อนที่เราจะเอาไป verified การแสดงผลแล้ว
