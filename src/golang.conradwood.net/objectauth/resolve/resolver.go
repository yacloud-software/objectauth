package resolve

import (
    "golang.conradwood.net/go-easyops/authremote"
	"fmt"
	af "golang.conradwood.net/apis/artefact"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/objectauth"
	ra "golang.conradwood.net/apis/rpcaclapi"
	"golang.conradwood.net/go-easyops/utils"
	"strconv"
	"strings"
)

func ResolveService(s string) (uint64, string, error) {
	ctx := authremote.Context()
	rpcapi := ra.GetRPCACLServiceClient()
	sir, err := rpcapi.ServiceNameToID(ctx, &ra.ServiceNameRequest{Name: s})
	if err != nil {
		fmt.Printf("ServiceNameToID(%s) failed: %s\n", s, utils.ErrorString(err))
		return 0, "", err
	}
	fmt.Printf("ServiceNameToID(%s) == %d\n", s, sir.ID)
	return sir.ID, s, err

	/*
		svs, err := rpcapi.GetServices(ctx, &common.Void{})
		if err != nil {
			return 0, "", err
		}
		matched := false
		var sr *ra.Service
		for _, sv := range svs.Services {
			if strings.Contains(sv.Name, s) {
				if matched {
					return 0, "", fmt.Errorf("multiple matches for servicename")
				}
				matched = true
				sr = sv
			}
		}
		if !matched {
			return 0, "", fmt.Errorf("no match for servicename")
		}
		return sr.ID, sr.Name, nil
	*/
}
func ResolveArtefact(s string) (uint64, string, error) {
	ctx := authremote.Context()
	afClient := af.GetArtefactServiceClient()
	afs, err := afClient.List(ctx, &common.Void{})
	if err != nil {
		return 0, "", err
	}
	var res *af.Contents
	for _, b := range afs.GetArtefacts() {
		if strings.Contains(b.Name, s) {
			if res != nil {
				return 0, "", fmt.Errorf("multiple matches for artefact name")
			}
			res = b
		}
	}
	if res == nil {
		return 0, "", fmt.Errorf("no match for artefact name")
	}
	if res.ArtefactID == nil {
		return 0, "", fmt.Errorf("artefact has no id")
	}
	return res.ArtefactID.ID, res.Name, nil
}

// resolves a string, e.g. a number or a servicename to a number
// ID, some name, or error
func ResolveToNumber(t pb.OBJECTTYPE, s string) (uint64, string, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		return id, fmt.Sprintf(" #%d", id), nil
	}
	if t == pb.OBJECTTYPE_Service {
		return ResolveService(s)
	}
	if t == pb.OBJECTTYPE_Artefact {
		return ResolveArtefact(s)
	}

	return 0, "", fmt.Errorf("Unsure how to convert \"%v\" into a number\n", t)

}





