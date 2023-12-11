package main

import (
	"context"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/objectauth"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/objectauth/composite"
	"golang.conradwood.net/objectauth/db"
)

var (
	compgroups *db.DBGroupToComposite
)

// this checks permissions. calling user must be also in the group and have execute permissions on the git repo id
func (e *objectAuthServer) GrantToSoftware(ctx context.Context, req *pb.IDGrantRequest) (*common.Void, error) {
	u := auth.GetUser(ctx)
	if u == nil {
		return nil, errors.Unauthenticated(ctx, "please login")
	}
	// check if user is in same group
	if !auth.IsInGroupByUser(u, req.GroupID) {
		return nil, errors.AccessDenied(ctx, "user \"%s\" must be in group \"%s\"", u.ID, req.GroupID)
	}
	// check if user has exe permissions
	ar, err := e.AskObjectAccess(ctx, &pb.AuthRequest{ObjectType: pb.OBJECTTYPE_GitRepository, ObjectID: req.ID})
	if err != nil {
		return nil, err
	}
	if !ar.Permissions.Execute {
		return nil, errors.AccessDenied(ctx, "user \"%s\" must have execute permissions to git repository #%d", u.ID, req.ID)
	}
	ls, err := compgroups.ByGroupID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	gtc := &pb.GroupToComposite{
		ObjectType: pb.COMPOSITETYPE_Software,
		ObjectID:   req.ID,
		GroupID:    req.GroupID,
		Active:     true,
		Read:       req.Read,
		Write:      req.Write,
		View:       true,
		Execute:    false,
	}

	for _, f := range ls {
		if f.ObjectType != pb.COMPOSITETYPE_Software {
			continue
		}
		if f.ObjectID != req.ID {
			continue
		}
		gtc.ID = f.ID
		break
	}
	if gtc.ID != 0 {
		err = compgroups.Update(ctx, gtc)
	} else {
		_, err = compgroups.Save(ctx, gtc)
	}
	if err != nil {
		return nil, err
	}
	return &common.Void{}, nil
}

// given an authrequest finds the matching composite ones
func composite_right(ctx context.Context, req *pb.AuthRequest) (*pb.AccessRightList, error) {
	c := composite.GetComposite(req.ObjectType)
	if c == nil {
		// no composite for this type -> empty list
		return &pb.AccessRightList{}, nil
	}
	return c.ForSingleObject(ctx, req)
}


