package main

import "context"

type auditLogger struct {
}

func (l *auditLogger) Log(ctx context.Context, actionKey string, actionInfo any) error {
	return nil
}

func New() *auditLogger {
	return &auditLogger{}
}

type SendSMSAuditLogInfo struct {
	CustomerName string
}

type ChangePasswordAuditLogInfo struct {
	Email string
}

type AuditLogAction interface {
	SendSMSAuditLogInfo | ChangePasswordAuditLogInfo
}

type Logger interface {
	Log(ctx context.Context, actionKey string, actionInfo any) error
}

type AuditLogger[T AuditLogAction] struct {
	logger Logger
}

func (l *AuditLogger[T]) Log(ctx context.Context, actionKey string, actionInfo T) error {
	return l.logger.Log(ctx, actionKey, actionInfo)
}

func NewAuditLogger[T AuditLogAction](l Logger) *AuditLogger[T] {
	return &AuditLogger[T]{
		logger: l,
	}
}

type Service struct {
	auditLogger Logger
}

func (s *Service) SendSMS(ctx context.Context) error {
	auditLogger := NewAuditLogger[SendSMSAuditLogInfo](s.auditLogger)
	if err := auditLogger.Log(ctx, "send_sms", SendSMSAuditLogInfo{
		CustomerName: "John",
	}); err != nil {
		return err
	}

	return nil
}
