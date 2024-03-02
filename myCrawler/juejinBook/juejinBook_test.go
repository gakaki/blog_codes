package juejinBook

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetAllXiaoces(t *testing.T) {
	GetAllBookListSortLatestSaveToJSON()
}

func TestDownload2Markdown(t *testing.T) {

	c := Config{
		Sessionid: "a1deccf2241693a29d8b7b9b316a8fb3",
		BookIDs: []string{
			//"7302990019642261567",
			//"6918979822425210891", // 0 打造通用型低
			//"7202598408815640631", //前端依赖治理
			//"7269673629964173331", // 前端可视化入门与实战
			//"7288940354408022074", //web动画之旅
			"7294082310658326565", //react 通关秘籍
			//"7306163555449962533", //前端工程体验优化实战
			//"6844733800300150797", //前端算法与数据结构面试：底层逻辑解读与大厂真题训练
		},
		SaveDir: "",
	}

	juejin, err := NewConfig(c)
	if err != nil {
		fmt.Println(err)
	}
	juejin.Download()
}

func TestGetMarkDownImageUrl(t *testing.T) {
	s := "# 前言\n\n性能优化不等于体验优化。\n\n就像吃饭不等于吃饱一样，吃饭和性能优化都只是**手段**，而吃饱和优化体验才是**目的**。\n\n谈到前端优化，许多人往往只关注性能优化，但性能优化只是改善前端项目的方式之一，**改善用户体验和开发体验**才应该是我们优化的根本目的。\n\n所以这本书不同于以往的前端性能优化书籍，不是简单的罗列优化方法，而是更多的关注**方法论**，引导读者从宏观的视角，关注前端优化开始前、实施中、生效后的全过程，最终极致地、有效地改善用户体验和开发体验。\n\n从而解决以往前端优化的诸多痛点：\n\n1.  **目标不明确：** 只会照本宣科，把别人的优化手段生搬硬套到自己的项目，找不到自己的优化目标；\n1.  **缺乏量化指标：** 无法评估优化效果，拿不出客观、可量化的指标证明优化效果；\n1.  **没有改善用户的主观体验：** 优化效果对用户主观体验影响甚微，无法让用户直观地感受到；\n1.  **欠缺长效化机制：** 无法保证优化效果长期稳定、不出现衰退；\n1.  **忽视开发体验：** 没有认识到开发体验和用户体验的正相关性；\n\n  \n这本书总结了作者6年多来优化和维护百万日访问量广告管理后台项目、千万日活信息流前端项目以及全球领先的浏览器平台音视频会议等众多项目的经验，将以具体场景和实践经验为例，搭配6000+行源码，深入浅出地讲解现代前端工程体验优化的方法论和具体措施。\n\n希望这本书能让读者有所收获，为用户带来更好体验，在工作中取得更多成就。\n\n更欢迎通过评论，邮件，微信交流群等方式和笔者、同学们交流探讨。\n\n  \n\n\n# 1. 为什么要建立量化指标？\n\n没有量化指标的优化是没有说服力的，不了解优化目标的现状更无法实现优化。\n\n  \n在实施优化时，一个常常陷入的误区是不了解现状，缺乏量化优化效果的监控指标就开始优化。这样的方式往往导致自欺欺人的优化，自以为做了效果显著的改进优化，实际上并没有改善用户体验。\n\n\n例如笔者在刚参加工作时就曾做过一次没有量化指标、生搬硬套的技术改进。\n\n当时计划为内部前端项目的 JS，CSS 等静态资源增加预加载`Prefetch` ，因为没有提前建立量化指标、监控优化效果的理念，所以在完成优化，部署上线后，没有得到优化效果的量化数据，只能通过 JS 加载命中了缓存来解释优化的收益，这次优化对用户体验的影响更是一无所知。\n\n这样没有反馈的技术改进显然不能算是成功的优化。\n\n> 注：https://react.dev/ 的预加载 Prefetch 示例：\n>\n> ![](img/1/1.png#?w=1920&h=942&s=188973&e=png&b=24262c)\n\n  \n\n所以，为了能真正的改善用户体验，我们需要在开始优化前，就建立量化指标，一方面帮助我们透彻的理解优化目标的现状，另一方面，也可以用量化指标监控优化前后的变化，评估优化效果。\n\n这就需要我们能把**主观**的用户体验或开发体验**量化**为**客观**的数据指标。\n\n  \n\n\n# 2. 将主观的体验量化为客观指标\n\n体验是主观的感受，同样的事物对不同的人，在不同的环境都会有不同的体验。\n\n我们以前端页面的加载速度为例，同一个页面，在不同的地理位置，不同的硬件设备上，加载速度都会有不同的表现，给用户的主观体验更是因人而异。\n\n所以要测量用户对页面的加载速度的主观体验，需要考虑的因素非常多，我们需要能抹平各种影响因素差异，稳定衡量用户体验的量化手段。\n  \n\n业界经过多年的实践，尝试过许多量化用户体验的方式，例如：\n\n-   测量页面白屏时长\n-   计算可交互耗时（Time to interact）\n-   计算总阻塞时间 (Total Blocking Time，TBT)\n-   计算首次有效绘制 (First Meaningful Paint，FMP)\n\n但这些指标往往逻辑复杂、难以测量，甚至定义都有显著的歧义，所以逐渐消亡。\n\n经过大浪淘沙，近年来最实用的用户体验量化指标是基于开源库`web-vitals`获取的页面渲染耗时，交互延迟等指标。\n\n> `web-vitals` GitHub 仓库: https://github.com/GoogleChrome/web-vitals\n\n\n## 1. web-vitals 各项指标简介\n\n`web-vitals`是谷歌的 Chrome 维护团队于 2020 年开源的工具库，它基于统一的浏览器 `Performance API` 获取标准化的用户体验指标。\n\n它主要测量6项指标，分别是：\n\n1.  首次内容绘制 (First Contentful Paint，FCP)\n1.  最大内容绘制 (Largest Contentful Paint，LCP)\n1.  首次输入延迟 (First Input Delay ，FID)\n1.  交互到绘制延迟（Interaction to Next Paint，INP）\n1.  累积布局偏移 (Cumulative Layout Shift，CLS)\n1.  第一字节时间 (Time to First Byte，TTFB)\n\n下面我们将逐一了解这些指标的测量目标、评分标准和具体示例。\n\n  \n### 1. 首次内容绘制 (First Contentful Paint，FCP)\n\nFCP测量从页面开始加载到页面中任意部分内容（文本、图像、`<svg/>`，`<canvas/>`等内容）完成渲染的时长。\n\n其值为浮点数，单位是毫秒。FCP值越小表示该指标状况越好、页面的初始内容渲染越快。\n\n页面中率先出现的文本图像等视觉可见内容，直接决定了用户对页面加载速度的主观体验，所以这一指标选择测量这些内容的渲染耗时，从而量化用户的主观体验。\n\n\n注意，FCP测量的是**任意部分**DOM完成渲染的耗时，而非**全部**内容完成渲染耗时，不等于`onLoad`事件。\n\n如下图中的例子，FCP指标的值为1439毫秒，在这个时刻页面中首次渲染出了文字和图像。\n\n![](img/1/2.png#?w=1920&h=1002&s=214561&e=png&b=232427)\n\n  \n\n\n按照Chrome官方的推荐标准，FCP指标3个等级的评分分别为：\n\n-   优：小于1.8秒\n-   待改进：大于1.8秒且小于3秒\n-   差：大于3秒\n\n![](img/1/3.png#?w=907&h=224&s=15461&e=png&b=d2d3d7)\n\n  \n\n\n### 2. 最大内容绘制 (Largest Contentful Paint，LCP)\n\nLCP测量从页面开始加载到可视区域内**尺寸最大**的文字或图像渲染完成的耗时。\n\n其值为浮点数，单位是毫秒。LCP值越小表示该指标状况越好、最大元素渲染越快。\n\n之所以测量最大的内容，是因为尺寸最大的内容往往最能吸引用户的注意力，其渲染耗时，直接影响了用户对页面整体渲染速度的体验。\n\n  \n\n\n我们可以用Chrome浏览器自带 DevTool 中的 Performance Insights 工具来判断页面中什么元素是最大内容，例如下图中的`img.banner-image`就是掘金首页的最大内容元素，这个元素渲染的耗时为1.55秒，即LCP的值。\n\n![](img/1/4.png#?w=1920&h=942&s=314104&e=png&b=222326)\n\n按照Chrome官方的推荐标准，LCP指标3个等级的评分分别为：\n\n-   优：小于2.5秒\n-   待改进：大于2.5秒且小于4秒\n-   差：大于4秒\n\n![](img/1/5.png#?w=888&h=256&s=16717&e=png&b=ffffff)\n\n  \n\n### 3. 首次输入延迟 (First Input Delay ，FID)\n\nFID 测量用户首次交互（点击、触摸）后到浏览器开始响应之间的时间间隔。\n\n其值为浮点数，单位是毫秒。FID值越小表示该指标状况越好，用户首次与页面交互时，浏览器响应的延迟越小。\n\n这一指标只关注页面中首次交互的原因是因为，首次交互时，页面往往处于尚未完全加载的状态，异步响应数据仍在等待响应、部分JS和CSS仍在执行和渲染的过程中，浏览器的主线程会短暂的处于忙碌状态，往往不能即时响应用户交互。\n\n但是第一次交互的延迟长短往往决定了用户对网页流畅度的第一印象，所以这一指标的测量目标，也能量化用户的主观体验。\n\n\n按照Chrome官方的推荐标准，FID指标3个等级的评分分别为：\n\n-   优：小于100毫秒\n-   待改进：大于100毫秒且小于300毫秒\n-   差：大于300毫秒\n\n![](img/1/6.png#?w=839&h=214&s=13670&e=png&b=ffffff)\n\n\n> 注：FID指标与下文将要提到的 INP 指标测量目标有所重叠，且普适性不及INP，未来可能会被INP替代。\n\n  \n\n\n### 4. 交互到绘制延迟（Interaction to Next Paint，INP）\n\nINP测量用户在页面浏览过程中的所有交互（点击、键盘输入、触摸等）与浏览器渲染响应的**整体**延迟情况。\n\n其值为浮点数，单位是毫秒。INP值越小表示该指标状况越好，用户的各类交互响应延迟越小。\n\n\n\n与FID只关注首次交互不同，INP会关注用户浏览网页全过程中的**所有**交互，所以`web-vitals`库中获取INP值的`onINP(FCPReportCallback)`方法，通常会在页面可视化状态变化或页面卸载时多次触发，综合统计一段时间内的多次交互，按特定算法，计算该时段内的INP指标值。\n\n\nINP指标3个等级的评分分别为：\n\n-   优：小于200毫秒\n-   待改进：大于200毫秒且小于500毫秒\n-   差：大于500毫秒\n\n![](img/1/7.png#?w=960&h=243&s=17232&e=png&b=202124)\n\n> INP是新近加入`web-vitals`的一项指标，仍处于实验状态，其标准可能会有调整，目前描述的是其2023年5月的状况。\n\n  \n\n### 5. 累积布局偏移 (Cumulative Layout Shift，CLS)\n\nCLS测量页面中所有**意外**布局变化的累计分值。\n\n其值为浮点数，**无单位，** 值的大小表示意外布局变化的多少和影响范围的大小。\n\nCLS值的计算类似INP，会统计一段时间内的所有意外布局变化，按特定算法，计算出分值。\n\n\n所谓**意外布局变化**是指 DOM 元素在前后绘制的2帧之间，非用户交互引起DOM元素尺寸、位置的变化。\n\n请看示例视频：\n\n![](img/1/8.png#?w=880&h=794&s=560666&e=gif&f=402&b=fefefe)\n\n这段视频中用户本想点击取消按钮，但是页面元素的布局位置突然产生了变化，出现了**非用户交互导致**的**意外布局变化**，原本取消按钮的位置被确认按钮替代，导致了用户本想点击取消，却触发了购买的误操作，严重损害了用户体验。\n\n> [《意外布局变化》在线DEMO](https://codesandbox.io/p/devbox/cls-demo-qfu8g5?file=%2Fsrc%2Findex.js%3A6%2C53)\n\n\n引入`web-vitals`库后调用`onCLS`API就能获取CLS的值，同时获取到对应的意外布局变化的具体来源，如下图中`sources`字段的2个对象就通过DOM元素引用，明确地告诉了我们引起布局变化的来源，以及变化前后的尺寸位置等详细数据`sources[i].currentRect, sources[i].previousRect`：\n\n![](img/1/9.png#?w=1919&h=990&s=313395&e=png&b=222326)\n\n按照Chrome官方的推荐标准，CLS指标3个等级的评分分别为：\n\n-   优：小于0.1\n-   待改进：大于0.1且小于0.25\n-   差：大于0.25\n\n![](img/1/10.png#?w=862&h=241&s=16329&e=png&b=ffffff)\n\n  \n\n\n### 6. 第一字节时间 (Time to First Byte，TTFB)\n\nTTFB测量前端页面（Document）的HTTP请求发送后，到接收到第一字节数据响应的耗时，通常包括重定向、DNS查询、服务器响应延迟等耗时。\n\n其值为浮点数，单位是毫秒。值越小表示该项指标状况越好，页面HTTP响应的耗时越短，也就是页面的加载更快。\n\n\n\nTTFB指标值的大小直接决定着页面初始内容渲染耗时的长短，往往和`FCP`、`LCP`指标有明显的相关关系，对用户体验有直接影响，所以`web-viatals`也将其当做了量化用户体验的指标之一。\n\n  \n\n除了可以通过`web-vitals`库的`onTTFB()`API获取，也可以使用 Chrome 自带的 DevTool Network 网络面板计算得出。\n\n如下图的例子知乎首页的`TTFB`指标值即为：\n\n-   `文档响应的整体耗时` 减去 `内容下载耗时（Content Download）`\n-   391毫秒 - 57毫秒 = 335毫秒\n\n![](img/1/11.png#?w=1920&h=942&s=317512&e=png&b=242528)\n\n  \n\n\nTTFB指标3个等级的评分分别为：\n\n-   优：小于800毫秒\n-   待改进：大于800毫秒且小于1800毫秒\n-   差：大于1800毫秒\n\n![](img/1/12.png#?w=909&h=219&s=14280&e=png&b=d2d3d7)\n\n尽管以上指标都可以通过原生Performance API计算获得，但仍然推荐使用的`web-vitals`库，因为它能帮助我们处理了许多细节问题，例如标签页处于后台时的计算、指标获取时机、浏览器兼容性等等，能确保我们测量出标准、稳定的指标数值。\n\n## 2. 六类指标对比\n\n| 名称                                         | 含义                                                            | 注意事项                                                                                                                                  | 值单位 | WebVitals 库获取结果示例 |\n| ------------------------------------------ | ------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- | --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |\n| 首次内容绘制(First Contentful Paint，**FCP**)     | 测量从页面开始加载到页面中**任意**部分内容（文本、图像、`<svg/>`，`<canvas/>`等内容）完成渲染的时长 | 测量任意**部分**DOM渲染的耗时，而非全部内容，不等于页面所有内容完全加载完成的`onLoad`事件。                                                                                 | 毫秒  | ```{   \"name\": \"FCP\",   \"value\": 463.20000076293945,   \"rating\": \"good\",   \"delta\": 463.20000076293945,   \"entries\": [     {       \"name\": \"first-contentful-paint\",       \"entryType\": \"paint\",       \"startTime\": 463.20000076293945,       \"duration\": 0     }   ],   \"id\": \"v3-1695054859140-2991050486027\",   \"navigationType\": \"reload\" }``` |\n| 最大内容绘制 (Largest Contentful Paint，**LCP**)  | 测量从页面开始加载到可视区域内尺寸最大的文字或图像渲染完成的耗时。                             | 对于UI渲染逻辑复杂的前端应用，不同优化可能会有不同的最大元素，统计获得的最大元素可能有多个。                                                                                       | 毫秒  | ``` {   \"name\": \"LCP\",   \"value\": 463.20000076293945,   \"rating\": \"good\",   \"delta\": 463.20000076293945,   \"entries\": [     {       \"name\": \"\",       \"entryType\": \"largest-contentful-paint\",       \"startTime\": 463.20000076293945,       \"duration\": 0,       \"size\": 8985,       \"renderTime\": 463.2,       \"loadTime\": 0,       \"firstAnimatedFrameTime\": 0,       \"id\": \"\",       \"url\": \"\"     }   ],   \"id\": \"v3-1695054859140-6431611124119\",   \"navigationType\": \"reload\" }``` |\n| 首次输入延迟(First Input Delay ，**FID**)         | 测量用户首次交互（点击、触摸）后到浏览器开始响应用户交互之间的时间间隔。 | 未来可能会被INP替代。 | 毫秒  | ```{   \"name\": \"FID\",   \"value\": 1.7999992370605469,   \"rating\": \"good\",   \"delta\": 1.7999992370605469,   \"entries\": [     {       \"name\": \"pointerdown\",       \"entryType\": \"first-input\",       \"startTime\": 1954,       \"duration\": 0,       \"processingStart\": 1955.7999992370605,       \"processingEnd\": 1955.7999992370605,       \"cancelable\": true     }   ],   \"id\": \"v3-1695054999447-8039144919554\",   \"navigationType\": \"reload\" }``` |\n| 交互到绘制延迟(Interaction to Next Paint，**INP**) | 测量用户在页面浏览过程中的所有交互（点击、键盘输入、触摸等）与浏览器绘制对应响应的**整体**延迟情况。          | 通常会在页面可视化状态变化或页面卸载时进行计算。 | 毫秒  | ``` {   \"name\": \"INP\",   \"value\": 8,   \"rating\": \"good\",   \"delta\": 8,   \"entries\": [     {       \"name\": \"pointerdown\",       \"entryType\": \"first-input\",       \"startTime\": 36701.30000114441,       \"duration\": 8,       \"processingStart\": 36702.80000114441,       \"processingEnd\": 36702.80000114441,       \"cancelable\": true     }   ],   \"id\": \"v3-1695054859140-4977365532114\",   \"navigationType\": \"reload\" }``` |\n| 累积布局偏移(Cumulative Layout Shift，**CLS**)    | 测量页面中，一定时间段内所有意外布局变化的累计分值。                                    | -   通常会在页面可视化状态变化或页面卸载时进行计算。<br/>-   `web-viatals`提供的`onCLS()`方法会多次触发。<br/>-   `onCLS()`获取到的`sources`字段可能会因为元素卸载而变成`null`，统计时可以使用xpath进行特殊处理。 | 分值  | ```{   \"name\": \"CLS\",   \"value\": 0.0007529577629112682,   \"rating\": \"good\",   \"delta\": 0.0007529577629112682,   \"entries\": [     {       \"entryType\": \"layout-shift\",       \"value\": 0.0007529577629112682,       // ...       \"sources\": [         {           \"previousRect\": {             \"x\": 128,             \"y\": 553,             \"width\": 20,             \"height\": 20,             \"top\": 553,             \"right\": 148,             \"bottom\": 573,             \"left\": 128           },           \"currentRect\": {               // ...           }         }       ]     }   ],   \"id\": \"v3-1695054859142-8118655247179\",   \"navigationType\": \"reload\" }``` |\n| 第一字节时间(Time to First Byte，**TTFB**)        | 测量页面本身（Document）的HTTP请求发送后，到接收到第一字节数据响应的耗时                    | 往往和FCP、LCP指标有相关关系。 | 毫秒  | ```{   \"name\": \"TTFB\",   \"value\": 369.20000076293945,   \"rating\": \"good\",   \"delta\": 369.20000076293945,   \"entries\": [     {       \"name\": \"https://output.jsbin.com/bizanep\",       \"entryType\": \"navigation\",       \"startTime\": 0,       \"duration\": 820.5,       \"initiatorType\": \"navigation\",       \"nextHopProtocol\": \"h2\",       \"renderBlockingStatus\": \"non-blocking\",       \"workerStart\": 0,       // ...       \"activationStart\": 0,       \"criticalCHRestart\": 0     }   ],   \"id\": \"v3-1695054859140-2231742211102\",   \"navigationType\": \"reload\" }``` |\n\n  \n\n\n## 3.  `web-vitals`使用示例\n\n以上6项指标均可通过`web-vitals`库内置的API方便的获取，将`web-vitals`库集成到用户访问的前端页面，即可方便地获取用户的真实体验数据，例如：\n\n> 获取`web-vitals`数据在线 DEMO: https://output.jsbin.com/bizanep\n\n![](img/1/13.png#?w=1919&h=985&s=246939&e=png&b=212225)\n\n``` html\n<!DOCTYPE html>\n<html>\n<head>\n  <meta charset=\"utf-8\">\n  <meta name=\"viewport\" content=\"width=device-width\">\n  <title>获取 web-vitals 数据 DEMO</title>\n</head>\n<body>\n  <h2 id=\"fcp\">FCP:</h2>\n  <h2 id=\"lcp\">LCP:</h2>\n  <h2 id=\"ttfb\">TTFB:</h2>\n  <p>首次交互（例如：点击任意位置）后可获取：</p>\n  <h2 id=\"fid\">FID:</h2>\n  <p>页面可视化状态变化为隐藏（例如：切换标签页）后可获取：</p>\n  <h2 id=\"inp\">INP:</h2>\n  <h2 id=\"cls\">CLS:</h2>\n  \n  \n  <a href=\"https://github.com/JuniorTour\">Author: https://github.com/JuniorTour</a>\n  \n  <script type=\"module\">\n    import {onFCP, onLCP, onFID, onCLS, onINP, onTTFB} from 'https://unpkg.com/web-vitals@3?module';\n\n    function setInnerHtml(id, html) {\n      if (!id || !html) {\n        return\n      }\n      const el = document.querySelector(`#${id}`)\n      if (el) {\n        el.innerHTML = html\n      }\n    }\n    \n    function onGetWebVitalsData(data) {\n      console.log(data)\n      if (!data?.name) {\n        return\n      }\n      const name = data.name\n      const value = data.value\n      const rating = data.rating\n      const msg = (`${name}: value=${value}, rating=${rating}`)\n      console.log(msg)\n      setInnerHtml(name?.toLowerCase(), msg)\n    }\n    \n    onFCP(onGetWebVitalsData);\n    onLCP(onGetWebVitalsData); \n    onFID(onGetWebVitalsData); \n    onCLS(onGetWebVitalsData);\n    onINP(onGetWebVitalsData);\n    onTTFB(onGetWebVitalsData);\n  </script>\n\n</body>\n</html>\n```\n\n要注意的细节是，这些指标中：\n\n-   `onFCP, onLCP, onTTFB` 均为在页面初始化时自动触发。\n-   `onFID`是在用户第一次与页面交互时触发。\n-   `onCLS, onINP`则因为要测量页面的全生命周期，往往无固定触发时间点，在实践中通常会在交互停止一段时间后，或页面可视状态变化（例如切换标签页）后多次触发。\n\n  \n\n\n`web-vitals`的这些指标是Chrome维护团队基于海量用户数据、经过大量实践后设计出来的，能科学地将主观的用户体验量化为客观的指标，是我们进行体验优化的必备工具。\n\n  \n\n\n大量的收集这些指标数据，加以汇总分析便可以实现针对用户体验的“真实用户监控”（https://en.wikipedia.org/wiki/Real_user_monitoring） ，从用户客户端收集到**海量**数据，要比我们在内部的测试开发环境上测量出的**少量**实验室数据更全面、更客观、更有说服力，更有助于我们做出数据驱动的优化决策。\n\n  \n\n\n# 小结\n\n这一节中我们主要学习了建立量化指标的意义，是为了能真正的改善用户体验。\n\n并详细介绍了`web-viatals`库，FCP、LCP、FID、INP、TTFB等6项用户体验指标的含义、细节和具体用法。"
	fmt.Println(GetMarkDownImageUrl(s))
}
func TestParseMarkdownImagePath(t *testing.T) {
	// 定义正则表达式
	text := " Electron 应用场景的分布\n\n使用 `Electron` 开发的应用品类非常丰富，" +
		"我们看看官网的一些案例展示中的统计数据：\n\n" +
		"<p align=center><img src=\"https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/94cb5ee0174642acb4a24cf2c0fe1fad~tplv-k3u1fbpfcp-jj-mark:0:0:0:0:q75.image#?w=1486&h=1193&s=53343&e=png&a=1&b=546ec6\" alt=\"image.png\"  /></p>\n\n> 数据来源：[Electron ShowCase](https://www.electronjs.org/apps)\n\n可以看到，在使用 Electron 开发的 APP 中，开发者工具、效率应用占据了大半江山。\n"
	//img_pattern := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
	//matches := img_pattern.FindAllStringSubmatch(text, -1)
	//for _, match := range matches {
	//	url := match[1]
	//	fmt.Println(url)
	//}
	images := FindImageUrls(1, text)
	fmt.Println(images)
}

func TestRenderPDF(t *testing.T) {
	//brew install pandoc
	//brew install --cask basictex

	// 遍历文件夹
	filepath.Walk("./book", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() == false && filepath.Ext(info.Name()) == ".md" {
			// 获取文件名
			fileName := filepath.Base(info.Name())
			// 切换到当前文件夹
			os.Chdir(filepath.Dir(path))

			newFileName := strings.Replace(fileName, filepath.Ext(fileName), "", 1)
			// 打印转换开始信息
			fmt.Println("转换开始：" + "pandoc " + fileName + " -o " + newFileName + ".pdf")

			// 调用 pandoc 进行格式转换
			cmd := exec.Command(fmt.Sprintf("pandoc %s -o %s", fileName, newFileName))
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("combined out:\n%s\n", string(out))
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}

			// 打印转换完成信息
			fmt.Println("转换完成...")
		}
		return err
	})

}
