package pick_word

import (
	"github.com/kevin-zx/site-info-crawler/sitethrougher"
	"github.com/kevin-zx/wordproperty"
	"strings"
)

var stopWords = []string{"img", "首页", "联系我们", "网站地图", "关于我们", "登录", "更多", "注册", "视频", "意见反馈", "友情链接", "图片", "2", "资讯", "常见问题", "3", "广告服务", "下一页", "地图", "4", "帮助中心", "新闻", "5", "汽车", "百度首页", "帮助", "北京", "6", "隐私政策", "上海", "网页", "联系方式", "其他", "法律声明", "旅游", "个人中心", "查看更多", "天津", "使用百度前必读", "重庆", "7", "生活", "知道", "文库", "8", "手机版", "9", "贴吧", "专题", "体育", "深圳", "广州", "教育", "论坛", "网站首页", "科技", "采购", "招聘信息", "武汉", "版权声明", "杭州", "游戏", "加入我们", "成都", "南京", "房产", "公司简介", "10", "全部", "音乐", "娱乐", "文化", "财经", "免费注册", "商城", "问答", "济南", "苏州", "免责声明", "长沙", "诚聘英才", "新闻中心", "沈阳", "金融", "西安", "人才招聘", "教育培训", "海南", "我的收藏", "历史", "English", "郑州", "吉林", "政策法规", "1", "合肥", "健康", "青岛", "网站导航", "退出", "哈尔滨", "石家庄", "手机", "四川", "家电", "新浪微博", "福州", "互联网", "厦门", "长春", "设为首页", "南宁", "无锡", "全国", "昆明", "空调", "广东", "南昌", "太原", "广西", "东莞", "用户协议", "江苏", "美食", "家居", "新疆", "在线客服", "贵阳", "母婴", "山东", "企业文化", "贵州", "详情", "兰州", "浙江", "云南", "河北", "宁波", "保险", "福建", "河南", "营业执照", "洗衣机", "湖南", "安徽", "大连", "陕西", "直播", "冰箱", "湖北", "数码", "时尚", "辽宁", "呼和浩特", "0", "乌鲁木齐", "江西", "山西", "甘肃", "客服中心", "黑龙江", "不限", "媒体报道", "企业", "社会", "银川", "内蒙古", "宠物", "宁夏", "电视", "军事", "行情", "广告", "新闻资讯", "二手房", "招贤纳士", "佛山", "笔记本", "经济", "公司介绍", "招聘", "人物", "西藏", "香港", "租房", "排行榜", "艺术", "家具", "青海", "软件下载", "评论", "西宁", "二手市场", "我的订单", "惠州", "珠海", "博客", "平板电脑", "查看详情", "联系客服", "建材", "产品", "服务", "商务服务", "家电维修", "生活服务", "常州", "服务条款", "通知公告", "网站律师", "微博", "在线咨询", "烟台", "温州", "三星", "中山", "投资者关系", "新闻动态", "手机维修", "三亚", "公司新闻", "公司", "热水器", "最新", "台式电脑", "秦皇岛", "动漫", "唐山", "泉州", "二手车", "站点地图", "商务合作", "徐州", "打印机", "嘉兴", "保定", "招商加盟", "洛阳", "交通", "南通", "扬州", "昆山", "绍兴", "价格", "经销商", "网站声明", "廊坊", "活动", "电影", "潍坊", "拉萨", "法律", "开放平台", "小米", "澳门", "合作伙伴", "微信", "汽车服务", "渠道招商", "高级搜索", "亲子", "维修保养", "切换城市", "医疗健康", "百度知道", "租车", "国际", "桂林", "威海", "银行", "绵阳", "江门", "相机", "评测", "包头", "台湾", "酒店", "版权所有", "美国", "搬家", "盐城", "济宁", "装修建材", "C", "装修", "星座", "我要提问", "泰州", "柳州", "宜昌", "健身教练", "海口", "装修效果图", "市场", "理财", "X", "智能家居", "淄博", "A", "帐号设置", "写字楼", "社区", "电脑维修", "自然", "股票", "地板", "镇江", "基金", "百科", "S", "M", "设计", "金华", "B", "购物", "地理", "L", "H", "衡水", "新房", "G", "W", "芜湖", "电脑", "Q", "J", "导购", "K", "T", "F", "湖州", "D", "11", "R", "临沂", "国内", "Y", "免费设计", "Z", "北海", "N", "信息公开", "商铺", "数码产品", "E", "尾页", "P", "汕头", "V", "松下", "赣州", "装修攻略", "清远", "现代", "婚庆摄影", "湛江", "产品中心", "O", "政策", "下载中心", "数据", "I", "媒体", "科学", "连云港", "台州", "新加坡", "淮安", "默认排序", "收藏", "U", "漳州", "日本", "九江", "育儿", "写字楼出租", "蚌埠", "南充", "游戏机", "空气净化器", "鞍山", "株洲", "餐饮", "写字楼出售", "免费发布信息", "舟山", "商铺出租", "产业", "南阳", "襄阳", "往期回顾", "要闻", "湘潭", "商铺出售", "邯郸", "下载", "泸州", "咸阳", "娱乐休闲", "新乡", "建筑", "岳阳", "我的", "朝阳", "黄石", "常德", "乐山", "房屋维修", "在线留言", "行业动态", "12", "历史上的今天", "泰安", "团购", "遵义", "商业", "衡阳", "肇庆", "沧州", "主页", "服务声明", "服务器", "读书", "举报不良信息", "土地", "客服", "政务公开", "客户端", "服务协议", "推广服务", "个人房源", "河源", "景点", "东营", "网站合作", "德阳", "搞笑", "苹果", "梅州", "房产知识", "网站建设", "德州", "原创", "开封", "承德", "数据恢复", "政府", "美容", "会计", "宝鸡", "医疗", "百科协议", "效果图", "发现", "阳江", "产品库", "百科任务", "行业资讯", "办公", "更多城市", "淮南", "滁州", "大庆", "管道疏通", "百科商城", "城市百科", "许昌", "百度百科合作平台", "报价", "蝌蚪团", "孝感", "商丘", "物流", "数字博物馆", "合作模式", "投诉侵权信息", "邢台", "张家口", "休闲娱乐", "郴州", "阜阳", "行业新闻", "吉安", "手机APP", "路由器", "渭南", "技术服务", "茂名", "十堰", "海尔", "宿迁", "家政", "信阳", "成长任务", "揭阳", "安庆", "聊城", "韶关", "搜索", "家政服务", "宜宾", "六安", "手机客户端", "驾校", "推荐", "热点", "大同", "天气预报", "百科冷知识", "IT", "编辑规则", "封禁查询与解封", "历史版本", "编辑入门", "官方贴吧", "安阳", "荆州", "临汾", "赤峰", "秒懂大师说", "本人编辑", "秒懂全视界", "APP下载", "咸宁", "图解百科", "集成灶", "电工", "未通过词条申诉", "秒懂看瓦特", "发展历程", "非遗百科", "濮阳", "懂啦", "了解更多", "潮州", "政策解读", "秒懂五千年", "翻译", "会员登录", "恐龙百科", "装修日记", "家居资讯", "多肉百科", "电动车", "看房团", "漯河", "火车票", "娄底", "电磁炉", "公益", "摄影", "菏泽", "学前教育", "齐齐哈尔", "庆阳", "客户服务", "TCL", "音响", "找小区", "宁德", "太仓", "售后服务", "人力资源", "别墅", "下载客户端", "一级建造师", "CPU", "工作动态", "楼盘导购", "汉中", "周口", "化妆品", "吸尘器", "眉山", "日照", "防城港", "执业药师", "平顶山", "房产百科", "内江", "龙岩", "经济师", "机票", "电饭煲", "广安", "益阳", "衢州", "中国", "综合", "摩托车", "LG", "立即注册", "房产问答", "鄂州", "营口", "微波炉", "抽油烟机", "我要卖房", "本田", "大家电", "燃气灶", "资料下载", "主板", "二手", "长治", "外汇", "玉林", "本月开盘", "上饶", "智能穿戴", "怀化", "业主论坛", "显卡", "永州", "联想", "平面设计", "消防工程师", "汕尾", "曲靖", "依申请公开", "品牌推广", "奥迪", "房产快讯", "地图找房", "会员中心", "中关村在线", "耳机", "职业", "马鞍山", "投影机", "RSS订阅", "加湿器", "荆门", "广告合作", "燃梦计划", "交换机", "海外房产", "教师资格证", "网络设备", "装修家居", "百度营销", "月嫂", "超极本", "排行", "家庭影院", "抚州", "三门峡", "浴霸", "保姆", "宝马", "房贷计算器", "投资理财", "美的", "商办云", "大众", "遂宁", "培训", "SUV", "枣庄", "丽水", "二手房排行榜", "网上有害信息举报专区", "合伙人", "美发", "心理咨询师", "别克", "期货", "兼职", "电子书", "家居云", "豆浆机", "查房价", "数码相机", "文库协议", "三农", "立即开通", "儿童", "保洁", "沙发", "汽车用品", "空调移机", "文库报告", "三居", "我的简历", "食品", "整租房源", "丹东", "张家界", "行情中心", "黄山", "乌兰察布", "宣城", "免费发布出租", "牡丹江", "热门楼盘", "在线支付", "黄金", "焦作", "讨论", "金融财经", "网络安全", "门窗", "小家电", "全部问题", "亳州", "我的房天下", "宿州", "驻马店", "营养师", "新房排行榜", "宜春", "丰田", "除湿机", "房企研究", "签证", "入驻房产圈", "加盟房天下", "达州", "别墅房源", "隐私保护", "留学", "装修报价", "托福考试", "在售房源", "专栏", "家电数码", "经纪云", "景德镇", "英国", "广元", "运动健身", "洗碗机", "加拿大", "公告", "开发云", "美国留学", "台式机", "厨房", "销售", "健康管理师", "小户型", "天水", "地产数据", "二手手机", "雅思考试", "奇瑞", "维修", "晋中", "人民网", "黄冈", "找别墅", "钦州", "律师", "生产制造", "PDF转换", "服装", "会员介绍", "企业文库", "购房知识", "快递", "运城", "设计策划", "解决方案", "剃须刀", "电暖器", "文化历史", "楼盘新动态", "电风扇", "电源", "58同城", "品牌", "注册会计师", "跑车", "礼品", "U盘", "个人护理", "消毒柜", "扫地机器人", "家教", "先行赔付", "思维导图", "液晶电视", "长安", "返回首页", "锦州", "卫浴", "影音家电", "电吹风", "餐饮美食", "VIP新客立减2元", "硬盘", "家用电器", "家装", "行业", "酒店预订", "家具维修", "奔驰", "会议中心", "申请认证", "机构认证", "教育文库", "印刷包装", "电熨斗", "产品展示", "装修公司", "日产", "每日任务", "内存", "学术专区", "海信", "数据中心", "VIP免费专区", "VIP福利专区", "国统调查", "BIM工程师", "通辽", "账号申诉", "政务服务", "企业应用软件", "兑换VIP", "百度题库", "精品文库", "悬赏文档", "详细", "电压力锅", "会议平台", "床", "忘记密码", "法律专区", "莆田", "代理招募", "教育云平台", "自考学历", "瓷砖选材", "在职考研", "专业认证", "云知识", "在租房源", "中高考学习", "福特", "环保家装", "服务项目", "悬赏任务", "少儿钢琴陪练", "CG设计", "烧烤", "滨州", "中央空调", "扫描仪", "认证团队", "新车", "MPV", "日报", "晋城", "新加坡房产", "新股", "时尚美容", "企业采购", "家居圈", "互动交流", "本月交房", "其它", "知道协议", "如何答题", "知道商城", "情感心理", "家纺", "汽车电子", "淮北", "在问", "美国房产", "首尔", "服务中心", "帮助更多人", "使用财富值", "橱柜", "日报广场", "邵阳", "电动牙刷", "商用", "手机答题", "婚庆", "知道团队", "获取采纳", "日报作者", "生活家电", "高质量问答", "新华网", "抚顺", "滑雪", "图赏", "芝麻团", "合伙人认证", "其他组织", "钟点工", "日报精选", "认证用户", "兑换商品", "澳大利亚房产", "体温计", "餐厅", "百城价格指数", "酸奶机", "希腊房产", "电火锅", "电热水器", "跑步机", "盘锦", "一居", "免费发布", "专题专栏", "比亚迪", "按摩器", "外语培训", "越南房产", "电饼铛", "办公软件", "辽阳", "日本房产", "厨师", "菲律宾房产", "移动硬盘", "足浴盆", "酒柜", "云浮", "英语", "中指数据库", "港股", "烤箱", "英国房产", "阿联酋房产", "助听器", "梧州", "空调扇", "车辆", "综艺", "长虹", "自贡", "EN", "销售代表", "美容保健", "柬埔寨房产", "成功案例", "天津房地产", "大连房地产", "清洁机", "产品排行", "家装案例", "常州房地产", "空调维修", "向TA提问", "购物流程", "南宁房地产", "葫芦岛", "南京房地产", "笔记本电脑", "无线上网卡", "模拟攒机", "西安房地产", "二手家电", "公司动态", "了解58同城", "定西", "加入58同城", "计步器", "配件", "新浪首页", "VR", "海南房地产", "苏州房地产", "马来西亚房产", "时间", "垃圾处理器", "医疗卫生", "重庆房地产", "湘西", "沈阳房地产", "点评", "进入房天下家族", "反欺诈联盟", "报案平台", "杭州房地产", "标致", "脱毛器", "房天下家族", "软件", "货车", "深圳房地产", "更多类似问题", "违章查询", "蒸箱", "高级会计师", "福州房地产", "债券", "私信TA", "东莞房地产", "无锡房地产", "燃气热水器", "轿车", "安顺", "北京房地产"}
var containStopWords = []string{"案例"}
var nums = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-"}

