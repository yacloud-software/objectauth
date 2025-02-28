package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.conradwood.net/go-easyops/auth"
)

var (
	dedup_log      = make(map[string]time.Time)
	dedup_log_lock sync.Mutex
)

func logAccessDenied(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	u := auth.CurrentUserString(ctx)
	svc := auth.GetService(ctx)
	svcs := "noservice"
	if svc != nil {
		svcs = fmt.Sprintf("service %s/%s", svc.ID, svc.Email)
	}
	txt := fmt.Sprintf("[%s @ %s] %s", u, svcs, msg)
	print := false
	dedup_log_lock.Lock()
	last_time, found := dedup_log[txt]
	if found {
		if time.Since(last_time) > time.Duration(60)*time.Second {
			print = true
		}
	} else {
		print = true
	}
	if print {
		dedup_log[txt] = time.Now()
	}
	dedup_log_lock.Unlock()
	if print {
		fmt.Println(txt)
	}
}
