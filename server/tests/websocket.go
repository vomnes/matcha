package tests

import (
	"context"
	"net/http"

	"../lib"
)

func WithContextWS(r *http.Request, ctxData ContextData) *http.Request {
	ctx := context.WithValue(r.Context(), lib.UserID, ctxData.UserID)
	if ctxData.DB != nil {
		ctx = context.WithValue(ctx, lib.Database, ctxData.DB)
	}
	if ctxData.Client != nil {
		ctx = context.WithValue(ctx, lib.Redis, ctxData.Client)
	}
	if ctxData.MailjetClient != nil {
		ctx = context.WithValue(ctx, lib.MailJet, ctxData.MailjetClient)
	}
	ctx = context.WithValue(ctx, lib.UserID, ctxData.UserID)
	ctx = context.WithValue(ctx, lib.Username, ctxData.Username)
	ctx = context.WithValue(ctx, lib.UUID, ctxData.UUID)
	return r.WithContext(ctx)
}
