---
title: "[Rust] lifetime annotation มีไว้ช่วย compiler เช็คความถูกต้อง reference ตอน compile time"
date: 2023-03-04T06:50:01+07:00
draft: fase
---

Rust พยายามไม่ใช้เกิดการใช้งาน reference ที่ผิดๆ เช่นอ้างอิงถึงตัวแปรที่ถูกคืนหน่วยความจำไปแล้ว (Dangling Reference)

ทีนี้มีรูปแบบการใช้งานบางรูปแบบที่ compiler ไม่สามารถตรวจสอบได้โดยตรงว่ามีการใช้งาน reference ของตัวแปรที่คืนหน่วยความจำไปแล้วหรือไม่ ทำให้ต้องเราต้องใช้ lifetime annotation ช่วยเพิ่มข้อมูลให้กับ compiler รู้ได้ว่า reference ที่เราใช้งานนั้นมีช่วงชีวิต lifetime อยู่แค่ไหน

โพสต์นี้จะอธิบายสองรูปแบบคือ

- การส่งค่า reference เข้าไปเป็น parameter ของฟังก์ชัน และ return ค่า reference กลับมา
- การเก็บค่า reference เป็น field หนึ่งของ struct

<!--more-->

ตัวอย่างเคสของการส่ง reference ไปให้กับฟังก์ชัน แล้ว ฟังก์ชันตอบกับ reference กลับคืนมา

```rust
fn ref_of_x(x: &str) -> &str {
        return x;
}

fn main() {
    let string1 = String::from("abcd");

    let result = ref_of_x(&string1);
    println!("{}", result);
}
```

โค้ดนี้ compile ผ่านเพราะเพราะเรามีแค่ 1 parameter คือ `x` ที่เป็น reference ของ `str` แล้วก็ตอบกับ `x` กลับไป

ใน `main` เราก็เรียก `ref_of_x` แล้วเก็บค่าที่ตอบกลับมาใน `result`

ซึ่งจะเห็นว่าตัวแปร `string1` และ `result` อยู่ใน scope เดียวกันจะถูกคืนค่าหน่วยความจำหลังจากจบการทำงาน block เดียวกันไม่มีโอกาสที่จะถูกเอาไปใช้ผิดๆหลังจากนั้น จึงทำให้ compiler compile ผ่าน เพราะวิเคราะห์แล้วไม่มีทางที่จะเกิด Dangling reference เกิดขึ้น

ทีนี้ถ้าโค้ดเราแบบนี้แทน

```rust
fn ref_of_x(x: &str) -> &str {
        return x;
}

fn main() {
    let result;
    {
        let string1 = String::from("abcd");
        result = ref_of_x(&string1);
    }
    println!("{}", result);
}
```

compiler จะ compile error แบบนี้

```txt
error[E0597]: `string1` does not live long enough
  --> src/main.rs:9:27
   |
9  |         result = ref_of_x(&string1);
   |                           ^^^^^^^^ borrowed value does not live long enough
10 |     }
   |     - `string1` dropped here while still borrowed
11 |     println!("{}", result);
   |                    ------ borrow later used here

For more information about this error, try `rustc --explain E0597`.
error: could not compile `playground` due to previous error
```

เพราะว่า `string1` ถูกคืนหน่วยความจำไปก่อน (drop) เนื่องจากว่าถูกประกาศตัวแปรไว้ใน block {} อีกอันที่มีช่วงชีวิตที่สั้นกว่า `result` ที่เป็น reference หา `string1`

จากโค้ดตัวอย่างเมื่อกี้ไม่เห็นมี lifetime annotation เลย เพราะ Rust มีกด lifetime elission คือถ้าไม่ใส่จะเติมให้เองได้ถ้าเข้าเงื่อนไขของ lifetime elission เช่นตัวอย่างเมื่อกี้ Rust จะมองว่าเหมือนมี lifetime annotation ให้แบบนี้

```rust
fn ref_of_x<'a>(x: &'a str) -> &'a str {
        return x;
}

fn main() {
    let string1 = String::from("abcd");

    let result = ref_of_x(&string1);
    println!("{}", result);
}
```

ซึ่ง lifetime annotation เราต้องลิสต์เป็น type parameter แบบเดียวกับ generic type parameter อื่นๆ แต่ว่าชื่อ parameter ต้องขึ้นต้นด้วย `'` แล้วตาม convention จะใช้อักษรตัวเลขไล่จาก `a` เป็นต้นไป

ตอนเอามากำหนดให้ตัวแปร reference ก็ใส่หลัง `&` ก่อนแล้วค่อยเว้นวรรค ตามด้วยชื่อ type

ทีนี้มาดูเคสที่ไม่ตรงเงื่อนไข lifetime elission แล้วต้องใส่เองกันเช่น

```rust
fn longest(x: &str, y: &str) -> &str {
    if x.len() > y.len() {
        x
    } else {
        y
    }
}

fn main() {
    let string1 = String::from("long string is long");

    {
        let string2 = String::from("xyz");
        let result = longest(string1.as_str(), string2.as_str());
        println!("The longest string is {}", result);
    }
}
```

ถ้าโค้ดแบบนี้ Rust จะงงว่า lifetime ของ return type จะเป็นอะไร เพราะมีทั้ง `x` และ `y` แล้วจะ compile error มาแบบนี้

