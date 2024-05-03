package helper

import (
	"embed"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/ua-parser/uap-go/uaparser"
)

//go:embed db/*
var DbFs embed.FS

var IpSearch *xdb.Searcher

func IpParser(ip string) (string, error) {
	if IpSearch == nil {
		file, err := DbFs.ReadFile("db/ip2region.xdb")
		if err != nil {
			return "", err
		}
		searcher, err := xdb.NewWithBuffer(file)
		if err != nil {
			return "", err
		}
		IpSearch = searcher
	}

	region, err := IpSearch.SearchByStr(ip)
	if err != nil {
		return "", err
	}

	return region, nil
}

var UapParser *uaparser.Parser

func UaParser(ua string) (*uaparser.Client, error) {
	if UapParser == nil {
		file, err := DbFs.ReadFile("db/regexes.yaml")
		if err != nil {
			return nil, err
		}
		parser, err := uaparser.NewFromBytes(file)
		if err != nil {
			return nil, err
		}
		UapParser = parser
	}
	client := UapParser.Parse(ua)
	return client, nil
}
