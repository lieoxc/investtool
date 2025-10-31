// 关键词搜索

package eastmoney

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/corpix/uarand"
	"github.com/sirupsen/logrus"
)

// SearchFundInfo 关键词搜索基金结构
type SearchFundInfo struct {
	Code string
	Name string
	Type string
}

// SearchFund 关键词搜索， 股票、代码、拼音
func (e EastMoney) SearchFund(ctx context.Context, kw string) (results []SearchFundInfo, err error) {
	count := 10
	apiurl := fmt.Sprintf("https://fundsuggest.eastmoney.com/FundCodeNew.aspx?input=%s&count=%d&cb=x", kw, count)
	logrus.WithContext(ctx).Debug("EastMoney SearchFund " + apiurl + " begin")
	beginTime := time.Now()
	header := map[string]string{
		"user-agent": uarand.GetRandom(),
	}
	resp, err := goutils.HTTPGETRaw(ctx, e.HTTPClient, apiurl, header)
	strresp := string(resp)
	latency := time.Now().Sub(beginTime).Milliseconds()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"latency(ms)": latency}).Debug("EastMoney SearchFund " + apiurl + " end")
	if err != nil {
		return nil, err
	}

	if len(strresp) < 6 {
		logrus.WithContext(ctx).Warnf("SearchFund invalid resp: %s", strresp)
		return nil, fmt.Errorf("无法找到相关基金")
	}
	reg, err := regexp.Compile(`"(?P<code>\d{6}),.+?,(?P<name>.+?),(?P<type>.+?),"`)
	if err != nil {
		logrus.WithContext(ctx).Error("regexp error:" + err.Error())
		return nil, err
	}
	matched := reg.FindAllStringSubmatch(strresp, -1)
	for _, m := range matched {
		results = append(results, SearchFundInfo{
			Code: m[1],
			Name: m[2],
			Type: m[3],
		})
	}
	return
}
