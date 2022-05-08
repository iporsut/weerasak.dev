---
title: "Elasticsearch max_expansions applies to shard level"
date: 2022-05-08T19:47:01+07:00
draft: false
---

โพสต์นี้สั้นๆ อยากจะโน้ตเอาไว้หน่อยว่า `max_expansions` option ของ Elasticsearch query นั้นมัน apply ที่ shard level

<!--more-->

หมายความว่าถ้าเรากำหนดค่า `max_expansions` เป็น 50 แล้วเรามี shards ทั้งหมด 5 shards เวลามันจะเลือก documents มันก็จะเลือกมาทั้งหมด shards ละ 50 เป็น 250 documents นั่นเอง
