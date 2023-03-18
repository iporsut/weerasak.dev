---
title: "[Go] ใช้ bytes.Buffer เป็น io.Writer เพื่อเก็บ output ใน memory ก่อนแปลงเป็น string หรือ []byte"
date: 2023-03-18T16:30:00+07:00
draft: false
---

มีฟังก์ชันการทำงานหลายตัวที่รับค่าเป็น interface `io.Writer` เพื่อเขียน output ออกไป ซึ่งช่วยให้คนใช้งานเลือกได้ว่าจะเขียน output การทำงานออกไปทางไหนก็ได้ตามต้องการขอแค่การทำงานนั้นถูก implements ด้วย interface `io.Writter` ที่มี method `Write([]byte) (int, error)` อยู่ และ `bytes.Buffer` ก็เป็นหนึ่งใน `io.Writer` ที่ช่วยให้เราเอาผลลัพธ์การทำงานมาเก็บไว้ใน memory ก่อนแล้วค่อยเอาค่าไปใช้ต่อได้เช่นแปลงเป็น `string` แปลงเป็น `[]byte` หรือแม้แต่ส่งไปให้ฟังก์ชันที่ต้องการ `io.Reader` เพราะ `bytes.Buffer` ก็ implement interface `io.Reader` ไว้เช่นกัน

<!--more-->

ตัวอย่างเช่นผมต้องการใช้ package `html/template` เพื่อสร้าง HTML content ทีนี้การใช้งาน template จะมี method `Execute(wr io.Writer, data any) error` อยู่ให้เราส่ง `io.Writer` เข้าไปเพื่อให้ template write output ที่ได้จากการ render `data`

แต่ว่าเราต้องการให้ output เก็บไว้ใน memory ก่อนเราก็ใช้ `bytes.Buffer` ช่วยได้ แบบนี้

```go
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

var tpl = template.Must(template.New("html").Parse(`<!doctype html>
<html>
<body>
<h1>{{.}}</h1>
</body>
</html>
`))

func main() {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, "Hello, World")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())

	// Output:
	// <!doctype html>
	// <html>
	// <body>
	// <h1>Hello World</h1>
	// </body>
	// </html>
}
```

- เราส่ง pointer ของ `bytes.Buffer` เข้าไปซึ่ง implements `io.Writer` อยู่
- จากนั้นเราก็เรียก method `String()` เพื่อให้แปลงข้อมูลใน `buf` ออกมาเป็น String
