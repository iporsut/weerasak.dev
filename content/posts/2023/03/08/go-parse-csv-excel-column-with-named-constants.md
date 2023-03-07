---
title: "[Go] Parse CSV/Excel columns ด้วย named constant และ map แทนที่จะใช้ index โดยตรง"
date: 2023-03-08T06:45:00+07:00
draft: false
---

ช่วงนี้มี usecase ที่ต้อง parse ข้อมูลจาก Excel แต่ว่าบ้างครั้งลำดับของ column ก็มีเปลี่ยนไป มีเพิ่ม มีลดกันบ้าง เลยได้ลองคิดท่าใหม่เพื่อให้โค้ดในการ parse ไม่ต้องแก้เยอะเพราะลำดับของ column เปลี่ยน ซึ่งก็คือใช้ constant และ map มาช่วยนั่นเอง

<!--more-->

สมมติเรามีข้อมูลใน CSV/Excel แบบนี้

| Name     | Address  | BirthDay   |
| -------- | -------- | ---------- |
| John Doe | Bangkok  | 20/02/2000 |
| Jane Doe | New York | 12/03/2001 |
| Jame Doe | London   | 01/02/2002 |

ผมจะละโค้ดตอนอ่านไฟล์ไปนะ สมมติว่าเราอ่านไฟล์มาแล้วเก็บไว้ใน slice ของ slice ของ string แบบนี้ และทำการตัด header row ออกแล้วด้วยนะเหลือแต่ data rows

```go
var rows [][]string = loadData()
```

แล้วเราจะ parse เพื่อให้เก็บใน slice ของ User แบบนี้

```go
type User struct {
        Name string
        Address string
        BirthDay string
}

func parseData(rows [][]string) []*User{
        var users []*User

        for _, row := range rows {
                var user User
                user.Name = row[0]
                user.Address = row[1]
                user.BirthDay = row[2]

                users = append(users, &user)
        }

        return users
}
```

จะเห็นแบบ version แรกนี้เราใช้ index ตรงๆอยู่ เราทำให้ดีกว่านี้ได้โดยใช้ constant ช่วยแบบนี้

```go
type UserParsingIndex int

const (
        UserNameColumn UserParsingIndex = iota
        UserAddressColumn
        UserBirthdayColumn
)
```

จะเห็นว่าเราใช้ iota ช่วยเพื่อไล่เลขจาก 0 และ define type ใหม่สำหรับ parsing index โดยเฉพาะ หลังจากนั้นปรับโค้ดเป็นแบบนี้

```go
func parseData(rows [][]string) []*User{
        var users []*User

        for _, row := range rows {
                var user User
                user.Name = row[UserNameColumn]
                user.Address = row[UserAddressColumn]
                user.BirthDay = row[UserBirthdayColumn]

                users = append(users, &user)
        }

        return users
}
```

แบบนี้เราสามารถปรับ ลำดับโค้ดจากตรง define constant ได้เลยเช่นสลับ Address กับ BirthDay เราก็แค่ปรับตรง const แบบนี้

```go
const (
        UserNameColumn UserParsingIndex = iota
        UserBirthdayColumn
        UserAddressColumn
)
```

โจทย์ต่อไปถ้าช่วงเปลี่ยนผ่านเช่นต้องการแยก column name เป็น FirstName กับ LastName ล่ะ แต่ว่า Struct ยังเก็บเป็น Name เหมือนเดิม จะทำยังไงให้รองรับได้ทั้งสองแบบ

ดังนั้นแทนที่จะใช้ constant โดยตรง เราสามารถใช้ constant และ map เข้าช่วย แทนที่จะใช้ constant แทน index เราใช้ constant เพื่อ lookup index อีกทีใน map ได้แบบนี้

```go
type UserColumnParse struct {
        row []string
        columnIndexMap map[UserParsingIndex]int
}

func (p *UserColumnParse) get(column UserParsingIndex) (string, bool) {
        if index, ok := columnIndexMap[column]; ok {
                return p.row[index], true
        }

        return "", false
}
```

สร้าง type ใหม่ที่รับ row และ columnIndexMap ที่เป็น map ระหว่าง column constant ไปหา int ที่เป็น index จริงๆ เสร็จแล้วสร้าง method get เพื่อรับ column แล้วไป lookup หา index ที่แท้จริง ถ้ามี ก็เอาไป get จาก row แล้วตอบกับไปพร้อม true เพื่อบอกว่ามีค่านี้อยู่

ทีนี้การสร้าง columnIndexMap ถ้ามาสร้างเองทีละค่าก็จะยุ่งยาก เลยสร้างฟังก์ชันช่วยในการสร้างแบบนี้

```go
func makeColumnIndexMap(columns ...UserParsingIndex) map[UserParsingIndex]int {
        m := map[UserParsingIndex]int{}
	for index, column := range columns {
		m[column] = index
	}

	return m
}
```

จะเห็นว่าฟังก์ชันรับเป็น variadic parameter ซึ่งทำให้เราจัดลำดับของ column ตามลำดับ parameter ได้เลย ถ้าเปลี่ยนก็แค่สร้างเป็นลำดับใหม่

ทีนี้มาดูโค้ดทั้งหมดหลังจากใชั้ constant และ map ช่วยกัน

```go
const (
        UserNameColumn UserParsingIndex = iota
        UserFirstNameColumn
        UserLastNameColumn
        UserBirthdayColumn
        UserAddressColumn
)

type UserColumnParse struct {
        row []string
        columnIndexMap map[UserParsingIndex]int
}

func (p *UserColumnParse) get(column UserParsingIndex) (string, bool) {
        if index, ok := columnIndexMap[column]; ok {
                return p.row[index], true
        }

        return "", false
}

func parseData(rows [][]string, columnIndexMap map[UserParsingIndex]int) []*User{
        var users []*User

        for _, row := range rows {
                columnParse := UserColumnParse{
                        row: row,
                        columnIndexMap: columnIndexMap,
                }

                var user User
                if name, ok := columnParse.get(UserNameColumn); ok {
                        user.Name = name
                }
                if firstName, ok := columnParse.get(UserFirstNameColumn); ok {
                        user.Name = firstName
                }
                if lastName, ok := columnParse.get(UserLastNameColumn); ok {
                        user.Name += " " + lastName
                }

                if address, ok := columnParse.get(UserAddressColumn); ok {
                        user.Address = address
                }
                if birthDay, ok := columnParse.get(UserBirthdayColumn); ok {
                        user.BirthDay = birthDay
                }

                users = append(users, &user)
        }

        return users
}

// ตอนเรียกใช้ parseData โดยใช้ mapping แบบเดิม
parseData(rows, makeColumnIndexMap(
        UserNameColumn,
        UserAddressColumn,
        UserBirthdayColumn
))

// ตอนเรียกใช้ parseData โดยใช้ mapping แบบใหม่
parseData(rows, makeColumnIndexMap(
        UserFirstNameColumn,
        UserLastNameColumn
        UserAddressColumn,
        UserBirthdayColumn
))
```

จะเห็นว่าตรงประกาศ constant จะไม่ได้มีผลกับลำดับของ column อีกต่อไปแล้ว จะเป็นแค่การประกาศชื่อ column ที่เป็นไปได้แทน แล้วเราค่อยไป mapping ชื่อกับ index จริงๆอีกทีผ่านฟังก์ชัน makeColumnIndexMap โดยเรียงลำดับตาม parameter ที่ส่งให้ฟังก์ชันนั่นเอง
