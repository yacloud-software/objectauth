package main

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/objectauth"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/objectauth/db"
	"time"
)

func (e *objectAuthServer) AllowAllServiceAccess(ctx context.Context, req *pb.AllAccessRequest) (*pb.AllAccessResponse, error) {
	if req.ReadAccess == false && req.WriteAccess == false {
		return nil, errors.InvalidArgs(ctx, "no rights specified", "no rights specified")
	}
	svc := auth.GetService(ctx)
	if svc == nil {
		return &pb.AllAccessResponse{ReadAccess: false, WriteAccess: false}, nil
	}
	ls, err := db.DefaultDBServiceAccess().ByCallingService(ctx, svc.ID)
	if err != nil {
		return nil, err
	}
	for _, sa := range ls {
		if sa.SubjectService == "" || sa.CallingService == "" {
			return nil, fmt.Errorf("DATABASE objectauth broken, entry with ID %d in table %s has no service", sa.ID, db.DefaultDBServiceAccess().Tablename())
		}
		if sa.SubjectService == req.ServiceID {
			return &pb.AllAccessResponse{
				ReadAccess:  sa.ReadAccess,
				WriteAccess: sa.WriteAccess,
			}, nil
		}
	}
	res := &pb.AllAccessResponse{ReadAccess: false, WriteAccess: false}
	return res, nil
}
func (e *objectAuthServer) GrantAllServiceAccess(ctx context.Context, req *pb.GrantAllAccessRequest) (*common.Void, error) {
	err := errors.NeedsRoot(ctx)
	if err != nil {
		return nil, err
	}
	u := auth.GetUser(ctx)
	if u == nil {
		return nil, errors.Unauthenticated(ctx, "login please")
	}
	if req.CallingService == "" {
		return nil, errors.InvalidArgs(ctx, "missing calling service id", "missing calling service id")
	}
	if req.SubjectService == "" {
		return nil, errors.InvalidArgs(ctx, "missing calling subject id", "missing calling subject id")
	}
	sa := &pb.ServiceAccess{
		CallingService: req.CallingService,
		SubjectService: req.SubjectService,
		CreatedBy:      u.ID,
		Created:        uint32(time.Now().Unix()),
	}
	_, err = db.DefaultDBServiceAccess().Save(ctx, sa)
	if err != nil {
		return nil, err
	}
	return &common.Void{}, nil
}
