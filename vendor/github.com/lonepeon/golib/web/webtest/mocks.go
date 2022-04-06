package webtest

//go:generate mockgen -destination=web.go -package webtest github.com/lonepeon/golib/web AuthenticationFrontendStorer,AuthenticationBackendStorer,Handler,Context
