package main

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/auth"
)

func logAccessDenied(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	u := auth.CurrentUserString(ctx)
	svc := auth.GetService(ctx)
	svcs := "noservice"
	if svc != nil {
		svcs = fmt.Sprintf("service %s/%s", svc.ID, svc.Email)
	}
	fmt.Printf("[%s @ %s] %s\n", u, svcs, msg)
}





