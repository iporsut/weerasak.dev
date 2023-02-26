---
title: "[Go] สรุปวิธีที่ Go ใช้ implements type http.ResponseController ใน Go 1.20"
date: 2023-02-26T16:12:32+07:00
draft: false
---

ขอสรุปวิธีที่ Go 1.20 ใช้ implements type `http.ResponseController` ของ package `http` เอาไว้หน่อย เพราะเป็น pattern ที่น่าสนใจดีนะ ว่าทำยังไงให้สามารถเพิ่มความสามารถของ ResponseWriter interface โดยที่ไม่ต้องปรับ type ResponseWriter โดยตรงเพราะจะทำให้คนที่เคย implements ต้องแก้โค้ดเพื่อ implements method เพิ่มขึ้น

<!--more-->

Go 1.20 เลือกที่จะสร้าง type ใหม่ชื่อว่า `http.ResponseController` โดยภายใน http.ResponseController ก็มี field ที่เก็บค่าของ tyep ResponseWriter เอาไว้แบบนี้

```go
type ResponseController struct {
	rw ResponseWriter
}
```

และสร้าง constructor function เอาไว้แบบนี้

```go
func NewResponseController(rw ResponseWriter) *ResponseController {
	return &ResponseController{rw}
}
```

ตัว ResponseController เองนั้นมี 4 methods คือ

```go
func (c *ResponseController) Flush() error

func (c *ResponseController) Hijack() (net.Conn, *bufio.ReadWriter, error)

func (c *ResponseController) SetReadDeadline(deadline time.Time) error

func (c *ResponseController) SetWriteDeadline(deadline time.Time) error
```

ซึ่งถ้าเราดูการ implements method ทั้ง 4 ตัวนี้จะพบว่า ResponseController ไม่ได้ implements การทำงานด้วยตัวมันเอง แต่จะทำการเช็คว่า rw ที่เป็น ResponseWriter นั้นตัว concret type มัน implements method พวกนี้ไว้อยู่หรือเปล่าโดยใช้ type assertion เช่น Flush implements ไว้แบบนี้

```go
// Flush flushes buffered data to the client.
func (c *ResponseController) Flush() error {
	rw := c.rw
	for {
		switch t := rw.(type) { // ใช้ switch type assertion
		case interface{ FlushError() error }: // ใช้ interface literal เลยไม่สร้าง type ใหม่
			return t.FlushError()
		case Flusher: // เช็คจาก interface http.Flusher
			t.Flush()
			return nil
		case rwUnwrapper:
			rw = t.Unwrap()
		default:
			return errNotSupported()
		}
	}
}
```

ถ้าค่าที่เก็บใน rw implements FlushError method ก็เรียก FlushError ถ้าไม่ใช่ก็เช็คต่อว่า implements Flusher หรือไม่ถ้าใช่ก็เรียก method Flush ของ interface Flusher

ถ้ายังไม่ได้ implements ทั้งสอง interface เมื่อกี้ ResponseController ยังรองรับให้ type implements rwUnwrapper ซึ่งคือ implements method Unwrap ที่จะแกะค่าของ type ResponseWriter อื่นๆที่ wrap กันอยู่ออกมาได้เพื่อเอาไปวนซ้ำหา interface FlushError หรือ Flusher ที่ต้องการ

สุดท้ายถ้าไม่ตรงตามเงื่อนไขเลยก็จะ return error กับไปว่า ErrNotSupported

3 methods ที่เหลือก็ใช้ pattern เดียวกัน

```go
// Hijack lets the caller take over the connection.
// See the Hijacker interface for details.
func (c *ResponseController) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	rw := c.rw
	for {
		switch t := rw.(type) {
		case Hijacker:
			return t.Hijack()
		case rwUnwrapper:
			rw = t.Unwrap()
		default:
			return nil, nil, errNotSupported()
		}
	}
}

// SetReadDeadline sets the deadline for reading the entire request, including the body.
// Reads from the request body after the deadline has been exceeded will return an error.
// A zero value means no deadline.
//
// Setting the read deadline after it has been exceeded will not extend it.
func (c *ResponseController) SetReadDeadline(deadline time.Time) error {
	rw := c.rw
	for {
		switch t := rw.(type) {
		case interface{ SetReadDeadline(time.Time) error }:
			return t.SetReadDeadline(deadline)
		case rwUnwrapper:
			rw = t.Unwrap()
		default:
			return errNotSupported()
		}
	}
}

// SetWriteDeadline sets the deadline for writing the response.
// Writes to the response body after the deadline has been exceeded will not block,
// but may succeed if the data has been buffered.
// A zero value means no deadline.
//
// Setting the write deadline after it has been exceeded will not extend it.
func (c *ResponseController) SetWriteDeadline(deadline time.Time) error {
	rw := c.rw
	for {
		switch t := rw.(type) {
		case interface{ SetWriteDeadline(time.Time) error }:
			return t.SetWriteDeadline(deadline)
		case rwUnwrapper:
			rw = t.Unwrap()
		default:
			return errNotSupported()
		}
	}
}

// errNotSupported returns an error that Is ErrNotSupported,
// but is not == to it.
func errNotSupported() error {
	retu
```

จะเห็นว่าการทำแบบนี้ทำให้ Go 1.20 ไม่ต้องเกิดการ breaking change type `ResponseWriter` คนที่เคย implements ResponseWriter เอาไว้ก็ไม่ต้องแก้โค้ด

ส่วนคนที่อยากให้ type ตัวเองรองรับการทำงานอื่นๆนอกจาก method ที่ implements ให้กับ ResponseWriter ก็สามารถ implements เพิ่มได้แล้วสร้าง ResponseController ผ่าน http.NewResponseController ก็สามารถเรียกใช้ได้เลยไม่ต้องเขียนโค้ด assert type เองเพราะมัน wrap เอาไว้ให้แล้ว

ref: ไปดู sourcecode เต็มๆได้ที่นี่นะ https://cs.opensource.google/go/go/+/refs/tags/go1.20.1:src/net/http/responsecontroller.go
