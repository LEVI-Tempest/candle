// 港股通 100 股 - 拉取最近 30 天数据并生成 HTML
// 港股通100股 - Fetch 30 days data and generate HTML
package main

import (
	"fmt"
	"html"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/LEVI-Tempest/Candle/pkg/datasource"
	v1 "github.com/LEVI-Tempest/Candle/pkg/proto"
)

// HKStock 港股通股票
type HKStock struct {
	Rank int
	Name string
	Code string
}

// HKResult 拉取结果
type HKResult struct {
	Stock   HKStock
	Candles []*v1.Candlestick
	Err     error
}

var hk100 = []HKStock{
	{1, "腾讯控股", "00700"}, {2, "阿里巴巴-W", "09988"}, {3, "美团-W", "03690"},
	{4, "建设银行", "00939"}, {5, "小米集团-W", "01810"}, {6, "中国移动", "00941"},
	{7, "中芯国际", "00981"}, {8, "快手-W", "01024"}, {9, "中国海洋石油", "00883"},
	{10, "比亚迪股份", "01211"}, {11, "工商银行", "01398"}, {12, "汇丰控股", "00005"},
	{13, "中国银行", "03988"}, {14, "京东集团-SW", "09618"}, {15, "友邦保险", "01299"},
	{16, "香港交易所", "00388"}, {17, "中国平安", "02318"}, {18, "招商银行", "03968"},
	{19, "药明康德", "02359"}, {20, "李宁", "02331"}, {21, "安踏体育", "02020"},
	{22, "联想集团", "00992"}, {23, "哔哩哔哩-W", "09626"}, {24, "网易-S", "09999"},
	{25, "百度集团-SW", "09888"}, {26, "京东健康", "06618"}, {27, "商汤-W", "00020"},
	{28, "优必选", "09880"}, {29, "招金矿业", "01818"}, {30, "建滔积层板", "01888"},
	{31, "阜博集团", "03738"}, {32, "中国神华", "01088"}, {33, "中海油服", "02883"},
	{34, "中国石化", "00386"}, {35, "中国石油股份", "00857"}, {36, "紫金矿业", "02899"},
	{37, "中国电信", "00728"}, {38, "中国联通", "00762"}, {39, "华润啤酒", "00291"},
	{40, "舜宇光学科技", "02382"}, {41, "瑞声科技", "02018"}, {42, "海尔智家", "06690"},
	{43, "中国生物制药", "01177"}, {44, "石药集团", "01093"}, {45, "金斯瑞生物科技", "01548"},
	{46, "信达生物", "01801"}, {47, "药明生物", "02269"}, {48, "碧桂园服务", "06098"},
	{49, "龙湖集团", "00960"}, {50, "华润置地", "01109"}, {51, "中国海外发展", "00688"},
	{52, "长城汽车", "02333"}, {53, "吉利汽车", "00175"}, {54, "理想汽车-W", "02015"},
	{55, "小鹏汽车-W", "09868"}, {56, "蔚来-SW", "09866"}, {57, "零跑汽车", "09863"},
	{58, "申洲国际", "02313"}, {59, "创科实业", "00669"}, {60, "新奥能源", "02688"},
	{61, "华润燃气", "01193"}, {62, "中国燃气", "00384"}, {63, "恒安国际", "01044"},
	{64, "蒙牛乳业", "02319"}, {65, "百威亚太", "01876"}, {66, "农夫山泉", "09633"},
	{67, "周大福", "01929"}, {68, "泡泡玛特", "09992"}, {69, "海底捞", "06862"},
	{70, "呷哺呷哺", "00520"}, {71, "阿里健康", "00241"}, {72, "平安好医生", "01833"},
	{73, "微创医疗", "00853"}, {74, "绿叶制药", "02186"}, {75, "中国太保", "02601"},
	{76, "中国人寿", "02628"}, {77, "交通银行", "03328"}, {78, "中信银行", "00998"},
	{79, "光大银行", "06818"}, {80, "邮储银行", "01658"}, {81, "中金公司", "03908"},
	{82, "中信证券", "06030"}, {83, "华泰证券", "06886"}, {84, "银河娱乐", "00027"},
	{85, "金沙中国", "01928"}, {86, "澳博控股", "00880"}, {87, "新鸿基地产", "00016"},
	{88, "恒基地产", "00012"}, {89, "长实集团", "01113"}, {90, "恒隆地产", "00101"},
	{91, "中银香港", "02388"}, {92, "恒生银行", "00011"}, {93, "渣打集团", "02888"},
	{94, "思摩尔国际", "06969"}, {95, "中教控股", "00839"}, {96, "中国民航信息网络", "00696"},
	{97, "金蝶国际", "00268"}, {98, "万国数据-SW", "09698"},
}

