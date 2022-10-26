package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func RandomVerCode() string {
	rand.Seed(time.Now().UnixMilli())
	str := fmt.Sprintf("%.6f", rand.Float32())[2:]
	return str
}

func RandStr(length int) string {
	arr := make([]byte, length)
	for i := 0; i < length; i++ {
		rand.Seed(time.Now().UnixMilli() + int64(i))
		n := byte(rand.Intn(26))
		arr[i] = 'a' + n
	}
	return string(arr)
}

func Struct2Map(obj interface{}) map[string]interface{} {
	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	var res map[string]interface{}
	json.Unmarshal(bytes, &res)
	return res
}

func Map2Struct(m map[string]interface{}, obj interface{}) {
	bytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bytes, obj)
}
