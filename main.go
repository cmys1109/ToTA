package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"strings"
	"time"
)

type MysqlConf struct {
	IP       string `json:"ip"`
	Post     int    `json:"post"`
	DBName   string `json:"db_name"`
	UserName string `json:"user_name"`
	Passwd   string `json:"passwd"`
}

type ServerConf struct {
	Post string `json:"post"`
}

type Conf struct {
	ServerConfig ServerConf `json:"server_conf"`
	MysqlConfig  MysqlConf  `json:"mysql_conf"`
}

type Table struct {
	Text string
	Time string
}

func main() {
	ConfFile, err := ioutil.ReadFile("./Config.json")
	if err != nil {
		return
	}
	var Conf Conf
	err = json.Unmarshal(ConfFile, &Conf)
	if err != nil {
		return
	}

	HTMLMap := LoadHTML()

	DBC := Conf.MysqlConfig.UserName + ":" + Conf.MysqlConfig.Passwd + "@tcp(" + Conf.MysqlConfig.IP + ")/" + Conf.MysqlConfig.DBName
	DB, err := sql.Open("mysql", DBC)
	if err != nil {
		panic(err)
		return
	}

	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		panic(err)
		return
	}

	App := iris.New()

	App.RegisterView(iris.HTML("./views", ".html"))

	App.Get("/", func(ctx iris.Context) {

		_, err := ctx.HTML(HTMLMap["index.html"])
		if err != nil {
			panic(err)
			return
		}

		App.Logger().Info(ctx.Path(), ctx.Request().RemoteAddr)
	})

	App.Post("/PushText", func(ctx iris.Context) {
		var (
			Name = ctx.FormValue("name")
			Text = ctx.FormValue("text")
			Time = []uint8(time.Now().Format("2006-01-02 15:04:05"))
			IP   = ctx.Request().RemoteAddr
		)

		// 检测是否为空然后添加至数据库
		if Name == "" || Text == "" {
			_, err := ctx.JSON(map[string]string{"ERR": "name OR text is NULL"})
			if err != nil {
				return
			}
			return
		}
		//将换行符替换成<br>实现换行
		Text = strings.ReplaceAll(Text, "\n", "<br>")
		SqlStr := "INSERT INTO ToTA(name,text,time,ip) VALUE (?,?,?,?)"
		_, err := DB.Exec(SqlStr, Name, Text, Time, IP)
		if err != nil {
			panic(err)
			return
		}

		_, err = ctx.HTML(HTMLMap["Push.html"])
		if err != nil {
			return
		}
		App.Logger().Info(ctx.Path(), IP, Name)
	})

	App.Post("/GetText", func(ctx iris.Context) {
		//key := ctx.FormValue("name")
		key := ctx.FormValue("name")
		// 执行查询
		rows, err := DB.Query("SELECT text,time from ToTA WHERE name = '" + key + "'")
		if err != nil {
			panic(err)
			return
		}
		// 关闭结果集
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				panic(err)
			}
		}(rows)

		var TextList []Table

		for rows.Next() {
			var DBTable Table
			var Time []uint8
			err = rows.Scan(&DBTable.Text, &Time)
			if err != nil {
				fmt.Println(err.Error())
				panic(err)
				return
			}
			DBTable.Time = string(Time)
			TextList = append(TextList, DBTable)
		}

		err = rows.Err()
		if err != nil {
			panic(err)
			return
		}

		// wuwuwu,难搞，模板不能直接插入文本，html的解析会出问题
		var body = "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>写给 " + key + " 的留言</title>\n</head>\n<body>\n<h1>写给  " + key + "  的留言:</h1>\n<hr>"
		for _, k := range TextList {
			body += "<hr>\n<table>\n    <tr>\n        <td>" + k.Time + "</td>\n    </tr>\n    <tr>\n        <td>" + k.Text + "</td>\n    </tr>\n</table>"
		}
		body += "</body>\n</html>"
		_, err = ctx.HTML(body)
		if err != nil {
			panic(err)
			return
		}

		App.Logger().Info(ctx.Path(), ctx.Request().RemoteAddr, key)
	})

	App.Get("/Say", func(ctx iris.Context) {
		_, err := ctx.HTML(HTMLMap["Say.html"])
		if err != nil {
			panic(err)
			return
		}

		App.Logger().Info(ctx.Path(), ctx.Request().RemoteAddr)
	})

	App.Get("/Look", func(ctx iris.Context) {
		_, err := ctx.HTML(HTMLMap["Look.html"])
		if err != nil {
			panic(err)
			return
		}

		App.Logger().Info(ctx.Path(), ctx.Request().RemoteAddr)
	})

	// 启动服务
	err = App.Run(iris.Addr(":" + Conf.ServerConfig.Post))
	if err != nil {
		panic(err)
		return
	}

}

func LoadHTML() map[string]string {
	var HTMLFileList = "index.html;;Say.html;;Look.html;;Push.html"
	HTMLMap := make(map[string]string)
	for _, v := range strings.Split(HTMLFileList, ";;") {
		html, err := ioutil.ReadFile("./views/" + v)
		if err != nil {
			return map[string]string{}
		}
		HTMLMap[v] = string(html)
	}
	return HTMLMap
}
