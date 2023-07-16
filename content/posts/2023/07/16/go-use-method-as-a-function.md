---
title: "Go: ใช้ method เป็นค่าของ function"
date: 2023-07-16T10:31:00+07:00
draft: false
---

ใน Go function เป็น value แบบนึง ที่เราสามารถส่งไปให้ function อื่นๆ ได้ และเราสามารถเก็บ function ไว้ในตัวแปรได้ เช่น ที่เราใช้กันบ่อยๆก็คือตอนที่เราเขียน HTTP handler

<!--more-->

```go
func main() {
        http.HandleFunc("/todos", func (w http.ResponseWriter, r *http.Request) {
                r.Write([]byte("Hello, World"))
        })
}
```

เราเขียน anonymous function เพื่อส่งค่าฟังก์ชันไปให้ http.HandleFunc เพื่อเขียนการทำงานเมื่อมี request เข้ามาที่ path `/todos`

จากที่เห็น handler function เราถูกกำหนดว่าต้องมี 2 parameters แต่ถ้าเราต้องการให้ handler เราใช้งานค่าอื่นๆเช่น มีตัวแปรของ database connection ด้วย เราจะส่งเข้าไปให้ handler function ไม่ได้ เราต้องประกาศเป็นตัวแปรไว้ที่อื่นเช่นใน main ก่อนเรียก http.HandleFunc แต่แบบนี้จะทำให้เราทดสอบ handler function เราได้ยาก ดังนั้นท่านึงที่เป็นไปได้ คือใช้ method เป็น function แล้วให้ type ที่ implements method เก็บค่าของ dependencies อื่นๆที่ต้องใช้ใน handler เอาไว้ให้ แบบนี้ (จริงๆทำได้หลายท่าเช่นสร้าง builder function ที่ return handler function อีกที แต่วันนี้อยากยกตัวอย่างการใช้ method เป็นค่าของ function)

```go
type Handler struct {
        db Database
}

func (h *Handler) GetAllTodo(w http.ResponseWriter, r *http.Request) {
        todos := h.db.GetAllTodo()
        json.NewEncoder(w).Encode(todos)
}

func main() {
        db := ConnectDB()
        h := Handler{
                db: db,
        }

        http.HandleFunc("/todos", h.GetAllTodo)
}
```

จะเห็นว่าเราส่ง h.GetAllTodo ซึ่งเป็น method ของ Handler ของตัวแปร h ไปให้ http.HandleFunc แทนที่จะส่ง anonymous หรือส่ง function ปกติตรงๆ ได้เหมือนกัน

เมื่อมี request วิ่งเข้ามา มันก็จะไปเรียก method GetAllTodo ของ h ให้นั่นเอง ซึ่งก็จะใช้งาน db ที่เรา config ไว้ในกับ h ตอนเราสร้าง Handler ใน main นั่นเอง

แบบนี้ตอนนี้เราต้องการทดสอบ GetAllTodo handler ก็สามารถทดสอบที่ method GetAllTodo ได้เลยโดย config ค่า db ได้ตามต้องการ
