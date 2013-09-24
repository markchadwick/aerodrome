package net

import (
	"github.com/markchadwick/spec"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func Test(t *testing.T) {
	spec.Run(t)
}
