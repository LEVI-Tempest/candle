#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
港股通 100 股 - 使用 AkShare 拉取最近 30 天数据并生成 HTML
港股通100股 - Fetch 30 days data via AkShare and generate HTML

依赖: pip install -r scripts/requirements.txt
运行:
  python scripts/hk100_akshare.py [output.html]
  # 或使用 venv:
  python3 -m venv .venv && source .venv/bin/activate  # Linux/macOS
  pip install -r scripts/requirements.txt
  python scripts/hk100_akshare.py docs/港股通100股_akshare.html

注意: 如遇 ProxyError，请关闭代理或取消 HTTP_PROXY/HTTPS_PROXY 环境变量后再试。
"""

import json
import sys
from datetime import datetime, timedelta
from html import escape

try:
    import akshare as ak
    import pandas as pd
except ImportError:
    print("请先安装: pip install akshare pandas")
    sys.exit(1)

HK100 = [
    (1, "腾讯控股", "00700"), (2, "阿里巴巴-W", "09988"), (3, "美团-W", "03690"),
    (4, "建设银行", "00939"), (5, "小米集团-W", "01810"), (6, "中国移动", "00941"),
    (7, "中芯国际", "00981"), (8, "快手-W", "01024"), (9, "中国海洋石油", "00883"),
    (10, "比亚迪股份", "01211"), (11, "工商银行", "01398"), (12, "汇丰控股", "00005"),
    (13, "中国银行", "03988"), (14, "京东集团-SW", "09618"), (15, "友邦保险", "01299"),
    (16, "香港交易所", "00388"), (17, "中国平安", "02318"), (18, "招商银行", "03968"),
    (19, "药明康德", "02359"), (20, "李宁", "02331"), (21, "安踏体育", "02020"),
    (22, "联想集团", "00992"), (23, "哔哩哔哩-W", "09626"), (24, "网易-S", "09999"),
    (25, "百度集团-SW", "09888"), (26, "京东健康", "06618"), (27, "商汤-W", "00020"),
    (28, "优必选", "09880"), (29, "招金矿业", "01818"), (30, "建滔积层板", "01888"),
    (31, "阜博集团", "03738"), (32, "中国神华", "01088"), (33, "中海油服", "02883"),
    (34, "中国石化", "00386"), (35, "中国石油股份", "00857"), (36, "紫金矿业", "02899"),
    (37, "中国电信", "00728"), (38, "中国联通", "00762"), (39, "华润啤酒", "00291"),
    (40, "舜宇光学科技", "02382"), (41, "瑞声科技", "02018"), (42, "海尔智家", "06690"),
    (43, "中国生物制药", "01177"), (44, "石药集团", "01093"), (45, "金斯瑞生物科技", "01548"),
    (46, "信达生物", "01801"), (47, "药明生物", "02269"), (48, "碧桂园服务", "06098"),
    (49, "龙湖集团", "00960"), (50, "华润置地", "01109"), (51, "中国海外发展", "00688"),
    (52, "长城汽车", "02333"), (53, "吉利汽车", "00175"), (54, "理想汽车-W", "02015"),
    (55, "小鹏汽车-W", "09868"), (56, "蔚来-SW", "09866"), (57, "零跑汽车", "09863"),
    (58, "申洲国际", "02313"), (59, "创科实业", "00669"), (60, "新奥能源", "02688"),
    (61, "华润燃气", "01193"), (62, "中国燃气", "00384"), (63, "恒安国际", "01044"),
    (64, "蒙牛乳业", "02319"), (65, "百威亚太", "01876"), (66, "农夫山泉", "09633"),
    (67, "周大福", "01929"), (68, "泡泡玛特", "09992"), (69, "海底捞", "06862"),
    (70, "呷哺呷哺", "00520"), (71, "阿里健康", "00241"), (72, "平安好医生", "01833"),
    (73, "微创医疗", "00853"), (74, "绿叶制药", "02186"), (75, "中国太保", "02601"),
    (76, "中国人寿", "02628"), (77, "交通银行", "03328"), (78, "中信银行", "00998"),
    (79, "光大银行", "06818"), (80, "邮储银行", "01658"), (81, "中金公司", "03908"),
    (82, "中信证券", "06030"), (83, "华泰证券", "06886"), (84, "银河娱乐", "00027"),
    (85, "金沙中国", "01928"), (86, "澳博控股", "00880"), (87, "新鸿基地产", "00016"),
    (88, "恒基地产", "00012"), (89, "长实集团", "01113"), (90, "恒隆地产", "00101"),
    (91, "中银香港", "02388"), (92, "恒生银行", "00011"), (93, "渣打集团", "02888"),
    (94, "思摩尔国际", "06969"), (95, "中教控股", "00839"), (96, "中国民航信息网络", "00696"),
    (97, "金蝶国际", "00268"), (98, "万国数据-SW", "09698"),
]


def fetch_one(symbol: str, start: str, end: str) -> pd.DataFrame | None:
    try:
        df = ak.stock_hk_hist(symbol=symbol, start_date=start, end_date=end, adjust="")
        return df if df is not None and not df.empty else None
    except Exception:
        return None


def main():
    output = sys.argv[1] if len(sys.argv) > 1 else "港股通100股_30日.html"

    end_d = datetime.now()
    start_d = end_d - timedelta(days=40)
    start_s = start_d.strftime("%Y%m%d")
    end_s = end_d.strftime("%Y%m%d")

    print("🕯️  港股通 100 股 - 最近 30 天数据 (AkShare)")
    print("=" * 50)
    print(f"数据源: AkShare (东方财富)")
    print(f"输出: {output}\n")

    results = []
    for i, (rank, name, code) in enumerate(HK100):
        df = fetch_one(code, start_s, end_s)
        if df is not None and len(df) > 0:
            df = df.tail(30)
            results.append({
                "rank": rank, "name": name, "code": code,
                "df": df, "rows": len(df)
            })
        else:
            results.append({"rank": rank, "name": name, "code": code, "df": None, "rows": 0})

        if (i + 1) % 20 == 0:
            print(f"  已拉取 {i + 1}/98")

    got = sum(1 for r in results if r["rows"] > 0)
    print(f"\n📊 成功获取 {got}/100 只股票数据\n")

    # 排序：有数据的放前面
    results.sort(key=lambda x: (-x["rows"], x["rank"]))

    # 生成 HTML
    table_rows = []
    chart_scripts = []
    chart_count = 0

    for r in results:
        row = f'<tr><td>{r["rank"]}</td><td>{escape(r["name"])}</td><td>{r["code"]}</td>'
        if r["df"] is not None and len(r["df"]) > 0:
            df = r["df"]
            last = df.iloc[-1]
            first = df.iloc[0]
            close = float(last["收盘"])
            high = float(df["最高"].max())
            low = float(df["最低"].min())
            pct = ((close - float(first["收盘"])) / float(first["收盘"]) * 100) if float(first["收盘"]) else 0
            row += f'<td>{close:.2f}</td><td>{pct:+.2f}%</td><td>{high:.2f}</td><td>{low:.2f}</td><td>✅</td>'

            if chart_count < 10:
                dates = [d.strftime("%m-%d") for d in df["日期"]]
                values = []
                for _, rw in df.iterrows():
                    values.append([float(rw["开盘"]), float(rw["收盘"]), float(rw["最低"]), float(rw["最高"])])
                dates_js = json.dumps(dates)
                values_js = json.dumps(values)
                div_id = f"chart{chart_count}"
                chart_scripts.append(f'''
<div class="chart">
<h3>{escape(r["name"])} ({r["code"]}) - 30日 K 线</h3>
<div id="{div_id}" class="chart-div"></div>
</div>
<script>
(function(){{
var opt={{xAxis:{{type:'category',data:{dates_js}}},yAxis:{{scale:true}},
series:[{{type:'candlestick',data:{values_js}}}],
tooltip:{{trigger:'axis'}},grid:{{left:'3%',right:'4%',bottom:'3%',containLabel:true}}}};
var c=echarts.init(document.getElementById("{div_id}"));
c.setOption(opt);
window.addEventListener('resize',function(){{c.resize();}});
}})();
</script>''')
                chart_count += 1
        else:
            row += "<td>-</td><td>-</td><td>-</td><td>-</td><td>❌</td>"
        table_rows.append(row + "</tr>")

    html_body = f"""<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<title>港股通 100 股 - 最近 30 天数据</title>
