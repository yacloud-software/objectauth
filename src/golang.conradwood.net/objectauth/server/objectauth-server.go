package main

import (
	"context"
	"flag"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/objectauth"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/sql"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/objectauth/db"
	//	"golang.conradwood.net/objectauth/shared"
	"google.golang.org/grpc"
	"os"
)

var (
	debug     = flag.Bool("debug", false, "debug mode")
	port      = flag.Int("port", 4100, "The grpc server port")
	allow_all = flag.Bool("allow_all", false, "allow all requests")
	objects   *db.DBUserToObject
	gobjects  *db.DBGroupToObject
	psql      *sql.DB
	accesses  = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "objectauth_total_requests",
			Help: "V=1 UNIT=ops DESC=incremented each time a request is received",
		},
	)
	accessesFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "objectauth_failed_requests",
			Help: "V=1 UNIT=ops DESC=incremented each time a request failed",
		},
	)
)

func init() {
	prometheus.MustRegister(accesses, accessesFailed)
}

type objectAuthServer struct {
}

func main() {
	flag.Parse()
	fmt.Printf("Starting ObjectAuthServiceServer...\n")
	xpsql, err := sql.Open()
	utils.Bail("failed to open sql database", err)
	psql = xpsql
	objects = db.NewDBUserToObject(psql)
	gobjects = db.NewDBGroupToObject(psql)
	compgroups = db.DefaultDBGroupToComposite() // newer, better way of doing it
	sd := server.NewServerDef()
	sd.Port = *port
	//	migratedb()
	sd.Register = server.Register(
		func(server *grpc.Server) error {
			e := new(objectAuthServer)
			pb.RegisterObjectAuthServiceServer(server, e)
			return nil
		},
	)
	err = server.ServerStartup(sd)
	utils.Bail("Unable to start server", err)
	os.Exit(0)
}
func migratedb() {
	ctx := authremote.Context()
	tables := []string{"usertoobject", "grouptoobject"}
	for _, t := range tables {
		rows, err := psql.QueryContext(ctx, "mig1_"+t, "select id from "+t+" where read = false")
		utils.Bail("mig "+t+" failed", err)
		for rows.Next() {
			var id uint64
			err = rows.Scan(&id)
			utils.Bail("Scan "+t+" failed", err)
			_, err = psql.ExecContext(ctx, "mig2_"+t, "update "+t+" set view=true,read=true,write=true,execute=true where id = $1", id)
			utils.Bail("update "+t+" failed", err)
		}
		rows.Close()
	}
}

/************************************
* grpc functions
************************************/
func (e *objectAuthServer) ResolveForUser(ctx context.Context, req *pb.ResolveRequest) (*pb.ResolveResponse, error) {
	err := errors.NeedsRoot(ctx)
	if err != nil {
		return nil, err
	}
	u := auth.GetUser(ctx)
	if u == nil {
		return nil, errors.Unauthenticated(ctx, "login please")
	}
	fmt.Printf("Resolving rights for user %s\n", req.UserID)
	nctx, err := authremote.ContextForUserID(req.UserID)
	if err != nil {
		return nil, err
	}
	res, err := resolve_user_rights(nctx, req.AuthRequest)
	if err != nil {
		return nil, err
	}
	u = auth.GetUser(nctx)
	r := &pb.ResolveResponse{AccessRightList: res, User: u}
	return r, nil
}

// ask if a user has access to a specific object
func (e *objectAuthServer) AskObjectAccess(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	accesses.Inc()
	resp := &pb.AuthResponse{Granted: false}
	u := auth.GetUser(ctx)
	if u == nil {
		u = auth.GetService(ctx)
		if u == nil {
			return resp, nil
		}
	}
	// TODO: HACK FOR USERAPPREPORIGHTS FLAGS
	if req.ObjectType == pb.OBJECTTYPE_UserAppFlags {
		req.ObjectType = pb.OBJECTTYPE_GitRepository
	}
	if HasAllAccess(u.ID, req.ObjectType) {
		return &pb.AuthResponse{Granted: true, Permissions: &pb.PermissionSet{Read: true, Write: true, Execute: true, View: true}}, nil
	}
	if *debug {
		fmt.Printf("Access request for user #%s(%s) for objecttype=%v, objectid=%d\n", u.ID, auth.Description(u), req.ObjectType, req.ObjectID)
	}

	if *allow_all {
		resp.Granted = true
		return resp, nil
	}
	arl, err := resolve_user_rights(ctx, req)
	if err != nil {
		return nil, err
	}
	res := &pb.AuthResponse{
		Permissions: arl.EffectivePermissions,
	}
	res.Granted = res.Permissions.Read || res.Permissions.Write || res.Permissions.Execute || res.Permissions.View
	return res, nil
	/*
		// we actually have to hit the database...
		q := "select " + objects.SelectCols() + " from " + objects.Tablename() + " where objecttype = $1 and userid=$2 and objectid = $3"
		r, err := psql.QueryContext(ctx, "getuserobjectaccess", q, req.ObjectType, u.ID, req.ObjectID)
		if err != nil {
			return resp, err
		}
		obs, err := objects.FromRows(ctx, r)
		r.Close()
		if err != nil {
			return nil, err
		}
		p := &pb.PermissionSet{}
		resp.Permissions = p
		for _, db := range obs {
			if db.Active {
				mergePerm(p, db)
				break
			}
		}
		// do we need to check groups?

		for _, g := range u.Groups {
			ga, err := getGroupACL(ctx, g.ID, req.ObjectType, req.ObjectID)
			if err != nil {
				return nil, err
			}
			if ga != nil && ga.Active {
				mergePerm(p, ga)
				break
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
		perms := bestPermsFromList(al)
		resp.Permissions = mergePerm(resp.Permissions, perms)
		//fmt.Printf("Composite perms: %#v\n", perms)
		resp.Granted = p.Read || p.Write || p.Execute || p.View
		if !resp.Granted {
			accessesFailed.Inc()
			logDenied(ctx, u, req.ObjectType, req.ObjectID)
		}

		return resp, nil
	*/
}

