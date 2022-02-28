package webtest

//go:generate mockgen -destination=web.go -package webtest github.com/lonepeon/golib/web CurrentAuthenticatedUserStorage,Handler,Context
