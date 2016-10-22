package main

import (
	"fmt"
	//"image"
	//"image/jpeg"
	//"io/ioutil"
	//"log"
	//"mime/multipart"
	"net/http"
	//"os"
	"strconv"
	"time"

	"github.com/codegangsta/martini-contrib/render"
	"github.com/martini-contrib/sessions"

	//"path/filepath"

	mf "github.com/mixamarciv/gofncstd3000"

	//"math/rand"

	//"github.com/nfnt/resize"
	"html"
)

var itoa = strconv.Itoa
var atoi = strconv.Atoi

func http_get_messages(r render.Render, session sessions.Session) {
	var m = map[string]interface{}{}
	post := GetSessJson(session, "post", "{}")
	user := GetSessJson(session, "user", "{}")
	m["user"] = user
	m["post"] = post
	if _, ok := post["uuid"]; !ok {
		post["uuid"] = mf.StrUuid()
		post["time"] = mf.CurTimeStrShort()
		SetSessJson(session, "post", post)
	}

	m["s"] = GetSessJson(session, "s", default_session_data)

	{ //загружаем общую информацию о сообщениях:
		query := `SELECT COUNT(*) FROM tmessage p WHERE p.ishide=0`
		rows, err := db.Query(query)
		LogPrintErrAndExit("db.Query error: \n"+query+"\n\n", err)
		rows.Next()
		var cnt int
		err = rows.Scan(&cnt)
		LogPrintErrAndExit("rows.Scan error: \n"+query+"\n\n", err)
		m["messages_count"] = cnt
	}

	if imgs, ok := post["imagesuploaded"]; ok {
		post["imagesuploaded_jsonstr"] = mf.ToJsonStr(imgs)
	}

	r.HTML(200, "messages", m)
}

func http_post_messages(req *http.Request, session sessions.Session) string {
	p := ParseBodyParams(req)

	var skip float64 = 0

	fmt.Printf("%#v", p)
	if t, ok := p["skip"]; ok {
		skip = t.(float64)
	}

	cnt_skip := itoa(int(skip))
	cnt_msgs := itoa(cfg_cnt_messages_on_page)

	query := "SELECT FIRST " + cnt_msgs + " SKIP " + cnt_skip
	query += "  uuid,uuid_parent,userdata,name,email,text,upddate,tdate,tdatet,"
	query += "  (SELECT COUNT(*) FROM timage t WHERE t.uuid_message=p.uuid)"
	query += "FROM tmessage p "
	query += "WHERE ishide=0 "

	// uuid_parent:
	if s, ok := p["uuid_parent"]; ok {
		if s.(string) != "-" {
			query += " AND p.uuid='" + s.(string) + "' "
		}
	}
	query += "ORDER BY tdatet"

	rows, err := db.Query(query)
	if err != nil {
		LogPrintErrAndExit("ERROR db.Query(query): \n"+query+"\n\n", err)
	}

	ret := make([]map[string]interface{}, 0)

	cfg := map[string]interface{}{"cnt_msgs": cfg_cnt_messages_on_page, "skip": skip}
	ret = append(ret, cfg) //первая строка всегда будет информационная

	cnt_rows := 0
	for rows.Next() {
		cnt_rows++
		var uuid, uuid_parent, userdata, name, email, text, upddate, tdate NullString
		var tdatet time.Time
		var imgcnt int
		if err := rows.Scan(&uuid, &uuid_parent, &userdata, &name, &email, &text, &upddate, &tdate, &tdatet, &imgcnt); err != nil {
			LogPrintErrAndExit("ERROR rows.Scan: \n"+query+"\n\n", err)
		}
		m := map[string]interface{}{"uuid": uuid.get(""), "uuid_parent": uuid_parent.get("-"), "text": post_text_to_html(text.get("-"))}
		m["userdata"] = mf.FromJsonStr([]byte(userdata.get("{}")))
		m["name"] = name.get("аноним")
		m["datefmt"] = tdatet.Format("02.01.2006 15:04")
		m["imgcnt"] = imgcnt
		if imgcnt > 0 {
			m["images"] = load_images_arr(m["uuid"].(string), cnt_rows)
		}
		ret = append(ret, m)
	}
	rows.Close()

	retstr := mf.ToJsonStr(ret)

	return retstr
}

func http_post_messagenewsavesession(req *http.Request, session sessions.Session) string {
	js := ParseBodyParams(req)
	SetSessJson(session, "post", js)

	user := map[string]interface{}{"name": js["name"], "email": js["email"]}
	SetSessJson(session, "user", user)
	return "{\"success\":1}"
}

func http_post_messagenew(req *http.Request, session sessions.Session) string {
	js := ParseBodyParams(req)
	save_message_data(js)

	var ret = map[string]interface{}{}
	{ //затираем данные сессии, что бы пользователь дважды не создал один и тот же пост
		var post = map[string]interface{}{}
		post["uuid"] = mf.StrUuid()
		post["time"] = mf.CurTimeStrShort()
		SetSessJson(session, "post", post)
		ret["new"] = post
	}
	ret["uuid"] = js["uuid"]
	ret["success"] = "ваше сообщение успешно отправлено"

	return mf.ToJsonStr(ret)
}

func save_message_data(p map[string]interface{}) {
	fmt.Printf("%#v", p)
	posttime := mf.CurTimeStrShort()
	query := "INSERT INTO tmessage(uuid_parent,uuid,userdata,name,email,text,upddate,tdate,tdatet) "
	query += "VALUES(?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP)"
	_, err := db.Exec(query, p["uuid_parent"], p["uuid"],
		mf.ToJsonStr(p["userdata"]), p["name"], p["email"], p["text"], posttime, posttime)
	LogPrintErrAndExit("ERROR db.Exec: \n"+query+"\n\n", err)

	imgs := p["imagesuploaded"].([]interface{})
	for _, imgi := range imgs {
		img, ok := imgi.(map[string]interface{})
		if !ok {
			continue
		}
		//log.Printf("%#v\n", img)
		query := "INSERT INTO timage(uuid_message,uuid,hash,title,path,pathmin,imgdate,imgdatet) "
		query += "VALUES(?,?,?,?,?,?,?,CURRENT_TIMESTAMP)"
		_, err := db.Exec(query, p["uuid"], mf.StrUuid(), "nohash", img["text"], img["path"], img["pathmin"], posttime)
		LogPrintErrAndExit("ERROR db.Exec: \n"+query+"\n\n", err)
	}
	//SendMailNewPostsToWork()
}

func load_images_arr(uuid_post string, i_post int) []map[string]string {
	ret := make([]map[string]string, 0)
	query := "SELECT uuid,title,path,pathmin "
	query += " FROM timage WHERE uuid_message=? ORDER BY uuid"
	rows, err := db.Query(query, uuid_post)
	if err != nil {
		LogPrintErrAndExit("ERROR db.Query(query): \n"+query+"\n\n", err)
	}
	for rows.Next() {
		var uuid, title, path, pathmin NullString
		if err := rows.Scan(&uuid, &title, &path, &pathmin); err != nil {
			LogPrintErrAndExit("ERROR rows.Scan: \n"+query+"\n\n", err)
		}
		m := map[string]string{"uuid": uuid.get("-"), "title": title.get(""), "path": path.get(""), "pathmin": pathmin.get(""), "ipost": strconv.Itoa(i_post)}
		ret = append(ret, m)
	}
	rows.Close()
	return ret
}

func post_text_to_html(text string) string {
	text = html.EscapeString(text)
	//text = mf.StrRegexpReplace(text, "\\n", "<br>")
	return text
}
