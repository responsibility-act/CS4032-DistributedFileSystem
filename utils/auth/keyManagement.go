package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/KillianDavitt/CS4032-DistributedFileSystem/utils/ticket"
	"github.com/kataras/iris"
	"gopkg.in/redis.v5"
	"log"
)

func StoreRedis(pubKey *rsa.PublicKey, identifier string) {
	PubASN1, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		log.Fatal(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: PubASN1,
	})

	pubKeyClient := getPubkeyRedis()
	err = pubKeyClient.Set(identifier, string(pubBytes), 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func RetrieveKey(identifier string) (pubKey *rsa.PublicKey) {
	pubKeyClient := getPubkeyRedis()
	pubKeyString, err := pubKeyClient.Get(identifier).Result()
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode([]byte(pubKeyString))
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := pub.(*rsa.PublicKey)
	return publicKey

}

func getPubkeyRedis() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 2})
}

func IsAllowed(ctx *iris.Context) bool {
	token := ctx.FormValue("token")
	pubKey := RetrieveKey("authserver")
	_ = ticket.GetTicketMap(token, pubKey)
	return true
}
