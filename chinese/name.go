package chinese

import (
	"math/rand"
	"strconv"
	"time"
)

func NewName() string {
	FirstNames := []string{
		"李", "王", "张", "阮", "成",
		"刘", "陈", "杨", "黄", "赵", "周", "吴", "徐", "孙", "朱", "马", "胡", "郭", "林",
		"何", "高", "梁", "郑", "罗", "宋", "谢", "唐", "韩", "曹", "许", "邓", "萧", "冯",
		"曾", "程", "蔡", "彭", "潘", "袁", "於", "董", "余", "苏", "叶", "吕", "魏", "蒋",
		"田", "杜", "丁", "沈", "姜", "范", "江", "傅", "钟", "卢", "汪", "戴", "崔", "任",
		"陆", "廖", "姚", "方", "金", "邱", "夏", "谭", "韦", "贾", "邹", "石", "熊", "孟",
		"秦", "阎", "薛", "侯", "雷", "白", "龙", "段", "郝", "孔", "邵", "史", "毛", "常",
		"万", "顾", "赖", "武", "康", "贺", "严", "尹", "钱", "施", "牛", "洪", "龚", "东方",
		"夏侯", "诸葛", "尉迟", "皇甫", "宇文", "鲜于", "西门", "司马", "独孤", "公孙", "慕容", "轩辕",
		"左丘", "欧阳", "皇甫", "上官", "闾丘", "令狐",
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := r.Intn(len(FirstNames))
	chineseName := FirstNames[idx]
	nu := 1 + r.Intn(2)
	for i := 0; i < nu; i++ {
		chineseName += string(30000 + r.Intn(3000))
		// chineseName += string(19968 + r.Intn(20901))
		//chineseName += string(19968 + r.Intn(1000))
	}
	return chineseName
}

//随机生成手机号码
func MobilePhone() (phone string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	array := [15]string{"130", "131", "133", "136", "139", "150", "151", "156", "177", "178", "185", "186", "188", "189", "158"}
	phone += array[r.Intn(15)]
	for i := 0; i < 8; i++ {
		phone += strconv.Itoa(r.Intn(10))
	}
	return phone
}
