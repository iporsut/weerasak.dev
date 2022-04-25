---
title: "ทดสอบโค้ด Go ที่อยู่ในรูปแบบ message loop goroutine"
date: 2022-04-25T22:14:56+07:00
draft: false
---

วันนี้เจอโค้ด Go ที่อยู่ในรูปแบบ message loop แล้วต้องเขียนเทสเพื่อทดสอบ goroutine เลยได้ท่าแบบนี้ในการเขียนเทสออกมา

<!--more-->

โค้ดอยู่ในรูปแบบประมาณนี้

```go
package sample

import (
	"context"
	"fmt"
)

var println = fmt.Println

func loop(ctx context.Context, msgCh chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgCh:
			println(msg)
		}
	}
}
```

จะเห็นว่าตัวฟังก์ชันรับค่า context กับ channel ซึ่งในตัวฟังก์ชันนั้นจะวนลูปเพื่อเช็คว่ามี message ส่งมาทาง channel หรือไม่ ถ้ามีก็ปริ้น message นั้นออกไป

หรือ ถ้า context นั้นถูก cancel/timeout ไปก่อนแล้ว (เรียก method Done() เพื่อรับค่า channel แล้วพยายามดึงค่าออกมาจาก channel นั้นเพื่อเช็คว่า context นั้นถูก cancel/timeout แล้วหรือยัง) ถ้า context cancel/timeout ไปก่อนแล้วก็จบฟังก์ชัน

ส่วนโค้ดที่เทสสำหรับฟังก์ชันนี้ก็จะได้ประมาณนี้

```go
package sample

import (
	"context"
	"fmt"
	"testing"
)

func TestLoop(t *testing.T) {
	t.Cleanup(func() { println = fmt.Println })

	// stub println function
	var msg string
	println = func(a ...any) (n int, err error) {
		msg = a[0].(string)
		return 1, nil
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	msgCh := make(chan string)

	// run loop in separated goroutine
	go loop(ctx, msgCh)

	// send message
	msgCh <- "Hello, World"

	// trigger done context
	cancelFn()

	if msg != "Hello, World" {
		t.Fatalf("expect \"Hello World\" but got %q", msg)
	}
}
```

- เรา stub ฟังก์ชัน println ก่อนแล้วเก็บค่า parameter ที่ส่งมาไว้ในตัวแปร msg เพื่อเอาไว้เช็คว่าโดนเรียกหรือไม่
- จากนั้นก็สร้าง context ใหม่ โดยใช้ WithCancel เพื่อให้ได้ cancelFn กลับมาด้วยเพื่อสั่ง cancel context ได้
- สร้าง message channel
- รัน loop แยกในอีก goroutine นึงพร้อม context กับ message channel ไป
- ทดสอบส่ง message ไปทาง msgCh
- สั่ง cancelFn() เพื่อ cancel context แล้วทำให้ goroutine จบการทำงาน
- เช็คว่า msg มีค่าเท่ากับ message ที่ส่งไปหรือไม่

ก็เท่านี้ น่าจะพอเห็น pattern เอาไว้ทดสอบถ้าเจอฟังก์ชันแบบนี้อีก
