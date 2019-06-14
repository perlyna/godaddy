package godaddy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Domain domain info
type Domain struct {
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
	Renewable           bool      `json:"renewable"`
	Status              string    `json:"status"`
	TransferProtected   bool      `json:"transferProtected"`
}

// ListDomain  Retrieve a list of Domains for the specified Shopper
// https://api.godaddy.com/v1/domains
func (c *Context) ListDomain() ([]Domain, error) {
	val := url.Values{}
	val.Set("limit", strconv.Itoa(limit))
	res := []Domain{}
	for {
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/v1/domains?%s", c.Host, val.Encode()), nil)
		if err != nil {
			return res, err
		}
		resItem := []Domain{}
		if err := c.RequestExecute("", request, &resItem); err != nil {
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

// ContactInfo domain contact info
type ContactInfo struct {
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

// DomainInfo domain detail info
type DomainInfo struct {
	Domain
	AuthCode          string      `json:"authCode"`
	RenewDeadline     time.Time   `json:"renewDeadline"`
	ContactAdmin      ContactInfo `json:"contactAdmin"`
	ContactBilling    ContactInfo `json:"contactBilling"`
	ContactRegistrant ContactInfo `json:"contactRegistrant"`
	ContactTech       ContactInfo `json:"contactTech"`
}

// GetDomain Retrieve details for the specified Domain
// https://api.godaddy.com/v1/domains/{domain}
func (c *Context) GetDomain(domain string) (info DomainInfo, err error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/v1/domains/%s", c.Host, domain), nil)
	if err != nil {
		return info, err
	}
	if err = c.RequestExecute("", request, &info); err != nil {
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

const (
	// AType A type
	AType = "A"
	// AAAAType AAAA type
	AAAAType = "AAAA"
	// CNameType CName type
	CNameType = "CNAME"
	// MXType MX type
	MXType = "MX"
	// NSType NS type
	NSType = "NS"
	// SOAType SOA type
	SOAType = "SOA"
	// TXTType TXT type
	TXTType = "TXT"
)

func (rt RecordType) String() string {
	switch rt {
	case A:
		return AType
	case AAAA:
		return AAAAType
	case CNAME:
		return CNameType
	case MX:
		return MXType
	case NS:
		return NSType
	case SOA:
		return SOAType
	case TXT:
		return TXTType
	}
	return ""
}

// ListRecords Retrieve DNS Records for the specified Domain, optionally with the specified Type and/or Name
// https://api.godaddy.com/v1/domains/{domain}/records/{type}/{name}
func (c *Context) ListRecords(domain string, recordType RecordType, name string) ([]Record, error) {
	res := []Record{}
	val := url.Values{}
	val.Set("limit", strconv.Itoa(limit))
	for offset := 1; true; offset++ {
		val.Set("offset", strconv.Itoa(offset))
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/v1/domains/%s/records/%s/%s?%s", c.Host, domain, recordType.String(), name, val.Encode()), nil)
		if err != nil {
			return res, err
		}
		resItem := []Record{}
		if err := c.RequestExecute("", request, &resItem); err != nil {
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

// AddRecords Add the specified DNS Records to the specified Domain
// https://api.godaddy.com/v1/domains/{domain}/records
func (c *Context) AddRecords(domain string, records []Record) error {
	b, err := json.Marshal(records)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("https://%s/v1/domains/%s/records", c.Host, domain), bytes.NewReader(b))
	if err != nil {
		return err
	}
	return c.RequestExecute("", request, nil)
}

// UpdateRecords Replace all DNS Records for the specified Domain with the specified Type
// https://api.godaddy.com/v1/domains/{domain}/records/{type}
func (c *Context) UpdateRecords(domain string, recordType RecordType, records []Record) error {
	b, err := json.Marshal(records)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", b)
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://%s/v1/domains/%s/records/%s", c.Host, domain, recordType.String()), bytes.NewReader(b))
	if err != nil {
		return err
	}
	return c.RequestExecute("", request, nil)
}
