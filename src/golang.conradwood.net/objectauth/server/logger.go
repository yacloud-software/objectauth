package main

import (
	"context"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	pb "golang.conradwood.net/apis/objectauth"
	"golang.conradwood.net/go-easyops/auth"
)

func logDenied(ctx context.Context, u *apb.User, object_type pb.OBJECTTYPE, object_id uint64) {
	if u == nil {
		return
	}
	fmt.Printf("User %s [%s]: Denied access to object #%d of type %d (%s)\n", auth.Description(u), u.ID, object_id,
		object_type, pb.OBJECTTYPE_name[int32(object_type)],
	)
}
