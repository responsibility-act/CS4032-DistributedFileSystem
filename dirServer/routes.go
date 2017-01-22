package main
import (
	"github.com/kataras/iris"
	"log"
	"encoding/json"
	"gopkg.in/redis.v5"
	"github.com/KillianDavitt/CS4032-DistributedFileSystem/utils/ticket"
	"github.com/KillianDavitt/CS4032-DistributedFileSystem/utils/auth"
)
func getFile(ctx *iris.Context){
	//filename := ctx.URLParam("filename")
	//f := readFiles()
	//ip := f.getFile(filename).String()
	ctx.HTML(iris.StatusOK, "hi")
}

func putFile(ctx *iris.Context){
	ctx.HTML(iris.StatusOK, "ok")
}


func listFiles(ctx *iris.Context){
	if !isAllowed(ctx){
		ctx.HTML(iris.StatusForbidden, "Invalid token")
	}
	//	f := readFiles()
	files := [2]string{"hi.pem", "test.txt"}
	jsonFiles, err := json.Marshal(files)
	if err != nil {
		log.Fatal(err)
	}
	ctx.HTML(iris.StatusOK, string(jsonFiles))//f.getFile("test.txt").String())
}

func registerToken(ctx *iris.Context){
	pubKey := auth.RetrieveKey("authserver")
	token := ctx.FormValue("token")
	ticket := ticket.GetTicketMap(token, pubKey)
	token_client := getTokenRedis()
	expiryString, err := ticket.Expiry_date.MarshalText()
	if err != nil {
		log.Fatal(err)
	}
	err = token_client.Set(string(ticket.Token), string(expiryString), 0).Err()
	if err != nil {
		log.Fatal(err)
	}
	ctx.HTML(iris.StatusOK, "Register token succ")
}

func isAllowed(ctx *iris.Context)(bool){
	token := ctx.FormValue("token")
	pubKey := auth.RetrieveKey("authserver")
	ticket := ticket.GetTicketMap(token, pubKey)
	client := getTokenRedis()
	_, err := client.Get(string(ticket.Token)).Result()
	if err != nil{
		return false
	}
	return true
}

func getTokenRedis() (*redis.Client){
	return redis.NewClient(&redis.Options{ Addr: "localhost:6379", Password: "", DB: 1})
}
