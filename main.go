// @file main.go
// @brief Catswords Community CLI (catswords-cli)
// @author Go Namhyeon <gnh1201@gmail.com>
// @created_at 2019-01-09
// @updated_at 2019-02-05

package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
    "io/ioutil"
    "github.com/urfave/cli"
    "gopkg.in/resty.v1"
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
)

type User struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type AuthToken struct {
    Token string `json:"token"`
}

type AuthSuccess struct {
    Data AuthToken `json:"data"`
    Public bool `json:"public"`
}

type MessageContext struct {
    Status string `json:"status"`
    Message string `json:"message"`
    Agent string `json:"agent"`
    Format string `json:"format"`
    Delimiter string `json:"delimiter"`
    Encoding string `json:"encoding"`
    Mime string `json:"mime"`
    Label string `json:"label"`
    Encryption string `json:"encryption"`
    EncryptionKey string `json:"encryption_key"`
    EncryptionIv string `json:"encryption_iv"`
    PrivateKey string `json:"private_key"`
    PublicKey string `json:"public_key"`
    HashFunction string `json:"hash_function"`
    HashValue string `json:"hash_value"`
    Mnemonic string `json:"mnemonic"`
    IntNetwork string `json:"int_network"`
    IntAddress string `json:"int_address"`
    ExtNetwork string `json:"ext_network"`
    ExtAddress string `json:"ext_address"`
    Host string `json:"host"`
    Lang string `json:"lang"`
    NetworkId string `json:"network_id"`
    AccessKey string `json:"access_key"`
    AccessSecret string `json:"access_secret"`
    Protocol string `json:"protocol"`
    TimeToLive int64 `json:"time_to_live"`
    Latitude float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Geohash string `json:"geohash"`
    Planet string `json:"planet"`
}