```txt
error[E0106]: missing lifetime specifier
 --> src/main.rs:1:33
  |
1 | fn longest(x: &str, y: &str) -> &str {
  |               ----     ----     ^ expected named lifetime parameter
  |
  = help: this function's return type contains a borrowed value, but the signature does not say whether it is borrowed from `x` or `y`
help: consider introducing a named lifetime parameter
  |
1 | fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
  |           ++++     ++          ++          ++

For more information about this error, try `rustc --explain E0106`.
error: could not compile `playground` due to previous error
```

ซึ่งก็ guide มาให้เรียบร้อย เราก็เลยต้องแก้โค้ดเป็นแบบนี้

```rust
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() > y.len() {
        x
    } else {
        y
    }
}

fn main() {
    let string1 = String::from("long string is long");

    {
        let string2 = String::from("xyz");
        let result = longest(string1.as_str(), string2.as_str());
        println!("The longest string is {}", result);
    }
}
```

ทีนี้จุดที่น่าสงสัยคือ ทำไมถึงใช้ `'a` กับทั้ง `x` และ `y`

จริงๆแล้วเราใช้คนละตัวก็ได้ แต่พอโค้ดเรามีโอกาสตอบกลับได้ทั้ง `x` และ `y` ทำให้ Rust compiler ไม่รู้ว่า return type เราจะเป็น `x` หรือ `y` กันแน่

Rust เลยให้เราใช้ `'a` ได้กับทั้ง `x` และ `y`

แม้ว่า syntax จะเป็นแบบนี้ ไม่ได้หมายความว่า `'a` จะเป็นข้อมูล lifetime ของทั้ง `x` และ `y` ได้ในเวลาเดียวกัน

แต่ Rust compiler จะเลือกแค่ lifetime ที่สั้นที่สุด อายุน้อยที่สุดระหว่าง `x` และ `y`

เช่นจากโค้ดการใช้งานใน main เมื่อกี้ เราส่ง ref ของ `string1` เข้าไปที่ `x`, ส่ง ref `string2` ไปที `y` ซึ่งช่วงชีวิตของ `string2` สั้นกว่าเพราะอยู่ใน scope ที่จบก่อน `string1` ทำให้ `'a` ก็จะแทนด้วย lifetime ของ `y` นั่นเอง ทำให้ lifetime ของค่าที่ตอบกลับมาของ `longest` ที่เก็บใน `result` ก็จะมีช่วงชีวิตเดียวกันกับ `string2` ไปด้วย

ดังนั้นเราจะเขียนโค้ดแบบ จะ compile ไม่ผ่าน

```rust
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() > y.len() {
        x
    } else {
        y
    }
}

fn main() {
    let string1 = String::from("long string is long");
    let result;
    {
        let string2 = String::from("xyz");
        result = longest(string1.as_str(), string2.as_str());
    }

    println!("The longest string is {}", result);
}
```

เพราะเรามีใช้งาน `result` หลังจากค่าที่ return นั้นจบชีวิต (drop) ไปแล้วนั่นเอง แม้ว่าฟังก์ชันมีโอกาสตอบกลับค่า `x` ตอน runtime แต่ว่า compiler ไม่สามารถรู้ได้ตอน compile time จึงทำให้ compile error นั่นเอง

ต่อไปคือเคสที่มี reference ใน data structure อื่นๆเช่น struct เราก็ต้องกำหนด lifetime annotation ให้กับ field ที่เป็น reference ด้วยแบบนี้

```rust
struct ImportantExcerpt<'a> {
    part: &'a str,
}

fn main() {
    let novel = String::from("Call me Ishmael. Some years ago...");
    let first_sentence = novel.split('.').next().expect("Could not find a '.'");
    let i = ImportantExcerpt {
        part: first_sentence,
    };
    println!("{}", i.part);
}
```

เพราะว่า Rust ไม่อยากปล่อยให้มีกรณีที่ช่วงชีวิตของตัวแปร struct เองนั้นยาวนานกว่าช่วงชีวิตของ reference ที่ struct เก็บไว้อยู่เพราะมันทำให้เกิด Dangling ได้นั่นเอง

เช่นแบบนี้ที่จะ compile ไม่ผ่าน เพราะ field part ไม่มีช่วงชีวิตที่ยาวพออีกแล้ว

```rust
struct ImportantExcerpt<'a> {
    part: &'a str,
}

fn main() {
    let i;
    {
        let novel = String::from("Call me Ishmael. Some years ago...");
        let first_sentence = novel.split('.').next().expect("Could not find a '.'");
        i = ImportantExcerpt {
            part: first_sentence,
        };
    }
    println!("{}", i.part);
}
```

สรุปก็คือ Rust compiler จะรู้ช่วงชีวิต lifetime ของตัวแปรต่างๆตั้งแต่ตอน compile time ส่วน lifetime annotation มีก็เพื่อช่วยกำหนดความสัมพันธ์ของช่วงชีวิต ว่า reference ที่ตอบกลับไปจาก function จะมีช่วงชีวิตสั้นยาวแค่ไหนโดยดูจาก reference parameter ที่ส่งเข้ามา

ส่วนของ Struct field ก็ดูจากค่า reference ที่ assign ให้กับ field นั้นๆเช่นกัน

Ref:
อ่านเพิ่มเติมได้จากที่นี่ https://doc.rust-lang.org/book/ch10-03-lifetime-syntax.html ซึ่งโค้ดตัวอย่างในโพสต์นี้ก็เอามาจากลิ้งนี้เช่นกัน
