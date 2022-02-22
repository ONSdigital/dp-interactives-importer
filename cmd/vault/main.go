package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	vault "github.com/ONSdigital/dp-vault"
	"github.com/ONSdigital/log.go/v2/log"
)

const maxRetries = 3

func main() {
	log.Namespace = "vault-example"

	devAddress := "http://localhost:8200"
	token := "0000-0000-0000-0000"
	path := "secret/shared/psk/testing"

	ctx := context.Background()

	client, err := vault.CreateClient(token, devAddress, maxRetries)

	// In production no tokens should be logged
	logData := log.Data{"address": devAddress, "token": token}
	log.Info(ctx, "Created vault client", logData)

	if err != nil {
		log.Error(ctx, "failed to connect to vault", err, logData)
	}

	raw := base64.StdEncoding.EncodeToString(createKey())
	//logData["raw"] = raw
	//logData["rawstr"] = string(raw[:])
	logData["raw"] = raw

	err = client.WriteKey(path, "key", raw)
	if err != nil {
		log.Error(ctx, "failed to write to vault", err, logData)
	}

	key, err := client.ReadKey(path, "key")
	if err != nil {
		log.Error(ctx, "failed to read key from vault", err, logData)
	}

	sEnc := base64.StdEncoding.EncodeToString([]byte(key))
	sDec, _ := base64.StdEncoding.DecodeString(sEnc)

	logData["key"] = key
	logData["sEnc"] = sEnc
	logData["sDec"] = sDec

	log.Info(ctx, "successfully written and read a key from vault", logData)
}

func createKey() []byte {
	key := make([]byte, 16)
	rand.Read(key) //nolint
	return key
}
