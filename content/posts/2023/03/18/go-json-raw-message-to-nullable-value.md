---
title: "[Go] สร้างฟังก์ชัน unmarshal json.RawMessage แบบ generic type และเช็คค่า null"
date: 2023-03-18T16:16:00+07:00
draft: false
---

มีเคสต้องกำหนดค่า request body field นึง ให้เป็น type `json.RawMessage` ไปก่อนแล้วค่อยไปเช็คเงื่อนไขอีกทีว่าจะ unmarshaling ไปใส่ตัวแปร type ไหน แต่ก็อยากรู้ด้วยว่าค่าใน `json.RawMessage` เป็น `null` value ของ JSON หรือไม่ เลยได้ลองสร้างฟังก์ชันที่เป็น generic รองรับการ parse ค่าจาก `json.RawMessage` เป็น type ใดๆดูพร้อมกับ return flag เพื่อบอกว่าเป็น `null` หรือไม่

<!--more-->

โค้ดเป็นแบบนี้

```go
func rawMessageToValue[T any](raw json.RawMessage) (T, bool, error) {
	var v T
	if string(raw) == `null` {
		return v, true, nil
	}
	err := json.Unmarshal(raw, &v)
	return v, false, err
}
```

- มี type parameter T มี constraint เป็น any คือเป็น type อะไรก็ได้
- รับค่า `raw` เป็น `json.RawMessage`
- และ return ออกไป 3 ค่าซึ่งก็คือ `T` เองหลังจาก unmarshaling แล้ว, ค่า `bool` เพื่อบอกว่าเป็น `null` หรือไม่ และ `error` ถ้าเกิดมี error ตอน unmarshaling
- ในโค้ดก็แค่ประกาศตัวแปร `v` type `T` ตาม type parameter
- เช็คว่า `raw` มีค่าเป็น `null` หรือเปล่าโดยแปลงเป็น `string` ก่อนแล้วค่อยเทียบกับ `string` ค่า `null` ถ้าใช่ก็ return `v` ออกไปเลยแบบไม่ต้องทำอะไร พร้อมค่าที่สองเป็น `true` และ error เป็น `nil`
- ถ้าไม่ใช่ก็แค่เรียก `json.Unmarshal(raw, &v)` แล้ว `return v, false, err` เพื่อบอกว่า Unmarshal ได้อะไรและบอกว่าไม่เป็น `null` นะ

เวลาใช้งานเราก็สามารถเช็คจากค่าที่สองได้ว่าเป็น `null` หรือเปล่า ถ้าไม่ใช่ก็เช็คว่า `error` หรือเปล่า ถ้าไม่ก็เอา ค่าแรกที่ return กลับไปใช้งานได้เลย เช่น

```go
func main() {
	raw := json.RawMessage("null")
	n, isNull, err := rawMessageToValue[int](raw)
	if isNull {
		fmt.Println(isNull)
	} else if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(n)
	}

        // Output:
        // true
}
```

```go
func main() {
	raw := json.RawMessage("10")
	n, isNull, err := rawMessageToValue[int](raw)
	if isNull {
		fmt.Println(isNull) // false
	} else if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(n)
	}

        // Output:
        // 10
}
```

```go
func main() {
	raw := json.RawMessage(`"Hello, World"`)
	str, isNull, err := rawMessageToValue[string](raw)
	if isNull {
		fmt.Println(isNull) // false
	} else if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(str)
	}

        // Output:
	// Hello, World
}
```

```go
func main() {
	raw := json.RawMessage(`"2023-03-18T16:15:00.000+07:00"`)
	t, isNull, err := rawMessageToValue[time.Time](raw)
	if isNull {
		fmt.Println(isNull) // false
	} else if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(t)
	}

        // Output:
	// 2023-03-18 16:15:00 +0700 +0700
}
```
