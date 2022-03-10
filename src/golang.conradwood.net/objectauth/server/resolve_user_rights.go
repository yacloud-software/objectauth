package main

import (
	"context"
	"fmt"
	//	au "golang.conradwood.net/apis/auth"
	pb "golang.conradwood.net/apis/objectauth"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/objectauth/shared"
)

func resolve_user_rights(ctx context.Context, req *pb.AuthRequest) (*pb.AccessRightList, error) {
	user := auth.GetUser(ctx)
	if user == nil {
		return nil, errors.Unauthenticated(ctx, "please log in")
	}
	userid := user.ID
	groups := user.Groups
	res := &pb.AccessRightList{
		ObjectType:           req.ObjectType,
		ObjectID:             req.ObjectID,
		UserPermissions:      make(map[string]*pb.PermissionSet),
		GroupPermissions:     make(map[string]*pb.PermissionSet),
		EffectivePermissions: &pb.PermissionSet{},
	}
	q := "select " + objects.SelectCols() + " from " + objects.Tablename() + " where objecttype = $1 and userid=$2 and objectid = $3"
	r, err := psql.QueryContext(ctx, "getuserobjectaccess", q, req.ObjectType, userid, req.ObjectID)
	if err != nil {
		return res, err
	}
	obs, err := objects.FromRows(ctx, r)
	r.Close()
	if err != nil {
		return nil, err
	}
	for _, db := range obs {
		p := res.UserPermissions[db.UserID]
		if db.Active {
			p = mergePerm(p, db)
		}
		res.UserPermissions[db.UserID] = p
	}
	// do we need to check groups?

	for _, g := range groups {
		ga, err := getGroupACL(ctx, g.ID, req.ObjectType, req.ObjectID)
		if err != nil {
			return nil, err
		}
		if ga != nil && ga.Active {
			p := res.GroupPermissions[ga.GroupID]
			p = mergePerm(p, ga)
			res.GroupPermissions[ga.GroupID] = p
		}
	}
	// now add composites...
	al, err := composite_right(ctx, req)
	if err != nil {
		return nil, err
	}
	if *debug {
		fmt.Printf("Composite perms result:\n")
		shared.PrintAccessRightList(al)
	}
	res = mergeAccessLists(res, al)
	res.EffectivePermissions = bestPermsFromList(res)

	return res, nil
}
