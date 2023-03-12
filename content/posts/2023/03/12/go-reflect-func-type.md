---
title: "[Go] ใช้ reflect ดึงข้อมูล parameter type และ return type ของ function"
date: 2023-03-12T09:57:09+07:00
draft: false
---

เมื่ออาทิตย์ที่ผ่านมาเจอ package นึงชื่อ [tonic](https://github.com/loopfz/gadgeto/tree/master/tonic) ช่วยในการเขียน Gin handler ให้มี request parameter และ return type ได้ เลยไปขุดดูว่าทำได้ยังไง เจอว่าใช้ reflect ในการหาข้อมูลของ handler function นั่นเองว่ารับ type อะไร หรือ return type อะไร

<!--more-->

วิธีการใช้งานก็คือเรียก function `reflect.TypeOf` โดยส่งชื่อฟังก์ชันที่เราต้องการดูข้อมูล type เข้าไป

```go
package main

import (
	"fmt"
	"reflect"
)

func add(a int, b int) int {
	return a + b
}

func main() {
	ft := reflect.TypeOf(add)
	fmt.Println(ft) // output: func(int, int) int
}
```

ถ้าอยากรู้ชื่อฟังก์ชันต้องใช้ฟังก์ชัน `runtime.FuncForPC` จาก package `runtime` เข้ามาช่วยโดยเราจะส่ง uintptr ซึ่งเป็น raw pointer ที่ชี้ไปที่โค้ดของฟังก์ชันตอน runtime โดยเรียก method `.Pointer()` ของ `relect.Value` เข้าไปแทนแบบนี้ แล้วเรียก method `Name()` อีกที ก็จะได้ชื่อฟังก์ชันในรูปแบบที่มีชื่อ package dot ข้างหน้าด้วย

```go
package main

import (
	"fmt"
	"reflect"
	"runtime"
)

func add(a int, b int) int {
	return a + b
}

func main() {
	funcValue := reflect.ValueOf(add) // ดึง reflect.Value ของ function add
	runtimeFunc := runtime.FuncForPC(funcValue.Pointer()) // หา runtime.Func โดยส่ง Pointer() ให้กับ FuncForPC
	fmt.Println(runtimeFunc.Name()) // main.add
}
```

ต่อไปถ้าเราอยากได้ข้อมูลของ parameter ว่ารับ type อะไรบ้าง เราจะใช้ method `NumIn()` และ `In` ของ `reflect.Type` ที่ได้ตอนเรียก `reflect.TypeOf` (หรือหา `reflect.Type` จาก method `Type()` ของ `reflect.Value` ได้เช่นกัน) ตัวอย่างเช่น

```go
package main

import (
	"fmt"
	"reflect"
)

func add(a int, b int) int {
	return a + b
}

func main() {
	funcValue := reflect.ValueOf(add)
	funcType := funcValue.Type()

	funcNumParameter := funcType.NumIn()
	fmt.Println(funcNumParameter) // 2

	for i := 0; i < funcNumParameter; i++ {
		parameterType := funcType.In(i)
		fmt.Println(parameterType)
	}

	// Output:
	// 2
	// int
	// int
}
```

ซึ่งเราจะได้ข้อมูลของ type แต่เราจะไม่ได้ข้อมูลว่า paramter name ชื่ออะไรเพราะสำหรับ reflect.Type นั้นมองว่า parameter name นั้นไม่ใช่ส่วนหนึ่งของ type นั่นเอง (ส่วนนี้จะเกิดขอบเขตการทำงานของ relect ของ Go ถ้าอยากได้ละเอียดขนาดนั้นต้องทำการ parse sourcecode)

ส่วนของ output หรือ return type ก็เช่นกัน แค่เปลี่ยนจาก NumIn เป็น NumOut และ In เป็น Out แบบนี้

```go
package main

import (
	"fmt"
	"reflect"
)

func add(a int, b int) int {
	return a + b
}

func main() {
	funcValue := reflect.ValueOf(add)
	funcType := funcValue.Type()

	funcNumReturn := funcType.NumOut()
	fmt.Println(funcNumReturn) // 1

	for i := 0; i < funcNumReturn; i++ {
		returnType := funcType.Out(i)
		fmt.Println(returnType)
	}

	// Output:
	// 1
	// int
}
```

ตัว package tonic ก็อาศัยข้อมูลจาก reflect พวกนี้แหละครับเพื่อสร้าง handler ที่รับ type สำหรับ request และ return type สำหรับ response แล้วสร้าง gin Handler ให้ที่มีการทำงานในการ parse http request ไปเป็น request type แล้วก็ marshaling response type กลับไปเป็น http response ให้เอง ลองดูโค้ดของ tonic เต็มๆได้ที่นี่ https://github.com/loopfz/gadgeto/blob/master/tonic/handler.go
