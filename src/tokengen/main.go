//
// Copyright © 2023 wzlu@amazon.com
//
package main

import (
    "schedrestd/common"
    "schedrestd/config"
    "schedrestd/common/jwt"

    rawjwt "github.com/dgrijalva/jwt-go"
    "github.com/boltdb/bolt"
    "os/user"
    "os"
    "log"
    "time"
    "fmt"
    "strconv"
)

func main() {
    var timeOut int64
    var err error
    conf := config.NewConfig()
    dbpath := conf.WorkDir + "/" + common.BoltDBName
    currentUser, err := user.Current()
    if err != nil {
        log.Fatalf(err.Error())
        os.Exit(1)
    }

    userName := currentUser.Username
    timeOutStr := ""
    if userName == "root" {
        if len(os.Args) < 2 {
            fmt.Println("Usage: " + os.Args[0] + " username [duration_minutes]")
            os.Exit(0)
        }
        userName = os.Args[1]
        if len(os.Args) == 3 {
            timeOutStr = os.Args[2]
        }
    } else {
        if len(os.Args[2]) == 2 {
            timeOutStr = os.Args[1]
        }
    }
    if timeOutStr != "" {
        timeOut, err = strconv.ParseInt(os.Args[2], 10, 64)
        if err != nil {
            log.Fatalf(err.Error())
        }
    } else {
        timeOut, err = strconv.ParseInt(conf.Timeout, 10, 64)
        if err != nil {
            timeOut = 5256000 // 10 years
        }
    }

    claims := jwt.AClaims {
           userName,
           rawjwt.StandardClaims {
               Issuer: "Schedrestd",
               IssuedAt: time.Now().Unix(),
           },
    }

    j := jwt.NewJWT()
    token, err := j.GenerateToken(claims)
    if err != nil {
        log.Fatalf(err.Error())
        os.Exit(1)
    }

    key := fmt.Sprintf("%v-%v", token, userName)
    expireTime := time.Now().Add(time.Minute * time.Duration(timeOut)).Unix()

    db, err := bolt.Open(dbpath, 0600, nil)
    if err != nil {
         log.Fatalf(err.Error())
         os.Exit(1)
    }
    err = db.Update(func(tx *bolt.Tx) error {
         _, err = tx.CreateBucketIfNotExists([]byte(common.BoltDBJWTTable))
         return err
    })
    err = db.Update(func(tx *bolt.Tx) error {
         b := tx.Bucket([]byte(common.BoltDBJWTTable))
         return b.Put([]byte(key),
                []byte(strconv.FormatInt(expireTime, 10)))
    })
    db.Close()
    if err != nil {
         log.Fatalf(err.Error())
         os.Exit(1)
    }
    fmt.Println(token)
}
