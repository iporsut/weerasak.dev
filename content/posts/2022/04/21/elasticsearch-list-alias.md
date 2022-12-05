---
title: "List Elasticsearch Alias"
date: 2022-04-21T21:56:39+07:00
draft: false
---

Elasticsearch Alias ช่วยให้เราตั้งชื่อแทนชื่อ index จริงๆได้ ทำให้เวลาเราเขียนโค้ดหรือเรียก API ไม่จำเป็นต้องใช้ชื่อ index ตรงๆ ช่วยให้เราสามารถสร้าง index ใหม่แล้วชี้ alias ไปที่ index ใหม่ได้ โดยที่เราไม่ต้องไปแก้โค้ดเพื่อเปลี่ยนชื่อ index ใหม่

โพสต์นี้บันทึกวิธีลิสต์ alias ที่มีอยู่แล้วออกมาไว้หน่อย

<!--more-->

เราสามารถลิสต์ alias ออกมาได้โดยเรียก method GET ไปที่ path นี้

```txt
GET /_alias
```

ซึ่งจะได้ผลลัพธ์ออกมาเป็น JSON หน้าตาประมาณนี้

```json
{
  "a-index-000001" : {
    "aliases" : {
      "a-index" : { }
    }
  },
  "a-index-000002" : {
    "aliases" : {
      "a-index" : {
        "is_write_index" : true
      }
    }
  },
  "b-index-000001" : {
    "aliases" : {
      "b-index" : { }
    }
  },
  "b-index-000002" : {
    "aliases" : {
      "index-b" : {
        "is_write_index" : true
      }
    }
  },
}
```

คือจะเป็น object ที่ key เป็นชื่อ index แล้วข้างในก็มี object ที่มี key ชื่อ aliases ที่ข้างในก็จะเป็น ชื่อ alias map กับ configuration ของ alias นั้นๆอีกที
