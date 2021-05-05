package server

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"os"
	"strconv"
	"time"
	"strings"
	"io"
	"io/ioutil"
	"encoding/json"
	"sync"

)

type Entry struct {
	AccessCount int `json:"count"`
	Id string `json:"id"`
	Name string `json:"name"`
	MaxAccess int `json:"max"`
	ContentType string `json:"content-type"`
}

func BuildRoutes(app *iris.Application) {
	accessToken := os.Getenv("ACCESS_TOKEN")
	mtx := &sync.Mutex{}
	os.Mkdir("data", os.ModePerm)
	app.Post("/add/{count string}", func(ctx context.Context) {
		token := ctx.GetHeader("Authorization")
		if token != accessToken {
			return
		}
		var count = ctx.Params().Get("count")
		if len(count) == 0  {
			ctx.StatusCode(iris.StatusBadRequest)
			_, _ = ctx.JSON(map[string]string{"message": "malformed"})
			return
		}
		n, err := strconv.Atoi(count)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			_, _ = ctx.JSON(map[string]string{"message": "invalid count"})
			return
		}

		file, info, err := ctx.FormFile("datafile")
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.JSON(map[string]string{"message": "internal_error_on_upload"})
			return
		}
		defer file.Close()
		id := strconv.FormatInt(time.Now().Unix(), 10)
		parts := strings.Split(info.Filename, ".")
		ext := parts[len(parts)-1]
		var e = Entry{
			AccessCount: 0,
			Id: id,
			Name: ext,
			MaxAccess: n,
			ContentType: info.Header["Content-Type"][0],
		}
		out, err:= os.OpenFile("./data/" + id + "." + ext, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.JSON(map[string]string{"message": "internal_error_on_upload"})
			return
		}
		defer out.Close()
		io.Copy(out, file)
		var entries []Entry
		if _, err := os.Stat("./data/files.json"); err == nil {
			content, err:= ioutil.ReadFile("./data/files.json")
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				_, _ = ctx.JSON(map[string]string{"message": "internal_error_on_upload"})
				return
			}
			err = json.Unmarshal(content, &entries)
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				_, _ = ctx.JSON(map[string]string{"message": "internal_error_on_upload"})
				return
			}
		}
		entries = append(entries, e)

		jsonout, err := json.Marshal(entries)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.JSON(map[string]string{"message": "internal_error_on_upload"})
			return
		}
    err = ioutil.WriteFile("./data/files.json", jsonout, 0644)
		ctx.Header("Content-type", "text/plain")
		_, _ = ctx.Write([]byte(id + "." + ext))

	})

	app.Get("/d/{id string}", func(ctx context.Context) {
		mtx.Lock()
		token := ctx.GetHeader("Authorization")

		var id = ctx.Params().Get("id")
		if len(id) == 0  {
			ctx.StatusCode(iris.StatusBadRequest)
			_, _ = ctx.JSON(map[string]string{"message": "malformed"})
			mtx.Unlock()
			return
		}
		var admin = token == accessToken
		var entries []Entry
		if _, err := os.Stat("./data/files.json"); err == nil {
			content, err:= ioutil.ReadFile("./data/files.json")
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				_, _ = ctx.JSON(map[string]string{"message": "internal_error_on_upload"})
				mtx.Unlock()
				return
			}
			err = json.Unmarshal(content, &entries)

		}

		for index, entry := range entries {
			if (entry.Id + "." + entry.Name) == id {
				if admin || entry.MaxAccess == 0 || entry.AccessCount < entry.MaxAccess {
					content, err := ioutil.ReadFile("./data/" + entry.Id + "." + entry.Name)
					if err != nil {
						mtx.Unlock()
						return
					}
					ctx.Header("Content-type", entry.ContentType)
					ctx.Write(content)

				if !admin && entry.MaxAccess != 0 {
					entries[index].AccessCount+=1

					jsonout, err := json.Marshal(entries)
					if err != nil {
						ctx.StatusCode(iris.StatusInternalServerError)
						_, _ = ctx.JSON(map[string]string{"message": "internal_error_on_upload"})
						return
					}
					err = ioutil.WriteFile("./data/files.json", jsonout, 0644)
				}
					mtx.Unlock()
				return
			}
			}
		}
		ctx.StatusCode(iris.StatusNotFound)
		_, _ = ctx.JSON(map[string]string{"message": "not_found"})
		mtx.Unlock()

	})
}