func main() {
	output := "港股通100股_30日.html"
	if len(os.Args) > 1 {
		output = os.Args[1]
	}
	token := "demo"
	if t := os.Getenv("TSANGHI_TOKEN"); t != "" {
		token = t
	}

	fmt.Println("🕯️  港股通 100 股 - 最近 30 天数据")
	fmt.Println("=====================================")
	fmt.Printf("数据源: Tsanghi(XHKG) + East Money\n")
	fmt.Printf("输出: %s\n\n", output)

	tsanghi := datasource.NewTsanghiClient(token)
	eastmoney := datasource.NewEastMoneyKlineClient()

	results := make([]HKResult, 0, len(hk100))
	for i, s := range hk100 {
		var candles []*v1.Candlestick
		var err error
		// 1. 先试 Tsanghi XHKG
		candles, err = tsanghi.Fetch(datasource.TsanghiXHKG, s.Code, &datasource.FetchOptions{Limit: 30, Order: 2})
		if err != nil || len(candles) == 0 {
			// 2. 备选 East Money
			candles, err = eastmoney.FetchHK(s.Code, 30)
		}
		results = append(results, HKResult{Stock: s, Candles: candles, Err: err})
		if (i+1)%20 == 0 {
			fmt.Printf("  已拉取 %d/%d\n", i+1, len(hk100))
		}
		time.Sleep(100 * time.Millisecond)
	}

	got := 0
	for _, r := range results {
		if len(r.Candles) > 0 {
			got++
		}
	}
	fmt.Printf("\n📊 成功获取 %d/100 只股票数据\n", got)

	// 若全部失败，用示例数据填充前 5 只，便于查看 HTML 结构
	if got == 0 {
		fmt.Println("⚠️  未获取到数据，使用示例数据生成前 5 只股票的图表结构")
		baseTime := time.Now().AddDate(0, 0, -30)
		for i := 0; i < 5 && i < len(results); i++ {
			candles := make([]*v1.Candlestick, 30)
			price := 100.0 + float64(i)*20
			for j := 0; j < 30; j++ {
				change := (float64(j%7) - 3) * 2
				candles[29-j] = &v1.Candlestick{
					Timestamp: baseTime.AddDate(0, 0, j).Unix(),
					Open:      price,
					High:      price + 3,
					Low:       price - 2,
					Close:     price + change*0.5,
					Volume:    1000000,
				}
				price = candles[29-j].Close
			}
			results[i].Candles = candles
		}
	}

	htmlPath, _ := filepath.Abs(output)
	if err := generateHTML(htmlPath, results); err != nil {
		log.Fatalf("❌ 生成 HTML 失败: %v", err)
	}
	fmt.Printf("✅ HTML 已保存: %s\n", htmlPath)
}

