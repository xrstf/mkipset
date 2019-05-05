package ipset

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"os/exec"
)

type execIpset struct {
}

func NewExec() (*execIpset, error) {
	return &execIpset{}, nil
}

func (e *execIpset) Create(setname string, t SetType) error {
	if !validateSetname(setname) {
		return fmt.Errorf("invalid set name: '%s'", setname)
	}

	_, err := e.run([]string{"create", setname, string(t), "-exist"})
	if err != nil {
		return err
	}

	return nil
}

func (e *execIpset) Add(setname string, entry string) error {
	if !validateSetname(setname) {
		return fmt.Errorf("invalid set name: '%s'", setname)
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
	if !validateSetname(setname) {
		return fmt.Errorf("invalid set name: '%s'", setname)
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
	if !validateSetname(setname) {
		return nil, fmt.Errorf("invalid set name: '%s'", setname)
	}

	output, err := e.run([]string{"list", setname})
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
		Type:     ipset.Type,
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
	if !validateSetname(oldname) {
		return fmt.Errorf("invalid old set name: '%s'", oldname)
	}

	if !validateSetname(newname) {
		return fmt.Errorf("invalid new set name: '%s'", newname)
	}

	_, err := e.run([]string{"swap", oldname, newname})
	if err != nil {
		return err
	}

	return nil
}

func (e *execIpset) Destroy(setname string) error {
	if !validateSetname(setname) {
		return fmt.Errorf("invalid set name: '%s'", setname)
	}

	_, err := e.run([]string{"destroy", setname})
	if err != nil {
		return err
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