type HashResult struct {
    HashFunction string `json:"hash_function"`
    HashValue string `json:"hash_value"`
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func authenticate(email string, password string, host string, protocol string) string {
    auth := new(AuthSuccess)
    user := &User{Email: email, Password: password}
    body, err := json.Marshal(user)
    check(err)

    resty.R().
          SetHeader("Content-Type", "application/json").
          SetBody(body).
          SetResult(&auth).
          Post(fmt.Sprintf("%s://%s/_/auth/authenticate", protocol, host))

    tokenBytes := []byte(auth.Data.Token)
    err2 := ioutil.WriteFile("token.dat", tokenBytes, 0644)
    check(err2)

    result := string(tokenBytes)

    return result
}

func sendMessage(context MessageContext, token string, host string, protocol string) {
    body, err := json.Marshal(context)
    check(err)

    // post data to server
    resp, err := resty.R().
          SetHeader("Content-Type", "application/json").
          SetBody(body).
          SetAuthToken(token).
          Post(fmt.Sprintf("%s://%s/_/items/catswords_cli", protocol, host))
    check(err)

    fmt.Println(resp)
}

func recvMessages(networkId string, token string, host string, protocol string, limit int64) {
    resp, err := resty.R().
          SetHeader("Accept", "application/json").
          SetAuthToken(token).
          Get(fmt.Sprintf("%s://%s/_/items/catswords_cli?filter[network_id][eq]=%s&limit=%d&sort=-id", protocol, host, networkId, limit))
    check(err)

    fmt.Println(resp)
}

func getHashResult(msgBytes []byte, hashFunction string) *HashResult {
    result := new(HashResult)
    checksum := ""

    if hashFunction == "md5" {
        hasher := md5.New()
        hasher.Write(msgBytes)
        checksum = fmt.Sprintf("%x", hasher.Sum(nil))
    }

    if hashFunction == "sha1" {
        hasher := sha1.New()
        hasher.Write(msgBytes)
        checksum = fmt.Sprintf("%x", hasher.Sum(nil))
    }

    if hashFunction == "sha256" {
        hasher := sha256.New()
        hasher.Write(msgBytes)
        checksum = fmt.Sprintf("%x", hasher.Sum(nil))
    }

    if hashFunction == "sha512" {
        hasher := sha512.New()
        hasher.Write(msgBytes)
        checksum = fmt.Sprintf("%x", hasher.Sum(nil))
    }

    result.HashFunction = hashFunction
    result.HashValue = checksum

    return result
}

func main() {
    app := cli.NewApp()
    app.Name = "Catswords Community CLI (catwords-cli)"
    app.Usage = "The awesome command line tool to keep work-life balance for developers"
    app.Version = "1.4-dev"
    app.Compiled = time.Now()
    app.Authors = []cli.Author{
        cli.Author{
            Name:  "Go Namhyeon",
            Email: "gnh1201@gmail.com",
        },
        cli.Author{
            Name:  "Catswords Community",
            Email: "support@exts.kr",
        },
        cli.Author{
            Name:  "Catswords Research",
            Email: "support@2s.re.kr",
        },
    }
    app.Copyright = "(c) 2019 Catswords Research."

    app.Flags = []cli.Flag {
        cli.StringFlag{
            Name: "host",
            Value: "catswords.re.kr",
            Usage: "set server hostname",
        },
        cli.StringFlag{
            Name: "lang",
            Value: "english",
            Usage: "set language",
        },
        cli.StringFlag{
            Name: "email",
            Value: "",
            Usage: "set user email",
        },
        cli.StringFlag{
            Name: "password",
            Value: "",
            Usage: "set user password",
        },
        cli.StringFlag{
            Name: "token",
            Value: "",
            Usage: "set access token",
            FilePath: "token.dat",
        },
        cli.StringFlag{
            Name: "protocol",
            Value: "https",
            Usage: "set protocol: https, or http, or more",
        },
        cli.StringFlag{
            Name: "action",
            Value: "send",
            Usage: "set action: send, or recv, or refresh",
        },
        cli.StringFlag{
            Name: "message",
            Value: "",
            Usage: "set message it will send to server",
        },
        cli.StringFlag{
            Name: "file",
            Value: "",
            Usage: "set file path instead of message",
        },
        cli.StringFlag{
            Name: "format",
            Value: "text",
            Usage: "set message type: text, json, xml, rfc5424(syslog), or more",
        },
        cli.StringFlag{
            Name: "delimiter",
            Value: "comma",
            Usage: "set delimiter: comma, or pipeline, or more",
        },
        cli.StringFlag{
            Name: "encoding",
            Value: "utf-8",
            Usage: "set encoding: character set, or encapsulation, or more",
        },
        cli.StringFlag{
            Name: "mine",
            Value: "text/plain",
            Usage: "set media type: text/plain, or application/json, or more",
        },
        cli.StringFlag{
            Name: "label,labels",
            Value: "",
            Usage: "set label(s) with delimiter",
        },
        cli.StringFlag{
            Name: "agent",
            Value: "",
            Usage: "set custom agent name",
        },
        cli.StringFlag{
            Name: "encryption,enc",
            Value: "",
            Usage: "set encryption algorithm: des, or aes, or more",
        },
        cli.StringFlag{
            Name: "encryption-key,ekey",
            Value: "",
            Usage: "set encryption key",
        },
        cli.StringFlag{
            Name: "encryption-iv,eiv",
            Value: "",
            Usage: "set encryption IV",
        },
        cli.StringFlag{
            Name: "private-key,privkey",
            Value: "",
            Usage: "set private key",
        },
        cli.StringFlag{
            Name: "public-key,pubkey",
            Value: "",
            Usage: "set public key",
        },
        cli.StringFlag{
            Name: "hash-function,hasher",
            Value: "",
            Usage: "set hash function(s) with delimiter: md5, or sha1, or sha256, or more",
        },
        cli.StringFlag{
            Name: "hash-value,hash",
            Value: "",
            Usage: "set hash value(s) with delimiter",
        },
        cli.StringFlag{
            Name: "mnemonic",
            Value: "",
            Usage: "set mnemonic",
        },
        cli.StringFlag{
            Name: "int-network,innet",
            Value: "",
            Usage: "set internal network name",
        },
        cli.StringFlag{
            Name: "int-address,inaddr",
            Value: "",
            Usage: "set address of specified internal network",
        },
        cli.StringFlag{
            Name: "ext-network,exnet",
            Value: "",
            Usage: "set external network name",
        },
        cli.StringFlag{
            Name: "ext-address,exaddr",
            Value: "",
            Usage: "set address of specified external network",
        },
        cli.StringFlag{
            Name: "network-id,netid",
            Value: "",
            Usage: "set network ID",
        },
        cli.StringFlag{
            Name: "access-key,akey",
            Value: "",
            Usage: "set access key",
        },
        cli.StringFlag{
            Name: "access-secret,asec",
            Value: "",
            Usage: "set access secret",
        },
        cli.StringFlag{
            Name: "time-to-live,ttl",
            Value: "",
            Usage: "set limits the number of reading devices",
        },
        cli.StringFlag{
            Name: "limit",
            Value: "128",
            Usage: "set limit the number of receving messages",
        },
        cli.StringFlag{
            Name: "latitude,lat",
            Value: "",
            Usage: "set latitude",
        },
        cli.StringFlag{
            Name: "longitude,lon",
            Value: "",
            Usage: "set longitude",
        },
        cli.StringFlag{
            Name: "geohash",
            Value: "",
            Usage: "set geohash",
        },
        cli.StringFlag{
            Name: "planet",
            Value: "",
            Usage: "set planet",
        },
    }

    app.Action = func(c *cli.Context) error {
        token := c.String("token")
        if (token == "" || c.String("action") == "refresh") {
            token = authenticate(c.String("email"), c.String("password"), c.String("host"), c.String("protocol"))
            fmt.Println(token)
        }

        if token == "" {
            fmt.Println("Could not find access token")
            return nil
        }

        if c.String("email") != "" {
            fmt.Println("Done authenticate. You have to add flag '--action refresh' when the token expired")
            return nil
        }

        if c.String("action") == "recv" {
            if c.String("network-id") == "" {
                fmt.Println("You must be set network ID '--network-id [your network ID]'")
                return nil
            }
            recvMessages(c.String("network-id"), token, c.String("host"), c.String("protocol"), c.Int64("limit"))
        } else {
            // set message
            message := c.String("message")
            msgBytes := []byte(message)

            // set message by file
            if c.String("file") != "" {
                bytes, err := ioutil.ReadFile(c.String("file"))
                check(err)
                msgBytes = bytes
                message = fmt.Sprintf("%s", bytes)
            }

            // get hash
            if c.String("action") == "hash" {
                hashResult := getHashResult(msgBytes, c.String("hash-function"))
                body, err := json.Marshal(hashResult)
                check(err)
                fmt.Println(fmt.Sprintf("%s", body))
            }

            // set message context
            if c.String("action") == "send" {
                msgContext := MessageContext{
                    Status: "published",
                    Message: message,
                    Agent: c.String("agent"),
                    Format: c.String("format"),
                    Delimiter: c.String("delimiter"),
                    Encoding: c.String("encoding"),
                    Mime: c.String("mime"),
                    Label: c.String("label"),
                    Encryption: c.String("encryption"),
                    EncryptionKey: c.String("encryption-key"),
                    EncryptionIv: c.String("encryption-iv"),
                    PrivateKey: c.String("private-key"),
                    PublicKey: c.String("public-key"),
                    HashFunction: c.String("hash-function"),
                    HashValue: c.String("hash-value"),
                    Mnemonic: c.String("mnemonic"),
                    IntNetwork: c.String("int-network"),
                    IntAddress: c.String("int-address"),
                    ExtNetwork: c.String("ext-network"),
                    ExtAddress: c.String("ext-address"),
                    Host: c.String("host"),
                    Lang: c.String("lang"),
                    NetworkId: c.String("network-id"),
                    AccessKey: c.String("access-key"),
                    AccessSecret: c.String("access-secret"),
                    Protocol: c.String("protocol"),
                    TimeToLive: c.Int64("time-to-live"),
                    Latitude: c.Float64("latitude"),
                    Longitude: c.Float64("longitude"),
                    Geohash: c.String("geohash"),
                    Planet: c.String("planet"),
                }
                sendMessage(msgContext, token, c.String("host"), c.String("protocol"))
            }
        }
        return nil
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
