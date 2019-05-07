package ipset

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/xrstf/mkipset/pkg/ip"
)

type execIpset struct {
}

func NewExec() (*execIpset, error) {
	return &execIpset{}, nil
}

func (e *execIpset) Create(setname string, t SetType) error {
	if err := ValidateSetName(setname); err != nil {
		return fmt.Errorf("invalid set name: %v", err)
	}

	_, err := e.run([]string{"create", setname, string(t), "-exist"})
	if err != nil {
		return err
	}

	return nil
}

func (e *execIpset) Add(setname string, entry string) error {
	if err := ValidateSetName(setname); err != nil {
		return fmt.Errorf("invalid set name: %v", err)
	}

	if len(entry) == 0 {
		return errors.New("no entry given")
	}

	_, err := e.run([]string{"add", setname, entry})
	if err != nil {
		return err
	}

	return nil
}

func (e *execIpset) Delete(setname string, entry string) error {
	if err := ValidateSetName(setname); err != nil {
		return fmt.Errorf("invalid set name: %v", err)
	}

	if len(entry) == 0 {
		return errors.New("no entry given")
	}

	_, err := e.run([]string{"delete", setname, entry})
	if err != nil {
		return err
	}

	return nil
}

type ipsetsXML struct {
	IPSets []struct {
		Name     string `xml:"name,attr"`
		Type     string `xml:"type"`
		Revision int    `xml:"revision"`
		Header   struct {
			Family      string `xml:"family"`
			HashSize    int    `xml:"hashsize"`
			MaxElements int    `xml:"maxelem"`
			MemorySize  int    `xml:"memsize"`
			References  int    `xml:"references"`
		} `xml:"header"`
		Members struct {
			Member []struct {
				Elem string `xml:"elem"`
			} `xml:"member"`
		} `xml:"members"`
	} `xml:"ipset"`
}

func (e *execIpset) Show(setname string) (*Set, error) {
	if err := ValidateSetName(setname); err != nil {
		return nil, fmt.Errorf("invalid set n%v '%s'", err)
	}

	output, err := e.run([]string{"list", setname, "-o", "xml"})
	if err != nil {
		return nil, err
	}

	var data ipsetsXML
	if err := xml.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("failed to parse ipset XML: %v", err)
	}

	if len(data.IPSets) == 0 {
		return nil, errors.New("could not find set")
	}

	ipset := data.IPSets[0]

	s := &Set{
		Name:     ipset.Name,
		Type:     SetType(ipset.Type),
		Revision: ipset.Revision,
		Header: SetHeader{
			Family:      ipset.Header.Family,
			HashSize:    ipset.Header.HashSize,
			MaxElements: ipset.Header.MaxElements,
			MemorySize:  ipset.Header.MemorySize,
			References:  ipset.Header.References,
		},
		Members: make([]string, 0),
	}

	for _, member := range ipset.Members.Member {
		s.Members = append(s.Members, member.Elem)
	}

	return s, nil
}

func (e *execIpset) Swap(oldname string, newname string) error {
	if err := ValidateSetName(oldname); err != nil {
		return fmt.Errorf("invalid old set name: %v", err)
	}

	if err := ValidateSetName(newname); err != nil {
		return fmt.Errorf("invalid new set name: %v", err)
	}

	_, err := e.run([]string{"swap", oldname, newname})
	if err != nil {
		return err
	}

	return nil
}

func (e *execIpset) Destroy(setname string) error {
	if err := ValidateSetName(setname); err != nil {
		return fmt.Errorf("invalid set name: %v", err)
	}

	_, err := e.run([]string{"destroy", setname})
	if err != nil {
		return err
	}

	return nil
}

func (e *execIpset) Synchronize(set Set, ips ip.Slice) error {
	// remember that names cannot be longer tan 31 characters
	// and set names in mkipset are allowed to be up to 20 chars
	tempSetName := fmt.Sprintf("%s_%s", set.Name, time.Now().Format("150405.999"))
	tempSetName = strings.Replace(tempSetName, ".", "_", -1)

	err := e.Create(tempSetName, set.Type)
	if err != nil {
		return fmt.Errorf("failed to create temp set: %v", err)
	}
	defer func() {
		e.Destroy(tempSetName)
	}()

	for _, ip := range ips {
		err := e.Add(tempSetName, ip.String())
		if err != nil {
			return fmt.Errorf("failed to add %v to temp set: %v", ip, err)
		}
	}

	err = e.Swap(set.Name, tempSetName)
	if err != nil {
		return fmt.Errorf("failed to enable new temp set: %v", err)
	}

	return nil
}

func (e *execIpset) run(args []string) ([]byte, error) {
	command := exec.Command("ipset", args...)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	return stdout.Bytes(), nil
}