// get all objects available for a user
func (e *objectAuthServer) AvailableObjects(ctx context.Context, req *pb.ObjectType) (*pb.ObjectIDList, error) {
	resp := &pb.ObjectIDList{}
	u := auth.GetUser(ctx)
	if u == nil {
		return resp, nil
	}
	fmt.Printf("Getting Available Objects for user #%s\n", u.ID)
	q := "select " + objects.SelectCols() + " from " + objects.Tablename() + " where objecttype = $1 and userid=$2 and active=true"
	r, err := psql.QueryContext(ctx, "getuserobjectall", q, req.ObjectType, u.ID)
	if err != nil {
		return resp, err
	}
	obs, err := objects.FromRows(ctx, r)
	r.Close()
	if err != nil {
		return nil, err
	}

	for _, ob := range obs {
		resp.ObjectIDs = append(resp.ObjectIDs, ob.ObjectID)
	}

	// now add groups to the mix

	for _, g := range u.Groups {
		q = "select " + gobjects.SelectCols() + " from " + gobjects.Tablename() + " where objecttype = $1 and groupid=$2 and active=true"
		r, err := psql.QueryContext(ctx, "getgroupall", q, req.ObjectType, g.ID)
		if err != nil {
			return resp, err
		}
		gobs, err := gobjects.FromRows(ctx, r)
		r.Close()
		if err != nil {
			return nil, err
		}
		for _, gob := range gobs {
			f := false
			for _, tst := range resp.ObjectIDs {
				if tst == gob.ObjectID {
					f = true
					break
				}
			}
			if !f {
				resp.ObjectIDs = append(resp.ObjectIDs, gob.ObjectID)
			}
		}

	}
	fmt.Printf("Got %d Available Objects for user #%s\n", len(resp.ObjectIDs), u.ID)
	return resp, nil
}
func (e *objectAuthServer) GrantToMe(ctx context.Context, req *pb.GrantUserRequest) (*common.Void, error) {
	if !auth.IsRoot(ctx) && !extraService(ctx, req.ObjectType) {
		return nil, errors.AccessDenied(ctx, "privileged access required")
	}
	u := auth.GetUser(ctx)
	if u == nil {
		return nil, errors.Unauthenticated(ctx, "access denied")
	}
	req.UserID = u.ID
	return e.GrantToUser(ctx, req)
}
func (e *objectAuthServer) GrantToGroup(ctx context.Context, req *pb.GrantGroupRequest) (*common.Void, error) {
	if !auth.IsRoot(ctx) && !extraService(ctx, req.ObjectType) {
		return nil, errors.AccessDenied(ctx, "only root can grant stuff atm")
	}
	if req.GroupID == "" {
		return nil, errors.InvalidArgs(ctx, "missing userid", "Missing userid")
	}
	if req.ObjectID == 0 {
		return nil, errors.InvalidArgs(ctx, "objectid of 0 is not valid", "objectid of 0 is not valid")
	}

	g, err := authremote.GetAuthManagerClient().GetGroupByID(ctx, &apb.GetGroupRequest{ID: req.GroupID})
	if err != nil {
		return nil, err
	}
	fmt.Printf("Granting access to group %s (%s)\n", g.ID, g.Name)

	uto, err := getGroupACL(ctx, req.GroupID, req.ObjectType, req.ObjectID)
	if err != nil {
		return nil, err
	}
	if uto != nil {
		uto.Active = true
		uto.Read = req.Read
		uto.Write = req.Write
		uto.Execute = req.Execute
		uto.View = req.View
		err = gobjects.Update(ctx, uto)
	} else {
		uto = &pb.GroupToObject{
			ObjectType: req.ObjectType,
			ObjectID:   req.ObjectID,
			GroupID:    req.GroupID,
			Read:       req.Read,
			Write:      req.Write,
			Execute:    req.Execute,
			View:       req.View,
			Active:     true,
		}
		_, err = gobjects.Save(ctx, uto)
	}
	if err != nil {
		return nil, err
	}
	return &common.Void{}, nil
}

