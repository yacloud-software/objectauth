package shared

import (
	"fmt"
	pb "golang.conradwood.net/apis/objectauth"
	"golang.conradwood.net/go-easyops/utils"
)

func PrintAccessRightList(arl *pb.AccessRightList) {
	PrintRequestResponse(&pb.ResolveResponse{AccessRightList: arl})
}
func PrintRequestResponse(rr *pb.ResolveResponse) {
	arl := rr.AccessRightList
	t := utils.Table{}
	gmap := make(map[string]string)
	if rr.User != nil {
		for _, g := range rr.User.Groups {
			gmap[g.ID] = g.Name
		}
	}
	t.AddHeaders(" ", "Execute", "View", "Read", "Write")
	for k, u := range arl.UserPermissions {
		t.AddString(fmt.Sprintf("User \"%s\"", k))
		t.AddBool(u.Execute)
		t.AddBool(u.View)
		t.AddBool(u.Read)
		t.AddBool(u.Write)
		t.NewRow()
	}
	for k, u := range arl.GroupPermissions {
		t.AddString(fmt.Sprintf("Group \"%s\" (%s)", k, gmap[k]))
		t.AddBool(u.Execute)
		t.AddBool(u.View)
		t.AddBool(u.Read)
		t.AddBool(u.Write)
		t.NewRow()

	}
	u := arl.EffectivePermissions
	t.AddString("Effective")
	if u == nil {
		t.AddString("not set")
	} else {
		t.AddBool(u.Execute)
		t.AddBool(u.View)
		t.AddBool(u.Read)
		t.AddBool(u.Write)
	}
	t.NewRow()

	fmt.Println(t.ToPrettyString())

}


