---
title: "[React] จะ dispatch event ให้ DOM ที่สร้างจาก React ต้องกำหนด bubbles เป็น true เสมอ"
date: 2023-03-09T19:55:21+07:00
draft: false
---

เพิ่งเจอว่า DOM ที่ได้จาก React แม้ว่าเราจะ binding Event ผ่าน JSX มาแล้วแต่เราจะ query DOM มาแล้วสั่ง dispatchEvent เลยไม่ได้ ต้องกำหนด option bubbles ของ Event object ให้เป็น true เสมอ

<!--more-->

ตัวอย่างเช่นเรามี React Component ง่ายๆแบบนี้

```jsx
export function Hello() {
  const onClick = () => {
    alert("Hello");
  };

  return (
    <div>
      <div data-testid="hello" onClick={onClick}>
        Click
      </div>
    </div>
  );
}
```

ถ้าเราเปิด browser console ของ inspector ขึ้นมา แล้วใช้ document.querySelector ช่วย select element เรียก method click แบบนี้ onClick callback ก็จะถูกเรียกใช้งานตามที่คิดไว้

```js
const dom = document.querySelector("[data-testid=hello]");
dom.click();
```

ที่เป็นแบบนี้เพราะพฤติกรรมของ method `click()` มันจะสร้าง MouseEvent object ที่กำหนด bubbles true ให้แล้ว

แต่ถ้าเราใช้ method `dispatchEvent` สำหรับ trigger event อื่นๆเช่น mousedown หรือกับ click เองก็ได้ แบบนี้ โดย default bubbles จะเป็น false ทำให้ callback เราไม่ถูกเรียก

```js
const dom = document.querySelector("[data-testid=hello]");
dom.dispatchEvent(new MouseEvent("click"));
```

ซึ่งแตกต่างกับการ binding event โดยใช้ native DOM โดยตรงแบบนี้

```html
<!DOCTYPE html>
<body>
  <button id="elem" onclick="alert('click!');">Autoclick</button>

  <script>
    let event = new MouseEvent("click");
    elem.dispatchEvent(event);
  </script>
</body>
```

ที่ต่อให้ไม่กำหนด bubbles เท่ากับ true ก็จะคง trigger onclick callback อยู่ดี

ที่มาของโพสต์นี้เพราะว่าพยายามเขียนเทส App ที่สร้างด้วย React แล้วเจอว่าไม่สามารถ trigger event click ได้เพราะ Component binding onMouseDown ไม่ใช่ onClick แต่พอไปลอง dispatchEvent ผ่าน inspector ก็ไม่ขึ้น สุดท้ายเพิ่งรู้ว่ามี option bubbles ลองเซตเป็น true เท่านั้นละได้เลย

Bonus: ถ้าใช้ Playwright อยู่แล้ว Locator object มี method dispatchEvent เช่นกันแต่เราใส่ชื่อ event ได้เลยแล้วมันจะสร้าง event object ให้เองพร้อมกำหนด bubbles กับ cancelable เป็น true ให้เอง