func (e *objectAuthServer) GrantToUser(ctx context.Context, req *pb.GrantUserRequest) (*common.Void, error) {
	if !auth.IsRoot(ctx) && !extraService(ctx, req.ObjectType) {
		return nil, errors.AccessDenied(ctx, "only root can grant stuff atm")
	}
	if req.UserID == "" {
		return nil, errors.InvalidArgs(ctx, "missing userid", "Missing userid")
	}
	if req.ObjectID == 0 {
		return nil, errors.InvalidArgs(ctx, "objectid of 0 is not valid", "objectid of 0 is not valid")
	}
	if req.ObjectType == 0 {
		return nil, errors.InvalidArgs(ctx, "objecttype of 0 is not valid", "objecttype of 0 is not valid")
	}
	uto, err := getUserACL(ctx, req.UserID, req.ObjectType, req.ObjectID)
	if err != nil {
		return nil, err
	}
	if uto != nil {
		uto.Active = true
		uto.Read = req.Read
		uto.Write = req.Write
		uto.Execute = req.Execute
		uto.View = req.View
		err = objects.Update(ctx, uto)
	} else {
		uto = &pb.UserToObject{ObjectType: req.ObjectType,
			ObjectID: req.ObjectID,
			UserID:   req.UserID,
			Active:   true,
			Read:     req.Read,
			Write:    req.Write,
			Execute:  req.Execute,
			View:     req.View,
		}
		_, err = objects.Save(ctx, uto)
	}
	if err != nil {
		return nil, err
	}
	return &common.Void{}, nil
}

// get rights of a particular object (all user rights - not current user)
func (e *objectAuthServer) GetRights(ctx context.Context, req *pb.AuthRequest) (*pb.AccessRightList, error) {
	if *debug {
		fmt.Printf("Getting rights  on object %d (type %s)\n", req.ObjectID, req.ObjectType)
	}

	g, err := e.AskObjectAccess(ctx, req)
	if err != nil {
		return nil, err
	}
	if !g.Permissions.View {
		return nil, errors.NotFound(ctx, "view access required to get rights of service")
	}
	res := &pb.AccessRightList{
		ObjectType:           req.ObjectType,
		ObjectID:             req.ObjectID,
		UserPermissions:      make(map[string]*pb.PermissionSet),
		GroupPermissions:     make(map[string]*pb.PermissionSet),
		EffectivePermissions: &pb.PermissionSet{},
	}

	// now add the users
	q := "select " + objects.SelectCols() + " from " + objects.Tablename() + " where objecttype = $1 and objectid = $2"
	r, err := psql.QueryContext(ctx, "getuserobjectaccess2", q, req.ObjectType, req.ObjectID)
	if err != nil {
		return nil, err
	}
	obs, err := objects.FromRows(ctx, r)
	r.Close()
	if err != nil {
		return nil, err
	}
	for _, o := range obs {
		if !o.Active {
			continue
		}
		res.UserPermissions[o.UserID] = mergePerm(res.UserPermissions[o.UserID], o)
	}

	// now add the groups
	q = "select " + gobjects.SelectCols() + " from " + gobjects.Tablename() + " where objecttype = $1 and objectid = $2"
	r, err = psql.QueryContext(ctx, "getgroupobjectaccess2", q, req.ObjectType, req.ObjectID)
	if err != nil {
		return nil, err
	}
	xobs, err := gobjects.FromRows(ctx, r)
	r.Close()
	if err != nil {
		return nil, err
	}
	for _, o := range xobs {
		if !o.Active {
			continue
		}
		res.GroupPermissions[o.GroupID] = mergePerm(res.GroupPermissions[o.GroupID], o)
	}

	/*
		// now add composites...

		al, err := composite_right(ctx, req)
		if err != nil {
			return nil, err
		}
		res = mergeAccessLists(res, al)
		res.EffectivePermissions = bestPermsFromList(res)
	*/
	return res, nil

}

