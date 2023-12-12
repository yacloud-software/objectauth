package composite

import (
	pb "golang.conradwood.net/apis/objectauth"
)

type embeddedPermissions interface {
	GetRead() bool
	GetWrite() bool
	GetExecute() bool
	GetView() bool
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



