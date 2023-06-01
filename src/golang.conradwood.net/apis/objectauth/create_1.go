// client create: ObjectAuthServiceClient
/*
  Created by /srv/home/cnw/devel/go/go-tools/src/golang.conradwood.net/gotools/protoc-gen-cnw/protoc-gen-cnw.go
*/

/* geninfo:
   filename  : protos/golang.conradwood.net/apis/objectauth/objectauth.proto
   gopackage : golang.conradwood.net/apis/objectauth
   importname: ai_0
   clientfunc: GetObjectAuthService
   serverfunc: NewObjectAuthService
   lookupfunc: ObjectAuthServiceLookupID
   varname   : client_ObjectAuthServiceClient_0
   clientname: ObjectAuthServiceClient
   servername: ObjectAuthServiceServer
   gscvname  : objectauth.ObjectAuthService
   lockname  : lock_ObjectAuthServiceClient_0
   activename: active_ObjectAuthServiceClient_0
*/

package objectauth

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ObjectAuthServiceClient_0 sync.Mutex
  client_ObjectAuthServiceClient_0 ObjectAuthServiceClient
)

func GetObjectAuthClient() ObjectAuthServiceClient { 
    if client_ObjectAuthServiceClient_0 != nil {
        return client_ObjectAuthServiceClient_0
    }

    lock_ObjectAuthServiceClient_0.Lock() 
    if client_ObjectAuthServiceClient_0 != nil {
       lock_ObjectAuthServiceClient_0.Unlock()
       return client_ObjectAuthServiceClient_0
    }

    client_ObjectAuthServiceClient_0 = NewObjectAuthServiceClient(client.Connect(ObjectAuthServiceLookupID()))
    lock_ObjectAuthServiceClient_0.Unlock()
    return client_ObjectAuthServiceClient_0
}

func GetObjectAuthServiceClient() ObjectAuthServiceClient { 
    if client_ObjectAuthServiceClient_0 != nil {
        return client_ObjectAuthServiceClient_0
    }

    lock_ObjectAuthServiceClient_0.Lock() 
    if client_ObjectAuthServiceClient_0 != nil {
       lock_ObjectAuthServiceClient_0.Unlock()
       return client_ObjectAuthServiceClient_0
    }

    client_ObjectAuthServiceClient_0 = NewObjectAuthServiceClient(client.Connect(ObjectAuthServiceLookupID()))
    lock_ObjectAuthServiceClient_0.Unlock()
    return client_ObjectAuthServiceClient_0
}

func ObjectAuthServiceLookupID() string { return "objectauth.ObjectAuthService" } // returns the ID suitable for lookup in the registry. treat as opaque, subject to change.
