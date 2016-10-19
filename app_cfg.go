package main

//глобальные переменные на весь проект
var db_pass = "masterkey"
var run_on_addr = ":8090"
var sitedomain string = "127.0.0.1" + run_on_addr
var secret_email_pass string = "AsPeefW2m42i03yqVB9f123"

var cookie_store_name = "s"
var secret_cookie_store = "qwer1234"

//емеил на который будут уходить сообщения с сайта
var work_emails = []string{"mixamarciv@gmail.com"}

var default_session_data = `{"white":1,"dark":0}`
