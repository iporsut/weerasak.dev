---
title: "Go io.Discard ใช้ตอนไหน"
date: 2022-04-30T16:16:32+07:00
draft: false
---

วันนี้เจอคำถามใน FB Group [Golang Thailand](https://www.facebook.com/groups/584867114995854) ว่า `ioutil.Discard` (ตั้งแต่ 1.16 ioutil.Discard ถูกย้ายไป `io.Discard` แล้ว) ควรใช้ตอนไหน เลยไปค้นใน standard lib ว่าเขาใช้ตอนไหนกันบ้าง

<!--more-->

ผมเลือกไปค้นใน package path "net" แล้วก็ path อื่นๆที่อยู่ภายใต้ "net" ด้วย ค้นเจอว่ามีการใช้ `io.Discard` แบบนี้

```txt
$ pwd
/Users/weerasak/src/go/src/net
$ ag "io\.Discard"
http/httputil/dump.go
135:			io.Copy(io.Discard, req.Body)

http/server.go
1371:			_, err := io.CopyN(io.Discard, w.reqBody, maxPostHandlerReadBytes+1)
3494:		io.Copy(io.Discard, mb)

http/fcgi/child.go
316:	io.CopyN(io.Discard, body, 100<<20)


http/client.go
697:				io.CopyN(io.Discard, resp.Body, maxBodySlurpSize)

http/transfer.go
375:			nextra, err = t.doBodyCopy(io.Discard, body)
997:			n, err = io.CopyN(io.Discard, bodyLocked{b}, maxPostHandlerReadBytes)
1008:		_, err = io.Copy(io.Discard, bodyLocked{b})
```

จะเห็นว่าเขาใช้ `io.Discard` โดยการ copy (เรียกใช้ `io.Copy`) ข้อมูลจาก `io.Reader` (type ที่ implements interface `io.Reader`) แล้วเขียนไปที่ `io.Discard` ที่เป็น `io.Writer` (type ที่ implements `io.Writer`)

ใน [document]() ของ `io.Discard` เขียนเอาไว้ว่า

```go
var Discard Writer = discard{}

// Discard is a Writer on which all Write calls succeed without doing anything.
```

นั่นคือเราเขียนข้อมูลไปที่ `io.Discard` มันจะเขียนได้สำเร็จ แต่มันไม่ได้ทำอะไรทั้งนั้นล่ะ

ดังนั้นที่เราเห็นโค้ด copy ข้อมูลไปใส่ `io.Discard` มันคือการ flush ข้อมูล `io.Reader` ทิ้งนั้นเอง

ถ้าใครเขียน shell คงเคยเห็น pattern แบบนี้

```txt
$ grep -r hello /sys/ 2> /dev/null
```

ที่หมายถึงเอาข้อมูลจาก standard error (file descriptor number 2) redirect ไปที่ `/dev/null` ซึ่งมันก็จะ redirect ไปทิ้งเฉยๆ ไม่มีอะไรเกิดขึ้น ซึ่งการใช้งาน `io.Copy` ไปที่ `io.Discard` จะเหมือน pattern ที่ใช้การ redirect ไปที่ `/dev/null` ของ shell นั่นเอง
