CALL "%~dp0/set_path.bat"



::@CLS
::@pause

@echo === install ===================================================================

go get "github.com/go-martini/martini"
go get "github.com/martini-contrib/binding"
go get "github.com/martini-contrib/sessions"
go get "github.com/codegangsta/martini-contrib/render"
go get "github.com/mixamarciv/gofncstd3000"
go get "github.com/go-gomail/gomail"
go get "github.com/nakagami/firebirdsql"
go get "github.com/nfnt/resize"

go install

@echo ==== end ======================================================================
@PAUSE
