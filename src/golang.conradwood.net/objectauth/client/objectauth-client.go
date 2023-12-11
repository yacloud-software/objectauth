package main

import (
	"flag"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/objectauth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/objectauth/grant"
	"golang.conradwood.net/objectauth/resolve"
	"golang.conradwood.net/objectauth/shared"
	"os"
	"strconv"
	"strings"
)

var (
	echoClient       pb.ObjectAuthServiceClient
	repo_yacloud_dev = flag.Int("repo_for_yacloud", 0, "composite software access. if this is provided and is a valid repository id, then that will be accessible for yacloud devs")
	exe              = flag.Bool("execute", false, "grant execute access")
	view             = flag.Bool("view", true, "grant write access")
	write            = flag.Bool("write", true, "grant write access")
	read             = flag.Bool("read", true, "grant read access")
	service          = flag.String("NOT_VALID_service", "", "DEPRECATED! request access to a service by name for me (short for -userid=[me] -objectid=service -objecttype=[service])")
	callingservice   = flag.Uint64("calling_service_id", 0, "service by id for allaccess, this is the id of the service which makes the call to objectauth")
	subjectservice   = flag.Uint64("subject_service_id", 0, "service by id for allaccess, this is the id of the service for which all access shall be granted")
	userid           = flag.String("userid", "", "userid to grant access to")
	groupid          = flag.String("groupid", "", "groupid to grant access to")
	objectid         = flag.String("objectid", "", "object id to grant access for")
	objecttype       = flag.String("objecttype", "", "object type to grant access")
	allaccess        = flag.Bool("allaccess", false, "set all accessmode for -serviceid on type -objecttype")
	check            = flag.Bool("check", false, "check access for user specified by userid (needs root)")
)

func main() {
	flag.Parse()
	if *allaccess {
		utils.Bail("failed to set AllAccess()", AllAccess())
		os.Exit(0)
	}
	if *repo_yacloud_dev != 0 {
		do_repo_yacloud_dev()
		return
	}
	if *service != "" {
		gid := ""
		if *groupid != "" {
			gid = getGroupID()
		}
		err := grant.GrantService(authremote.Context(), *service, gid) // gid is optional
		utils.Bail("failed to grant", err)
		fmt.Printf("Done\n")
		os.Exit(0)
	}

	// a context with authentication
	ctx := authremote.Context()
	var err error
	echoClient = pb.GetObjectAuthServiceClient()
	if *check {
		uid := getUserID()
		fmt.Printf("Getting permissions for user \"%s\"\n", uid)
		r, err := echoClient.ResolveForUser(ctx, &pb.ResolveRequest{
			AuthRequest: &pb.AuthRequest{ObjectType: getObjectType(), ObjectID: getObjectID()},
			UserID:      uid,
		})

		utils.Bail("failed to get access rights", err)
		shared.PrintRequestResponse(r)
		os.Exit(0)
	}

	if *userid != "" {
		f := &pb.GrantUserRequest{
			ObjectType: getObjectType(),
			ObjectID:   getObjectID(),
			UserID:     *userid,
			Read:       true,
			Write:      true,
			Execute:    *exe,
			View:       true,
		}
		_, err := echoClient.GrantToUser(ctx, f)
		utils.Bail("Failed to grant access", err)
		fmt.Printf("granted access.\ndone\n")
		os.Exit(0)
	}
	if *groupid != "" {
		f := &pb.GrantGroupRequest{
			ObjectType: getObjectType(),
			ObjectID:   getObjectID(),
			GroupID:    getGroupID(),
			Read:       *read,
			Write:      *write,
			Execute:    *exe,
			View:       *view,
		}
		_, err := echoClient.GrantToGroup(ctx, f)
		utils.Bail("Failed to grant group access", err)
		fmt.Printf("done\n")
		os.Exit(0)
	}
	ctx = authremote.Context()
	aor := &pb.ObjectType{ObjectType: getObjectType()}
	fmt.Printf("Getting rights for type #%d...\n", aor.ObjectType)
	response, err := echoClient.AvailableObjects(ctx, aor)
	utils.Bail("Failed to ping server", err)
	fmt.Printf("Got %d objects\n", len(response.ObjectIDs))
	for _, r := range response.ObjectIDs {
		fmt.Printf("Object #%d\n", r)
		ri, err := echoClient.GetRights(ctx, &pb.AuthRequest{})
		utils.Bail("failed to get rights", err)
		fmt.Printf(" Permissions: %s\n", permToString(ri))
	}
	fmt.Printf("Done.\n")
	os.Exit(0)
}

