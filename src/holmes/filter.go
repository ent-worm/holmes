package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

const (
	YES = iota // is a human
	NO         // is not a human
	UNKNOWN
)

func Filter(c chan int) {
	var accesslogLine string
	var accesslog AccessLog
	var filterResult int
	var i int
	redisConn := NewRedisConn()
	defer CloseRedisConn(redisConn)
	for {
		_, accesslogLine = redisConn.BlockListRightPop("accesslog", 5)
		if accesslogLine == "" {
			fmt.Printf("%s now list have no log to process,continue to wait others to add log to list\n", time.Now())
			continue
		}

		accesslog = GetLog(accesslogLine)
		// fmt.Printf("%dfilter==>%s\n", i, accesslogLine)
		if i%10000 == 0 {
			fmt.Printf("%s %d\n", time.Now(), i)
		}
		i++

		filterResult = DoFilter(redisConn, accesslog)
		if filterResult == YES {
			redisConn.ListLeftPush("accesslog_yes", accesslogLine)
		} else if filterResult == NO {
			redisConn.ListLeftPush("accesslog_no", accesslogLine)
		} else {
			redisConn.ListLeftPush("accesslog_unkown", accesslogLine)
		}
	}
	c <- 1
}

func DoFilter(redisConn RedisConn, accesslog AccessLog) int {
	//if matched, err := regexp.MatchString("^/prop/view", accesslog.RequestURI); err == nil && matched {
	//	if matched, err := regexp.MatchString("^2", accesslog.HttpCode); err == nil && matched {
	//		fmt.Println(accesslog.RequestURI, accesslog.HttpCode)
	//	}
	//}
	//FilterFlag := UNKNOWN
	//if TRUE == ValidClickFilter(redisConn, accesslog);{
	//    if FilterFlag =
	//} else if
	//switch {
	//case (UNKNOWN == GUIDFilter(redisConn, accesslog)):
	//	return FilterFlag
	//case (UNKNOWN == IPFilter(redisConn, accesslog)):
	//	return FilterFlag
	//}
	return URIFilter(redisConn, accesslog)
}

func URIFilter(redisConn RedisConn, accesslog AccessLog) int {
	if matched, err := regexp.MatchString("^/prop/view", accesslog.RequestURI); err == nil && matched {
		return HttpCodeFilter(redisConn, accesslog)
	} else {
		return UNKNOWN //Analysis(redisConn,accesslog)
	}
}

func HttpCodeFilter(redisConn RedisConn, accesslog AccessLog) int {
	if matched, err := regexp.MatchString("^2", accesslog.HttpCode); err == nil && matched {
		return SpiderFilter(redisConn, accesslog)
	} else {
		return UNKNOWN
	}
}

func SpiderFilter(redisConn RedisConn, accesslog AccessLog) int {
	res, err := http.Get("http://www.useragentstring.com/?usa=" + accesslog.UserAgent + "&getText=all")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success", res)
	}
	return UNKNOWN
}

func WhiteIpFilter(redisConn RedisConn, accesslog AccessLog) int {
	//if {
	//} else {

	//}
	return UNKNOWN
}

//func AddWatchingList(redisConn RedisConn, accesslog AccessLog) {
//
//}
//
//func DelWatchingList(redisConn RedisConn, accesslog AccessLog) {
//
//}
//
//func AddWhiteList(redisConn RedisConn, accesslog AccessLog){
//
//}
//
//func AddIgnoreList(redisConn RedisConn, accesslog Accesslog){
//
//}

//func GUIDFilter(redisConn RedisConn, accesslog AccessLog) int {
//	if accesslog.GUID == "-" {
//		return NO
//	} else {
//		redisConn.ListLeftPush("guid", accesslog.GUID)
//		redisConn.ListLeftPush(accesslog.GUID, "----"+accesslog.Referer)
//		uri := accesslog.LogTimeString() + "==>" + accesslog.RequestURI
//		redisConn.ListLeftPush(accesslog.GUID, uri)
//		return YES
//	}
//}
//
//func IPFilter(redisConn RedisConn, accesslog AccessLog) int {
//	redisConn.SetAdd("ip", accesslog.RemoteAddr)
//	redisConn.ListLeftPush(accesslog.RemoteAddr, "----"+accesslog.Referer)
//	uri := accesslog.LogTimeString() + "==>" + accesslog.RequestURI
//	redisConn.ListLeftPush(accesslog.RemoteAddr, uri)
//	return UNKNOWN
//}
