---
title: "[Go] any type argument satisfies comparable constraint in Go 1.20"
date: 2023-02-26T09:11:22+07:00
draft: false
---

ขอสรุปเรื่องของภาษาที่เปลี่ยนไปใน Go 1.20 เกี่ยวกับกระบวนการ satifying a type constraint ของ type `any` เอาไว้หน่อย

<!--more-->

(ป.ล. โพสต์นี้จะใช้คำภาษาอังกฤษไปเลย ไม่ได้อธิบายเพิ่ม เอาที่เจ้าของโพสต์เข้าใจ ไว้ว่างๆจะหาเวลามาอธิบายแต่ละคำทีหลังนะ)

ใน Go1.18 และ Go1.19 หลังจากที่รองรับ Generic type ก็จะรองรับการเขียนโค้ดที่เราสามารถมี type parameter ติดไว้ได้พร้อมระบุ constraint ว่า type argument ที่เป็นไปได้ที่จะ apply type parameter นั้นมีอะไรบ้าง (type set)

แต่มีเรื่องนึงที่ยังทำให้มีพฤติกรรมแปลกๆคือ type interface กับ `comparable` constraint

เพราะว่าตามสเปคของภาษา interface นั่นน่ะ compare ด้วย operator `==` ได้ โดยเงื่อนไขคือถ้า concrete type ที่เก็บใน interface เป็น type เดียวกัน และ type นั้น compare กันได้ และค่าที่เก็บเป็นค่าเดียวกัน ก็จะได้ true, ถ้าเป็นคนละ type จะได้ false, แต่ถ้า concret type เหมือนกัน แล้ว type นั้นไม่สามารถ compare กันได้ ก็จะเกิด panic แทน

ซึ่งกระบวนการเช็คเกิดตอน runtime ตอน compile time การเปรียบเทียบ interface Go จะยอมให้ compile ผ่านไปได้ (แน่ละ เพราะ interface จะรู้ concret type ก็ตอน runtime แล้ว)

ตัวอย่างง่ายๆเช่นถ้าเรามีโค้ดแบบนี้

```go
func Compare(a any, b any) bool {
        return a == b
}
```

เราใช้งานแบบนี้ได้ไม่เกิด panic แม้ว่า slice จะไม่สามารถ compare ได้

```go
func main() {
	fmt.Println(Compare("Hello", "Hello")) // true
	fmt.Println(Compare("Hello", []int{1})) // false
	fmt.Println(Compare([]int{1}, "Hello")) // false
}
```

แต่แบบนี้จะ panic

```go
func main() {
	fmt.Println(Compare([]int{1}, []int{1})) // panic: runtime error: comparing uncomparable type []int
}
```

แต่พอเป็นเรื่อง Generic type parameter กับ comparable constraint นั่นต่างกันออกไป เพราะ การ resolve type ของ type parameter นั้นเกิดตอน compile time เท่านั้น ทำให้แม้ว่า interface จะ compare กันได้ แต่ไม่สามารถเข้าเงื่อนไขของ comparable constraint ได้เช่น

```go
// Go 1.18 and Go 1.19 will compile error
// Go 1.20 can compile and print true

func Compare[T comparable](a, b T) bool {
	return a == b
}

func main() {
	var a any = 2
	var b any = 2
	fmt.Println(Compare[any](a, b)) // compile error: any does not implement comparable
}
```

นั่นทำให้ Go 1.20 เกิดการเปลี่ยนแปลงของ spec ภาษาเพื่อรองรับให้สามารถ apply type argument ที่เป็น interface ได้ถึงแม้ว่า interface type จะไม่เข้าเงื่อนไขของ comparable แต่ว่า interface type นั้น compare กันได้ (แม้ว่าจะเป็นตอน compile time ก็ตาม)

ใน Go 1.20 โค้ดเมื่อกี้จะทำงานได้และได้ค่าเป็น true ซึ่งก็พฤติกรรมก็จะเหมือน Compare ก่อนหน้าที่ใช้ type any ปกติไม่ที่จะมีโอกาส panic ได้ถ้าโค้ดเป็นแบบนี้

```go
func Compare[T comparable](a, b T) bool {
	return a == b
}

func main() {
	var a any = []int{1}
	var b any = []int{1}
	fmt.Println(Compare[any](a, b)) // panic: runtime error: comparing uncomparable type []int
}
```

เรื่องนี้ไม่ได้ apply แค่กับ type `any` ตรงๆเท่านั้น แต่ยังหมายถึง type interface อื่นๆ หรือ type อื่นๆที่มี element เป็น type interface เช่น struct, array อีกด้วย เช่น

```go
type T struct {
        D any
}

func Compare[T comparable](a, b T) bool {
	return a == b
}

func main() {
	var a any = T{D: 1}
	var b any = T{D: 1}
	fmt.Println(Compare[any](a, b)) // true in Go 1.20, compile error in Go 1.18, 1.19
}
```

ใครสนใจอ่านสเปคภาษาจริงๆลองอ่านดูได้ที่นี่ https://go.dev/ref/spec#Satisfying_a_type_constraint หรืออ่าน blog ที่ Robert Griesemer เขียนไว้แบบละเอียดๆได้ที่นี่ https://go.dev/blog/comparable