func getObjectType() pb.OBJECTTYPE {
	m := make(map[string]pb.OBJECTTYPE)
	for k, v := range pb.OBJECTTYPE_name {
		if k == 0 {
			continue
		}
		lv := strings.ToLower(v)
		m[lv] = pb.OBJECTTYPE(k)
	}

	for k, v := range m {
		if strings.ToLower(k) == strings.ToLower(*objecttype) {
			return v
		}
	}
	fmt.Printf("Invalid objecttype. choose one of:\n")
	for k, _ := range m {
		fmt.Printf(" \"%s\"\n", k)
	}
	os.Exit(10)
	return 0
}
func getObjectID() uint64 {
	id, err := strconv.ParseUint(*objectid, 10, 64)
	if err == nil {
		return id
	}
	num, s, err := resolve.ResolveToNumber(getObjectType(), *objectid)
	utils.Bail("failed to resolve to number", err)
	fmt.Printf("Operating on \"%s\" (id=%d)\n", s, num)
	return num

}

func getGroupID() string {
	if *groupid == "all" {
		return *groupid
	}
	gid, err := strconv.ParseUint(*groupid, 10, 64)
	if err == nil {
		return fmt.Sprintf("%d", gid)
	}
	am := authremote.GetAuthManagerClient()
	ctx := authremote.Context()
	lr, err := am.ListGroups(ctx, &common.Void{})
	utils.Bail("failed to list groups", err)
	var group *apb.Group
	for _, g := range lr.Groups {
		if g.Name == *groupid {
			if group != nil {
				fmt.Printf("Multiple Group matches: %s & %s\n", group.Name, g.Name)
				os.Exit(10)
			}
			group = g
		}
	}
	if group == nil {
		fmt.Printf("No match for group %s\n", *groupid)
		os.Exit(10)
	}
	fmt.Printf("Group \"%s\" has ID #%s\n", group.Name, group.ID)
	return group.ID
}

func permToString(pb *pb.AccessRightList) string {
	res := "users: "
	for k, v := range pb.UserPermissions {
		res = res + k + "(" + permString(v) + ") "
	}
	res = res + "groups: "
	for k, v := range pb.GroupPermissions {
		res = res + k + "(" + permString(v) + ") "
	}
	return res
}
func permString(pb *pb.PermissionSet) string {
	res := make([]byte, 4)
	if pb.Read {
		res[0] = 'r'
	}
	if pb.Write {
		res[1] = 'w'
	}
	if pb.Execute {
		res[2] = 'x'
	}
	if pb.View {
		res[3] = 'v'
	}
	return string(res)
}

func do_repo_yacloud_dev() {
	id := uint64(*repo_yacloud_dev)
	ctx := authremote.Context()
	req := &pb.AuthRequest{ObjectType: pb.OBJECTTYPE_GitRepository, ObjectID: id}
	ar, err := pb.GetObjectAuthServiceClient().AskObjectAccess(ctx, req)
	utils.Bail("failed to get access rights", err)
	if !ar.Permissions.Execute {
		fmt.Printf("You do not have execute rights to repo #%d\n", id)
		os.Exit(10)
	}
	igr := &pb.IDGrantRequest{ID: id, GroupID: "3", Read: true, Write: true}
	_, err = pb.GetObjectAuthServiceClient().GrantToSoftware(ctx, igr)
	utils.Bail("failed to set permissions", err)
	fmt.Printf("Done\n")

}

func getUserID() string {
	return *userid
}

func serviceid(id uint64) string {
	return fmt.Sprintf("%d", id)
}
func AllAccess() error {
	ar := &pb.GrantAllAccessRequest{
		ObjectType:     getObjectType(),
		CallingService: serviceid(*callingservice),
		SubjectService: serviceid(*subjectservice),
		ReadAccess:     *read,
		WriteAccess:    *write,
	}
	ctx := authremote.Context()
	_, err := pb.GetObjectAuthServiceClient().GrantAllServiceAccess(ctx, ar)
	if err != nil {
		return err
	}
	return nil
}


