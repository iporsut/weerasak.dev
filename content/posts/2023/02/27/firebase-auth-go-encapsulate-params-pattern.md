---
title: "[Go] pattern ในการ encapsulate payload params ของ Firebase Auth SDK"
date: 2023-02-27T17:15:13+07:00
draft: false
---

วันนี้ใช้งาน Firebase Auth Go SDK แล้วเจอ pattern น่าสนใจในการ encapsulate request payload params เลยมาเขียนสรุปเอาไว้หน่อย

<!--more-->

Firebase Auth คือ service นึงของ Firebase ที่เอาไว้สำหรับจัดการ Authentication โดยที่เราไม่ต้องเขียนระบบ Authen เอง

ทีนี้ตอนใช้งานก็จะมี library ของแต่ละภาษาครอบการเรียก HTTP ตรงๆอยู่ อย่าง Go ก็มี Firebase Admin SDK สำหรับ Auth service อยู่

จุดที่น่าสนใจคือการออกแบบ Go SDK ของ Firebase มีท่าที่ไม่ค่อยเหมือน SDK ของ API อื่นๆที่เคยเจอมา ตัวอย่างเช่นสมมติเรามี API สำหรับสร้าง User เราคงจะเจอ pattern นี้บ่อยๆ

```go
type User struct {
        Email    `json:"email"`
        Password `json:"password"`
}

func CreateUser(ctx context.Context, user *User) (*User, error) {
}
```

คือ payload ก็เป็น struct ที่มี field และ json tag ช่วยในการ mapping ระหว่าง Field กับ json text ที่จะถูกส่งไปหา API

แต่สำหรับ Firebase Auth SDK เลือกที่จะสร้าง type ใหม่สำหรับแต่ละ API action เช่นการสร้าง User ของ Firebase Auth SDK โค้ดจะเป็นแบบนี้

```go
createUserParam := &auth.UserToCreate{}
createUserParam.Email(email)
createUserParam.Password(password)

userRecord, err := client.CreateUser(ctx, createUserParam)
if err != nil {
        return nil, err
}
```

แทนที่จะรับ payload ที่เป็น type struct ตรงๆกลับสร้าง struct ที่ไม่สามารถระบุ field ได้ตรงๆ แต่ให้เรียก setter method เพื่อกำหนดค่าแต่ละ field แทน

ถ้าเราไปแกะโค้ดของ SDK ดูจะพวกว่าประกาศ struct UserToCreate ไว้แบบนี้

```go
// UserToCreate is the parameter struct for the CreateUser function.
type UserToCreate struct {
	params map[string]interface{}
}
```

ซึ่งมีแค่ field เดียวคือ params ที่เป็น type `map[string]interface{}` (any ใน Go 1.18+)

ส่วน setter method นั้นเป็นแบบนี้

```go
// Email setter.
func (u *UserToCreate) Email(email string) *UserToCreate {
	return u.set("email", email)
}

// Password setter.
func (u *UserToCreate) Password(pw string) *UserToCreate {
	return u.set("password", pw)
}

func (u *UserToCreate) set(key string, value interface{}) *UserToCreate {
	if u.params == nil {
		u.params = make(map[string]interface{})
	}
	u.params[key] = value
	return u
}
```

ซึ่งก็คือเอาค่าไป set ใน map params อีกทีนั่นเอง แล้วตอบกลับมาเป็น pointer ของ struct ตัวเดิมเพื่อให้เขียนแบบ chain ได้แบบนี้

```go
createUserParam := &auth.UserToCreate{}
createUserParam.Email(email).Password(password)

userRecord, err := client.CreateUser(ctx, createUserParam)
if err != nil {
        return nil, err
}
```

ซึ่งก็เป็นอีกท่านึงที่น่าสนใจเพราะพอทำเป็น Setter ก็ทำให้เปิดโอกาสในการแก้ไขโครงสร้างภายในของ type ได้โดยไม่กระทบกับคนเรียกใช้งานหลังจากอัพเกรด SDK

ลองดู Package document และโค้ดของ Auth SDK ทั้งหมดได้ที่นี่ https://pkg.go.dev/firebase.google.com/go/v4/auth
