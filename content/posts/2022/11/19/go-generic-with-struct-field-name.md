---
title: "Go generic with struct field name"
date: 2022-11-19T13:47:20+07:00
draft: false
---

Go [generic](https://go.dev/blog/intro-generics) type ที่มีมาตั้งแต่ Go version 1.18 นั้นช่วยให้เราลดโค้ดที่ซ้ำซ้อนกรณีที่ logic ของ function เหมือนกันหมด ต่างกันแค่ type ของ parameter ลงไปได้ แต่ถ้า code ของ function เรารับค่าเป็น struct ใช้งาน field เหมือนกันหมด ต่างกันแค่คนละ type เรายังไม่สามารถใช้ generic ช่วย refactor code ออกมาเป็น function เดียวได้ง่ายๆ ตัวอย่างเช่น

<!--more-->

```go
func FillEmployeeNameValue(e Employee) {
        e.Name = e.FirstName + " " + e.LastName
}

func FillCustomerNameValue(c Customer) {
        c.Name = c.FirstName + " " + c.LastName
}

// ใช้ Generic
func FillNameValue[T any](t T) {
        t.Name = t.FirstName + " " + t.LastName
}
```

โค้ดแบบนี้ยังไม่สามารถใช้งานได้ใน Go มีคนเปิด issue เอาไว้ให้แล้วด้วยที่นี่ https://github.com/golang/go/issues/48522 แต่ก็ยังไม่มีแผนว่าจะไปท่าไหน

ในระหว่างที่ยังไม่มี syntax ให้กำหนด constraint เกี่ยวกับ field name ได้ ตอนนี้ทางเลือกที่จะพอยุบโค้ดที่ซ้ำซ้อนตรงนี้ได้ก็คือ

1. ใช้ interface แทนที่จะใช้ generic type เราใช้ interface กำนหด method ที่จำเป็นสำหรับ function ที่เรา implement logic ต้องใช้ เช่นโค้ดด้านบนถ้าเปลี่ยนมาใช้ interface ธรรมดา ก็อาจจะได้แบบนี้

```go
type NameAccessor interface {
        SetName(name string)
        GetFirstName() string
        GetLastName() string
}

func FillNameValue(t NameAccessor) {
        t.SetName(t.GetFirstName() + " " + t.GetLastName())
}

// type ที่ต้องใช้ ต้อง implements NameAccessor เช่น
type Employee {
        Name, FirstName, LastName string
}

func (e *Employee) SetName(name string) { e.Name = name }
func (e *Employee) GetFirstName() string { return e.FirstName }
func (e *Employee) GetLastName() string { return e.LastName }
```

ข้อเสียคือ type จำเป็นต้อง implements method ซึ่งบางครั้งเราก็ไม่ได้อยากให้ type เรามี public method เหล่านี้ออกมาให้คนใช้งานเห็นเหมือนกัน จะทำเป็น public ก็ไม่ได้ถ้าเกิด FillNameValue และ NameAccessor อยู่คนละ package

2. มีอีกวิธีนึงที่ไม่ต้องใช้ interface คือ สร้าง type ขึ้นมาใหม่สำหรับเป็นตัวกลาง แล้วค่อยสร้าง function สำหรับแปลงค่าจาก type ที่เป็น input เป็น type ที่เราสร้างมาใหม่ ตัวอย่างเช่น

```go
type NameAccessor struct {
        Name, FirstName, LastName *string // ใช้ pointer เพื่อให้เราแก้ค่าของ Name ได้
}

func FillNameValue(t NameAccessor) {
        *t.Name = *t.FirstName + " " + *t.LastName
}

// สร้าง convert function สำหรับ type ที่ต้องการเช่น
type Employee struct {
	Name, FirstName, LastName string
}

func EmployeeToNameAccessor(e *Employee) NameAccessor {
	return NameAccessor{
		Name:      &e.Name,
		FirstName: &e.FirstName,
		LastName:  &e.LastName,
	}
}

func main() {
	e := Employee{
		FirstName: "John",
		LastName:  "Doe",
	}

        // ตอนเรียกใช้งานก็ใช้แบบนี้
	FillNameValue(EmployeeToNameAccessor(&e))

	fmt.Println(e.Name)
}
```

ทีนี้มาดูอีกตัวอย่างให้เห็นเวลาใช้ร่วมกับ generic สมมติเราต้องรับ slice ของ T ([]T) แล้ววนลูปเพื่อ fill ค่า name มาดูโค้ดแบบแรกที่คอมไพล์ไม่ผ่านกันก่อน

```go
func FillNameValue[T any](vals []*T) {
        for _, val := range vals {
                val.Name = val.FirstName + " " val.LastName
        }
}
```

ทีนี้เราจะเพิ่มให้ FillNameValue รับฟังก์ชัน convert จาก *T เป็น NameAccessor ด้วยแบบนี้

```go
func FillNameValue[T any](vals []*T, accFn func(*T) NameAccessor) {
	for _, val := range vals {
		accessor := accFn(val)
		*accessor.Name = *accessor.FirstName + " " + *accessor.LastName
	}
}
```

จะเห็นว่าเราทำการเรียก accFn(val) เพื่อให้ได้ accessor ก่อนเพื่อเราไปใช้งานในการกำหนดค่า Name
พอได้แบบนี้เราก็ใช้งาน FillNameValue กับ slice ของ *T ใดๆได้แล้ว เพียงแค่เขียน converter function ขึ้นมาเพื่อแปลงเป็น Accessor type สำหรับ field ที่ต้องการ ตัว Type เองก็ไม่ต้องมี method อื่นๆที่ไม่จำเป็น ตัว converter function เราจะเขียนเป็น annonymous ก็ยังได้ถ้าใช้แค่ครั้งเดียวไม่ได้ reuse อะไรบ่อยๆ

ตัวอย่างการเรียกใช้

```go
package main

import "fmt"

type Employee struct {
	Name, FirstName, LastName string
}

type NameAccessor struct {
	Name, FirstName, LastName *string // ใช้ pointer เพื่อให้เราแก้ค่าของ Name ได้
}

func EmployeeToNameAccessor(e *Employee) NameAccessor {
	return NameAccessor{
		Name:      &e.Name,
		FirstName: &e.FirstName,
		LastName:  &e.LastName,
	}
}

func FillNameValue[T any](vals []*T, accFn func(*T) NameAccessor) {
	for _, val := range vals {
		accessor := accFn(val)
		*accessor.Name = *accessor.FirstName + " " + *accessor.LastName
	}
}

func main() {
	empls := []*Employee{
		{
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			FirstName: "Jane",
			LastName:  "Doe",
		},
	}

	FillNameValue(empls, EmployeeToNameAccessor)

	for _, empl := range empls {
		fmt.Println(empl.Name)
	}

        // Output:
        // John Doe
        // Jane Doe
}
```

ตัวอย่างเขียน converter เป็น anonymous function

```go
func main() {
	empls := []*Employee{
		{
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			FirstName: "Jane",
			LastName:  "Doe",
		},
	}

	FillNameValue(empls, func(e *Employee) NameAccessor {
                return NameAccessor{
                        Name:      &e.Name,
                        FirstName: &e.FirstName,
                        LastName:  &e.LastName,
                }
        })

	for _, empl := range empls {
		fmt.Println(empl.Name)
	}

        // Output:
        // John Doe
        // Jane Doe
}
```

ลองเล่นใน Go playgroud ได้ที่นี่ https://go.dev/play/p/HAchEsjr4md