func generateHTML(path string, results []HKResult) error {
	// 按数据量排序，有数据的放前面
	sort.Slice(results, func(i, j int) bool {
		a, b := len(results[i].Candles), len(results[j].Candles)
		if a != b {
			return a > b
		}
		return results[i].Stock.Rank < results[j].Stock.Rank
	})

	var tableRows, chartScripts string
	chartCount := 0
	for _, r := range results {
		row := fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td>",
			r.Stock.Rank, html.EscapeString(r.Stock.Name), r.Stock.Code)
		if len(r.Candles) > 0 {
			last := r.Candles[0]
			first := r.Candles[len(r.Candles)-1]
			pct := 0.0
			if first.Close > 0 {
				pct = (last.Close - first.Close) / first.Close * 100
			}
			row += fmt.Sprintf("<td>%.2f</td><td>%.2f%%</td><td>%.2f</td><td>%.2f</td><td>✅</td>",
				last.Close, pct, last.High, last.Low)
			if chartCount < 10 {
				chartScripts += makeChartScript(r, chartCount)
				chartCount++
			}
		} else {
			row += "<td>-</td><td>-</td><td>-</td><td>-</td><td>❌</td>"
		}
		tableRows += row + "</tr>\n"
	}

	tpl := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<title>港股通 100 股 - 最近 30 天数据</title>
<script src="https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js"></script>
<style>
body{font-family:system-ui,-apple-system,sans-serif;margin:20px;background:#f5f5f5;}
h1{color:#333;}
table{border-collapse:collapse;background:#fff;box-shadow:0 2px 8px rgba(0,0,0,.1);border-radius:8px;overflow:hidden;}
th,td{padding:10px 14px;text-align:left;border-bottom:1px solid #eee;}
th{background:#1a73e8;color:#fff;}
tr:hover{background:#f8f9fa;}
.up{color:#e53935;}
.down{color:#43a047;}
.chart{margin:24px 0;padding:16px;background:#fff;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,.1);}
.chart h3{color:#333;margin-top:0;}
.chart-div{height:300px;}
</style>
</head>
<body>
<h1>🕯️ 港股通 100 股 - 最近 30 天数据</h1>
<p>数据来源: Tsanghi(XHKG) / 东方财富 | 生成时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>

<table>
<thead><tr><th>排名</th><th>股票名称</th><th>代码</th><th>最新价</th><th>30日涨跌幅</th><th>30日最高</th><th>30日最低</th><th>数据</th></tr></thead>
<tbody>
` + tableRows + `
</tbody>
</table>

<div class="chart-section">
` + chartScripts + `
</div>

<p style="margin-top:32px;color:#666;font-size:14px;">
说明：前 10 只有数据的股票显示 K 线图。若数据为 ❌，请检查 Tsanghi token 或网络。
</p>
</body>
</html>`

	return os.WriteFile(path, []byte(tpl), 0644)
}

func makeChartScript(r HKResult, id int) string {
	if len(r.Candles) == 0 {
		return ""
	}
	dates := make([]string, 0, len(r.Candles))
	values := make([]string, 0, len(r.Candles))
	for i := len(r.Candles) - 1; i >= 0; i-- {
		c := r.Candles[i]
		dates = append(dates, time.Unix(c.Timestamp, 0).Format("01-02"))
		values = append(values, fmt.Sprintf("[%.2f,%.2f,%.2f,%.2f]", c.Open, c.Close, c.Low, c.High))
	}
	datesJSON := "["
	for i, d := range dates {
		if i > 0 {
			datesJSON += ","
		}
		datesJSON += fmt.Sprintf("%q", d)
	}
	datesJSON += "]"
	valuesJSON := "[" + values[0]
	for i := 1; i < len(values); i++ {
		valuesJSON += "," + values[i]
	}
	valuesJSON += "]"

	divID := fmt.Sprintf("chart%d", id)
	return fmt.Sprintf(`
<div class="chart">
<h3>%s (%s) - 30日 K 线</h3>
<div id="%s" class="chart-div"></div>
</div>
<script>
(function(){
var opt = { xAxis:{type:'category',data:%s}, yAxis:{scale:true},
  series:[{type:'candlestick',data:%s}],
  tooltip:{trigger:'axis'},
  grid:{left:'3%%',right:'4%%',bottom:'3%%',containLabel:true}
};
var c = echarts.init(document.getElementById('%s'));
c.setOption(opt);
window.addEventListener('resize',function(){c.resize();});
})();
</script>
`, html.EscapeString(r.Stock.Name), r.Stock.Code, divID, datesJSON, valuesJSON, divID)
}
