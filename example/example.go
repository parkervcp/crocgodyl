package main

import (
	"github.com/parkervcp/crocgodyl"
)

type config struct {
	PanelURL string `json:"panel_url"`
	APIToken string `json:"api_token"`
}

func main() {
	crocgodyl.New("http://testing.synahost.com/", "", "8oGY8PHEmdNzsBjzqTPPb94eUFPsYLlocfhG1EU4fMyyhlP3")
}
