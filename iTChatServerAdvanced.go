package main

import (
	"encoding/json"
	"fmt"
	"time"

	// "html"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	// "strconv"
	"regexp"
	"strings"
)

type UserLogin struct {
	Flag int8
	User UserA
}

type UserInfo struct {
	Flag  int8
	Token string
	Time  string
	RTime string
	User  UserA
}

//flag 0:注册 1:登陆
type UserA struct {
	Name     string
	Password string
	Nickname string
	Sweet    string
	Level    string
}

type Message struct {
	Status     bool
	ReturnInfo int
	Token      string
	Time       string
	Dec        string
	// User       UserA
}

type UserCC struct {
	Src string
	Des string
}

type UserChat struct {
	Token string
	Mes   string
	Users UserCC
}

type Object interface{}

type ChatMessage struct {
	MseFrom []Message2
	MseTo   []Message3
}

type Message2 struct {
	// UserName string
	// FromUser string
	ToUser string
	Time   string
	Mes    string
}

type Message3 struct {
	// UserName string
	FromUser string
	// ToUser   string
	Time string
	Mes  string
}

type AUser struct {
	UName     string
	UNickName string
	ULevel    string
	USweet    string
}

type GetUsers struct {
	Users []AUser
}

type UpdateUserInfo struct {
	Token    string
	Nickname string
	Sweeth   string
	Level    string
}

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var coder = base64.NewEncoding(base64Table)

func base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}

func prerr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

