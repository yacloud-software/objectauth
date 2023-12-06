package composite

import (
	"context"
	pb "golang.conradwood.net/apis/objectauth"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/objectauth/db"
)

// given an id and compositetype, will get permissions for this object
// for example compositetype == software and id == 231 (gitrepo)
func accessrightlist_for(ctx context.Context, t pb.COMPOSITETYPE, id uint64) (*pb.AccessRightList, error) {
	Debugf("looking up %v, id #%d\n", t, id)
	u := auth.GetUser(ctx)
	if u == nil {
		return nil, errors.Unauthenticated(ctx, "login required")
	}
	cdb := db.DefaultDBGroupToComposite() // newer, better way of doing it
	res := &pb.AccessRightList{
		// objectid & type are a bit meaningless here, because the composite type does not match the objectid
		GroupPermissions: make(map[string]*pb.PermissionSet),
		UserPermissions:  make(map[string]*pb.PermissionSet),
	}
	for _, g := range u.Groups {
		cs, err := cdb.ByGroupID(ctx, g.ID)
		if err != nil {
			return nil, err
		}
		p := &pb.PermissionSet{}
		for _, c := range cs {
			if c.ObjectType != t {
				continue
			}
			if c.ObjectID != id {
				continue
			}
			Debugf("Merging permissions for group \"%s\"\n", g.ID)
			p = mergePerm(p, c)
		}
		res.GroupPermissions[g.ID] = p
	}

	return res, nil
}

