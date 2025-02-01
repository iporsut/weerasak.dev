---
title: "Basic Rust: Pattern Matching"
date: 2025-02-01T10:17:00+07:00
draft: false
---

Rust มีกลไก pattern matching ที่ให้เราสามารถเปรียบเทียบ pattern ของค่าต่าง ๆ แล้วให้ทำงานอะไรบางอย่าง ถ้าค่าที่เปรียบเทียบกันนั้น match กัน โดย Rust ใช้ keyword ชื่อว่า `match` ในการเขียน pattern matching

<!--more-->

ตัวอย่างโค้ด (จาก Rust Book https://doc.rust-lang.org/book/ch06-02-match.html)

```rust
enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter,
}

fn value_in_cents(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => 1,
        Coin::Nickel => 5,
        Coin::Dime => 10,
        Coin::Quarter => 25,
    }
}
```

ฟังก์ชัน `value_in_cents` ต้องการ เช็คค่า coin ว่าตรงกับ enum value ตัวไหน โดยใช้ `match coin` แล้วตามด้วยลิสต์ในแต่ละเคสโดยแต่ละเคสคือ enum value ที่เป็นไปได้ จากนั้นก็ ระบุค่า value ในหน่วย cents ของ Coin แต่ละแบบ

จริง ๆ โค้ดการทำงานในแต่ละเคสที่ match จากเขียนเป็น block ก็ได้ แต่ค่าที่ return ของ block นั้น ๆ ก็ต้องเป็น type เดียวกันกับทุก ๆ เคส เช่น

```rust
fn value_in_cents(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => {
            println!("Lucky penny!");
            1
        }
        Coin::Nickel => 5,
        Coin::Dime => 10,
        Coin::Quarter => 25,
    }
}
```

## Patterns with bind values (deconstruction)

เราสามารถเขียน matching value โดยให้ binding ค่าออกจาก value มาใส่ในตัวแปรได้ด้วย ตัวอย่างเช่น

```rust
#[derive(Debug)]
enum Message {
    Text(String),
    Number(f64),
}

fn main() {
    let msg = Message::Text(String::from("Hello"));
    
    match msg {
        Message::Text(ref s) => println!("Text message {:?}", s),
        Message::Number(n) => println!("Number message {:?}", n)
    };
    
    let msg = Message::Number(10.55);
    match msg {
        Message::Text(s) => println!("Text message {:?}", s),
        Message::Number(n) => println!("Number message {:?}", n)
    };
}
```

เรื่อง move semantic ก็ใช้เหมือนกันกับการ pass parameter เช่นตัวอย่างบน ถ้าเราเอา msg ที่ match Test(s) ไปแล้วมาใช้อีก ก็จะ compile error เพราะมันคือการ move ownership ไปแล้ว

```
   |
11 |         Message::Text(s) => println!("Text message {:?}", s),
   |                       - value partially moved here
...
15 |     println!("{:?}", msg);
   |                      ^^^ value borrowed here after partial move
   |
```

ถ้าเราต้องการที่จะ binding โดยต้องการ match แต่จะใช้เป็น reference ต้องใช้ keyword `ref` หน้าชื่อตัวแปรที่เราต้องการ binding ด้วย แบบตัวอย่างนี้ ถึงจะไม่เกิดการ move แต่จะเป็น borrow แทน

```rust
#[derive(Debug)]
enum Message {
    Text(String),
    Number(f64),
}

fn main() {
    let msg = Message::Text(String::from("Hello"));
    
    match msg {
        Message::Text(ref s) => println!("Text message {:?}", s),
        Message::Number(n) => println!("Number message {:?}", n)
    };
    
    println!("{:?}", msg);
    
    let msg = Message::Number(10.55);
    match msg {
        Message::Text(s) => println!("Text message {:?}", s),
        Message::Number(n) => println!("Number message {:?}", n)
    };
}
```

## Matches Are Exhaustive

คือ matching ใน Rust นั้นต้องระบุให้ cover ทุก ๆ เคสที่เป็นไปได้ เช่นถ้าเราเขียนโค้ด matching ค่าของ Option type แบบนี้จะ compile ไม่ผ่าน เพราะไม่ครอบคลุมเคส None

```rust
fn plus_one(x: Option<i32>) -> Option<i32> {
    match x {
        Some(i) => Some(i + 1),
    }
}
```

compile error

```
error[E0004]: non-exhaustive patterns: `None` not covered
 --> src/main.rs:3:15
  |
3 |         match x {
  |               ^ pattern `None` not covered
  |
```

เราสามารถ binding ค่าที่เหลือที่นอกจากที่เราต้องการแบบเฉพาะเจาะจงได้ หรือเลือกที่จะใช้ underscore (_) เพื่อ binding ค่าอื่น ๆ แบบที่เราไม่สนใจจะใช้ค่าที่ match นั้น ๆ ได้ เช่น

```rust
fn plus_one(x: Option<i32>) -> Option<i32> {
    match x {
        Some(i) => Some(i + 1),
        _ => None,
    }
}
```

จาก match keyword แล้วเรายังมี `if let` เพื่อเช็ค matching แบบสนใจแค่ pattern เดียวเท่านั้น ถ้า match ก็ทำงานตาม body ของ if หรือไม่ match ก็ผ่านไป เช่น

```rust
let config_max = Some(3u8);
if let Some(max) = config_max {
    println!("The maximum is configured to be {max}");
}
```

หรือ

```rust
    let mut count = 0;
    if let Coin::Quarter(state) = coin {
        println!("State quarter from {state:?}!");
    } else {
        count += 1;
    }
```

ซึ่งช่วยให้เราไม่ต้องเขียน cover ทุก ๆ เคสหรือใช้ `_`

ใน Rust นั้นมี std type Option และ Result ช่วยในการจัดการ error ซึ่ง 2 types นี้เป็น enum ที่มี value ข้างใน แล้วเราเลยจะเห็นการใช้ pattern matching เยอะ ๆ ในการจัดการกับค่าของ 2 types ในโค้ดที่เราเขียน