//
var weatherJsonData = []byte(`{
  "message" : "",
  "data" : {
    "flag" : "179481",
    "lts" : "1443357622",
    "items" : [
      {
        "width" : "100",
        "component" : {
          "monthColor" : "129,129,129,255",
          "showId" : "179498",
          "weekDayBgUrl" : "http:\/\/m2.pimg.cn\/images\/images\/20130518\/301d624c-cfc2-4467-956b-c9f7951de662.png",
          "picUrl" : "http:\/\/s2.pimg.cn\/group6\/M00\/D1\/80\/wKgBjVYFDD2AZ9GvAAFKmvdyJR4109.jpg?imageMogr2\/thumbnail\/314x%3E",
          "showTimeColor" : "255,56,141,255",
          "month" : "09\/15",
          "monthOnly" : "09",
          "backgroundUrl" : "",
          "componentType" : "hangtag",
          "weekDay" : "周日",
          "showTypeId" : "9730",
          "xingQi" : "星期日",
          "day" : "27",
          "year" : "2015",
          "publishColor" : "129,129,129,255",
          "showTime" : "20:30",
          "weekDayColor" : "129,129,129,255",
          "actions" : [
            {

            },
            {
              "unixtime" : 1443084822,
              "id" : "183027",
              "actionType" : "detail",
              "type" : "thread"
            }
          ],
          "showType" : "timedrop;;thread",
          "dayColor" : "85,85,85,255"
        },
        "height" : "150"
      },
      {
        "component" : {
          "itemsCount" : "1",
          "description" : "田园风碎花长袖百褶裙",
          "trackValue" : "star_latest_252857",
          "showTypeId" : "252857",
          "id" : "252857",
          "showId" : "179497",
          "picUrl" : "http:\/\/s4.pimg.cn\/group5\/M00\/88\/43\/wKgBf1X6ibyAYrl2AAKcKn5TIvs088.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "田园风碎花长袖百褶裙",
            "trackValue" : "star_latest_252857",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s6.pimg.cn\/group6\/M00\/EA\/84\/wKgBjFVeBBiAC_fxAAnp5GExsxM46.jpeg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "lianyiqun"
              }
            ],
            "id" : "252857",
            "width" : "644",
            "title" : "",
            "userName" : "玫瑰如我",
            "userId" : "93",
            "actionType" : "starDetail",
            "collectionCount" : "85",
            "height" : "867",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m2.pimg.cn\/images\/images\/54\/9\/83\/cfi_133699427880_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s5.pimg.cn\/group5\/M00\/88\/43\/wKgBf1X6ibyAYrl2AAKcKn5TIvs088.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月27日20点"
          },
          "showType" : "star;;",
          "collectionCount" : "85",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1443355207.00",
        "width" : "644",
        "height" : "867"
      },
      {
        "component" : {
          "itemsCount" : "18",
          "description" : "日韩达人，黑色棒球帽+白色长袖T恤+黑色背心+黑色运动裤+...",
          "trackValue" : "star_latest_250626",
          "showTypeId" : "250626",
          "id" : "250626",
          "showId" : "179496",
          "picUrl" : "http:\/\/s.pimg.cn\/group5\/M00\/7C\/BC\/wKgBf1XxQAyALNvNAAW1EiFmqEA816.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "日韩达人，黑色棒球帽+白色长袖T恤+黑色背心+黑色运动裤+白色运动鞋",
            "trackValue" : "star_latest_250626",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s1.pimg.cn\/group6\/M00\/CE\/9D\/wKgBjVYD8yKAQjDdAAPhT9ZQWnk131.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "shangyi"
              },
              {
                "picUrl" : "http:\/\/s2.pimg.cn\/group6\/M00\/CE\/A6\/wKgBjFYD9yiAEHCMAAS088GtxIY851.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "shangyi"
              },
              {
                "picUrl" : "http:\/\/s3.pimg.cn\/group6\/M00\/CE\/B6\/wKgBjFYD-8OAUZAZAAIfI3ov_6k197.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "ku"
              },
              {
                "picUrl" : "http:\/\/s4.pimg.cn\/group6\/M00\/CE\/9B\/wKgBjVYD8juAb-Z8AAF5ARvivPM074.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "xie"
              },
              {
                "picUrl" : "http:\/\/s5.pimg.cn\/group6\/M00\/CE\/A7\/wKgBjFYD91SACeStAAmbVAtR6ro182.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "peishi"
              }
            ],
            "id" : "250626",
            "width" : "500",
            "title" : "",
            "userName" : "黄小红",
            "userId" : "94",
            "actionType" : "starDetail",
            "collectionCount" : "210",
            "height" : "749",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m3.pimg.cn\/images\/images\/39\/38\/73\/hdg_133699427882_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s0.pimg.cn\/group5\/M00\/7C\/BC\/wKgBf1XxQAyALNvNAAW1EiFmqEA816.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月27日20点"
          },
          "showType" : "star;;",
          "collectionCount" : "210",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1443355210.00",
        "width" : "500",
        "height" : "749"
      },
      {
        "component" : {
          "itemsCount" : "14",
          "description" : "本土达人，蓝色牛仔外套+蓝色条纹上衣+白色系带短裤...",
          "trackValue" : "star_latest_249865",
          "showTypeId" : "249865",
          "id" : "249865",
          "showId" : "179495",
          "picUrl" : "http:\/\/s6.pimg.cn\/group6\/M00\/B6\/66\/wKgBjVXvmZeACrHOAAcQZ5v97uA851.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "本土达人，蓝色牛仔外套+蓝色条纹上衣+白色系带短裤+白色休闲鞋+白色棒球帽",
            "trackValue" : "star_latest_249865",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s0.pimg.cn\/group6\/M00\/BE\/27\/wKgBjFX2d7-APfJmAAS3J1l1BCg542.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "waitao"
              },
              {
                "picUrl" : "http:\/\/s1.pimg.cn\/group5\/M00\/81\/B5\/wKgBf1X2eISAEvdrAAGZBa0nLQg741.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "shangyi"
              },
              {
                "picUrl" : "http:\/\/s2.pimg.cn\/group5\/M00\/7F\/D6\/wKgBfVX2eXqAAK4RAAEz9f1Uy6g787.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "ku"
              },
              {
                "picUrl" : "http:\/\/s3.pimg.cn\/group6\/M00\/EA\/7D\/wKgBjVVdv02AdFSxAAVWbakQ3GA489.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "xie"
              },
              {
                "picUrl" : "http:\/\/s4.pimg.cn\/group6\/M00\/BE\/36\/wKgBjVX2evCACNHVAALtY3PJRmU689.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "peishi"
              }
            ],
            "id" : "249865",
            "width" : "595",
            "title" : "",
            "userName" : "孤单芭蕾",
            "userId" : "75",
            "actionType" : "starDetail",
            "collectionCount" : "893",
            "height" : "848",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m4.pimg.cn\/images\/images\/86\/50\/69\/rcn_133699423635_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s.pimg.cn\/group6\/M00\/B6\/66\/wKgBjVXvmZeACrHOAAcQZ5v97uA851.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月19日11点"
          },
          "showType" : "star;;",
          "collectionCount" : "893",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1442631605.00",
        "width" : "595",
        "height" : "848"
      },
      {
        "component" : {
          "itemsCount" : "8",
          "description" : "精品，卡其色翻领风衣，面料选用优质风衣料，经过水洗工艺更...",
          "trackValue" : "star_latest_251897",
          "showTypeId" : "251897",
          "id" : "251897",
          "showId" : "179494",
          "picUrl" : "http:\/\/s5.pimg.cn\/group6\/M00\/C1\/05\/wKgBjVX4C86ALKbvAAPKVoh9QK4505.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "精品，卡其色翻领风衣，面料选用优质风衣料，经过水洗工艺更有韧性；踩线细致工整，牢固且有质感，有型的大翻领设计，同系质感纽扣运用，大口袋的设计潇洒有型，向经典致敬。",
            "trackValue" : "star_latest_251897",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s.pimg.cn\/group6\/M00\/CE\/C9\/wKgBjFYEAcuAEuOtAATRnJEmSx0216.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "waitao"
              }
            ],
            "id" : "251897",
            "width" : "586",
            "title" : "",
            "userName" : "日坛左边",
            "userId" : "58",
            "actionType" : "starDetail",
            "collectionCount" : "136",
            "height" : "773",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m5.pimg.cn\/images\/images\/36\/55\/81\/ytd_133699423576_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s6.pimg.cn\/group6\/M00\/C1\/05\/wKgBjVX4C86ALKbvAAPKVoh9QK4505.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月27日20点"
          },
          "showType" : "star;;",
          "collectionCount" : "136",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1443355202.00",
        "width" : "586",
        "height" : "773"
      },
      {
        "component" : {
          "itemsCount" : "16",
          "description" : "闵孝琳，白色衬衫+黄色阔腿裤+裸色厚底坡跟鞋+黑色复古蛤蟆墨...",
          "trackValue" : "star_latest_249872",
          "showTypeId" : "249872",
          "id" : "249872",
          "showId" : "179493",
          "picUrl" : "http:\/\/s0.pimg.cn\/group5\/M00\/8C\/B8\/wKgBf1X_zT6ANADwAAP6Jp1YycY464.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "闵孝琳，白色衬衫+黄色阔腿裤+裸色厚底坡跟鞋+黑色复古蛤蟆墨镜",
            "trackValue" : "star_latest_249872",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s2.pimg.cn\/group6\/M00\/6E\/F8\/wKgBjFXEPpCATw3pAAMzwtNM_qQ093.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "shangyi"
              },
              {
                "picUrl" : "http:\/\/s3.pimg.cn\/group6\/M00\/CF\/59\/wKgBjVYEOjeAAI4KAAaiXkA_GIU053.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "ku"
              },
              {
                "picUrl" : "http:\/\/s4.pimg.cn\/group5\/M00\/92\/CD\/wKgBf1YEImyALuayAAVxdV-tKFo974.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "xie"
              },
              {
                "picUrl" : "http:\/\/s5.pimg.cn\/group5\/M00\/EF\/74\/wKgBf1WUtmiAVwcJAAE7PDsWgro862.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "peishi"
              }
            ],
            "id" : "249872",
            "width" : "600",
            "title" : "",
            "userName" : "恋物志",
            "userId" : "83",
            "actionType" : "starDetail",
            "collectionCount" : "103",
            "height" : "900",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m6.pimg.cn\/images\/images\/62\/10\/29\/jwe_133699427791_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s1.pimg.cn\/group5\/M00\/8C\/B8\/wKgBf1X_zT6ANADwAAP6Jp1YycY464.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月27日20点"
          },
          "showType" : "star;;",
          "collectionCount" : "103",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1443355204.00",
        "width" : "600",
        "height" : "900"
      },
      {
        "component" : {
          "itemsCount" : "19",
          "description" : "日韩达人，Naning9蓝色V领宽松针织衫+Nan...",
          "trackValue" : "star_latest_247658",
          "showTypeId" : "247658",
          "id" : "247658",
          "showId" : "179492",
          "picUrl" : "http:\/\/s6.pimg.cn\/group6\/M00\/C9\/1F\/wKgBjFX_zBGAe5wLAAItfiDj8KA766.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "日韩达人，Naning9蓝色V领宽松针织衫+Naning9蓝色阔腿牛仔裤+白色厚底休闲鞋+蓝色大帽沿礼帽\n",
            "trackValue" : "star_latest_247658",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s0.pimg.cn\/group5\/M00\/93\/FD\/wKgBf1YEfi2AB2sjAAWNRwO9DSg584.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "shangyi"
              },
              {
                "picUrl" : "http:\/\/s1.pimg.cn\/group5\/M00\/85\/4E\/wKgBf1X4yr-AKeyWAAGejuvl1dA468.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "ku"
              },
              {
                "picUrl" : "http:\/\/s2.pimg.cn\/group6\/M00\/CE\/E9\/wKgBjFYEC1SAFxshAAd0BnpWo9U300.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "xie"
              },
              {
                "picUrl" : "http:\/\/s3.pimg.cn\/group5\/M00\/90\/75\/wKgBfVYD-RaAJDU1AAO9Oc6Sppc624.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "peishi"
              }
            ],
            "id" : "247658",
            "width" : "580",
            "title" : "",
            "userName" : "豆子美眉",
            "userId" : "86",
            "actionType" : "starDetail",
            "collectionCount" : "366",
            "height" : "734",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m.pimg.cn\/images\/images\/60\/48\/24\/bue_133699427855_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s.pimg.cn\/group6\/M00\/C9\/1F\/wKgBjFX_zBGAe5wLAAItfiDj8KA766.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月27日20点"
          },
          "showType" : "star;;",
          "collectionCount" : "366",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1443355212.00",
        "width" : "580",
        "height" : "734"
      },
      {
        "component" : {
          "itemsCount" : "11",
          "description" : "江一燕，卡其色格子马甲+白色V领衬衫+卡其色格子短裤+黑色高...",
          "trackValue" : "star_latest_251183",
          "showTypeId" : "251183",
          "id" : "251183",
          "showId" : "179490",
          "picUrl" : "http:\/\/s4.pimg.cn\/group6\/M00\/C5\/83\/wKgBjFX7bxGAVqtfAAM4C_sPP2w636.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "江一燕，卡其色格子马甲+白色V领衬衫+卡其色格子短裤+黑色高跟短靴",
            "trackValue" : "star_latest_251183",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s6.pimg.cn\/group5\/M00\/89\/5C\/wKgBf1X7kNyANcQZAAFO9HiNfmk323.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "shangyi"
              },
              {
                "picUrl" : "http:\/\/s.pimg.cn\/group6\/M00\/CE\/41\/wKgBjFYD3buAJIVBAAUS3U_Weyc741.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "shangyi"
              },
              {
                "picUrl" : "http:\/\/s0.pimg.cn\/group6\/M00\/C8\/D7\/wKgBjVX_rAGAJkX2AAF3mXLDA34380.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "ku"
              },
              {
                "picUrl" : "http:\/\/s1.pimg.cn\/group6\/M00\/29\/F4\/wKgBjVWTShKAKJbbAAEA8e6od6w376.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "xie"
              }
            ],
            "id" : "251183",
            "width" : "514",
            "title" : "",
            "userName" : "水皮儿",
            "userId" : "55",
            "actionType" : "starDetail",
            "collectionCount" : "57",
            "height" : "792",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m0.pimg.cn\/images\/images\/37\/73\/86\/rsw_133699423567_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s5.pimg.cn\/group6\/M00\/C5\/83\/wKgBjFX7bxGAVqtfAAM4C_sPP2w636.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月27日20点"
          },
          "showType" : "star;;",
          "collectionCount" : "57",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1443355200.00",
        "width" : "514",
        "height" : "792"
      },
      {
        "component" : {
          "itemsCount" : "5",
          "description" : "藤井莉娜，Salire深灰色针织开衫+Salire白色无袖...",
          "trackValue" : "star_latest_250847",
          "showTypeId" : "250847",
          "id" : "250847",
          "showId" : "179488",
          "picUrl" : "http:\/\/s2.pimg.cn\/group5\/M00\/8B\/3F\/wKgBf1X9lIWAfRMyAAHs3cehVFI620.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "藤井莉娜，Salire深灰色针织开衫+Salire白色无袖连衣裙+黑色搭扣高跟鞋",
            "trackValue" : "star_latest_250847",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s4.pimg.cn\/group5\/M00\/92\/59\/wKgBf1YD-d2AfOFVAAFiKoUGj3Y535.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "waitao"
              },
              {
                "picUrl" : "http:\/\/s5.pimg.cn\/group5\/M00\/90\/63\/wKgBfVYD9LeAHRMgAAHy7CgZ5dg424.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "lianyiqun"
              },
              {
                "picUrl" : "http:\/\/s6.pimg.cn\/group5\/M00\/3F\/ED\/wKgBf1XMS1CAGzWfAANNpFUcXLE413.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "xie"
              }
            ],
            "id" : "250847",
            "width" : "510",
            "title" : "",
            "userName" : "栾小欠",
            "userId" : "343609",
            "actionType" : "starDetail",
            "collectionCount" : "85",
            "height" : "680",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m1.pimg.cn\/images\/images\/20140421\/a20954fa-d27a-49de-8d55-5dc34c22555b.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s3.pimg.cn\/group5\/M00\/8B\/3F\/wKgBf1X9lIWAfRMyAAHs3cehVFI620.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月27日20点"
          },
          "showType" : "star;;",
          "collectionCount" : "85",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1443355201.00",
        "width" : "510",
        "height" : "680"
      },
      {
        "component" : {
          "itemsCount" : "1",
          "description" : "中袖长款腰带连体衣",
          "trackValue" : "star_latest_252923",
          "showTypeId" : "252923",
          "id" : "252923",
          "showId" : "179487",
          "picUrl" : "http:\/\/s.pimg.cn\/group6\/M00\/C4\/AF\/wKgBjFX6ifyAbMt8AAMSedMs7wM388.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "中袖长款腰带连体衣",
            "trackValue" : "star_latest_252923",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s1.pimg.cn\/group5\/M00\/D1\/BE\/wKgBfVV5N7OATyuKAATejUTt0eo759.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "ku"
              }
            ],
            "id" : "252923",
            "width" : "538",
            "title" : "",
            "userName" : "下午茶",
            "userId" : "24",
            "actionType" : "starDetail",
            "collectionCount" : "42",
            "height" : "812",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m2.pimg.cn\/images\/images\/35\/40\/10\/xml_133699423448_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s0.pimg.cn\/group6\/M00\/C4\/AF\/wKgBjFX6ifyAbMt8AAMSedMs7wM388.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月27日20点"
          },
          "showType" : "star;;",
          "collectionCount" : "42",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1443355203.00",
        "width" : "538",
        "height" : "812"
      },
      {
        "component" : {
          "itemsCount" : "9",
          "description" : "时尚达人，白色印花夹克+白色印花长裤+黑色印花乐福鞋+银色圆...",
          "trackValue" : "star_latest_249864",
          "showTypeId" : "249864",
          "id" : "249864",
          "showId" : "179486",
          "picUrl" : "http:\/\/s2.pimg.cn\/group6\/M00\/B6\/5F\/wKgBjFXvmZKAYeKbAAbYWXqPofU482.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "时尚达人，白色印花夹克+白色印花长裤+黑色印花乐福鞋+银色圆框太阳镜",
            "trackValue" : "star_latest_249864",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s4.pimg.cn\/group5\/M00\/7F\/C8\/wKgBfVX2dPqABsCdAAOYqyapxck993.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "waitao"
              },
              {
                "picUrl" : "http:\/\/s5.pimg.cn\/group6\/M00\/BE\/25\/wKgBjVX2dSqARgn5AAXhH7U-46s717.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "ku"
              },
              {
                "picUrl" : "http:\/\/s6.pimg.cn\/group5\/M00\/81\/AB\/wKgBf1X2dZWAIYGjAAGkl5evd9g074.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "xie"
              },
              {
                "picUrl" : "http:\/\/s.pimg.cn\/group6\/M00\/7C\/53\/wKgBjVXMMpKAHHl5AAPn0SmYPxI852.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "peishi"
              }
            ],
            "id" : "249864",
            "width" : "619",
            "title" : "",
            "userName" : "花花公主",
            "userId" : "72",
            "actionType" : "starDetail",
            "collectionCount" : "647",
            "height" : "884",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m3.pimg.cn\/images\/images\/81\/37\/33\/cpd_133699423628_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s3.pimg.cn\/group6\/M00\/B6\/5F\/wKgBjFXvmZKAYeKbAAbYWXqPofU482.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月19日13点"
          },
          "showType" : "star;;",
          "collectionCount" : "647",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1442638800.00",
        "width" : "619",
        "height" : "884"
      },
      {
        "component" : {
          "itemsCount" : "17",
          "description" : "时尚达人，黑色镂空夹克外套+白色圆领针织上衣+黑色破洞...",
          "trackValue" : "star_latest_249868",
          "showTypeId" : "249868",
          "id" : "249868",
          "showId" : "179485",
          "picUrl" : "http:\/\/s0.pimg.cn\/group5\/M00\/77\/ED\/wKgBfVXvmt-AbYzXAATtXqjKF4k555.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "时尚达人，黑色镂空夹克外套+白色圆领针织上衣+黑色破洞小脚裤+白色板鞋+黑色手拿包",
            "trackValue" : "star_latest_249868",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s2.pimg.cn\/group6\/M00\/BE\/28\/wKgBjVX2djSAPxvsAAR-ypsMogs020.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "waitao"
              },
              {
                "picUrl" : "http:\/\/s3.pimg.cn\/group5\/M00\/81\/BF\/wKgBf1X2e9yAN26kAAKI7ibtMqg799.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "shangyi"
              },
              {
                "picUrl" : "http:\/\/s4.pimg.cn\/group6\/M00\/8A\/A1\/wKgBjVT7_QWATaF-AAG4m2r2xWc67.jpeg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "ku"
              },
              {
                "picUrl" : "http:\/\/s5.pimg.cn\/group6\/M00\/C5\/FD\/wKgBjVU4biOAYl2sAADTNmKfGVw410.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "bao"
              },
              {
                "picUrl" : "http:\/\/s6.pimg.cn\/group6\/M00\/5D\/D5\/wKgBjVXBsyCAeJpAAAC9UkTZBv8950.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "xie"
              }
            ],
            "id" : "249868",
            "width" : "514",
            "title" : "",
            "userName" : "德云社的粉",
            "userId" : "42",
            "actionType" : "starDetail",
            "collectionCount" : "1335",
            "height" : "765",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m4.pimg.cn\/images\/images\/46\/49\/65\/ikz_133699423511_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s1.pimg.cn\/group5\/M00\/77\/ED\/wKgBfVXvmt-AbYzXAATtXqjKF4k555.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月19日17点"
          },
          "showType" : "star;;",
          "collectionCount" : "1335",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1442653208.00",
        "width" : "514",
        "height" : "765"
      },
      {
        "component" : {
          "itemsCount" : "17",
          "description" : "霍思燕，黑色短款皮夹克+LaneCrawford黑色...",
          "trackValue" : "star_latest_250867",
          "showTypeId" : "250867",
          "id" : "250867",
          "showId" : "179483",
          "picUrl" : "http:\/\/s.pimg.cn\/group6\/M00\/C2\/F0\/wKgBjFX5RWmARoiwAAdAchkT02o496.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "霍思燕，黑色短款皮夹克+LaneCrawford黑色无袖T恤+粉色百褶半身裙+裸色单肩包",
            "trackValue" : "star_latest_250867",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s1.pimg.cn\/group5\/M00\/5E\/2A\/wKgBf1Xeg6-AKW2gAAMZUfiPEwc004.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "waitao"
              },
              {
                "picUrl" : "http:\/\/s2.pimg.cn\/group6\/M00\/C2\/F3\/wKgBjVX5QmiAado-AACofbKfJ8M020.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "shangyi"
              },
              {
                "picUrl" : "http:\/\/s3.pimg.cn\/group5\/M00\/86\/79\/wKgBf1X5QguAJa0XAAGNWhY8L_U224.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "qun"
              },
              {
                "picUrl" : "http:\/\/s4.pimg.cn\/group5\/M00\/92\/05\/wKgBf1YD5HeAKgyyAANPf4PS5ZE741.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "bao"
              }
            ],
            "id" : "250867",
            "width" : "626",
            "title" : "",
            "userName" : "美而立",
            "userId" : "2",
            "actionType" : "starDetail",
            "collectionCount" : "126",
            "height" : "939",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m5.pimg.cn\/images\/images\/69\/43\/82\/tkc_133699423336_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s0.pimg.cn\/group6\/M00\/C2\/F0\/wKgBjFX5RWmARoiwAAdAchkT02o496.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月27日20点"
          },
          "showType" : "star;;",
          "collectionCount" : "126",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1443355213.00",
        "width" : "626",
        "height" : "939"
      },
      {
        "component" : {
          "itemsCount" : "1",
          "description" : "真丝棉不规则连衣裙",
          "trackValue" : "star_latest_252927",
          "showTypeId" : "252927",
          "id" : "252927",
          "showId" : "179482",
          "picUrl" : "http:\/\/s5.pimg.cn\/group5\/M00\/86\/63\/wKgBfVX6ijCAf8jBAAT_U9fdxVE654.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "真丝棉不规则连衣裙",
            "trackValue" : "star_latest_252927",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s.pimg.cn\/group5\/M00\/EC\/BF\/wKgBf1WSQVyANXCIAAgTOAKRxug255.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "lianyiqun"
              }
            ],
            "id" : "252927",
            "width" : "542",
            "title" : "",
            "userName" : "走路撞到猫",
            "userId" : "60",
            "actionType" : "starDetail",
            "collectionCount" : "85",
            "height" : "819",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m6.pimg.cn\/images\/images\/29\/19\/37\/dze_133699423582_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s6.pimg.cn\/group5\/M00\/86\/63\/wKgBfVX6ijCAf8jBAAT_U9fdxVE654.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月27日20点"
          },
          "showType" : "star;;",
          "collectionCount" : "85",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1443355205.00",
        "width" : "542",
        "height" : "819"
      },
      {
        "component" : {
          "itemsCount" : "28",
          "description" : "精品，几何印花衬衫连衣裙。个性简约的图案印花，大方...",
          "trackValue" : "star_latest_250654",
          "showTypeId" : "250654",
          "id" : "250654",
          "showId" : "179481",
          "picUrl" : "http:\/\/s0.pimg.cn\/group5\/M00\/7A\/F7\/wKgBfVXxTJSAbEkeAAf-ZEg-dXE135.jpg?imageMogr2\/thumbnail\/314x%3E",
          "componentType" : "waterfallStarCell",
          "action" : {
            "videoUrl" : "",
            "description" : "精品，几何印花衬衫连衣裙。个性简约的图案印花，大方得体的OL气质。时尚简约明亮而充满活力的风格。",
            "trackValue" : "star_latest_250654",
            "itemPicUrlList" : [
              {
                "picUrl" : "http:\/\/s2.pimg.cn\/group6\/M00\/A1\/04\/wKgBjVXkH9eAMFewAANnlLByKME628.jpg?imageMogr2\/thumbnail\/100x%3E",
                "part" : "lianyiqun"
              }
            ],
            "id" : "250654",
            "width" : "790",
            "title" : "",
            "userName" : "月光泡泡",
            "userId" : "50",
            "actionType" : "starDetail",
            "collectionCount" : "112",
            "height" : "942",
            "commentCount" : "0",
            "userPicUrl" : "http:\/\/m.pimg.cn\/images\/images\/44\/29\/77\/ckq_133699423545_avatar.jpg?imageMogr2",
            "normalPicUrl" : "http:\/\/s1.pimg.cn\/group5\/M00\/7A\/F7\/wKgBfVXxTJSAbEkeAAf-ZEg-dXE135.jpg?imageMogr2\/thumbnail\/640x%3E",
            "dateTime" : "9月27日20点"
          },
          "showType" : "star;;",
          "collectionCount" : "112",
          "commentCount" : "0",
          "hasVideo" : "0"
        },
        "timestamp" : "1443355206.00",
        "width" : "790",
        "height" : "942"
      }
    ],
    "pin" : "179498",
    "tip" : "首页更新了116850条内容。",
    "appApi" : ""
  }
}
`)

