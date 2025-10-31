// 获取选股器中的行业列表数据

package eastmoney

import (
	"context"
	"fmt"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/sirupsen/logrus"
)

// RespIndustryList 接口返回的 json 结构
type RespIndustryList struct {
	Result struct {
		Count int `json:"count"`
		Pages int `json:"pages"`
		Data  []struct {
			Industry    string `json:"INDUSTRY"`
			FirstLetter string `json:"FIRST_LETTER"`
		} `json:"data"`
	} `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// QueryIndustryList 获取行业列表
func (e EastMoney) QueryIndustryList(ctx context.Context) ([]string, error) {
	apiurl := "https://datacenter.eastmoney.com/stock/selection/api/data/get/"
	reqData := map[string]string{
		"source": "SELECT_SECURITIES",
		"client": "APP",
		"type":   "RPTA_APP_INDUSTRY",
		"sty":    "ALL",
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{"reqData": reqData}).Debug("EastMoney IndustryList " + apiurl + " begin")
	beginTime := time.Now()
	req, err := goutils.NewHTTPMultipartReq(ctx, apiurl, reqData)
	if err != nil {
		return nil, err
	}
	resp := RespIndustryList{}
	err = goutils.HTTPPOST(ctx, e.HTTPClient, req, &resp)
	latency := time.Now().Sub(beginTime).Milliseconds()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"latency(ms)": latency}).Debug("EastMoney IndustryList " + apiurl + " end")
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("%#v", resp)
	}
	result := []string{}
	for _, i := range resp.Result.Data {
		result = append(result, i.Industry)
	}
	return result, nil
}
