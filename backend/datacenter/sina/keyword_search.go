// 关键词搜索

package sina

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// SearchResult 搜索结果
type SearchResult struct {
	// 数字代码
	SecurityCode string
	// 带后缀的代码
	Secucode string
	// 股票名称
	Name string
	// 股市类型: 11=A股 31=港股 41=美股 103=英股
	Market int
}

// KeywordSearch 关键词搜索， 股票、代码、拼音
func (s Sina) KeywordSearch(ctx context.Context, kw string) (results []SearchResult, err error) {
	apiurl := fmt.Sprintf("https://suggest3.sinajs.cn/suggest/key=%s", kw)
	logrus.WithContext(ctx).Debug("Sina KeywordSearch " + apiurl + " begin")
	beginTime := time.Now()
	resp, err := goutils.HTTPGETRaw(ctx, s.HTTPClient, apiurl, nil)
	latency := time.Now().Sub(beginTime).Milliseconds()
	logrus.WithContext(ctx).WithFields(logrus.Fields{"latency(ms)": latency, "resp": string(resp)}).Debug("Sina KeywordSearch " + apiurl + " end")
	if err != nil {
		return nil, err
	}

	trans := transform.NewReader(bytes.NewReader(resp), simplifiedchinese.GBK.NewDecoder())
	utf8resp, err := io.ReadAll(trans)
	if err != nil {
		logrus.WithContext(ctx).Error("transform ReadAll error:" + err.Error())
	}
	ds := strings.Split(string(utf8resp), "=")
	if len(ds) != 2 {
		return nil, errors.New("search resp invalid:" + string(utf8resp))
	}
	data := strings.Trim(ds[1], `"`)
	for _, line := range strings.Split(data, ";") {
		lineitems := strings.Split(line, ",")
		if len(lineitems) < 9 {
			continue
		}
		market, err := strconv.Atoi(lineitems[1])
		if err != nil {
			logrus.WithContext(ctx).Errorf("market:%s atoi error:%v", lineitems[1], err)
		}
		secucode := lineitems[3][2:] + "." + lineitems[3][:2]
		result := SearchResult{
			SecurityCode: lineitems[2],
			Secucode:     secucode,
			Name:         lineitems[6],
			Market:       market,
		}
		results = append(results, result)
	}
	// 按股市编号排序确保A股在前面
	sort.Slice(results, func(i, j int) bool {
		return results[i].Market < results[j].Market
	})
	return
}