// 公钥和私钥可以从文件中读取
var weatherJsonData_2 = []byte(`{"pub": "2013-06-29 22:59",
		   "name": "郑州",
		   "wind": {
		       "chill": 81,
		       "direction": 140,
		       "speed": 7
		   },
		   "astronomy": {
		       "sunrise": "6:05",
		       "sunset": "19:34"
		   },
		   "atmosphere": {
		       "humidity": 89,
		       "visibility": 6.21,
		       "pressure": 29.71,
		       "rising": 1
		   },
		   "forecasts": [
		       {
		           "date": "2013-06-29",
		           "day": 6,
		           "code": 29,
		           "text": "局部多云",
		           "low": 26,
		           "high": 32,
		           "image_large": "http://weather.china.xappengine.com/static/w/img/d29.png",
		           "image_small": "http://weather.china.xappengine.com/static/w/img/s29.png"
		       },
		       {
		           "date": "2013-06-30",
		           "day": 0,
		           "code": 30,
		           "text": "局部多云",
		           "low": 25,
		           "high": 33,
		           "image_large": "http://weather.china.xappengine.com/static/w/img/d30.png",
		           "image_small": "http://weather.china.xappengine.com/static/w/img/s30.png"
		       },
		       {
		           "date": "2013-07-01",
		           "day": 1,
		           "code": 37,
		           "text": "局部雷雨",
		           "low": 24,
		           "high": 32,
		           "image_large": "http://weather.china.xappengine.com/static/w/img/d37.png",
		           "image_small": "http://weather.china.xappengine.com/static/w/img/s37.png"
		       },
		       {
		           "date": "2013-07-02",
		           "day": 2,
		           "code": 38,
		           "text": "零星雷雨",
		           "low": 25,
		           "high": 32,
		           "image_large": "http://weather.china.xappengine.com/static/w/img/d38.png",
		           "image_small": "http://weather.china.xappengine.com/static/w/img/s38.png"
		       },
		       {
		           "date": "2013-07-03",
		           "day": 3,
		           "code": 38,
		           "text": "零星雷雨",
		           "low": 25,
		           "high": 32,
		           "image_large": "http://weather.china.xappengine.com/static/w/img/d38.png",
		           "image_small": "http://weather.china.xappengine.com/static/w/img/s38.png"
		       }
		   ]
		}
`)

var listJsonData_1 = []byte(`{
    "result": [
        {
            "name": "行前事项",
            "orderList": 0
        },
        {
            "name": "文件备份",
            "orderList": 0
        },
        {
            "name": "资金",
            "orderList": 0
        },
        {
            "name": "服装",
            "orderList": 0
        },
        {
            "name": "个护/化妆",
            "orderList": 0
        },
        {
            "name": "医疗/健康",
            "orderList": 0
        },
        {
            "name": "电子/数码",
            "orderList": 0
        },
        {
            "name": "潜水装备",
            "orderList": 0
        },
        {
            "name": "杂项",
            "orderList": 0
        },
        {
            "name": "旅途备忘",
            "orderList": 0
        },
        {
            "name": "购买机票",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "预定酒店",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "购买旅行保险",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "办理签证",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "换钱",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "打疫苗（特殊国家",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "结算必要账单",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "配行李牌",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "给电子设备充电",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "确认护照尾页的签名",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "手机设定目的地时区",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "U盘预装中文输入法",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "安装其他旅行APP",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "设置电子邮件自动回复",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "把行程表给家人/朋友",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "锁好贵重物品",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "妥善安置宠物和植物",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "关水电煤气",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "关好门窗",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "倒掉垃圾",
            "orderList": 1,
            "parentId": 1
        },
        {
            "name": "机票及复印件",
            "orderList": 1,
            "parentId": 2
        },
        {
            "name": "护照签证及复印件",
            "orderList": 1,
            "parentId": 2
        },
        {
            "name": "行程单",
            "orderList": 1,
            "parentId": 2
        },
        {
            "name": "酒店预订单及复印件",
            "orderList": 1,
            "parentId": 2
        },
        {
            "name": "邀请函及复印件",
            "orderList": 1,
            "parentId": 2
        },
        {
            "name": "保险单及复印件",
            "orderList": 1,
            "parentId": 2
        },
        {
            "name": "车票及复印件",
            "orderList": 1,
            "parentId": 2
        },
        {
            "name": "驾照及复印件",
            "orderList": 1,
            "parentId": 2
        },
        {
            "name": "学生证",
            "orderList": 1,
            "parentId": 2
        },
        {
            "name": "青旅会员卡",
            "orderList": 1,
            "parentId": 2
        },
        {
            "name": "2寸护照照片2张",
            "orderList": 1,
            "parentId": 2
        },
        {
            "name": "紧急联系人信息",
            "orderList": 1,
            "parentId": 2
        },
        {
            "name": "现金",
            "orderList": 1,
            "parentId": 3
        },
        {
            "name": "信用卡",
            "orderList": 1,
            "parentId": 3
        },
        {
            "name": "银联卡",
            "orderList": 1,
            "parentId": 3
        },
        {
            "name": "T恤衫",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "长袖衬衫",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "毛衣",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "冲锋衣",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "外套",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "短裤",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "牛仔裤",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "快干裤",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "袜子",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "内衣裤",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "凉鞋",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "运动鞋",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "登山鞋",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "拖鞋",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "帽子",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "围巾",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "泳装",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "墨镜",
            "orderList": 1,
            "parentId": 4
        },
        {
            "name": "洗发水/护发素",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "沐浴液",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "牙膏牙刷",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "毛巾",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "防晒霜/晒后修复",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "护肤/化妆用品",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "剃须刀",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "香水",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "指甲刀",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "隐形眼镜",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "眼药水",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "润唇膏",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "棉签",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "卫生巾",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "梳子",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "镜子",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "电吹风",
            "orderList": 1,
            "parentId": 5
        },
        {
            "name": "止疼药",
            "orderList": 1,
            "parentId": 6
        },
        {
            "name": "感冒药",
            "orderList": 1,
            "parentId": 6
        },
        {
            "name": "止泻/通便药",
            "orderList": 1,
            "parentId": 6
        },
        {
            "name": "晕车/船药",
            "orderList": 1,
            "parentId": 6
        },
        {
            "name": "维生素",
            "orderList": 1,
            "parentId": 6
        },
        {
            "name": "消炎药",
            "orderList": 1,
            "parentId": 6
        },
        {
            "name": "跌打损伤药",
            "orderList": 1,
            "parentId": 6
        },
        {
            "name": "创可贴",
            "orderList": 1,
            "parentId": 6
        },
        {
            "name": "防虫驱蚊",
            "orderList": 1,
            "parentId": 6
        },
        {
            "name": "避孕套/避孕药",
            "orderList": 1,
            "parentId": 6
        },
        {
            "name": "仁丹",
            "orderList": 1,
            "parentId": 6
        },
        {
            "name": "电源转接器",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "插头转换器",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "各类充电器",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "各类数据线",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "移动电源",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "数码伴侣",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "MP3",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "耳机",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "笔记本/平板电脑",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "电池",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "相机/DV",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "镜头",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "胶卷",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "三脚架",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "存储卡",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "相机包",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "读卡器",
            "orderList": 1,
            "parentId": 7
        },
        {
            "name": "毯子",
            "orderList": 1,
            "parentId": 8
        },
        {
            "name": "潜水面镜",
            "orderList": 1,
            "parentId": 8
        },
        {
            "name": "呼吸管",
            "orderList": 1,
            "parentId": 8
        },
        {
            "name": "脚蹼",
            "orderList": 1,
            "parentId": 8
        },
        {
            "name": "浮潜鞋",
            "orderList": 1,
            "parentId": 8
        },
        {
            "name": "潜水衣",
            "orderList": 1,
            "parentId": 8
        },
        {
            "name": "旅行指南",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "充气枕/耳塞/眼罩",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "多用军刀(需托运)",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "记事本/笔",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "行李密码锁",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "雨伞/雨衣",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "防水分类袋",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "望远镜",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "水壶",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "针线包",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "手电筒",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "纸巾",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "小礼品",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "火柴",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "香烟",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "睡袋",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "绳子",
            "orderList": 1,
            "parentId": 9
        },
        {
            "name": "检查手机漫游状态",
            "orderList": 1,
            "parentId": 10
        },
        {
            "name": "给家人报平安",
            "orderList": 1,
            "parentId": 10
        },
        {
            "name": "兑换当地货币",
            "orderList": 1,
            "parentId": 10
        },
        {
            "name": "去本地找旅客询问处",
            "orderList": 1,
            "parentId": 10
        },
        {
            "name": "给自己和朋友寄明信片",
            "orderList": 1,
            "parentId": 10
        },
        {
            "name": "与当地人聊聊",
            "orderList": 1,
            "parentId": 10
        },
        {
            "name": "确认回程机票",
            "orderList": 1,
            "parentId": 10
        },
        {
            "name": "别忘了退税",
            "orderList": 1,
            "parentId": 10
        }
    ],
    "struts": "1"
}
`)

