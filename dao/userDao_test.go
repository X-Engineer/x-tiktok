package dao

import (
	"log"
	"testing"
)

func TestGetUserBasicInfoByIdt(t *testing.T) {
	res, err := GetUserBasicInfoById(5)
	if err == nil {
		log.Println(res)
	}
}

func TestGetUserBasicInfoByName(t *testing.T) {
	res, err := GetUserBasicInfoByName("qcj")
	if err == nil {
		log.Println(res)
	}
}