<script src="https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js"></script>
<style>
body{{font-family:system-ui,-apple-system,sans-serif;margin:20px;background:#f5f5f5;}}
h1{{color:#333;}}
table{{border-collapse:collapse;background:#fff;box-shadow:0 2px 8px rgba(0,0,0,.1);border-radius:8px;overflow:hidden;}}
th,td{{padding:10px 14px;text-align:left;border-bottom:1px solid #eee;}}
th{{background:#1a73e8;color:#fff;}}
tr:hover{{background:#f8f9fa;}}
.chart{{margin:24px 0;padding:16px;background:#fff;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,.1);}}
.chart h3{{color:#333;margin-top:0;}}
.chart-div{{height:300px;}}
</style>
</head>
<body>
<h1>🕯️ 港股通 100 股 - 最近 30 天数据</h1>
<p>数据来源: AkShare (东方财富) | 生成时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}</p>

<table>
<thead><tr><th>排名</th><th>股票名称</th><th>代码</th><th>最新价</th><th>30日涨跌幅</th><th>30日最高</th><th>30日最低</th><th>数据</th></tr></thead>
<tbody>
{"".join(table_rows)}
</tbody>
</table>

<div class="chart-section">
{"".join(chart_scripts)}
</div>

<p style="margin-top:32px;color:#666;font-size:14px;">
说明：前 10 只有数据的股票显示 K 线图。
</p>
</body>
</html>"""

    with open(output, "w", encoding="utf-8") as f:
        f.write(html_body)

    print(f"✅ HTML 已保存: {output}")


if __name__ == "__main__":
    main()
