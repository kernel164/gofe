package main

import (
	fe "./fe"
	models "./models"
	settings "./settings"
	"fmt"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"log"
	"runtime"
	"strings"
)

var DEFAULT_API_ERROR_RESPONSE = models.GenericResp{
	models.GenericRespBody{false, "Not Supported"},
}

type SessionInfo struct {
	User         string
	Password     string
	FileExplorer fe.FileExplorer
	Uid          string
}

func main() {
	configRuntime()
	startServer()
}

func configRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func startServer() {
	settings.Load()
	macaron.Classic()
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Static("static", macaron.StaticOptions{
		Prefix:      "static",
		SkipLogging: true,
	}))
	m.Use(cache.Cacher())
	m.Use(session.Sessioner())
	m.Use(macaron.Renderer())
	m.Use(Contexter())

	m.Post("/api/_", binding.Bind(models.GenericReq{}), apiHandler)
	m.Post("/bridges/php/handler.php", binding.Bind(models.GenericReq{}), apiHandler)
	m.Get("/", mainHandler)
	m.Get("/login", loginHandler)
	m.Post("/api/download", defaultHandler)
	m.Post("/api/upload", defaultHandler)

	m.Run()
}

func mainHandler(ctx *macaron.Context) {
	ctx.HTML(200, "index")
}

func loginHandler(ctx *macaron.Context) {
	ctx.HTML(200, "login")
}

func defaultHandler(ctx *macaron.Context) {
	ctx.JSON(200, DEFAULT_API_ERROR_RESPONSE)
}

func apiHandler(c *macaron.Context, json models.GenericReq) {
	if json.Params.Mode == "list" {
		ls, err := c.Data["SessionInfo"].(SessionInfo).FileExplorer.ListDir(json.Params.Path)
		if err == nil {
			c.JSON(200, models.ListDirResp{ls})
		} else {
			ApiErrorResponse(c, 400, err)
		}
	} else if json.Params.Mode == "rename" { // path, newPath
		err := c.Data["SessionInfo"].(SessionInfo).FileExplorer.Move(json.Params.Path, json.Params.NewPath)
		if err == nil {
			ApiSuccessResponse(c, "")
		} else {
			ApiErrorResponse(c, 400, err)
		}
	} else if json.Params.Mode == "copy" { // path, newPath
		err := c.Data["SessionInfo"].(SessionInfo).FileExplorer.Copy(json.Params.Path, json.Params.NewPath)
		if err == nil {
			ApiSuccessResponse(c, "")
		} else {
			ApiErrorResponse(c, 400, err)
		}
	} else if json.Params.Mode == "delete" { // path
		err := c.Data["SessionInfo"].(SessionInfo).FileExplorer.Delete(json.Params.Path)
		if err == nil {
			ApiSuccessResponse(c, "")
		} else {
			ApiErrorResponse(c, 400, err)
		}
	} else if json.Params.Mode == "savefile" { // content, path
		c.JSON(200, DEFAULT_API_ERROR_RESPONSE)
	} else if json.Params.Mode == "editfile" { // path
		c.JSON(200, DEFAULT_API_ERROR_RESPONSE)
	} else if json.Params.Mode == "addfolder" { // name, path
		err := c.Data["SessionInfo"].(SessionInfo).FileExplorer.Mkdir(json.Params.Path, json.Params.Name)
		if err == nil {
			ApiSuccessResponse(c, "")
		} else {
			ApiErrorResponse(c, 400, err)
		}
	} else if json.Params.Mode == "changepermissions" { // path, perms, permsCode, recursive
		err := c.Data["SessionInfo"].(SessionInfo).FileExplorer.Chmod(json.Params.Path, json.Params.Perms)
		if err == nil {
			ApiSuccessResponse(c, "")
		} else {
			ApiErrorResponse(c, 400, err)
		}
	} else if json.Params.Mode == "compress" { // path, destination
		c.JSON(200, DEFAULT_API_ERROR_RESPONSE)
	} else if json.Params.Mode == "extract" { // path, destination, sourceFile
		c.JSON(200, DEFAULT_API_ERROR_RESPONSE)
	}
}

func IsApiPath(url string) bool {
	return strings.HasPrefix(url, "/api/") || strings.HasPrefix(url, "/bridges/php/handler.php")
}

func Contexter() macaron.Handler {
	return func(c *macaron.Context, cache cache.Cache, session session.Store, f *session.Flash) {
		IsSigned := false
		UserInfo := SessionInfo{}
		uid := session.Get("uid")

		if uid == nil {
			IsSigned = false
		} else {
			userInfo := cache.Get(uid.(string))
			if userInfo == nil {
				IsSigned = false
			} else {
				UserInfo = userInfo.(SessionInfo)
				if UserInfo.User == "" || UserInfo.Password == "" {
					IsSigned = false
				} else {
					IsSigned = true
					c.Data["UserInfo"] = UserInfo
					if UserInfo.FileExplorer == nil {
						fe, err := BackendConnect(UserInfo.User, UserInfo.Password)
						UserInfo.FileExplorer = fe
						if err != nil {
							IsSigned = false
							if IsApiPath(c.Req.URL.Path) {
								ApiErrorResponse(c, 500, err)
							} else {
								AuthError(c, f, err)
							}
						}
					}
				}
			}
		}

		if IsSigned == false {
			if strings.HasPrefix(c.Req.URL.Path, "/login") {
				if c.Req.Method == "POST" {
					username := c.Query("username")
					password := c.Query("password")
					fe, err := BackendConnect(username, password)
					if err != nil {
						AuthError(c, f, err)
					} else {
						uid := username // TODO: ??
						UserInfo = SessionInfo{username, password, fe, uid}
						cache.Put(uid, UserInfo, 100000000000)
						session.Set("uid", uid)
						c.Data["UserInfo"] = UserInfo
						c.Redirect("/")
					}
				}
			} else {
				c.Redirect("/login")
			}
		} else {
			if strings.HasPrefix(c.Req.URL.Path, "/logout") {
				UserInfo.FileExplorer.Close()
				session.Delete("uid")
				cache.Delete(uid.(string))
				c.SetCookie("MacaronSession", "")
				c.Redirect("/login")
			}
		}
	}
}

func BackendConnect(username string, password string) (fe.FileExplorer, error) {
	fe := fe.NewSSHFileExplorer(settings.SshHost, username, password)
	err := fe.Init()
	if err == nil {
		return fe, nil
	}
	log.Println(err)
	return nil, err
}

func ApiErrorResponse(c *macaron.Context, code int, obj interface{}) {
	var message string
	if err, ok := obj.(error); ok {
		message = err.Error()
	} else {
		message = obj.(string)
	}
	c.JSON(code, models.GenericResp{models.GenericRespBody{false, message}})
}

func ApiSuccessResponse(c *macaron.Context, message string) {
	c.JSON(200, models.GenericResp{models.GenericRespBody{true, message}})
}

func AuthError(c *macaron.Context, f *session.Flash, err error) {
	f.Set("ErrorMsg", err.Error())
	c.Data["Flash"] = f
	c.Data["ErrorMsg"] = err.Error()
	c.Redirect("/login")
}