func getUserACL(ctx context.Context, userid string, object_type pb.OBJECTTYPE, object_id uint64) (*pb.UserToObject, error) {
	q := "select " + objects.SelectCols() + " from " + objects.Tablename() + " where objecttype = $1 and userid=$2 and objectid = $3"
	r, err := psql.QueryContext(ctx, "getuserobjectaccess", q, object_type, userid, object_id)
	if err != nil {
		return nil, err
	}
	obs, err := objects.FromRows(ctx, r)
	r.Close()
	if err != nil {
		return nil, err
	}
	if len(obs) == 0 {
		return nil, nil
	}
	return obs[0], nil
}
func getGroupACL(ctx context.Context, groupid string, object_type pb.OBJECTTYPE, object_id uint64) (*pb.GroupToObject, error) {
	/*
		if *debug {
			fmt.Printf("Getting acl for group=%s, type=%d, id=%d\n", groupid, object_type, object_id)
		}
	*/
	q := "select " + gobjects.SelectCols() + " from " + gobjects.Tablename() + " where objecttype = $1 and groupid=$2 and objectid = $3"
	r, err := psql.QueryContext(ctx, "getgroupobjectaccess", q, object_type, groupid, object_id)
	if err != nil {
		return nil, err
	}
	obs, err := gobjects.FromRows(ctx, r)
	r.Close()
	if err != nil {
		return nil, err
	}
	if len(obs) == 0 {
		return nil, nil
	}
	return obs[0], nil
}

type embeddedPermissions interface {
	GetRead() bool
	GetWrite() bool
	GetExecute() bool
	GetView() bool
}

func toPerm(ep embeddedPermissions) *pb.PermissionSet {
	res := &pb.PermissionSet{
		Read:    ep.GetRead(),
		Write:   ep.GetWrite(),
		Execute: ep.GetExecute(),
		View:    ep.GetView(),
	}
	return res
}

func addToAccessList(src, add map[string]*pb.PermissionSet) {
	for k, p := range add {
		old, exists := src[k]
		if !exists {
			// no need to merge, doesn't exist yet, so add it
			src[k] = p
			continue
		}
		src[k] = mergePerm(old, p)
	}
}

// given two accesslists, merges their permission sets
func mergeAccessLists(a, b *pb.AccessRightList) *pb.AccessRightList {
	res := &pb.AccessRightList{
		ObjectType:       a.ObjectType,
		ObjectID:         b.ObjectID,
		GroupPermissions: make(map[string]*pb.PermissionSet),
		UserPermissions:  make(map[string]*pb.PermissionSet),
	}
	addToAccessList(res.GroupPermissions, a.GroupPermissions)
	addToAccessList(res.GroupPermissions, b.GroupPermissions)
	addToAccessList(res.UserPermissions, a.UserPermissions)
	addToAccessList(res.UserPermissions, b.UserPermissions)

	return res
}

// given a map of permissions, works out the "best" (that is most permissive) set
func bestPermsFromMap(m map[string]*pb.PermissionSet) *pb.PermissionSet {
	res := &pb.PermissionSet{}
	for _, v := range m {
		res = mergePerm(res, v)
	}
	return res
}

// given an accessrightlist, works out the "best" (that is most permissive) set
func bestPermsFromList(a *pb.AccessRightList) *pb.PermissionSet {
	if a == nil {
		return &pb.PermissionSet{}
	}
	pg := bestPermsFromMap(a.GroupPermissions)
	pu := bestPermsFromMap(a.UserPermissions)
	return mergePerm(pg, pu)
}

func mergePerm(p *pb.PermissionSet, ep embeddedPermissions) *pb.PermissionSet {
	res := p
	if res == nil {
		res = &pb.PermissionSet{}
	}
	res.Read = res.Read || ep.GetRead()
	res.Write = res.Write || ep.GetWrite()
	res.Execute = res.Execute || ep.GetExecute()
	res.View = res.View || ep.GetView()
	return res
}

func HasAllAccess(userid string, obj pb.OBJECTTYPE) bool {
	if obj == pb.OBJECTTYPE_SingingCatModule {
		if userid == "6139" { // scapply
			return true
		}
	}
	return false
}

func extraService(ctx context.Context, t pb.OBJECTTYPE) bool {
	svc := auth.GetService(ctx)
	if svc == nil {
		fmt.Printf("No extra service at all\n")
		return false
	}
	ruid := auth.GetServiceIDByName("repobuilder.RepoBuilder")
	if svc.ID == ruid && t == pb.OBJECTTYPE_GitRepository { //repobuilder
		return true
	}
	if svc.ID == ruid && t == pb.OBJECTTYPE_Artefact { //repobuilder
		return true
	}
	if svc.ID == "145" && t == pb.OBJECTTYPE_SingingCatModule { //scweb
		return true
	}
	fmt.Printf("not an extra service: \"%s\" for %v\n", svc.ID, t)
	return false
}
