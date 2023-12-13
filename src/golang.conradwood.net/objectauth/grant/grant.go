package grant

import (
	"context"
	"fmt"
	oa "golang.conradwood.net/apis/objectauth"
	"golang.conradwood.net/apis/registry"
	raa "golang.conradwood.net/apis/rpcaclapi"
	"golang.conradwood.net/go-easyops/client"
)

// groupid is optional, if empty, will grant access to current user
func GrantService(ctx context.Context, servicename string, groupid string) error {
	l := &registry.V2ListRequest{NameMatch: servicename}
	lr, err := client.GetRegistryClient().ListRegistrations(ctx, l)
	if err != nil {
		return err
	}
	if len(lr.Registrations) == 0 {
		return fmt.Errorf("No such service: \"%s\"", servicename)
	}
	sn := lr.Registrations[0]
	//	fmt.Printf("Getting ID for service \"%s\"\n", sn)
	sr, err := raa.GetRPCACLServiceClient().ServiceNameToID(ctx, &raa.ServiceNameRequest{Name: sn.Target.ServiceName})
	if err != nil {
		return err
	}
	fmt.Printf("ID: %d\n", sr.ID)
	if groupid != "" {
		f := &oa.GrantGroupRequest{
			ObjectType: oa.OBJECTTYPE_Service,
			ObjectID:   sr.ID,
			GroupID:    groupid,
			Read:       true,
			Write:      true,
			Execute:    true,
			View:       true,
		}
		_, err = oa.GetObjectAuthServiceClient().GrantToGroup(ctx, f)
	} else {
		gr := &oa.GrantUserRequest{
			ObjectType: oa.OBJECTTYPE_Service,
			ObjectID:   sr.ID,
			Read:       true,
			Write:      true,
			Execute:    true,
			View:       true,
		}
		_, err = oa.GetObjectAuthServiceClient().GrantToMe(ctx, gr)
	}
	if err != nil {
		return err
	}
	fmt.Printf("Granted access to service (groupid=%s) \"%s\"\n", groupid, sn.Target.ServiceName)
	return nil
}





