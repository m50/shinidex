package logger

import (
	"github.com/gookit/slog"
	"github.com/m50/shinidex/pkg/web/session"
)

type ContextProcessor struct{}

func (fn ContextProcessor) Process(record *slog.Record) {
	ctx := record.Ctx
	user, _ := session.GetAuthedUserContext(ctx)
	if user != nil {
		record.AddData(slog.M{KeyUserID: user.ID})
	}
}