// 公钥和私钥可以从文件中读取
var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDvMazRTUjkfJpE7DvldfHZUrZVLgdx6hco520btXsr1mdjqybp
fhJODJUcIssYGDNLRcjh0S6RyGT+qh6FDEvQ98I/yYq66PCjRTPgUfJdItncBv8S
jcs3Okefv2CJmhocCK11FEozZPEd3P2bF0XYjbyLbcnq3oUiZtenRGeeVwIDAQAB
AoGASrT0Kgb+bkawlDhIWNmmqN7Zje8rahvYEfF+NXpQNxfnAM0zARhcNT5e0APZ
9POSCb+JB2ajKesyCAwwLhPyFX68VQtH86AazqEjLR3e0z/F3vE1FReV8KR1bsyi
i91mmpVb3fTzhsi9dE0wNVGaxRl/sZaFp4TkL3SwwHpGm9ECQQDxPdzTn0vKLjEL
ZzzLZn2lUfK9Ztdsck0vv4s+l5pbd6FqJ9276nMkkaJttEdqySoxHilGd2oEzHVp
0iKyCjV5AkEA/dO+k8EmZ8fnBtZA1aCIWTwi+uCHUBWjQrkqbS67pIIv1FuNWsgi
pqZh2bZUPSnj7LlCXw/Kvu+J0i72iFiOTwJAOovZ6N3zBck6C9ttLKvd+F4v+/lW
dLI0u07QG0utoV8iJGIydOWMNibF9bvXzTmu7Ka2O6zFZQ69vAXMd8r0eQJBAMLO
UgOgR9N6rqqmoRfTnxGtf8M/s1oZYTWCWzd0mHrHl+HJahF0bHOuWob20mwmzFEQ
VgoTWq1ztjjj5j36iS0CQAkYVhPj7oVE5HmxJ5FRgTlY12fWjGNPXEz3H+nu6+Px
p0ideQeiI4ubemYfQIIrAx7pJaoPfgNnqxF70QmrEAk=
-----END RSA PRIVATE KEY-----
`)

// 解密
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

var decrypted string

func init() {
	flag.StringVar(&decrypted, "d", "", "加密过的数据")
	flag.Parse()
}

//正则匹配 邮箱
func isEmail(email string) bool {
	if email != "" {
		if isOk, _ := regexp.MatchString("^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$", email); isOk {
			return true
		}
		return false
	}
	return false
}

//打开数据库连接
func opendb(dbstr string) *sql.DB {
	db, err := sql.Open("mysql", dbstr)
	prerr(err)
	return db
}

//查询数据库
func query(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM user")
	prerr(err)
	defer rows.Close()
	// 普通demo
	// fmt.Println("")
	// cols, _ := rows.Columns()
	// for i := range cols {
	// 	fmt.Print(cols[i])
	// 	fmt.Print("\t")
	// }
	// fmt.Println("")

	// for rows.Next() {
	// 	var name string
	// 	var passw string
	// 	var token string
	// 	var statu int8
	// 	err = rows.Scan(&name, &passw, &token, &statu)
	// 	prerr(err)
	// 	fmt.Println(name)
	// 	fmt.Print("\t")
	// 	fmt.Println(passw)
	// 	fmt.Print("\t")
	// 	fmt.Println(token)
	// 	fmt.Print("\t")
	// 	fmt.Println(statu)
	// 	fmt.Print("\t\r\n")
	// }
	fmt.Println("=============================查询==============================")
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		fmt.Println(record)
	}
	fmt.Println("=================================================================")
}

//查询数据库 GET all user list
func getUserInfo(db *sql.DB, uName string) AUser {
	var ui AUser

	rows, err := db.Query("SELECT * FROM user where name=?", uName)

	prerr(err)
	defer rows.Close()

	fmt.Println("=============================查询好友==============================")
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	//name=?,passw=?,nickname=?,sweetheart=?,level=?,token=?,statu=?,rtime=?,ltime=?")

	if rows.Next() {
		err = rows.Scan(scanArgs...)
		prerr(err)
		// isLogined = true

		for i, col := range values {
			if col != nil {
				switch os := columns[i]; os {
				case "name":
					// fmt.Println(os, ":", string(col.([]byte)))
					ui.UName = string(col.([]byte))
				case "nickname":
					// fmt.Println(os, ":", string(col.([]byte)))
					ui.UNickName = string(col.([]byte))
				case "level":
					ui.ULevel = string(col.([]byte))
				case "sweetheart":
					// fmt.Println(os, ":", string(col.([]byte)))
					ui.USweet = string(col.([]byte))
				default:
					// fmt.Printf("%s.", os)
				}
			}
		}

	}

	return ui
}

//查询数据库 GET all user list
func getAllUsers(db *sql.DB) GetUsers {

	var user1 GetUsers

	rows, err := db.Query("SELECT * FROM user")
	prerr(err)
	defer rows.Close()

	fmt.Println("=============================查询好友==============================")
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		var u AUser
		for i, col := range values {
			if col != nil {
				// fmt.Println(columns[i], ":", string(col.([]byte)))
				switch os := columns[i]; os {
				case "name":
					// fmt.Println(os, ":", string(col.([]byte)))
					u.UName = string(col.([]byte))
				case "nickname":
					// fmt.Println(os, ":", string(col.([]byte)))
					u.UNickName = string(col.([]byte))
				case "level":
					// 	// fmt.Println(os, ":", string(col.([]byte)))
					u.ULevel = string(col.([]byte))
				case "sweetheart":
					// fmt.Println(os, ":", string(col.([]byte)))
					u.USweet = string(col.([]byte))
				default:
					// fmt.Printf("%s.", os)
				}

			}
		}

		user1.Users = append(user1.Users, u)
	}
	fmt.Println("=================================================================")
	return user1

}

//插入指定数据
func insert(db *sql.DB, u UserInfo) int64 {
	fmt.Println("=============================Insert==============================")
	// fmt.Println(u.User.Name, u.User.Password, u.Token, u.Flag, u.RTime, u.Time)

	stmt, err := db.Prepare("INSERT INTO user SET name=?,passw=?,nickname=?,sweetheart=?,level=?,token=?,statu=?,rtime=?,ltime=?")
	prerr(err)

	res, err := stmt.Exec(u.User.Name, u.User.Password, u.User.Nickname, u.User.Sweet, u.User.Level, u.Token, u.Flag, u.RTime, u.Time)
	prerr(err)

	id, err := res.LastInsertId()
	prerr(err)

	fmt.Println(id)
	fmt.Println("==================================================================")
	return id

}

//插入指定数据
func insertMessage(db *sql.DB, um UserChat) int64 {

	stmt, err := db.Prepare("INSERT INTO message SET f=?,t=?,mes=?,time=?")
	prerr(err)

	t_str := time.Now().Format("2006-01-02 15:04:05")

	res, err := stmt.Exec(um.Users.Src, um.Users.Des, um.Mes, t_str)
	prerr(err)

	id, err := res.LastInsertId()
	prerr(err)

	fmt.Println(id)

	return id

}

func checkMessageDB(db *sql.DB, um UserChat) ChatMessage {
	var mess ChatMessage
	rows, err := db.Query("SELECT * FROM message where f=?", um.Users.Src)
	prerr(err)

	fmt.Println("=============================查询==============================")
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		// record := make(map[string]string)
		// for i, col := range values {
		// 	if col != nil {
		// 		record[columns[i]] = string(col.([]byte))
		// 	}
		// }
		// fmt.Println(record)

		var mesTo Message2

		for i, col := range values {
			if col != nil {
				// fmt.Println(columns[i], ":", string(col.([]byte)))
				switch os := columns[i]; os {
				// case "f":
				// fmt.Println(os, ":", string(col.([]byte)))
				case "t":
					// fmt.Println(os, ":", string(col.([]byte)))
					mesTo.ToUser = string(col.([]byte))
				case "mes":
					// fmt.Println(os, ":", string(col.([]byte)))
					mesTo.Mes = string(col.([]byte))
				case "time":
					// fmt.Println(os, ":", string(col.([]byte)))
					mesTo.Time = string(col.([]byte))
				default:
					// fmt.Printf("%s.", os)
				}

			}
		}

		mess.MseFrom = append(mess.MseFrom, mesTo)

	}

	// fmt.Println(mess)

	rows, err = db.Query("SELECT * FROM message where t=?", um.Users.Des)
	prerr(err)

	columns, _ = rows.Columns()
	scanArgs = make([]interface{}, len(columns))
	values = make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)

		var mesFrom Message3

		for i, col := range values {
			if col != nil {
				// fmt.Println(columns[i], ":", string(col.([]byte)))
				switch os := columns[i]; os {
				// case "f":
				// fmt.Println(os, ":", string(col.([]byte)))
				case "f":
					// fmt.Println(os, ":", string(col.([]byte)))
					mesFrom.FromUser = string(col.([]byte))
				case "mes":
					// fmt.Println(os, ":", string(col.([]byte)))
					mesFrom.Mes = string(col.([]byte))
				case "time":
					// fmt.Println(os, ":", string(col.([]byte)))
					mesFrom.Time = string(col.([]byte))
				default:
					// fmt.Printf("%s.", os)
				}

			}
		}

		mess.MseTo = append(mess.MseTo, mesFrom)

	}

	// fmt.Println(mess)
	// r, err2 := json.Marshal(mess)
	// if err2 != nil {
	// 	fmt.Println(err2)
	// }
	// fmt.Println(r)
	fmt.Println(um.Users.Src, "talk to:", um.Users.Des)
	fmt.Println("=================================================================")

	return mess
}

func checkQuery(db *sql.DB, userName string) bool {
	rows, err := db.Query("SELECT * FROM user where name=?", userName)

	prerr(err)
	defer rows.Close()

	// var isLogined = false
	// var uname string
	// var passw string
	// var token string
	// var status int8
	// var rTime string
	// var lTime string
	// var nickName string
	// var level int
	// var sweetheart string
	// if rows.Next() {
	// 	err = rows.Scan(&uname, &passw, &token, &nickName, &level, &sweetheart, &status, &rTime, &lTime)
	// 	prerr(err)
	// 	isLogined = true
	// }

	// fmt.Println("=============================Check==============================")
	// fmt.Println(uname, passw, token)
	// fmt.Println("================================================================")

	var isLogined = false
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	if rows.Next() {
		err = rows.Scan(scanArgs...)
		prerr(err)
		isLogined = true
	}

	fmt.Println(isLogined)

	return isLogined
}

func checkPasswQuery(db *sql.DB, user UserInfo) bool {
	rows, err := db.Query("SELECT * FROM user where passw=?", user.User.Password)
	prerr(err)
	defer rows.Close()

	var isLogined = false
	var uname string
	var passw string
	var token string
	var status int8
	var rTime string
	var lTime string
	var nickName string
	var level string
	var sweetheart string

	if rows.Next() {
		err = rows.Scan(&uname, &passw, &token, &nickName, &level, &sweetheart, &status, &rTime, &lTime)
		prerr(err)
		isLogined = true
	}

	// for rows.Next() {
	// 	if err := rows.Scan(&passw); err == nil {

	// 	} else {
	// 		prerr(err)
	// 	}
	// }
	return isLogined
}

func checkTokenQuery(db *sql.DB, user UserInfo) bool {
	rows, err := db.Query("SELECT * FROM user where token=?", user.Token)
	prerr(err)
	defer rows.Close()

	var isLogined = false
	var uname string
	var passw string
	var token string
	var status int8
	var rTime string
	var lTime string
	var nickName string
	var level string
	var sweetheart string

	if rows.Next() {
		err = rows.Scan(&uname, &passw, &token, &nickName, &level, &sweetheart, &status, &rTime, &lTime)
		prerr(err)
		isLogined = true
	}

	fmt.Println("===========================Token================================")
	fmt.Println(uname, passw, token)
	fmt.Println("================================================================")
	// for rows.Next() {
	// 	if err := rows.Scan(&token); err == nil {
	// 		isLogined = true
	// 		return isLogined
	// 	} else {
	// 		prerr(err)
	// 	}
	// }
	return isLogined
}

//更新数据
func update(db *sql.DB, uname string) {
	rows, err := db.Prepare("update user set statu=?,ltime=? where name=?")
	prerr(err)
	defer rows.Close()

	t_str := time.Now().Format("2006-01-02 15:04:05")
	res, err := rows.Exec(1, t_str, uname)
	prerr(err)

	affect, err := res.RowsAffected()
	prerr(err)
	fmt.Println(affect)
}

//更新数据
func updateUserToken(db *sql.DB, user UserInfo) bool {
	var isUpdataed = true

	rows, err := db.Prepare("update user set statu=?,token=?,ltime=? where name=?")
	prerr(err)
	if err != nil {
		isUpdataed = false
	}

	defer rows.Close()

	t_str := time.Now().Format("2006-01-02 15:04:05")
	res, err := rows.Exec(1, user.Token, t_str, user.User.Name)
	prerr(err)
	if err != nil {
		isUpdataed = false
	}

	affect, err := res.RowsAffected()
	prerr(err)
	fmt.Println(affect)
	if err != nil {
		isUpdataed = false
	}

	return isUpdataed
}

//退出登录 刷新 state 以及token
func updateUserExit(db *sql.DB, user UserInfo) bool {
	var isUpdataed = true

	rows, err := db.Prepare("update user set statu=?,token=?,ltime=? where name=?")
	prerr(err)
	if err != nil {
		isUpdataed = false
	}

	defer rows.Close()

	user.Token = ""

	t_str := time.Now().Format("2006-01-02 15:04:05")
	res, err := rows.Exec(0, user.Token, t_str, user.User.Name)
	prerr(err)
	if err != nil {
		isUpdataed = false
	}

	affect, err := res.RowsAffected()
	prerr(err)
	fmt.Println(affect)
	if err != nil {
		isUpdataed = false
	}

	return isUpdataed
}

//更新用户数据
func updateUserInfoDB(db *sql.DB, user UpdateUserInfo) bool {
	var isUpdataed = true
	//name=?,passw=?,nickname=?,sweetheart=?,level=?,token=?,statu=?,rtime=?,ltime=?
	rows, err := db.Prepare("update user set nickname=?,sweetheart=?,level=? where token=?")
	prerr(err)
	if err != nil {
		isUpdataed = false
	}

	defer rows.Close()

	// t_str := time.Now().Format("2006-01-02 15:04:05")
	res, err := rows.Exec(user.Nickname, user.Sweeth, user.Level, user.Token)
	prerr(err)
	if err != nil {
		isUpdataed = false
	}

	affect, err := res.RowsAffected()
	prerr(err)
	fmt.Println(affect)
	if err != nil {
		isUpdataed = false
	}

	return isUpdataed
}

//判断是否是文件夹
func isDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
	return false
}

// var staticHandler http.Handler

// 静态文件处理
func StaticServer(w http.ResponseWriter, req *http.Request) {
	// fmt.Println("path:" + req.URL.Path)
	// staticHandler.ServeHTTP(w, req)

	// fmt.Println("path:" + req.URL.Path)
	// if req.URL.Path != "/down/" {
	// 	staticHandler.ServeHTTP(w, req)
	// 	return
	// }

	// fmt.Fprintf(w, "no file.")
}

// func init() {
// 	staticHandler = http.StripPrefix("./", http.FileServer(http.Dir("radio")))
// }

// const (
//   regular ="^(13[0-9]|14[57]|15[0-35-9]|18[07-9])\d{8}$"
// )

func checkEmail(name string) {
	match, _ := regexp.MatchString(name, "^([w-.]+)@(([[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.)|(([w-]+.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(]?)$")
	fmt.Print(match)
}

func main() {

	//设置访问的路由
	//登陆
	http.HandleFunc("/api/login", handler)
	//注册
	http.HandleFunc("/api/register", handlerR)
	//主页数据
	http.HandleFunc("/api/home/data", handlerH)
	//主页数据2
	http.HandleFunc("/api/home/data2", handlerH2)
	//用户退出
	http.HandleFunc("/api/exit", handlerE)
	//用户上传文件
	http.HandleFunc("/api/upload", handleUploadFile)
	//用户下载指定文件
	http.HandleFunc("/api/download", handleDownloadFile)
	//用户聊天数据
	http.HandleFunc("/api/chat", handleChat)

	http.HandleFunc("/api/friends/chat", handleChatTest)
	//用户获取好友信息
	http.HandleFunc("/api/friends", handleFriends)
	//更新用户个人信息
	http.HandleFunc("/api/update/userinfo", handleUpdateUserInfo)
	http.HandleFunc("/api/assets/", StaticServer)
	//设置监听的端口
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func handleUpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("client:", r.RemoteAddr, "method:", r.Method)
	//解析参数 默认是不会解析的
	r.ParseForm()
	if r.Method == "GET" {
		fmt.Println("method :", r.Method)

		var returnD Message
		returnD.Status = true
		returnD.ReturnInfo = 200
		returnD.Token = ""
		returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		returnD.Dec = "Warning!!! The Request DOT SUPPORT GET , SUPPORT POST!!!!"

		r, err2 := json.Marshal(returnD)
		if err2 != nil {
			fmt.Println(err2)
		}

		fmt.Println(string(r))
		fmt.Fprintf(w, "%s", r)
	} else if r.Method == "POST" {
		fmt.Println("method :", r.Method)

		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		var ui UpdateUserInfo
		json.Unmarshal([]byte(result), &ui)
		//b64解密
		tDec, _ := base64.URLEncoding.DecodeString(ui.Token)

		// fmt.Println(string(tDec))
		ui.Token = string(tDec)
		// fmt.Println(ui.Nickname)
		// fmt.Println(ui.Level)
		// fmt.Println(ui.Sweeth)
		// fmt.Println(ui.Token)

		// query(db)
		//解析token
		user_token := string(tDec)
		// fmt.Printf("%q\n", strings.Split(user_token, " "))
		str := strings.Split(user_token, " ")
		// fmt.Println("=================", str[0], len(str))

		user_name := str[0]
		if len(user_name) < 0 {
			// defer db.Close()
			return
		}
		//查询数据库有没有此人
		db := opendb("root:@tcp(localhost:3307)/boxDB?charset=utf8")
		ret := checkQuery(db, user_name)
		if ret == true {
			//验证token
			var user UserInfo
			user.Token = user_token
			user.User.Name = user_name

			ret = checkTokenQuery(db, user)
			if ret == true {
				//更新用户信息

				ret = updateUserInfoDB(db, ui)
				if ret == true {
					//将用户信息返回

					aUser := getUserInfo(db, user_name)
					// if aUser == true {
					r2, err2 := json.Marshal(aUser)
					if err2 != nil {
						fmt.Println(err2)
					}

					fmt.Println(string(r2))
					fmt.Fprintf(w, "%s", r2)
					// }

				} else {
					var returnD Message
					returnD.Status = false
					returnD.ReturnInfo = 200
					returnD.Token = ""
					returnD.Time = time.Now().Format("2006-01-02 15:04:05")
					returnD.Dec = "更新失败,查看token或联系 Box!!!!!!!!"

					r2, err2 := json.Marshal(returnD)
					if err2 != nil {
						fmt.Println(err2)
					}

					fmt.Println(string(r2))
					fmt.Fprintf(w, "%s", r2)
				}

			} else {
				var returnD Message
				returnD.Status = false
				returnD.ReturnInfo = 200
				returnD.Token = ""
				returnD.Time = time.Now().Format("2006-01-02 15:04:05")
				returnD.Dec = "Token 过期!!!!!!!!"

				r2, err2 := json.Marshal(returnD)
				if err2 != nil {
					fmt.Println(err2)
				}

				fmt.Println(string(r2))
				fmt.Fprintf(w, "%s", r2)
			}

		} else {
			var returnD Message
			returnD.Status = false
			returnD.ReturnInfo = 200
			returnD.Token = ""
			returnD.Time = time.Now().Format("2006-01-02 15:04:05")
			returnD.Dec = "用户不存在!!!!!!!!"

			r2, err2 := json.Marshal(returnD)
			if err2 != nil {
				fmt.Println(err2)
			}

			fmt.Println(string(r2))
			fmt.Fprintf(w, "%s", r2)
		}

		// var returnD Message
		// returnD.Status = true
		// returnD.ReturnInfo = 200
		// returnD.Token = ""
		// returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		// returnD.Dec = "Warning!!! The Request DOT SUPPORT GET , SUPPORT POST!!!!"

		// r, err2 := json.Marshal(returnD)
		// if err2 != nil {
		// 	fmt.Println(err2)
		// }

		// fmt.Println(string(r))
		// fmt.Fprintf(w, "%s", r)

		defer db.Close()
	}
}

func handleFriends(w http.ResponseWriter, r *http.Request) {
	fmt.Println("client:", r.RemoteAddr, "method:", r.Method)
	//解析参数 默认是不会解析的
	r.ParseForm()

	if r.Method == "GET" {
		fmt.Println("method :", r.Method)
		fmt.Println("token :", r.Form["token"])
		//解析token
		tDec, _ := base64.URLEncoding.DecodeString(r.Form["token"][0])
		fmt.Println(string(tDec))

		user_token := string(tDec)
		// fmt.Printf("%q\n", strings.Split(user_token, " "))
		str := strings.Split(user_token, " ")
		fmt.Println("========用户名=========", str[0], ":", len(str))

		user_name := str[0]
		if len(user_name) < 0 {
			return
		}

		//判断token 是否过期
		//查询数据库
		db := opendb("root:@tcp(localhost:3307)/boxDB?charset=utf8")
		var user UserInfo
		user.Token = user_token
		user.User.Name = user_name

		ret := checkTokenQuery(db, user)
		if ret == true {

			//获取所有好友列表
			// var user1 GetUsers

			allPeople := getAllUsers(db)
			r2, err2 := json.Marshal(allPeople)
			if err2 != nil {
				fmt.Println(err2)
			}

			fmt.Println(string(r2))
			fmt.Fprintf(w, "%s", r2)

		} else { //token 错误
			var returnD Message
			returnD.Status = false
			returnD.ReturnInfo = 200
			returnD.Token = ""
			returnD.Time = time.Now().Format("2006-01-02 15:04:05")
			returnD.Dec = "Token 过期!!!!!!!!"

			r2, err2 := json.Marshal(returnD)
			if err2 != nil {
				fmt.Println(err2)
			}

			fmt.Println(string(r2))
			fmt.Fprintf(w, "%s", r2)
		}

		defer db.Close()

	} else if r.Method == "POST" {
		fmt.Println("method :", r.Method)

		var returnD Message
		returnD.Status = true
		returnD.ReturnInfo = 200
		returnD.Token = ""
		returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		returnD.Dec = "Warning!!! The Request DOT SUPPORT POST , SUPPORT GET!!!!"

		r, err2 := json.Marshal(returnD)
		if err2 != nil {
			fmt.Println(err2)
		}

		fmt.Println(string(r))
		fmt.Fprintf(w, "%s", r)
	}
}

func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("client:", r.RemoteAddr, "method:", r.Method)
	r.ParseForm()
	r.ParseMultipartForm(32 << 20) //最大内存为32M

	//读取参数
	fmt.Println("token", r.Form["token"][0])
	if len(r.Form["token"]) > 0 {

		//解析token
		tDec, _ := base64.URLEncoding.DecodeString(r.Form["token"][0])
		fmt.Println(string(tDec))

		user_token := string(tDec)
		// fmt.Printf("%q\n", strings.Split(user_token, " "))
		str := strings.Split(user_token, " ")
		// fmt.Println("=================", str[0], len(str))

		user_name := str[0]
		if len(user_name) < 0 {
			return
		}

		//判断token 是否过期
		//查询数据库
		db := opendb("root:@tcp(localhost:3307)/boxDB?charset=utf8")
		var user UserInfo
		user.Token = user_token
		user.User.Name = user_name

		ret := checkTokenQuery(db, user)

		if ret == false { //token 错误
			var returnD Message
			returnD.Status = false
			returnD.ReturnInfo = 200
			returnD.Token = ""
			returnD.Time = time.Now().Format("2006-01-02 15:04:05")
			returnD.Dec = "Token 过期!!!!!!!!"

			r2, err2 := json.Marshal(returnD)
			if err2 != nil {
				fmt.Println(err2)
			}

			fmt.Println(string(r2))
			fmt.Fprintf(w, "%s", r2)
		} else {
			path := "./userData/" + user_name + "/"

			ret = isDirExists(path)
			if ret == false {
				os.MkdirAll(path, 0700)
			}

			// os.MkdirAll(path, 0700)

			mp := r.MultipartForm
			if mp == nil {
				fmt.Println("not MultipartForm")
				w.Write(([]byte)("不是MultipartForm格式"))
				return
			}

			fileHeaders, findFile := mp.File["file"]
			if !findFile || len(fileHeaders) == 0 {
				fmt.Println("file count == 0.")
				w.Write(([]byte)("没有上传文件"))
				return
			}

			fmt.Println("file headers : ", fileHeaders)

			for _, v := range fileHeaders {
				fileName := v.Filename
				file, err := v.Open()

				if err != nil {
					fmt.Println("Open file error. + ", err)
				}

				defer file.Close()

				outputFilePath := path + fileName

				writer, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					fmt.Println("Open local file error. + ", err)
				}

				io.Copy(writer, file)
			}

			// msg := fmt.Sprintf("成功上传了%d个文件", len(fileHeaders))
			// w.Write(([]byte)(msg))
			var returnD Message
			returnD.Status = false
			returnD.ReturnInfo = 200
			returnD.Token = ""
			returnD.Time = time.Now().Format("2006-01-02 15:04:05")
			returnD.Dec = "write sucess!!!!!!!!"

			r2, err2 := json.Marshal(returnD)
			if err2 != nil {
				fmt.Println(err2)
			}

			fmt.Println(string(r2))
			fmt.Fprintf(w, "%s", r2)

		}
		defer db.Close()
	}

}

func handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("client:", r.RemoteAddr, "method:", r.Method)
	r.ParseForm()
	r.ParseMultipartForm(32 << 20) //最大内存为32M

	//读取参数
	fmt.Println("token", r.Form["token"][0])
	if len(r.Form["token"]) > 0 {

		//解析token
		tDec, _ := base64.URLEncoding.DecodeString(r.Form["token"][0])
		fmt.Println(string(tDec))

		user_token := string(tDec)
		// fmt.Printf("%q\n", strings.Split(user_token, " "))
		str := strings.Split(user_token, " ")
		// fmt.Println("=================", str[0], len(str))
		user_name := str[0]
		if len(user_name) < 0 {
			return
		}

		//判断token 是否过期
		//查询数据库
		db := opendb("root:@tcp(localhost:3307)/boxDB?charset=utf8")
		var user UserInfo
		user.Token = user_token
		user.User.Name = user_name

		ret := checkTokenQuery(db, user)

		if ret == false { //token 错误
			var returnD Message
			returnD.Status = false
			returnD.ReturnInfo = 200
			returnD.Token = ""
			returnD.Time = time.Now().Format("2006-01-02 15:04:05")
			returnD.Dec = "Token 过期!!!!!!!!"

			r2, err2 := json.Marshal(returnD)
			if err2 != nil {
				fmt.Println(err2)
			}

			fmt.Println(string(r2))
			fmt.Fprintf(w, "%s", r2)
		} else {
			// 获取文件名称以及存储路径
			path := "./userData/" + user_name + "/"
			filename := r.URL.Path[strings.LastIndex(path, "/"):]

			// 打开文件
			fin, err := os.Open(path + filename)
			defer fin.Close()
			if err != nil {
				// log.Fatal("static resource:", err)
				//没有文件
				var returnD Message
				returnD.Status = true
				returnD.ReturnInfo = 200
				returnD.Token = ""
				returnD.Time = time.Now().Format("2006-01-02 15:04:05")
				returnD.Dec = "文件不存在!!!!!!!!"

				r2, err2 := json.Marshal(returnD)
				if err2 != nil {
					fmt.Println(err2)
				}

				fmt.Println(string(r2))
				fmt.Fprintf(w, "%s", r2)
			} else {
				// 读文件
				fd, _ := ioutil.ReadAll(fin)

				// 通过http socket 写回文件
				w.Write(fd)
			}

		}
		defer db.Close()
	}

}

//退出登录
func handlerE(w http.ResponseWriter, r *http.Request) {
	fmt.Println("client:", r.RemoteAddr, "method:", r.Method)
	//解析参数 默认是不会解析的
	r.ParseForm()

	if r.Method == "GET" {
		fmt.Println("method :", r.Method)
		fmt.Println("token :", r.Form["token"])

		if len(r.Form["token"]) > 0 {
			fmt.Println(r.Form["token"][0])
			// base64解密
			tDec, _ := base64.URLEncoding.DecodeString(r.Form["token"][0])

			fmt.Println(string(tDec))

			user_token := string(tDec)
			str := strings.Split(user_token, " ")
			//获取用户名
			user_name := str[0]
			//用户名存在
			if len(user_name) > 0 {
				//查询数据库
				db := opendb("root:@tcp(localhost:3307)/boxDB?charset=utf8")

				var u_info UserInfo
				u_info.Token = user_token
				u_info.User.Name = user_name
				//验证token
				ret := checkTokenQuery(db, u_info)
				if ret == true {

					//将token 移除 重新登录
					ret := updateUserExit(db, u_info)
					if ret == true {

						var returnD Message
						returnD.Status = false
						returnD.ReturnInfo = 200
						returnD.Token = ""
						returnD.Time = time.Now().Format("2006-01-02 15:04:05")
						returnD.Dec = "退出成功，需要重新登录!!!"

						r, err2 := json.Marshal(returnD)
						if err2 != nil {
							fmt.Println(err2)
						}

						fmt.Println(string(r))
						fmt.Fprintf(w, "%s", r)

					} else {
						var returnD Message
						returnD.Status = false
						returnD.ReturnInfo = 200
						returnD.Token = ""
						returnD.Time = time.Now().Format("2006-01-02 15:04:05")
						returnD.Dec = "退出失败，找找问题!!!"

						r, err2 := json.Marshal(returnD)
						if err2 != nil {
							fmt.Println(err2)
						}

						fmt.Println(string(r))
						fmt.Fprintf(w, "%s", r)
					}

					// // arr := []interface{}{} //interface{} arr := []Object{}
					// var user Message2
					// user.UserName = "box"
					// user.FromUser = "super box"
					// user.ToUser = "super super box"
					// user.Time = "2015.09.23.15:44"

					// var u ChatMessage

					// u.Mse = append(u.Mse, user)
					// u.Mse = append(u.Mse, user)
					// u.Mse = append(u.Mse, user)
					// u.Mse = append(u.Mse, user)
					// u.Mse = append(u.Mse, user)
					// u.Mse = append(u.Mse, user)

					// // fmt.Println(u)

					// // arr = append(arr, user)
					// // arr = append(arr, user)
					// // arr = append(arr, arr)
					// // fmt.Println(arr)

					// r, err2 := json.Marshal(u)
					// if err2 != nil {
					// 	fmt.Println(err2)
					// }
					// // fmt.Println(r)

					// var f []interface{}
					// json.Unmarshal(r, &f)
					// fmt.Println("------------f : ", f)

					// fmt.Fprintf(w, "%s", r)

				} else {
					var returnD Message
					returnD.Status = false
					returnD.ReturnInfo = 200
					returnD.Token = ""
					returnD.Time = time.Now().Format("2006-01-02 15:04:05")
					returnD.Dec = "TOKEN_ERROR!!!"

					r, err2 := json.Marshal(returnD)
					if err2 != nil {
						fmt.Println(err2)
					}

					fmt.Println(string(r))
					fmt.Fprintf(w, "%s", r)
				}

				defer db.Close()

			} else {

				var returnD Message
				returnD.Status = false
				returnD.ReturnInfo = 200
				returnD.Token = ""
				returnD.Time = time.Now().Format("2006-01-02 15:04:05")
				returnD.Dec = "USER WAS DEAD!!!"

				r, err2 := json.Marshal(returnD)
				if err2 != nil {
					fmt.Println(err2)
				}

				fmt.Println(string(r))
				fmt.Fprintf(w, "%s", r)

			}
		}
		// var returnD Message
		// returnD.Status = true
		// returnD.ReturnInfo = 200
		// returnD.Token = ""
		// returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		// returnD.Dec = "Warning!!! The Request DOT SUPPORT GET , SUPPORT POST!!!!"

		// r, err2 := json.Marshal(returnD)
		// if err2 != nil {
		// 	fmt.Println(err2)
		// }

		// fmt.Println(string(r))
		// fmt.Fprintf(w, "%s", r)
	} else if r.Method == "POST" {
		fmt.Println("method :", r.Method)

		var returnD Message
		returnD.Status = true
		returnD.ReturnInfo = 200
		returnD.Token = ""
		returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		returnD.Dec = "Warning!!! The Request DOT SUPPORT POST , SUPPORT GET!!!!"

		r, err2 := json.Marshal(returnD)
		if err2 != nil {
			fmt.Println(err2)
		}

		fmt.Println(string(r))
		fmt.Fprintf(w, "%s", r)

		// result, _ := ioutil.ReadAll(r.Body)
		// r.Body.Close()
		// fmt.Println("结果: ", result)

		// var f interface{}
		// json.Unmarshal(result, &f)
		// fmt.Println("f : ", f)
		// m := f.(map[string]interface{})
		// fmt.Println("============================")
		// fmt.Println("m : ", m)
		// fmt.Println("============================")

		// for k, v := range m {
		// 	switch vv := v.(type) {
		// 	case string:
		// 		fmt.Println(k, "is string", vv)
		// 	case int:
		// 		fmt.Println(k, "is int", vv)
		// 	case float64:
		// 		fmt.Println(k, "is float64", vv)
		// 	case interface{}:
		// 		fmt.Println(k, "is an array----->", vv)
		// 	default:
		// 		fmt.Println(k, "is of a type I don't konw how to handle")
		// 	}
		// }

	}

}

//获取home2 数据
func handlerH2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("client:", r.RemoteAddr, "method:", r.Method)
	//解析参数 默认是不会解析的
	r.ParseForm()

	if r.Method == "GET" {
		fmt.Println("method :", r.Method)
		fmt.Println("Host :", r.Host)
		fmt.Fprintf(w, "%s", listJsonData_1)

	} else {
		var returnD Message
		returnD.Status = true
		returnD.ReturnInfo = 200
		returnD.Token = ""
		returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		returnD.Dec = "Warning!!! The Request DOT SUPPORT POST , SUPPORT GET!!!!"

		r, err2 := json.Marshal(returnD)
		if err2 != nil {
			fmt.Println(err2)
		}

		fmt.Println(string(r))
		fmt.Fprintf(w, "%s", r)
	}

}

//获取home 数据
func handlerH(w http.ResponseWriter, r *http.Request) {
	fmt.Println("client:", r.RemoteAddr, "method:", r.Method)
	//解析参数 默认是不会解析的
	r.ParseForm()

	if r.Method == "GET" {
		fmt.Println("method :", r.Method)
		// fmt.Println("token :", r.Form["token"])
		// fmt.Println("password", r.Form["password"])
		fmt.Println("Host :", r.Host)

		if len(r.Form["token"]) > 0 {
			fmt.Println(r.Form["token"][0])

			// base64解密
			tDec, _ := base64.URLEncoding.DecodeString(r.Form["token"][0])
			fmt.Println(string(tDec))

			user_token := string(tDec)

			// fmt.Printf("%q\n", strings.Split(user_token, " "))

			str := strings.Split(user_token, " ")
			fmt.Println("=================", str[0], len(str))

			user_name := str[0]

			fmt.Println("用户姓名:", user_name)

			if len(user_name) > 0 {
				//查询数据库
				db := opendb("root:@tcp(localhost:3307)/boxDB?charset=utf8")

				var u_info UserInfo
				u_info.Token = user_token
				//验证token
				ret := checkTokenQuery(db, u_info)
				if ret == true {
					fmt.Fprintf(w, "%s", weatherJsonData)
				} else {
					var returnD Message
					returnD.Status = false
					returnD.ReturnInfo = 200
					returnD.Token = ""
					returnD.Time = time.Now().Format("2006-01-02 15:04:05")
					returnD.Dec = "TOKEN_ERROR!!!"

					r, err2 := json.Marshal(returnD)
					if err2 != nil {
						fmt.Println(err2)
					}

					fmt.Println(string(r))
					fmt.Fprintf(w, "%s", r)
				}

				defer db.Close()

			} else {
				var returnD Message
				returnD.Status = false
				returnD.ReturnInfo = 200
				returnD.Token = ""
				returnD.Time = time.Now().Format("2006-01-02 15:04:05")
				returnD.Dec = "USER WAS DEAD!!!"

				r, err2 := json.Marshal(returnD)
				if err2 != nil {
					fmt.Println(err2)
				}

				fmt.Println(string(r))
				fmt.Fprintf(w, "%s", r)
			}

		} else {
			var returnD Message
			returnD.Status = false
			returnD.ReturnInfo = 200
			returnD.Token = ""
			returnD.Time = time.Now().Format("2006-01-02 15:04:05")
			returnD.Dec = "TOKEN MEI YOU!!!"

			r, err2 := json.Marshal(returnD)
			if err2 != nil {
				fmt.Println(err2)
			}

			fmt.Println(string(r))

			fmt.Fprintf(w, "%s", r)
		}

		//从token 获取用户名

		// db := opendb("root:@tcp(localhost:3307)/boxDB?charset=utf8")

		//先找用户

		// var token_str string
		// for k, v := range r.Form {
		// 	fmt.Println("key:", k, ";")
		// 	fmt.Println("val:", strings.Join(v, ""))
		// 	if k == "token" {
		// 		// fmt.Fprintln(w, "HI,Welcome To 微软IT学院, 小鲜肉:", strings.Join(v, ""), "恭喜登陆成功!!!!!,欢迎下次再来\n")
		// 		token_str = strings.Join(v, "")
		// 	}
		// }

		// var returnD Message
		// returnD.Status = true
		// returnD.ReturnInfo = 200
		// returnD.Token = "ABCDEFGHIJKLMN"
		// returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		// returnD.Dec = "Hello golang!!!!!"

		// r, err2 := json.Marshal(returnD)
		// if err2 != nil {
		// 	fmt.Println(err2)
		// }

		// fmt.Println(string(r))
		// fmt.Fprintf(w, "%s", r)

		// fmt.Fprintf(w, "%s", weatherJsonData)

	} else {
		var returnD Message
		returnD.Status = true
		returnD.ReturnInfo = 200
		returnD.Token = ""
		returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		returnD.Dec = "Warning!!! The Request DOT SUPPORT POST , SUPPORT GET!!!!"

		r, err2 := json.Marshal(returnD)
		if err2 != nil {
			fmt.Println(err2)
		}

		fmt.Println(string(r))
		fmt.Fprintf(w, "%s", r)
	}

}

//聊天测试
func handleChatTest(w http.ResponseWriter, r *http.Request) {
	//解析参数 默认是不会解析的
	r.ParseForm()
	if r.Method == "GET" {

		var returnD Message
		returnD.Status = true
		returnD.ReturnInfo = 200
		returnD.Token = ""
		returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		returnD.Dec = "Warning!!! The Request DOT SUPPORT GET , SUPPORT POST!!!!"

		r, err2 := json.Marshal(returnD)
		if err2 != nil {
			fmt.Println(err2)
		}

		fmt.Println(string(r))
		fmt.Fprintf(w, "%s", r)

	} else if r.Method == "POST" {

		fmt.Println("method :", r.Method)

		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		var chatM UserChat
		json.Unmarshal([]byte(result), &chatM)

		//b64解密
		tDec, _ := base64.URLEncoding.DecodeString(chatM.Token)
		mDec, _ := base64.URLEncoding.DecodeString(chatM.Mes)
		sDec, _ := base64.URLEncoding.DecodeString(chatM.Users.Src)
		dDec, _ := base64.URLEncoding.DecodeString(chatM.Users.Des)

		chatM.Token = string(tDec)
		chatM.Mes = string(mDec)
		chatM.Users.Src = string(sDec)
		chatM.Users.Des = string(dDec)

		fmt.Println("=============================CHAT================================")
		fmt.Println("token:", chatM.Token)
		fmt.Println("Mes:", chatM.Mes)
		fmt.Println("src:", chatM.Users.Src)
		fmt.Println("des:", chatM.Users.Des)
		fmt.Println("=================================================================")

		//打开数据库
		db := opendb("root:@tcp(localhost:3307)/boxDB?charset=utf8")
		//插入数据库
		id := insertMessage(db, chatM)

		fmt.Println("r : ", id)

		//查询数据库
		ret := checkMessageDB(db, chatM)

		r, err2 := json.Marshal(ret)
		if err2 != nil {
			fmt.Println(err2)
		}

		// var returnD Message
		// returnD.Status = true
		// returnD.ReturnInfo = 200
		// returnD.Token = ""
		// returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		// returnD.Dec = "Warning!!! The Request DOT SUPPORT GET , SUPPORT POST!!!!"

		// r, err2 := json.Marshal(returnD)
		// if err2 != nil {
		// 	fmt.Println(err2)
		// }

		fmt.Println(string(r))
		fmt.Fprintf(w, "%s", r)

		defer db.Close()
	}
}

//聊天
func handleChat(w http.ResponseWriter, r *http.Request) {
	//解析参数 默认是不会解析的
	r.ParseForm()
	if r.Method == "GET" {

		var returnD Message
		returnD.Status = true
		returnD.ReturnInfo = 200
		returnD.Token = ""
		returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		returnD.Dec = "Warning!!! The Request DOT SUPPORT GET , SUPPORT POST!!!!"

		r, err2 := json.Marshal(returnD)
		if err2 != nil {
			fmt.Println(err2)
		}

		fmt.Println(string(r))
		fmt.Fprintf(w, "%s", r)

	} else if r.Method == "POST" {

		fmt.Println("method :", r.Method)

		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		fmt.Println("结果: ", result)

		var f interface{}
		json.Unmarshal(result, &f)
		fmt.Println("f : ", f)
		m := f.(map[string]interface{})
		fmt.Println("============================")
		fmt.Println("m : ", m)
		fmt.Println("============================")

		for k, v := range m {
			switch vv := v.(type) {
			case string:
				fmt.Println(k, "is string", vv)
			case int:
				fmt.Println(k, "is int", vv)
			case float64:
				fmt.Println(k, "is float64", vv)
			case interface{}:
				fmt.Println(k, "is an array----->", vv)
			default:
				fmt.Println(k, "is of a type I don't konw how to handle")
			}
		}

	}
}

func handlerR(w http.ResponseWriter, r *http.Request) {
	//解析参数 默认是不会解析的
	r.ParseForm()

	//回给client
	// fmt.Fprintln(w, "HI,I love GO language! %s\n", html.EscapeString(r.URL.Path[1:]))

	// //json encode
	// j1 := make(map[string]interface{})
	// j1["name"] = "超级合子"
	// j1["url"] = "http://www.baidu.com"

	// js1, err := json.Marshal(j1)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(string(js1))

	// fmt.Fprintf(w, "%s", js1)

	//这些信息是解析客户端的request信息

	if r.Method == "GET" {
		// fmt.Println("method :", r.Method)
		// fmt.Println("username", r.Form["username"])
		// fmt.Println("password", r.Form["password"])
		// fmt.Println("method :", r.Host)

		// fmt.Println("r.from : ", r.Form)

		// var userName string

		// for k, v := range r.Form {
		// 	fmt.Println("key:", k, ";")
		// 	fmt.Println("val:", strings.Join(v, ""))
		// 	if k == "username" {
		// 		fmt.Fprintln(w, "HI,Welcome To 微软IT学院, 小鲜肉:", strings.Join(v, ""), "恭喜登陆成功!!!!!,欢迎下次再来\n")
		// 		// userName = "HI,Welcome To 微软IT学院, 小鲜肉:" + strings.Join(v, "") + "恭喜登陆成功!!!!!,欢迎下次再来\n"
		// 	}
		// }

		var returnD Message
		returnD.Status = true
		returnD.ReturnInfo = 200
		returnD.Token = ""
		returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		returnD.Dec = "Warning!!! The Request DOT SUPPORT GET , SUPPORT POST!!!!"

		r, err2 := json.Marshal(returnD)
		if err2 != nil {
			fmt.Println(err2)
		}

		fmt.Println(string(r))
		fmt.Fprintf(w, "%s", r)

	} else if r.Method == "POST" {

		// fmt.Println(r.Body)

		fmt.Println("method :", r.Method)

		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		// fmt.Println("结果: ", result)

		// var f interface{}
		// json.Unmarshal(result, &f)
		// fmt.Println("f : ", f)
		// m := f.(map[string]interface{})
		// fmt.Println("============================")
		// fmt.Println("m : ", m)
		// fmt.Println("============================")

		// for k, v := range m {
		// 	switch vv := v.(type) {
		// 	case string:
		// 		fmt.Println(k, "is string", vv)
		// 	case int:
		// 		fmt.Println(k, "is int", vv)
		// 	case float64:
		// 		fmt.Println(k, "is float64", vv)
		// 	case interface{}:
		// 		fmt.Println(k, "is an array----->", vv)
		// 	default:
		// 		fmt.Println(k, "is of a type I don't konw how to handle")
		// 	}
		// }

		// decName, err := base64Decode([]byte(result))
		// if err != nil {
		// 	fmt.Println("---------->", err.Error())
		// }

		// // if hello != string(decName) {
		// // 	fmt.Println("hello is not equal to enbyte")
		// // }

		// fmt.Println(" ================new----------------->: ", string(decName))

		//结构已知,解析到结构体
		var u UserLogin
		json.Unmarshal([]byte(result), &u)

		fmt.Println(u)

		fmt.Println("flag is : ", u.Flag)
		fmt.Println("name is : ", u.User.Name)
		fmt.Println("pass is : ", u.User.Password)
		//base64解密
		uDec, _ := base64.URLEncoding.DecodeString(u.User.Name)
		pDec, _ := base64.URLEncoding.DecodeString(u.User.Password)
		fmt.Println("*************************************************")

		u.User.Name = string(uDec)
		u.User.Password = string(pDec)
		u.User.Level = "100"
		u.User.Nickname = ""
		u.User.Sweet = ""
		fmt.Println("++++++++++++++++++++++++++", u.User.Name)

		retEmail := isEmail(u.User.Name)
		if retEmail {
			fmt.Println("is Email!!!")
		} else {
			fmt.Println("is not Email!!!!")

			var returnD Message
			returnD.Status = true
			returnD.ReturnInfo = 200
			returnD.Token = ""
			returnD.Time = time.Now().Format("2006-01-02 15:04:05")
			returnD.Dec = "Warning!!! your name don't email!!!!"

			r, err2 := json.Marshal(returnD)
			if err2 != nil {
				fmt.Println(err2)
			}

			fmt.Println(string(r))
			fmt.Fprintf(w, "%s", r)

			return
		}

		fmt.Println("*************************************************")

		// decName, err := base64Decode([]byte(u.User.Name))
		// if err != nil {
		// 	fmt.Println("---------->", err.Error())
		// }

		// if hello != string(decName) {
		// 	fmt.Println("hello is not equal to enbyte")
		// }

		// fmt.Println(" ================new name ----------------->: ", string(decName))

		// u.User.Name = string(decName)

		// var data []byte
		// if decrypted != "" {
		// 	data, err := base64.StdEncoding.DecodeString(decrypted)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// }

		// origData, err := RsaDecrypt(decName)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println("------------------------->", string(origData))

		db := opendb("root:@tcp(localhost:3307)/boxDB?charset=utf8")

		// query(db)

		ret := checkQuery(db, u.User.Name)
		if ret == false {
			fmt.Println("no")

			var newUser UserInfo

			newUser.User.Name = u.User.Name
			newUser.User.Password = u.User.Password

			//用 名+密+时间 生产 token
			t := time.Now().Format("2006-01-02 15:04:05")
			// tokenstr := newUser.User.Name + " " + u.User.Password + " " + t

			// fmt.Println(tokenstr)

			newUser.Token = ""
			newUser.Flag = 0
			newUser.RTime = t
			newUser.Time = ""
			newUser.User.Level = "0"
			newUser.User.Nickname = ""
			newUser.User.Sweet = ""

			//插入数据库
			id := insert(db, newUser)
			fmt.Println("------------------------------------------------->", id)

			var returnD Message
			returnD.Status = false
			returnD.ReturnInfo = 200
			returnD.Dec = "注册成功!!!!!请重新登录!!!!"

			// sEnc := base64.StdEncoding.EncodeToString([]byte(tokenstr))

			returnD.Token = ""
			returnD.Time = time.Now().Format("2006-01-02 15:04:05")

			r, err2 := json.Marshal(returnD)
			if err2 != nil {
				fmt.Println(err2)
			}

			fmt.Println(string(r))

			fmt.Fprintf(w, "%s", r)

		} else {

			fmt.Println("Yes")

			var returnD Message
			returnD.Status = false
			returnD.ReturnInfo = 200
			returnD.Dec = "哇,重名了,重来!!!!!"
			returnD.Token = ""
			returnD.Time = time.Now().Format("2006-01-02 15:04:05")

			r, err2 := json.Marshal(returnD)
			if err2 != nil {
				fmt.Println(err2)
			}

			fmt.Println(string(r))

			fmt.Fprintf(w, "%s", r)
		}

		query(db)

		defer db.Close()

		// var returnClientData UserLogin

		// returnClientData.Flag = 20
		// returnClientData.User.Name = u.User.Name
		// returnClientData.User.Password = "token_return"

		// r, err2 := json.Marshal(returnClientData)
		// if err2 != nil {
		// 	fmt.Println(err2)
		// }
		// // fmt.Println(string(r))
		// fmt.Fprintf(w, "%s", r)

	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	//解析参数 默认是不会解析的
	r.ParseForm()

	//这些信息是解析客户端的request信息

	if r.Method == "GET" {
		// fmt.Println("method :", r.Method)
		// fmt.Println("username", r.Form["username"])
		// fmt.Println("password", r.Form["password"])
		// fmt.Println("method :", r.Host)

		// var returnD Message
		// returnD.Status = true
		// returnD.ReturnInfo = 200
		// returnD.Token = "ABCDEFGHIJKLMN"
		// returnD.Time = time.Now().Format("2006-01-02 15:04:05")

		// r, err2 := json.Marshal(returnD)
		// if err2 != nil {
		// 	fmt.Println(err2)
		// }

		// fmt.Println(string(r))
		// fmt.Fprintf(w, "%s", r)
		var returnD Message
		returnD.Status = true
		returnD.ReturnInfo = 200
		returnD.Token = ""
		returnD.Time = time.Now().Format("2006-01-02 15:04:05")
		returnD.Dec = "Warning!!! The Request DOT SUPPORT GET , SUPPORT POST!!!!"

		r, err2 := json.Marshal(returnD)
		if err2 != nil {
			fmt.Println(err2)
		}

		fmt.Println(string(r))
		fmt.Fprintf(w, "%s", r)

	} else if r.Method == "POST" {
		// fmt.Println(r.Body)

		fmt.Println("method :", r.Method)

		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		// fmt.Println("结果: ", result)

		// var f interface{}
		// json.Unmarshal(result, &f)
		// fmt.Println("f : ", f)
		// m := f.(map[string]interface{})
		// fmt.Println("============================")
		// fmt.Println("m : ", m)
		// fmt.Println("============================")

		// for k, v := range m {
		// 	switch vv := v.(type) {
		// 	case string:
		// 		fmt.Println(k, "is string", vv)
		// 	case int:
		// 		fmt.Println(k, "is int", vv)
		// 	case float64:
		// 		fmt.Println(k, "is float64", vv)
		// 	case interface{}:
		// 		fmt.Println(k, "is an array----->", vv)
		// 	default:
		// 		fmt.Println(k, "is of a type I don't konw how to handle")
		// 	}
		// }

		//结构已知,解析到结构体

		var u_info UserLogin
		json.Unmarshal([]byte(result), &u_info)

		fmt.Println("*********************User Login****************************")
		fmt.Println("flag is : ", u_info.Flag)
		fmt.Println("name is : ", u_info.User.Name)
		fmt.Println("pass is : ", u_info.User.Password)
		// fmt.Println("token is : ", u_info.Token)
		// fmt.Println("time is : ", u_info.Time)

		fmt.Println("***********************************************************")

		fmt.Println("***********************B64 Dec*****************************")
		//base64解密
		uDec, _ := base64.URLEncoding.DecodeString(u_info.User.Name)
		pDec, _ := base64.URLEncoding.DecodeString(u_info.User.Password)
		// tDec, _ := base64.URLEncoding.DecodeString(u_info.Token)

		u_info.User.Name = string(uDec)
		u_info.User.Password = string(pDec)
		// u_info.Token = string(tDec)
		//判断用户名是否符合格式
		// fmt.Println("++++++++++++++++++++++++++", u_info.Token)

		retEmail := isEmail(u_info.User.Name)
		if retEmail {
			fmt.Println("is Email!!!")
		} else {
			fmt.Println("is not Email!!!!")

			var returnD Message
			returnD.Status = true
			returnD.ReturnInfo = 200
			returnD.Token = ""
			returnD.Time = time.Now().Format("2006-01-02 15:04:05")
			returnD.Dec = "Warning!!! your name don't email!!!!"

			r, err2 := json.Marshal(returnD)
			if err2 != nil {
				fmt.Println(err2)
			}

			fmt.Println(string(r))
			fmt.Fprintf(w, "%s", r)

			return
		}

		fmt.Println("++++++++++++++++++++++++++", u_info.User.Name, u_info.User.Password)

		fmt.Println("************************************************************")

		db := opendb("root:@tcp(localhost:3307)/boxDB?charset=utf8")

		// query(db)

		ret := checkQuery(db, u_info.User.Name)
		if ret == false {
			fmt.Println("no user!!!!!!!!!")

			var returnD Message
			returnD.Status = false
			returnD.ReturnInfo = 200
			returnD.Dec = "没有 " + u_info.User.Name + " 的小鲜肉!!"

			// sEnc := base64.StdEncoding.EncodeToString([]byte(tokenstr))

			returnD.Token = ""
			returnD.Time = time.Now().Format("2006-01-02 15:04:05")

			r, err2 := json.Marshal(returnD)
			if err2 != nil {
				fmt.Println(err2)
			}

			fmt.Println(string(r))

			fmt.Fprintf(w, "%s", r)

		} else {

			fmt.Println("Yes user alive!!!")

			//用户存在 首先要验证密码正确？
			//验证密码 如果成功就给token 服务器记录并send to client

			var u_info2 UserInfo

			u_info2.Flag = u_info.Flag
			u_info2.User.Name = u_info.User.Name
			u_info2.User.Password = u_info.User.Password
			u_info2.Time = ""
			u_info2.Token = ""

			ret := checkPasswQuery(db, u_info2)

			if ret == true {

				//此处已经不需要验证token了
				//此处要用来给用户登陆的时候 给个新token并更新数据库

				//用 名+密+时间 生产 token
				t := time.Now().Format("2006-01-02 15:04:05")
				tokenstr := u_info2.User.Name + " " + u_info2.User.Password + " " + t

				u_info2.Token = tokenstr

				//存数据库
				ret := updateUserToken(db, u_info2)

				if ret == true {
					var returnD Message
					returnD.Status = true
					returnD.ReturnInfo = 200
					returnD.Dec = "恭喜登陆成功!!!!"
					sEnc := base64.StdEncoding.EncodeToString([]byte(tokenstr))
					returnD.Token = sEnc
					returnD.Time = time.Now().Format("2006-01-02 15:04:05")

					r, err2 := json.Marshal(returnD)
					if err2 != nil {
						fmt.Println(err2)
					}

					fmt.Println(string(r))

					fmt.Fprintf(w, "%s", r)
				} else {
					var returnD Message
					returnD.Status = true
					returnD.ReturnInfo = 200
					returnD.Dec = "恭喜登陆成功!!!!但是数据库录入失败!!!!无法提供token!!!联系super box!!!!"
					returnD.Token = ""
					returnD.Time = time.Now().Format("2006-01-02 15:04:05")

					r, err2 := json.Marshal(returnD)
					if err2 != nil {
						fmt.Println(err2)
					}

					fmt.Println(string(r))

					fmt.Fprintf(w, "%s", r)
				}

				// //验证token用来测试
				// ret := checkTokenQuery(db, u_info)

				// if ret == true {
				// 	//token 正确
				// 	update(db, u_info.User.Name)

				// 	var returnD Message
				// 	returnD.Status = true
				// 	returnD.ReturnInfo = 200
				// 	returnD.Dec = "恭喜登陆成功!!!!"
				// 	returnD.Token = ""
				// 	returnD.Time = time.Now().Format("2006-01-02 15:04:05")

				// 	r, err2 := json.Marshal(returnD)
				// 	if err2 != nil {
				// 		fmt.Println(err2)
				// 	}

				// 	fmt.Println(string(r))

				// 	fmt.Fprintf(w, "%s", r)

				// } else {
				// 	var returnD Message
				// 	returnD.Status = false
				// 	returnD.ReturnInfo = 200
				// 	returnD.Dec = "token error!"
				// 	returnD.Token = ""
				// 	returnD.Time = time.Now().Format("2006-01-02 15:04:05")

				// 	r, err2 := json.Marshal(returnD)
				// 	if err2 != nil {
				// 		fmt.Println(err2)
				// 	}

				// 	fmt.Println(string(r))

				// 	fmt.Fprintf(w, "%s", r)
				// }
			} else {
				var returnD Message
				returnD.Status = false
				returnD.ReturnInfo = 200
				returnD.Dec = "智力为0,忘记密码!!记忆力<5s,鱼的记忆力7s!!"
				returnD.Token = ""
				returnD.Time = time.Now().Format("2006-01-02 15:04:05")

				r, err2 := json.Marshal(returnD)
				if err2 != nil {
					fmt.Println(err2)
				}

				fmt.Println(string(r))

				fmt.Fprintf(w, "%s", r)
			}

		}

		defer db.Close()

	}

}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	//解析参数 默认是不会解析的
// 	r.ParseForm()

// 	//回给client
// 	// fmt.Fprintln(w, "HI,I love GO language! %s\n", html.EscapeString(r.URL.Path[1:]))

// 	// //json encode
// 	// j1 := make(map[string]interface{})
// 	// j1["name"] = "超级合子"
// 	// j1["url"] = "http://www.baidu.com"

// 	// js1, err := json.Marshal(j1)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }

// 	// fmt.Println(string(js1))

// 	// fmt.Fprintf(w, "%s", js1)

// 	//这些信息是解析客户端的request信息

// 	if r.Method == "GET" {
// 		fmt.Println("method :", r.Method)
// 		fmt.Println("username", r.Form["username"])
// 		fmt.Println("password", r.Form["password"])
// 		fmt.Println("method :", r.Host)

// 		// fmt.Println("r.from : ", r.Form)

// 		// var userName string

// 		// for k, v := range r.Form {
// 		// 	fmt.Println("key:", k, ";")
// 		// 	fmt.Println("val:", strings.Join(v, ""))
// 		// 	if k == "username" {
// 		// 		// fmt.Fprintln(w, "HI,Welcome To 微软IT学院, 小鲜肉:", strings.Join(v, ""), "恭喜登陆成功!!!!!,欢迎下次再来\n")
// 		// 		userName = "HI,Welcome To 微软IT学院, 小鲜肉:" + strings.Join(v, "") + "恭喜登陆成功!!!!!,欢迎下次再来\n"
// 		// 	}
// 		// }

// 		// //json encode
// 		// j1 := make(map[string]interface{})
// 		// j1["name"] = "超级合子"
// 		// j1["url"] = "http://www.baidu.com"

// 		// js1, err := json.Marshal(j1)
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// }

// 		// fmt.Println(string(js1))

// 		// jsonStr := `{"host": "http://localhost:9090","port": 9090,"analytics_file": "","static_file_version": 1,"static_dir": "E:/Project/goTest/src/","templates_dir": "E:/Project/goTest/src/templates/","serTcpSocketHost": ":12340","serTcpSocketPort": 12340,"fruits": ["apple", "peach"]}`

// 		// fmt.Fprintf(w, "%s", js1)

// 		var returnD Message
// 		returnD.Status = true
// 		returnD.ReturnInfo = 200
// 		returnD.Token = "ABCDEFGHIJKLMN"
// 		returnD.Time = time.Now().Format("2006-01-02 15:04:05")
// 		// if userName != "" {
// 		// 	returnD.User.Name = userName
// 		// } else {
// 		// 	returnD.User.Name = ""
// 		// }
// 		// returnD.User.Password = ""

// 		r, err2 := json.Marshal(returnD)
// 		if err2 != nil {
// 			fmt.Println(err2)
// 		}

// 		fmt.Println(string(r))
// 		fmt.Fprintf(w, "%s", r)

// 	} else if r.Method == "POST" {

// 		fmt.Println(r.Body)

// 		fmt.Println("method :", r.Method)

// 		result, _ := ioutil.ReadAll(r.Body)
// 		r.Body.Close()
// 		fmt.Println("结果: ", result)

// 		var f interface{}
// 		json.Unmarshal(result, &f)
// 		fmt.Println("f : ", f)
// 		m := f.(map[string]interface{})
// 		fmt.Println("============================")
// 		fmt.Println("m : ", m)
// 		fmt.Println("============================")

// 		for k, v := range m {
// 			switch vv := v.(type) {
// 			case string:
// 				fmt.Println(k, "is string", vv)
// 			case int:
// 				fmt.Println(k, "is int", vv)
// 			case float64:
// 				fmt.Println(k, "is float64", vv)
// 			case interface{}:
// 				fmt.Println(k, "is an array----->", vv)
// 			default:
// 				fmt.Println(k, "is of a type I don't konw how to handle")
// 			}
// 		}

// 		//结构已知,解析到结构体
// 		var u UserLogin
// 		json.Unmarshal([]byte(result), &u)

// 		fmt.Println(u)

// 		fmt.Println("flag is : ", u.Flag)
// 		fmt.Println("name is : ", u.User.Name)
// 		fmt.Println("pass is : ", u.User.Password)

// 		retJson := make(map[string]interface{})
// 		retJson["name"] = "超级合子"
// 		retJson["url"] = "http://www.baidu.com"

// 		d, err := json.Marshal(retJson)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		fmt.Println(string(d))

// 		// fmt.Fprintf(w, "%s", d)

// 		var returnClientData UserLogin

// 		returnClientData.Flag = 20
// 		returnClientData.User.Name = "hello client!!!!!"
// 		returnClientData.User.Password = "token_return"

// 		r, err2 := json.Marshal(returnClientData)
// 		if err2 != nil {
// 			fmt.Println(err2)
// 		}

// 		fmt.Println(string(r))

// 		fmt.Fprintf(w, "%s", r)

// 		// jsonStr := `{"host": "http://localhost:9090","port": 9090,"analytics_file": "","static_file_version": 1,"static_dir": "E:/Project/goTest/src/","templates_dir": "E:/Project/goTest/src/templates/","serTcpSocketHost": ":12340","serTcpSocketPort": 12340,"fruits": ["apple", "peach"]}`

// 		// fmt.Fprintln(w, "HI,I love GO language!")

// 		// fmt.Println("name is : ", u.user.Name)
// 		// fmt.Println("password is : ", u.user.Password)

// 		// for i := 0; i < len(u.user); i++ {
// 		// 	fmt.Println(u.user[i].Name)
// 		// 	fmt.Println(u.user[i].Password)
// 		// }

// 		// 	b := []byte(`{"pub": "2013-06-29 22:59",
// 		//    "name": "南宁",
// 		//    "wind": {
// 		//        "chill": 81,
// 		//        "direction": 140,
// 		//        "speed": 7
// 		//    },
// 		//    "astronomy": {
// 		//        "sunrise": "6:05",
// 		//        "sunset": "19:34"
// 		//    },
// 		//    "atmosphere": {
// 		//        "humidity": 89,
// 		//        "visibility": 6.21,
// 		//        "pressure": 29.71,
// 		//        "rising": 1
// 		//    },
// 		//    "forecasts": [
// 		//        {
// 		//            "date": "2013-06-29",
// 		//            "day": 6,
// 		//            "code": 29,
// 		//            "text": "局部多云",
// 		//            "low": 26,
// 		//            "high": 32,
// 		//            "image_large": "http://weather.china.xappengine.com/static/w/img/d29.png",
// 		//            "image_small": "http://weather.china.xappengine.com/static/w/img/s29.png"
// 		//        },
// 		//        {
// 		//            "date": "2013-06-30",
// 		//            "day": 0,
// 		//            "code": 30,
// 		//            "text": "局部多云",
// 		//            "low": 25,
// 		//            "high": 33,
// 		//            "image_large": "http://weather.china.xappengine.com/static/w/img/d30.png",
// 		//            "image_small": "http://weather.china.xappengine.com/static/w/img/s30.png"
// 		//        },
// 		//        {
// 		//            "date": "2013-07-01",
// 		//            "day": 1,
// 		//            "code": 37,
// 		//            "text": "局部雷雨",
// 		//            "low": 24,
// 		//            "high": 32,
// 		//            "image_large": "http://weather.china.xappengine.com/static/w/img/d37.png",
// 		//            "image_small": "http://weather.china.xappengine.com/static/w/img/s37.png"
// 		//        },
// 		//        {
// 		//            "date": "2013-07-02",
// 		//            "day": 2,
// 		//            "code": 38,
// 		//            "text": "零星雷雨",
// 		//            "low": 25,
// 		//            "high": 32,
// 		//            "image_large": "http://weather.china.xappengine.com/static/w/img/d38.png",
// 		//            "image_small": "http://weather.china.xappengine.com/static/w/img/s38.png"
// 		//        },
// 		//        {
// 		//            "date": "2013-07-03",
// 		//            "day": 3,
// 		//            "code": 38,
// 		//            "text": "零星雷雨",
// 		//            "low": 25,
// 		//            "high": 32,
// 		//            "image_large": "http://weather.china.xappengine.com/static/w/img/d38.png",
// 		//            "image_small": "http://weather.china.xappengine.com/static/w/img/s38.png"
// 		//        }
// 		//    ]
// 		// }`)

// 		// fmt.Fprintf(w, "%s", b)
// 	}

// }