func GetWords(si *sitethrougher.SiteInfo) []string {
	sitethrougher.FillSiteLinksDetailHrefText(si)
	var r []string
	for _, link := range si.SiteLinks {
		for key, info := range link.DetailHrefTexts {
			key = clear(key)
			key = clearTailNums(key)
			klen := len(strings.Split(key, ""))
			if klen >= 20 || klen <= 1 {
				continue
			}
			stop := false
			for _, word := range stopWords {
				if key == word {
					stop = true
					break
				}
			}
			for _, word := range containStopWords {
				if strings.Contains(key, word) {
					stop = true
					break
				}
			}
			if stop {
				break
			}
			if info.Count >= 1 {
				if ok, _ := wordproperty.EnvWordProperty(key); ok {
					continue
				}
				r = append(r, key)

			}
		}
	}

	return removeDuplicate(r)
}

func clearTailNums(key string) string {
	for true {
		c := false
		for _, num := range nums {
			if strings.HasSuffix(key, num) {
				key = strings.ReplaceAll(key, num, "")
				c = true
			}
		}
		if !c {
			break
		}
	}
	return key
}

func clear(key string) string {
	key = strings.ReplaceAll(key, " ", "")
	key = strings.ReplaceAll(key, "\r", "")
	key = strings.ReplaceAll(key, "\n", "")
	key = strings.ReplaceAll(key, "\t", "")
	return key
}