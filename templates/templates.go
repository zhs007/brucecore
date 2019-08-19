package brucetemplates

import (
	"github.com/zhs007/adacore"
	adacorepb "github.com/zhs007/adacore/proto"
	brucecorebase "github.com/zhs007/brucecore/base"
	"github.com/zhs007/brucecore/ipgeo"
	"github.com/zhs007/brucecore/templates/spnormal"
	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

// TempFunc - template function
type TempFunc func(name string, url string, reply *jarviscrawlercore.ReplyAnalyzePage,
	km *adacore.KeywordMappingList, ipgeodb *ipgeo.DB) (*adacorepb.MarkdownData, error)

// TemplatesMgr - templates manager
type TemplatesMgr struct {
	mapTemp map[string]TempFunc
}

// RegTemplate -
func (mgr *TemplatesMgr) RegTemplate(name string, funcTemp TempFunc) error {
	mgr.mapTemp[name] = funcTemp

	return nil
}

// GetTemplate -
func (mgr *TemplatesMgr) GetTemplate(name string) TempFunc {
	f, isok := mgr.mapTemp[name]
	if isok {
		return f
	}

	return nil
}

var mgr TemplatesMgr

func init() {
	mgr.mapTemp = make(map[string]TempFunc)

	mgr.RegTemplate("spnormal", spnormal.GenMarkdown)
}

// GenMarkdown -
func GenMarkdown(tempname string, name string, url string, reply *jarviscrawlercore.ReplyAnalyzePage,
	km *adacore.KeywordMappingList, ipgeodb *ipgeo.DB) (*adacorepb.MarkdownData, error) {

	f := mgr.GetTemplate(tempname)
	if f != nil {
		return f(name, url, reply, km, ipgeodb)
	}

	return nil, brucecorebase.ErrNoTemplate
}
