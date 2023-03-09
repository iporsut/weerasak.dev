---
title: "[Rust] สรุปความต่างของ Box<T>, Rc<T> และ RefCell<T>"
date: 2023-03-06T11:33:00+07:00
draft: false
---

สรุปความต่างของ Box, Rc และ Refcell ของ Rust ที่เป็น smart pointer เอาไว้หน่อย

<!--more-->

ก่อนจะดูว่าต่างกันยังไง สรุปเรื่องกฎของ ownership และ borrow/reference checker เอาไว้หน่อย

## กฎของ ownership

- แต่ละค่าของตัวแปรใน Rust ต้องมี owner
- ต้องมีแค่ 1 owner ณ เวลาใดๆ (และ owner move ไปให้ owner อื่นได้)
- เมื่อจบการทำงาน scope ของ owner ค่านั้นจะถูก drop และคืนหน่วยความจำ

## กฎของ borrow/reference checker

- ณ เวลาใดๆ จะมีได้แค่ 1 mutable borrow และมีหลายๆ immbutable borrow ได้
- reference ต้อง valid ไม่สามารถ reference หาค่าที่ lifetime จบไปแล้วได้

กฎเหล่านี้จะถูกเช็คตอน compile time แต่ smart pointer อย่าง Rc, RefCell จะทำให้เราใช้งานโดยหลีกเลี่ยงกฎเหล่านี้ตอน compile time ได้

## ความต่างของ Box, Rc, RefCell

- Box ยึดตาม ownership/borrow ตอน compile time ตามปกติ สิ่งที่ Box ช่วยคือทำให้เรา เก็บ value ใน heap ได้
- Rc ทำให้เราแชร์ค่าเดียวกันที่อยู่ใน heap และมีหลายๆ ownership ได้ซึ่งค่าจะยังไม่ถูก drop จนกว่าทุก ownership นั้นจะ drop จนหมด แต่ว่ายังยึดกฎ borrow และ reference checker ตอน compile time เหมือนเดิม
- RefCell นั้นยึดกฎ ownership ตอน compile time เหมือนเดิม แต่ว่า ยืดหยุ่นให้สามารถ mutate ค่าข้างใน RefCel ได้ และเช็คกฎ borrow/reference ตอน runtime แม้ว่าตัวแปรของ RefCell เองจะ immutable ก็ตาม

ตัวอย่างโค้ด

```rust
// Box
fn main() {
    let box_v = Box::new(10); // immutable
    let mut box_mv = Box::new(10); // mutable
    println!("{}", box_v);
    println!("{}", box_mv);
    *box_mv = 20;
    println!("{}", box_mv);
}
```

```rust
// Rc
use std::rc::Rc;

fn main() {
    let rc_1 = Rc::new(10); // owner 1
    let rc_2 = Rc::clone(&rc_1) // owner 2
    println!("{}", rc_1);
    println!("{}", rc_2);
    println!("{}", Rc::strong_count(&rc1)) // 2
}
```

Box นั้นใช้ `*box_mv = 20;` แก้ไขค่า mutable ข้างในได้โดยตรงแต่ Rc ทำไม่ได้เพราะจะทำได้ต้อง implements `DerefMut` trait

เหตุผลคือ Rc เปิดโอกาสให้มี owner หลายคนที่ internal นั้นมี ref ที่ชี้ไปที่ข้อมูลเดียวกัน ถ้าเกิดทำให้ DerefMut ได้จะทำให้เบรกกฎ borrow checker

อย่างไรก็ตาม Rc นั้นมี function ที่ให้เรา get mutable reference ของค่าด้านในได้แต่จะได้กลับออกมาเป็น `Option<&mut T>` เพราะมีเงื่อนไขว่า ถ้ามี owner หลายค่าจะได้ค่า None ออกมาเพื่อไม่ให้มีการเปลี่ยนพร้อมกันขัดกับกฎ borrow checker

แต่ถ้ามีแค่ 1 owner ก็จะได้ &mut T ออกมา แก้ไขค่าได้ เช่น โค้ดแบบนี้ get_mut จะได้ None

```rust
// Rc
use std::rc::Rc;

fn main() {
    let mut rc_1 = Rc::new(10); // owner 1
    let rc_2 = Rc::clone(&rc_1); // owner 2
    println!("{:?}", Rc::get_mut(&mut rc_1));  // None
    println!("{}", rc_2); // 2
    println!("{}", Rc::strong_count(&rc_1)) // 2
}
```

