package core

import "context"

// Define a custom type for the context key
type contextKey string

const KeyRequester contextKey = "requester"

type Requester interface {
	GetSubject() string
	GetTokenId() string
}

type requesterData struct {
	Sub string `json:"user_id"`
	Tid string `json:"tid"`
}

func NewRequester(sub, tid string) *requesterData {
	return &requesterData{
		Sub: sub,
		Tid: tid,
	}
}

func (r *requesterData) GetSubject() string {
	return r.Sub
}

func (r *requesterData) GetTokenId() string {
	return r.Tid
}

func GetRequester(ctx context.Context) Requester {
	if requester, ok := ctx.Value(KeyRequester).(Requester); ok {
		return requester
	}

	return nil
}

func ContextWithRequester(ctx context.Context, requester Requester) context.Context {
	return context.WithValue(ctx, KeyRequester, requester)
}
