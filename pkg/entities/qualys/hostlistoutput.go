package qualys

import (
	"encoding/xml"
	"time"
)

type HOSTLISTOUTPUT struct {
	XMLName  xml.Name `xml:"HOST_LIST_OUTPUT"`
	Text     string   `xml:",chardata"`
	RESPONSE struct {
		Text     string `xml:",chardata"`
		DATETIME struct {
			Text string `xml:",chardata"`
		} `xml:"DATETIME"`
		HOSTLIST struct {
			Text string `xml:",chardata"`
			HOST []HOST `xml:"HOST"`
		} `xml:"HOST_LIST"`
	} `xml:"RESPONSE"`
}

type ADDIPOUTPUT struct {
	XMLName  xml.Name `xml:"SIMPLE_RETURN"`
	Text     string   `xml:",chardata"`
	RESPONSE struct {
		Text     string `xml:",chardata"`
		DATETIME struct {
			Text string `xml:",chardata"`
		} `xml:"DATETIME"`
	} `xml:"RESPONSE"`
}

type ADDIPOUTPUTJSON struct {
	Text     string   `json:"text"`
	IsPaid   bool     `json:"isPaid"`
	RESPONSE RESPONSE `json:"response"`
}

type RESPONSE struct {
	Text     string    `json:"text"`
	DATETIME time.Time `json:"createdAt"`
}
