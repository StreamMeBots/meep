package routes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/StreamMeBots/meep/pkg/bot"
	"github.com/StreamMeBots/meep/pkg/config"
	"github.com/StreamMeBots/meep/pkg/greetings"
	"github.com/StreamMeBots/meep/pkg/user"
	pkgBot "github.com/StreamMeBots/pkg/bot"
	"github.com/gin-gonic/gin"
)

// our running bouts
var bots = bot.NewBots()

var Debugf = func(string, ...interface{}) {}
var Debugln = func(...interface{}) {}

// Init adds the http routes to the gin engine
func Init(r *gin.Engine) {
	// setup debug logging
	if config.Conf.Debug {
		Debugf = func(f string, args ...interface{}) {
			log.Printf(f, args...)
		}
		Debugln = func(args ...interface{}) {
			log.Println(args...)
		}
	}

	// setup oauth config
	auth = newAuth(config.Conf)

	// oauth2 login routes
	r.GET("/login", auth.loginHandler)
	r.GET("/login-redirect", auth.redirectHandler)
	r.GET("/logout", logout)

	// API routes
	api := r.Group("/api", checkAuth)
	{
		// current user info
		api.GET("/me", loggedInUser)

		// Bot
		// Start bot
		api.POST("/bot", startBot)

		// Stop Bot
		api.DELETE("/bot", stopBot)

		// bot info
		api.GET("/bot", botInfo)

		// Grettings
		// get greeting messages
		api.GET("/greeting-templates", getGreetings)

		// save greeting messages
		api.POST("/greeting-templates", saveGreetings)

		// bot log
		api.GET("/bot/log-stream", logStream)

		// Commands
		// get commands
		api.GET("/commands", getCommands)

		// update commands list
		api.PUT("/commands", updateCommand)

		// get a single
		api.GET("/commands/:name", getCommand)

		// remove a command from the commands list
		api.DELETE("/commands/:name", deleteCommand)
	}
}

func loggedInUser(ctx *gin.Context) {
	ctx.JSON(200, getAuthedUser(ctx))
}

func logout(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:   "sessid",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	ctx.Redirect(302, "/")
}

func botInfo(ctx *gin.Context) {
	u := getAuthedUser(ctx)
	ctx.JSON(200, bots.Info(u.user.PublicId))
}

func logStream(ctx *gin.Context) {
	u := getAuthedUser(ctx)

	ch, err := bots.LogStream(u.user.PublicId)
	if err != nil {
		ctx.Stream(func(w io.Writer) bool {
			log.Println("botError", err.Error())
			ctx.SSEvent("botError", err.Error())
			return false
		})
		return
	}
	defer bots.CloseLogStream(u.user.PublicId)

	ctx.Stream(func(w io.Writer) bool {
		e, ok := <-ch
		if !ok {
			return false
		}
		switch t := e.(type) {
		case pkgBot.EventStateChange:
			ctx.SSEvent("stateChange", t)
		case pkgBot.EventReadCommand:
			ctx.SSEvent("read", t)
		case pkgBot.EventReadError:
			ctx.SSEvent("readError", t.Error())
		case pkgBot.EventWrite:
			ctx.SSEvent("write", t)
		case pkgBot.EventWriteError:
			ctx.SSEvent("writeError", t.Error())
		}
		return true
	})
}

func startBot(ctx *gin.Context) {
	u := getAuthedUser(ctx)

	if err := bots.Start(u.user.PublicId, u.client); err != nil {
		ctx.JSON(500, map[string]string{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]string{
		"message": "Bot has been started!",
	})
}

func stopBot(ctx *gin.Context) {
	u := getAuthedUser(ctx)

	bots.Stop(u.user.PublicId)

	ctx.JSON(200, map[string]string{
		"message": "Bot has been stopped",
	})
}

func getGreetings(ctx *gin.Context) {
	u := getAuthedUser(ctx)

	tmpl, err := greetings.Get(user.BucketName(u.user.PublicId))
	if err != nil {
		log.Printf("msg='json-decode-error', error='%v'\n", err)
		ctx.JSON(500, map[string]string{
			"message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(200, tmpl)
}

func saveGreetings(ctx *gin.Context) {
	u := getAuthedUser(ctx)

	tmpl := &greetings.Template{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&tmpl); err != nil {
		log.Printf("msg='json-decode-error', error='%v'\n", err)
		ctx.JSON(400, map[string]string{
			"message": "Invalid JSON body",
		})
		return
	}

	if err := tmpl.Validate(); err != nil {
		ctx.JSON(422, map[string]string{
			"message": err.Error(),
		})
		return
	}

	if err := tmpl.Save(user.BucketName(u.user.PublicId)); err != nil {
		log.Printf("msg='error-saving-greeting', userPublicId='%s', error='%v'\n", u.user.PublicId, err)
		ctx.JSON(500, map[string]string{
			"message": "Internal server error",
		})
	}

	ctx.JSON(200, tmpl)
}

func getCommands(ctx *gin.Context) {
	//u := getAuthedUser(ctx)

}

func getCommand(ctx *gin.Context) {
	//u := getAuthedUser(ctx)

}

func updateCommand(ctx *gin.Context) {
	//u := getAuthedUser(ctx)

}

func deleteCommand(ctx *gin.Context) {
	//u := getAuthedUser(ctx)

}