แต่แบบนี้จะได้ `Some(10)` แล้วจะแก้ไขค่าได้

```rust
// Rc
use std::rc::Rc;

fn main() {
    let mut rc_1 = Rc::new(10); // owner 1
    let ref_to_inside: &mut i32 = Rc::get_mut(&mut rc_1).unwrap();
    *ref_to_inside = 20;
    println!("{}", rc_1);
}
```

ส่วน RefCell ทำให้กฎเรื่อง borrow/reference ของ borrow checker เช็คตอน runtime ได้ นั่นคือ เราสามารถมีตัวแปร immutable ของ RefCell แต่สามารถแก้ไขค่าด้านใน mutate ค่าที่ RefCell reference ไว้อยู่ได้ ตัวอย่างเช่น

```rust
use std::cell::RefCell;

fn main() {
    let ref_cell_1 = RefCell::new(10);
    {
        let mut mutable_borrow = ref_cell_1.borrow_mut();
        *mutable_borrow = 20;
    }
    let immutable_borrow = ref_cell_1.borrow();
    println!("{}", immutable_borrow); // 20
}
```

เราใช้ method `borrow_mut()` เพื่อให้ได้ type ใหม่ที่ implement DerefMut ซึ่งเราจะเปลี่ยนค่าได้ และใช้ `borrow()` ธรรมดาเพื่อให้ได้ Deref เฉยๆที่เปลี่ยนค่าไม่ได้

อย่างไรก็ตามถ้าเราเรียก `borrow_mut()` โดยที่มีการเรียก `borrow()` เพื่อยืมไปก่อนหน้าและยังมีช่วงชีวิตอยู่จะทำให้ตอนเรียก `borrow_mut()` เกิด panic ตอน runtime เกิดขึ้น เช่นโค้ดด้านบนถ้าเราเอา block `{}` ออกเป็นแบบนี้

```rust
use std::cell::RefCell;

fn main() {
    let ref_cell_1 = RefCell::new(10);
    let mut mutable_borrow = ref_cell_1.borrow_mut();
    *mutable_borrow = 20;
    let immutable_borrow = ref_cell_1.borrow();
    println!("{}", immutable_borrow); // 20
}
```

จะได้ panic message ประมาณนี้

```txt
thread 'main' panicked at 'already mutably borrowed: BorrowError', src/main.rs:7:39
note: run with `RUST_BACKTRACE=1` environment variable to display a backtrace
```

ถ้าเราไม่อยากให้เกิด panic ให้ไปใช้ method `try_borrow` และ `try_borrow_mut` แทนซึ่งจะตอบกลับมาเป็น type Result ที่ wrap error ไว้ด้วย แล้วเอามาเช็คอีกทีได้

สุดท้าย เราสามารถเอา Rc และ RefCell มาประกอบกันเพื่อให้ได้ type ที่ผ่อนปรนกฎเรื่อง ownership และ borrowing จากตอน compile ไปเป็นตอน runtime ได้ ตัวอย่างเช่น

```rust
use std::cell::RefCell;
use std::rc::Rc;

fn main() {
    let rc_ref_cell_1 = Rc::new(RefCell::new(10));
    let rc_ref_cell_2 = Rc::clone(&rc_ref_cell_1); // เพิ่ม owner
    println!("{:?}", rc_ref_cell_1);
    println!("{:?}", rc_ref_cell_2);

    {
        // borrow ค่าจาก owner 2 เรียก method borrow_mut  ของ RefCell ได้เลยเพราะ Rc implements Deref
        let mut mut_v = rc_ref_cell_2.borrow_mut();
        *mut_v = 20; // เปลี่ยนค่าเป็น 20 แล้วทั้ง 2 owner เห็นค่าเปลี่ยนไปตาม
    }

    println!("{:?}", rc_ref_cell_1);
}

// Output:
// RefCell { value: 10 }
// RefCell { value: 10 }
// RefCell { value: 20 }
// RefCell { value: 20 }
```

สรุปสุดท้ายสั้นๆ

- Box เพิ่มความสามารถในการเก็บตัวแปรใน heap
- Rc เก็บใน heap แต่เพิ่ม owner หลายๆ owner ได้
- RefCell ทำให้ mutable ค่าที่อยู่ข้างใน RefCell แม้ว่า RefCell จะเป็น immutable ได้
- Rc+RefCell ก็ทำให้ทั้งมีหลายๆ owner และ แก้ไขค่าได้ตอน runtime ของค่าที่เก็บอยู่ใน heap ได้
