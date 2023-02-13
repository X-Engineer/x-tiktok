package redis

import (
	"testing"
)

func TestConnRedis(t *testing.T) {
	connRedis()
}

func TestSetValue(t *testing.T) {
	setValue("name", "zhicheng")
}

func TestGetValue(t *testing.T) {
	getValue("name")
}
