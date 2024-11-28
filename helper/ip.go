package helper

import (
	"embed"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/ua-parser/uap-go/uaparser"
)

//go:embed db/*
var DbFs embed.FS

var IpSearch *xdb.Searcher

type IpParserResult struct {
	Country  string
	Region   string
	Province string
	City     string
	Isp      string
}

func IpParser(ip string) (*IpParserResult, error) {
	if IpSearch == nil {
		file, err := DbFs.ReadFile("db/ip2region.xdb")
		if err != nil {
			return nil, err
		}
		searcher, err := xdb.NewWithBuffer(file)
		if err != nil {
			return nil, err
		}
		IpSearch = searcher
	}

	region, err := IpSearch.SearchByStr(ip)
	if err != nil {
		return nil, err
	}

	split := strings.Split(region, "|")

	return &IpParserResult{
		Country:  split[0],
		Region:   split[1],
		Province: split[2],
		City:     split[3],
		Isp:      split[4],
	}, nil
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
