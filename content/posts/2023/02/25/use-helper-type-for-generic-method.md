---
title: "[Go] ใช้ generic ครอบ method ที่รับ any เพื่อให้อ่านและใช้งานง่ายยิ่งขึ้น"
date: 2023-02-25T07:19:03+07:00
draft: false
---

เรามี method ที่รับค่าเป็น type `any` ซึ่งจะส่งอะไรให้ก็ได้แต่จริงๆเราอยากจำกัดให้ส่งได้แค่บาง type วันนี้จะลองใช้ generic type ช่วยสร้าง type ใหม่ครอบการทำงานของ method นี้ทำให้คนใช้งานเห็นชัดเจนว่าใช้กับ type ไหนได้บ้าง แถม compiler ช่วยเช็คให้ตั้งแต่ตอน compile และ IDE/Editor ช่วย autocomplete ขึ้นมาให้อีกด้วย

<!--more-->

สมมติว่าเรามี method ที่เขียน Audit Log ซึ่งมี method แบบนี้

```go
// package auditlog
package auditlog

type auditLogger struct {
        // ...
}

func (l *auditLogger) Log(ctx context.Context, actionKey string, actionInfo any) error {
        // ...
        return nil
}

func New() *auditLogger {
        return &auditLogger{}
}
```

ตอนใช้งานเราก็จะส่งค่าอะไรให้กับ parameter actionInfo ก็ได้แบบนี้

```go
type SendSMSAuditLogInfo struct {
	CustomerName string
}

type ChangePasswordAuditLogInfo struct {
	Email string
}

type Logger interface {
	Log(ctx context.Context, actionKey string, actionInfo any) error
}

type Service struct {
	auditLogger Logger
}

func (s *Service) SendSMS(ctx context.Context) error {
        if err := s.auditLogger.Log(ctx, "send_sms", SendSMSAuditLogInfo{
                CustomerName: "John",
        }); err != nil {
                return err
        }

        return nil
}
```

จะเห็นว่าก็ใช้งานง่ายดีแต่ว่าถ้าเราใช้ง่าย auditLogger.Log เราจะไม่รู้เลยว่า actionInfo เป็น type อะไรได้บ้าง ถ้าเราเปลี่ยนไปใช้ Generic แล้วสร้าง type constraint ขึ้นมาใหม่ให้มีแต่เซตของ type ที่เราต้องการ ก็จะทำให้ compiler ช่วยตรวจสอบได้แบบนี้

```go
// 1) เริ่มจากสร้าง type constraint เพื่อระบุเซตของ AuditLogInfo type ที่เราอนุญาติให้ใช้เท่านั้น
type AuditLogAction interface {
	SendSMSAuditLogInfo | ChangePasswordAuditLogInfo
}

// 2) เริ่มจากสร้าง type ขึ้นมาใหม่เพื่อครอบการทำงานของ auditLogger.Log method
// ที่รองรับ generic type parameter T ที่มี constraint เป็น AuditLogAction
// และมี field ที่เก็บ type ใดๆที่ implements Log method

type AuditLogger[T AuditLogAction] struct {
        logger Logger
}

// 3) จากนั้นสร้าง wrapper method ชื่อ Log สำหรับ type AuditLogger
// ที่รองรับ generic type parameter T เพื่อเอาไปใช้เป็น type ของ actionInfo parameter
func (l *AuditLogger[T]) Log(ctx context.Context, actionKey string, actionInfo T) error {
        // 4) เรียก l.logger.Log อีกทีโดยใช้ค่าที่รับมาผ่าน generic type
	return l.logger.Log(ctx, actionKey, actionInfo)
}

// 5) สร้าง constructor function ให้กับ AuditLogger[T AuditLogAction]
func NewAuditLogger[T AuditLogAction](l Logger) *AuditLogger[T] {
	return &AuditLogger[T]{
		logger: l,
	}
}
```

จากขั้นตอนที่ 3 และ 4 จะเห็นว่าเรา implement wrapper method โดยก็ไปเรียก method ที่รับเป็น type any อีกที

ตอนใช้งานแทนที่เราจะใช้ผ่าน Log ที่รับ type any เราก็ใช้ผ่าน AuditLogger แทนแบบนี้

```go
func (s *Service) SendSMS(ctx context.Context) error {
	auditLogger := NewAuditLogger[SendSMSAuditLogInfo](s.auditLogger)
	if err := auditLogger.Log(ctx, "send_sms", SendSMSAuditLogInfo{
		CustomerName: "John",
	}); err != nil {
		return err
	}

	return nil
}
```

ตอนใช้งานกับ IDE/Editor ที่รองรับ ก็จะมี autocomplete ที่ระบุ type ให้ชัดเจน ไม่ใส่ผิดพลาดแบบนี้

![generic autocompleted in VSCode](/2023-02-25-generic-autocomplete.png)
