package extraction

import (
	"net"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
)

type Extractor interface {
	Extract(candidate string) []Asset
}

type Asset struct {
	Type string
	ID   string
}

func ListMatches(candidate string) []Asset {
	assets := []Asset{}

	for _, extractor := range supportedAssets {
		assets = append(assets, extractor.Extract(candidate)...)
	}

	return assets
}

//nolint:gochecknoglobals
var supportedAssets = []Extractor{
	&IP{
		version: "4",
		regex:   regexp.MustCompile(`(\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`), //nolint:lll
	},
	&IP{
		version: "6",
		regex:   regexp.MustCompile(`(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`), //nolint:lll
	},
	&CIDR{
		version: "4",
		regex:   regexp.MustCompile(`(\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`), //nolint:lll
	},
	&CIDR{
		version: "6",
		regex:   regexp.MustCompile(`(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`), //nolint:lll
	},
	&ARN{
		regex: regexp.MustCompile(`arn:(?P<Partition>[^:\n]*):(?P<Service>[^:\n]*):(?P<Region>[^:\n]*):(?P<AccountID>[^:\n]*):(?P<Ignore>(?P<ResourceType>[^:\/\n]*)[:\/])?(?P<Resource>.*)`), //nolint:lll
	},
	&Regex{
		_type: "email",
		regex: regexp.MustCompile(`([\w\-_.]*[^.])(@\w+)(\.\w+(\.\w+)?[^.\W])`),
	},
}

type IP struct {
	regex   *regexp.Regexp
	version string
}

func (i *IP) Extract(candidate string) []Asset {
	assets := []Asset{}

	for _, match := range i.regex.FindAllString(candidate, -1) {
		if ip := net.ParseIP(match); ip != nil {
			assets = append(assets, Asset{Type: "net.ip." + i.version, ID: match})
		}
	}

	return assets
}

type MAC struct {
	regex *regexp.Regexp
}

func (m *MAC) Extract(candidate string) []Asset {
	assets := []Asset{}

	for _, match := range m.regex.FindAllString(candidate, -1) {
		if _, err := net.ParseMAC(match); err != nil {
			continue
		}

		assets = append(assets, Asset{Type: "net.mac", ID: match})
	}

	return assets
}

type CIDR struct {
	regex   *regexp.Regexp
	version string
}

func (c *CIDR) Extract(candidate string) []Asset {
	assets := []Asset{}

	for _, match := range c.regex.FindAllString(candidate, -1) {
		_, _, err := net.ParseCIDR(match)
		if err != nil {
			continue
		}

		assets = append(assets, Asset{Type: "net.cidr." + c.version, ID: match})
	}

	return assets
}

type ARN struct {
	regex *regexp.Regexp
}

func (a *ARN) Extract(candidate string) []Asset {
	assets := []Asset{}

	for _, match := range a.regex.FindAllString(candidate, -1) {
		arn, err := arn.Parse(match)
		if err != nil {
			continue
		}

		assets = append(assets, Asset{Type: "aws." + arn.Service, ID: match})
	}

	return assets
}

type Regex struct {
	regex *regexp.Regexp
	_type string
}

func (r *Regex) Extract(candidate string) []Asset {
	assets := []Asset{}

	for _, match := range r.regex.FindAllString(candidate, -1) {
		assets = append(assets, Asset{Type: r._type, ID: match})
	}

	return assets
}
