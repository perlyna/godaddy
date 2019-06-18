package godaddy

// record type

// RecordType is an enumeration of possible DNS record types
type RecordType int

const (
	// A is an address record type
	A RecordType = iota
	// AAAA is an IPv6 address record type
	AAAA
	// CNAME is a Canonical record name (alias) type
	CNAME
	// MX is a mail exchange record type
	MX
	// NS is a name server record type
	NS
	// SOA is a start of authority record type
	SOA
	// SRV is a service locator type
	SRV
	// TXT is a text record type
	TXT
)
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
