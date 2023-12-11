package composite

import (
	"context"
	pb "golang.conradwood.net/apis/objectauth"
	"golang.conradwood.net/go-easyops/utils"
)

type Software struct {
}

func (s *Software) ForSingleObject(ctx context.Context, req *pb.AuthRequest) (*pb.AccessRightList, error) {
	var err error
	var lookupid uint64
	Debugf("looking up single object %v id #%d\n", req.ObjectType, req.ObjectID)

	found := false
	if req.ObjectType == pb.OBJECTTYPE_GitRepository {
		return accessrightlist_for(ctx, pb.COMPOSITETYPE_Software, req.ObjectID)
	}
	if req.ObjectType == pb.OBJECTTYPE_Artefact {
		lookupid, err = get_git_repo_for_artefact(ctx, req.ObjectID)
		if err != nil {
			Debugf("unable to get git repo for artefact %i: %s\n", req.ObjectID, utils.ErrorString(err))
			return nil, err
		}
		found = true
	}
	if !found {
		return &pb.AccessRightList{}, nil
	}
	return accessrightlist_for(ctx, pb.COMPOSITETYPE_Software, lookupid)
}


