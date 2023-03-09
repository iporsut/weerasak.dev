---
title: "ทำไม Go return pointer ของตัวแปร local ได้ แต่ Rust return reference ของตัวแปร local ไม่ได้"
date: 2023-03-05T07:49:19+07:00
draft: false
---

ลองมาทำความเข้าใจกันว่าทำไม Go return pointer ของตัวแปร local ได้ แต่ Rust return reference ของตัวแปร local ไม่ได้

<!--more-->

ลองดูโค้ด Go กันก่อน ตัวอย่างเช่นเรามี struct User เรามักสร้าง function constructor ที่ initialize ค่าแล้วตอบกลับเป็น pointer ของ User แบบนี้

```go
package main

import "fmt"

type User struct {
	Name  string
	Email string
}

func NewUser(name, email string) *User {
	return &User{
		Name:  name,
		Email: email,
	}
}

func main() {
	u := NewUser("me", "me@email.com")
	fmt.Println(u)
}
```

ทีนี้ถ้าเราสร้างโค้ดคล้ายๆกันแบบนี้ด้วย Rust แบบนี้

```rust
struct User {
    name: String,
    email: String
}

impl User {
    fn new<'a> (name: String, email: String) -> &'a User {
        return &User{name: name, email: email}
    }
}

fn main() {
    let u = User::new(String::from("me"), String::from("me@email.com"));
    println!("{}, {}", u.name, u.email);
}
```

โค้ดนี้จะไม่สามารถ compile ผ่านได้เลย Rust จะขึ้น compile error แบบนี้

```txt
error[E0515]: cannot return reference to temporary value
 --> src/main.rs:8:16
  |
8 |         return &User{name: name, email: email}
  |                ^------------------------------
  |                ||
  |                |temporary value created here
  |                returns a reference to data owned by the current function

For more information about this error, try `rustc --explain E0515`.
error: could not compile `playground` due to previous error
```

Rust compiler บอกว่า `User{name: name, email: email}` เนี่ยถูกสร้างใน function ซึ่ง ownership อยู่ใน scope ของ function ซึ่ง Rust จะคืนค่าหน่วยความจำเมื่อจบการทำงานของฟังก์ชัน ดังนั้นการที่เราจะส่ง reference ของสิ่งที่จะถูกคืนหน่วยความจำกลับไป Rust compiler จึงไม่ยอมให้เกิดขึ้น

แล้วทำไม Go ถึงยอมให้เกิดขึ้น

เพราะว่า Go นั้นเลือกที่จะมี runtime และ Garbage Collector (GC) ช่วยในการเก็บกวาดสิ่งที่ไม่ได้ใช้แล้วให้เอง ทีนี้ถ้าเราส่ง pointer กลับไป แล้วมีตัวแปร pointer อื่นๆเก็บค่า address ของหน่วยความจำที่เราตอบกลับไป มันก็จะถือว่ามีคนใช้อยู่ GC ก็ไม่มาเก็บ มันจะมาเก็บเองเมื่อไม่มีใครอ้างอิงถึงหน่วยความจำนั้นๆแล้ว (ป.ล. Go compiler อาจจะเลือกเปลี่ยนโค้ดการเรียก NewUser เป็น inline ทำให้การ allocate memory ยังอยู่บน stack ไม่ออกไป heap ให้เองอัตโนมัต ก็เป็นได้ในบางกรณี)

ส่วน Rust นั้นไม่มี GC ใช้กลไก Ownership ในการจัดการหน่วยความจำโดยจะคืนก็ต่อเมื่อ scope ของตัวแปรที่เป็น owner จบการทำงาน และ ถ้าไม่อยากให้มันคืน สิ่งเดียวที่ทำได้คือ move ownership ให้ตัวแปรอื่นที่มีชีวิตยืนยาว (outlive) มากกว่า หรือ copy ค่าออกไปนั่นเอง

ดังนั้นท่าที่ควรจะเป็นของ Rust คือแบบนี้

```rust
struct User {
    name: String,
    email: String
}

impl User {
    fn new(name: String, email: String) -> User {
        return User{name: name, email: email}
    }
}

fn main() {
    let u = User::new(String::from("me"), String::from("me@email.com"));
    println!("{}, {}", u.name, u.email);
}
```

จะเห็นว่าฟังก์ชัน new เปลี่ยน return type เป็น User ไม่ใช่ reference ของ User ทำให้การ return ค่าจะไม่ได้การส่ง reference (borrow) ของตัวแปรที่จะถูกคืนค่าหลังจากจบ scope แต่เป็นการ move ownership ออกไปทำให้จบ scope ก็ไม่ถูกทำลาย หรือว่าไม่ผิดกฎ ตัวแปร u ใน main จะได้ ownership ไปแทน ซึ่งมีอายุยืนยาวกว่าตัวแปรในฟังก์ชัน new นั่นเอง (ป.ล. ถ้า type implement trait Copy ก็จะเป็นการ copy แต่ถ้า type ไม่ได้ implement trait Copy ก็จะเป็นการ move นั่นเอง)

แล้วนี่ก็คือความต่างของ Go กับ Rust ในการจัดการหน่วยความจำ และ วิธีเขียนโค้ดให้ตอบกลับ pointer หรือ reference ซึ่ง Go ตอบกลับ pointer ได้เลยส่วน Rust ตอบกลับ reference ของตัวแปรที่สร้างใน local ไม่ได้เพราะจะถูกทำลายก่อนต้องเป็นการ move ownership หรือ copy เท่านั้น แต่ถ้า Rust ตอบกลับเป็น reference นั่นคือตอบ reference ของ parameter ที่เข้ามาเท่านั้น เหมือนตัวอย่างโค้ดของโพสต์เมื่อวาน http://weerasak.dev/posts/2023/03/04/understand-rust-lifetime-annotation/
