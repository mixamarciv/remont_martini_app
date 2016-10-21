package main

import (
	//"fmt"
	//"image"
	//"image/jpeg"
	//"io/ioutil"
	//"log"
	//"mime/multipart"
	"net/http"
	//"os"
	//"strconv"

	"github.com/codegangsta/martini-contrib/render"
	"github.com/martini-contrib/sessions"

	//"path/filepath"

	mf "github.com/mixamarciv/gofncstd3000"

	//"math/rand"

	//"github.com/nfnt/resize"
)

func http_get_messages(r render.Render, session sessions.Session) {
	var m = map[string]interface{}{}
	post := GetSessJson(session, "post", "{}")
	user := GetSessJson(session, "user", "{}")
	m["user"] = user
	m["post"] = post
	m["s"] = GetSessJson(session, "s", default_session_data)

	if imgs, ok := post["imagesuploaded"]; ok {
		post["imagesuploaded_jsonstr"] = mf.ToJsonStr(imgs)
	}

	r.HTML(200, "messages", m)
}

func http_post_messages(req *http.Request, session sessions.Session) string {
	/**********
	p := ParseBodyParams(req)
	ret := make([]map[string]interface{}, 0)

	uuid_post, ok := p["uuid_post"]
	if !ok {
		return "{\"error\":\"uuid_post не задан\"}"
	}

	query := "SELECT uuid,uuid_parent,iif(ishideuser=1,'{}',userdata) AS userdata,text,commentdatet, "
	query += "(SELECT COUNT(*) FROM timage t WHERE t.uuid_comment=p.uuid) "
	query += "FROM tcomment p "
	query += "WHERE uuid_post=? AND isactive=1 AND ishide=0 "
	query += "ORDER BY commentdatet "
	rows, err := db.Query(query, uuid_post)
	if err != nil {
		LogPrintErrAndExit("ERROR db.Query(query): \n"+query+"\n\n", err)
	}

	cnt_rows := 0
	for rows.Next() {
		cnt_rows++
		var uuid, uuid_parent, userdata, text NullString
		var commentdatet time.Time
		var imgcnt int
		if err := rows.Scan(&uuid, &uuid_parent, &userdata, &text, &commentdatet, &imgcnt); err != nil {
			LogPrintErrAndExit("ERROR rows.Scan: \n"+query+"\n\n", err)
		}
		m := map[string]interface{}{"uuid": uuid.get(""), "uuid_parent": uuid_parent.get(""), "text": post_text_to_html(text.get("-"))}
		m["userdata"] = mf.FromJsonStr([]byte(userdata.get("{}")))
		for k, _ := range m["userdata"].(map[string]interface{}) {
			switch k {
			case "phone", "email":
				delete(m["userdata"].(map[string]interface{}), k)
			}
		}
		m["datefmt"] = commentdatet.Format("02.01.2006 15:04")
		m["imgcnt"] = imgcnt
		if imgcnt > 0 {
			m["images"] = load_posts_or_comment_images_arr(m["uuid"].(string), "comment", cnt_rows)
		}
		ret = append(ret, m)
	}
	rows.Close()

	test := map[string]interface{}{"uuid_post": uuid_post, "cnt_rows": cnt_rows, "query": query}
	ret = append(ret, test)

	SetSessStr(session, "post", string("")) //затираем данные сессии, что бы пользователь дважды не создал один и тот же пост
	retstr := mf.ToJsonStr(js)
	return retstr
	*******/
	return "{}"
}

func http_post_messagenewsavesession(req *http.Request, session sessions.Session) string {
	js := ParseBodyParams(req)
	SetSessJson(session, "post", js)
	return "{\"success\":1}"
}

func http_post_messagenew(req *http.Request, session sessions.Session) string {
	js := ParseBodyParams(req)

	SetSessStr(session, "post", string("")) //затираем данные сессии, что бы пользователь дважды не создал один и тот же пост
	retstr := mf.ToJsonStr(js)
	return retstr
}
