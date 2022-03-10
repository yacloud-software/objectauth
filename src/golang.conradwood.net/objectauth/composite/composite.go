package composite

import (
	"context"
	"fmt"
	pb "golang.conradwood.net/apis/objectauth"
)

type Composite interface {
	// return accessrights based on composite rights for a particular object (identified by ID)
	ForSingleObject(ctx context.Context, req *pb.AuthRequest) (*pb.AccessRightList, error)
}

func GetComposite(t pb.OBJECTTYPE) Composite {
	Debugf("Getting composite resolver for type %v\n", t)
	if t == pb.OBJECTTYPE_GitRepository || t == pb.OBJECTTYPE_Artefact || t == pb.OBJECTTYPE_Proto {
		return &Software{}
	}
	return nil
}

func Debugf(format string, args ...interface{}) {
	fmt.Printf("[composite] "+format, args...)
}
