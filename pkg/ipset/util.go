package ipset

import "regexp"

type SetType string

const (
	SetTypeBitmapIP       SetType = "bitmap:ip"
	SetTypeBitmapIPMAC    SetType = "bitmap:ip,mac"
	SetTypeBitmapPort     SetType = "bitmap:port"
	SetTypeHashIP         SetType = "hash:ip"
	SetTypeHashMAC        SetType = "hash:mac"
	SetTypeHashNet        SetType = "hash:net"
	SetTypeHashNetNet     SetType = "hash:net,net"
	SetTypeHashIPPort     SetType = "hash:ip,port"
	SetTypeHashNetPort    SetType = "hash:net,port"
	SetTypeHashIPPortIP   SetType = "hash:ip,port,ip"
	SetTypeHashIPPortNet  SetType = "hash:ip,port,net"
	SetTypeHashIPMark     SetType = "hash:ip,mark"
	SetTypeHashNetPortNet SetType = "hash:net,port,net"
	SetTypeHashNetIface   SetType = "hash:net,iface"
	SetTypeListSet        SetType = "list:set"
)

var setnameRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

func validateSetname(name string) bool {
	return setnameRegexp.MatchString(name)
}
