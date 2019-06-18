package godaddy

import (
	"fmt"
	"testing"
)

func Test_ListDomain(t *testing.T) {
	c := InitConfig()
	domains, err := c.ListDomain()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	t.Logf("domains size :%d", len(domains))
}

func Test_GetDomain(t *testing.T) {
	c := InitConfig()
	domains, err := c.ListDomain()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	body, err := c.GetDomain(domains[0].Domain)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	fmt.Printf("%+v\n", body)
}

func Test_ListRecords(t *testing.T) {
	c := InitConfig()
	domains, err := c.ListDomain()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	body, err := c.ListRecords(domains[0].Domain, A, "*")
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	fmt.Printf("%+v\n", body)
}
func Test_AddRecords(t *testing.T) {
	c := InitConfig()
	domains, err := c.ListDomain()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	records := []Record{Record{Data: "8.8.8.9", Name: "*", TTL: DefaultTTL, RecordType: "A"}, Record{Data: "114.114.114.115", Name: "*", TTL: DefaultTTL, RecordType: AType}}
	if err = c.AddRecords(domains[0].Domain, records); err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}
func Test_UpdateRecords(t *testing.T) {
	c := InitConfig()
	domains, err := c.ListDomain()
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	records := []Record{Record{Data: "8.8.8.8", Name: "*", TTL: DefaultTTL, RecordType: "A"}, Record{Data: "114.114.114.114", Name: "*", TTL: DefaultTTL, RecordType: AType}}
	if err = c.UpdateRecordsByType(domains[0].Domain, records, A); err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}
