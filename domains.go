package godaddy

// @see https://developer.godaddy.com/doc/endpoint/domains#/
//
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Contact domain contact info
type Contact struct {
	AddressMailing struct {
		Address1   string `json:"address1"`
		Address2   string `json:"address2"`
		City       string `json:"city"`
		Country    string `json:"country"`
		PostalCode string `json:"postalCode"`
		State      string `json:"state"`
	} `json:"addressMailing"`
	Email        string `json:"email"`
	Fax          string `json:"fax"`
	NameFirst    string `json:"nameFirst"`
	NameLast     string `json:"nameLast"`
	Organization string `json:"organization"`
	Phone        string `json:"phone"`
}

// Domain domain info
type Domain struct {
	AuthCode            string    `json:"authCode"`
	ContactAdmin        Contact   `json:"contactAdmin"`
	ContactBilling      Contact   `json:"contactBilling"`
	ContactRegistrant   Contact   `json:"contactRegistrant"`
	ContactTech         Contact   `json:"contactTech"`
	CreatedAt           time.Time `json:"createdAt"`
	Domain              string    `json:"domain"`
	DomainID            int64     `json:"domainId"`
	ExpirationProtected bool      `json:"expirationProtected"`
	Expires             time.Time `json:"expires"`
	HoldRegistrar       bool      `json:"holdRegistrar"`
	Locked              bool      `json:"locked"`
	NameServers         []string  `json:"nameServers"`
	Privacy             bool      `json:"privacy"` // 隐私保护
	RenewAuto           bool      `json:"renewAuto"`
	RenewDeadline       time.Time `json:"renewDeadline"`
	Renewable           bool      `json:"renewable"`
	Status              string    `json:"status"`
	TransferProtected   bool      `json:"transferProtected"`
}

// ListDomain  Retrieve a list of Domains for the specified Shopper
// GET https://api.godaddy.com/v1/domains
func (c *Context) ListDomain() ([]Domain, error) {
	val := url.Values{}
	val.Set("limit", strconv.Itoa(limit))
	res := []Domain{}
	for {
		URL := fmt.Sprintf("https://%s/v1/domains?%s", c.Host, val.Encode())
		request, err := http.NewRequest(http.MethodGet, URL, nil)
		if err != nil {
			return res, err
		}
		if len(c.XShopperID) != 0 {
			request.Header.Set("X-Shopper-Id", c.XShopperID)
		}
		resItem := []Domain{}
		if err := c.RequestExecute(request, &resItem); err != nil {
			return res, err
		}
		size := len(resItem)
		for _, item := range resItem {
			res = append(res, item)
		}
		if size < limit {
			break
		}
		val.Set("marker", resItem[size-1].Domain)
	}
	return res, nil
}

// GetDomain Retrieve details for the specified Domain
// GET https://api.godaddy.com/v1/domains/{domain}
func (c *Context) GetDomain(domain string) (info Domain, err error) {
	URL := fmt.Sprintf("https://%s/v1/domains/%s", c.Host, domain)
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return info, err
	}
	if len(c.XShopperID) != 0 {
		request.Header.Set("X-Shopper-Id", c.XShopperID)
	}
	if err = c.RequestExecute(request, &info); err != nil {
		return
	}
	return
}

// Record domain record
type Record struct {
	Data       string `json:"data"`
	Name       string `json:"name"`
	TTL        int    `json:"ttl"`                // 缓存时间
	RecordType string `json:"type,omitempty"`     // 记录类型 [A, AAAA, CNAME, MX, NS, SOA, SRV, TXT]
	Port       int    `json:"port,omitempty"`     // 服务端口号, 最小为1 最大为65535 仅限 SRV 类型
	Priority   int    `json:"priority,omitempty"` // 记录优先级, 仅限 MX 和 SRV 类型
	Protocol   string `json:"protocol,omitempty"` // 服务协议, 仅限 SRV 类型
	Service    string `json:"service,omitempty"`  // 服务类型, 仅限 SRV 类型
	Weight     string `json:"weight,omitempty"`   // 记录权重, 仅限 SRV 类型
}

// AddRecords Add the specified DNS Records to the specified Domain
// PATCH https://api.godaddy.com/v1/domains/{domain}/records
func (c *Context) AddRecords(domain string, records []Record) error {
	b, err := json.Marshal(records)
	if err != nil {
		return err
	}
	URL := fmt.Sprintf("https://%s/v1/domains/%s/records", c.Host, domain)
	request, err := http.NewRequest(http.MethodPatch, URL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	return c.RequestExecute(request, nil)
}

// UpdateAllRecords Replace all DNS Records for the specified Domain
// PUT https://api.godaddy.com/v1/domains/{domain}/records
func (c *Context) UpdateAllRecords(domain string, records []Record) error {
	b, err := json.Marshal(records)
	if err != nil {
		return err
	}
	URL := fmt.Sprintf("https://%s/v1/domains/%s/records", c.Host, domain)
	request, err := http.NewRequest(http.MethodPut, URL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	return c.RequestExecute(request, nil)
}

// ListRecords Retrieve DNS Records for the specified Domain, optionally with the specified Type and/or Name
// GET https://api.godaddy.com/v1/domains/{domain}/records/{type}/{name}
func (c *Context) ListRecords(domain string, recordType RecordType, name string) ([]Record, error) {
	res := []Record{}
	val := url.Values{}
	val.Set("limit", strconv.Itoa(limit))
	for offset := 1; true; offset++ {
		val.Set("offset", strconv.Itoa(offset))
		URL := fmt.Sprintf("https://%s/v1/domains/%s/records/%s/%s?%s", c.Host, domain, recordType.String(), name, val.Encode())
		request, err := http.NewRequest(http.MethodGet, URL, nil)
		if err != nil {
			return res, err
		}
		resItem := []Record{}
		if err := c.RequestExecute(request, &resItem); err != nil {
			return res, err
		}
		size := len(resItem)
		for _, item := range resItem {
			res = append(res, item)
		}
		if size < limit {
			break
		}
	}
	return res, nil
}

// UpdateRecordsByName Replace all DNS Records for the specified Domain with the specified Type and Name
// PUT https://api.godaddy.com/v1/domains/{domain}/records/{type}/{name}
func (c *Context) UpdateRecordsByName(domain string, records []Record, recordType RecordType, name string) error {
	b, err := json.Marshal(records)
	if err != nil {
		return err
	}
	URL := fmt.Sprintf("https://%s/v1/domains/%s/records/%s/%s", c.Host, domain, recordType.String(), name)
	request, err := http.NewRequest(http.MethodPut, URL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	return c.RequestExecute(request, nil)
}

// UpdateRecordsByType Replace all DNS Records for the specified Domain with the specified Type
// PUT https://api.godaddy.com/v1/domains/{domain}/records/{type}
func (c *Context) UpdateRecordsByType(domain string, records []Record, recordType RecordType) error {
	b, err := json.Marshal(records)
	if err != nil {
		return err
	}
	URL := fmt.Sprintf("https://%s/v1/domains/%s/records/%s", c.Host, domain, recordType.String())
	request, err := http.NewRequest(http.MethodPut, URL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	return c.RequestExecute(request, nil)
}
