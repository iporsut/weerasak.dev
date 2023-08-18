---
title: "Go: วิธีเขียน http handler ให้รู้ว่า request ถูก cancel ไปแล้ว"
date: 2023-08-18T14:00:00+07:00
draft: false
---

net/http (หรือ gin) handler จะมีจะรับ request type เข้ามาคือ type \*http.Request ซึ่งในนี้จะมีตัวแปร context ของ request ที่เราสามารถใช้เช็คได้ว่า request ถูก cancel ไปแล้วหรือยัง

<!--more-->

ตัวอย่างเช่นถ้าเรามี handler แบบนี้

```go
package main

import (
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                log.Println("start")
                defer log.Println("done")
		// simulate long running process
		time.Sleep(60 * time.Second)
		w.Write([]byte("Hello World"))
	})
	http.ListenAndServe(":8000", nil)
}
```

ถ้าเราเรียกไปที่ endpoint นี้ด้วย `curl http://localhost:8000` จะต้องรอ 60 วินาทีถึงจะได้ response กลับมาเป็น `"Hello World"` และเราจะเห็น log แบบนี้

```
2023/08/18 14:14:20 start
2023/08/18 14:15:20 done
```

และถ้าเราเรียกอีกครั้ง แต่กด ctrl+c ก่อนเพื่อยกเลิกการทำงานของ curl และปิด connection ไปแล้วนั้น จะเห็นว่า `time.Sleep` ก็ยังทำงานต่อไปจนครบ 60 วินาทีแล้วค่อยจบการทำงาน ยังเห็น log pattern เดิมแบบนี้

```
2023/08/18 14:14:20 start
2023/08/18 14:15:20 done
```

คือทำจนครบ 60 วินาที

ถ้าเราอยากจะเช็คได้ว่า request ถูก cancel แล้วหรือยังให้เราเขียนโค้ดเช็คจาก request.Context() ได้โดยใช้เมธอด Done() ของ request.Context() เพื่อดักจับ signal จาก channel นี้ว่า request cancel แล้วหรือยัง นอกจากนั้นเราจะเปลี่ยน time.Sleep เป็น time.NewTimer เพื่อที่จะสามารถยกเลิกการ sleep ได้ผ่าน เมธอด timer.Stop

```go
package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("start")
		defer log.Println("done")
		// simulate long running process with aware of context
		tm := time.NewTimer(60 * time.Second)
		select {
		case <-r.Context().Done():
			tm.Stop()
		case <-tm.C:
			w.Write([]byte("Hello World"))
		}
	})
	http.ListenAndServe(":8000", nil)
}
```

อย่างไรก็ตาม ต่อให้เราใช้ r.Context.Done() เพื่อเช็คแล้ว แต่พอ request ถูก disconnect จากฝั่ง client แล้วแต่กลับยังไม่มี signal ส่งมาที่ Done channel นั้นเป็นเพราะว่า client อาจจะส่ง request body มาด้วย แต่ handler ไม่ยอม read request body ที่ส่งมาให้หมดก่อน ตัวอย่างเช่นถ้าเราเรียก curl ด้วย POST และมี body ด้วย แบบนี้

```
 curl -XPOST http://localhost:8000 -d '{}'
```

แล้วกด ctrl-c ตัว handler ก็ยังจะรอ 60 วิอยู่ดี

วิธีที่จะทำให้เช็คได้คือ ต้อง flush request body ออกให้หมดก่อนด้วย ซึ่งการเขียน REST หรือ HTTP API ปกติ เราก็จะ unmarshaling request body แล้วทำให้อ่านจนหมดกันอยู่แล้ว ถ้าจะจำลองก็ใช้ io.Copy ได้แบบนี้

```go
package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("start")
		defer log.Println("done")
                // simulate unmarshaling request body
		io.Copy(io.Discard, r.Body)
		// simulate long running process with aware of context
		tm := time.NewTimer(60 * time.Second)
		select {
		case <-r.Context().Done():
			tm.Stop()
		case <-tm.C:
			w.Write([]byte("Hello World"))
		}
	})
	http.ListenAndServe(":8000", nil)
}
```

กลไกการเช็คแบบนี้เป็น pattern ที่ package context ใช้งาน และ pattern ที่เราเห็นได้เป็นปกติใน Go คือส่ง ctx contex.Context เป็น parameter แรกต่อๆกันไป เพื่อให้ process ที่ทำงานเช็คได้ว่าควรจะยกเลิกการทำงานที่เหลือเวลา request ถูกยกเลิกไปแล้วนั่นเอง
