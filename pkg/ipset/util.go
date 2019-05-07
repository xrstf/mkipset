package ipset

import (
	"errors"
	"regexp"
)

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

func ValidateSetName(name string) error {
	if !setnameRegexp.MatchString(name) {
		return errors.New("set names must be alphanumeric and not start with a number")
	}

	if len(name) > 31 {
		return errors.New("set names must be at most 31 characters long")
	}

	return nil
}
