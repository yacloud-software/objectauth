package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/objectauth"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/objectauth/db"
)

var (
	disallow_all_service_access_for_users = flag.Bool("disallow_all_service_access_for_users", true, "if true, only services get to access all objects. if a context has both, user and service, it will be disallowed")
	debug_service_access                  = flag.Bool("debug_service_access", false, "if true debug service access")
)

// ask for access
func (e *objectAuthServer) AllowAllServiceAccess(ctx context.Context, req *pb.AllAccessRequest) (*pb.AllAccessResponse, error) {
	resp, err := e.AllowAllServiceAccessErr(ctx, req)
	if err != nil || (resp.ReadAccess == false && resp.WriteAccess == false) {
		svc, _ := authremote.GetUserByID(ctx, req.ServiceID)
		svcs := fmt.Sprintf("%s", req.ServiceID)
		if svc != nil {
			svcs = auth.UserIDString(svc)
		}
		logAccessDenied(ctx, "all access denied for service %s for objecttype \"%v\" (%d)", svcs, req.ObjectType, req.ObjectType)
	}
	return resp, err
}
func (e *objectAuthServer) AllowAllServiceAccessErr(ctx context.Context, req *pb.AllAccessRequest) (*pb.AllAccessResponse, error) {
	svc := auth.GetService(ctx)
	if svc == nil {
		return &pb.AllAccessResponse{ReadAccess: false, WriteAccess: false}, nil
	}
	svc_debugf("access request from service %s for service %s for objecttype \"%v\" (%d)\n", svc.ID, req.ServiceID, req.ObjectType, req.ObjectType)

	if *disallow_all_service_access_for_users {
		u := auth.GetUser(ctx)
		if u != nil {
			svc_debugf("access request from service %s for service %s for objecttype \"%v\" (%d) denied because also running as user %s\n", svc.ID, req.ServiceID, req.ObjectType, req.ObjectType, auth.UserIDString(u))
			return nil, errors.AccessDenied(ctx, "access denied")
		}
	}

	ls, err := db.DefaultDBServiceAccess().ByCallingService(ctx, svc.ID)
	if err != nil {
		return nil, err
	}
	for _, sa := range ls {
		if sa.SubjectService == "" || sa.CallingService == "" {
			return nil, fmt.Errorf("DATABASE objectauth broken, entry with ID %d in table %s has no service", sa.ID, db.DefaultDBServiceAccess().Tablename())
		}
		if sa.ObjectType != req.ObjectType {
			continue
		}
		if sa.SubjectService == req.ServiceID {
			return &pb.AllAccessResponse{
				ReadAccess:  sa.ReadAccess,
				WriteAccess: sa.WriteAccess,
			}, nil
		}
	}
	svc_debugf("access request from service %s for service %s for objecttype %v: DENIED\n", svc.ID, req.ServiceID, req.ObjectType)
	res := &pb.AllAccessResponse{ReadAccess: false, WriteAccess: false}
	return res, nil
}
func (e *objectAuthServer) GrantAllServiceAccess(ctx context.Context, req *pb.GrantAllAccessRequest) (*common.Void, error) {
	if req.ReadAccess == false && req.WriteAccess == false {
		return nil, errors.InvalidArgs(ctx, "no rights specified", "no rights specified")
	}
	err := errors.NeedsRoot(ctx)
	if err != nil {
		return nil, err
	}
	u := auth.GetUser(ctx)
	if u == nil {
		return nil, errors.Unauthenticated(ctx, "login please")
	}
	if req.CallingService == "" || req.CallingService == "0" {
		return nil, errors.InvalidArgs(ctx, "missing calling service id", "missing calling service id")
	}
	if req.SubjectService == "" || req.SubjectService == "0" {
		return nil, errors.InvalidArgs(ctx, "missing calling subject id", "missing calling subject id")
	}
	sa := &pb.ServiceAccess{
		CallingService: req.CallingService,
		SubjectService: req.SubjectService,
		CreatedBy:      u.ID,
		ObjectType:     req.ObjectType,
		ReadAccess:     req.ReadAccess,
		WriteAccess:    req.WriteAccess,
		Created:        uint32(time.Now().Unix()),
	}
	_, err = db.DefaultDBServiceAccess().Save(ctx, sa)
	if err != nil {
		return nil, err
	}
	return &common.Void{}, nil
}
func svc_debugf(format string, args ...interface{}) {
	if !*debug_service_access {
		return
	}
	fmt.Printf(format, args...)
}
