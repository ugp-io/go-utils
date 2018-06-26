package utils


import(
    "crypto/sha256"
    "encoding/hex"
    "io"
    "crypto/rand"
    "encoding/base64"
    //"crypto/md5"
    //"io"
    //"gopkg.in/mgo.v2/bson"
    //"errors"
    //"strings"
    //"strconv"
    //"github.com/dgrijalva/jwt-go"
    //scope "github.com/chowt/chowt-api/chowt/scope"
)

func HashString(salt string, input string) (string){

    h256 := sha256.New()
    io.WriteString(h256, salt + input)
    hashedString := hex.EncodeToString(h256.Sum(nil))

    return hashedString
}

func StringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func GenerateRandomString(n int) (string, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), err
}
