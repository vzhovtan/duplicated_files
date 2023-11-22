package main

import (
	"fmt"
	"log"
	"testing"
)

var tRoot = "./testdata"

var tHash = "d41d8cd98f00b204e9800998ecf8427e"

var tStat = `&[{empty1.txt d41d8cd98f00b204e9800998ecf8427e} {empty2.txt d41d8cd98f00b204e9800998ecf8427e} {file.txt 122a10d6a32262217e0e79e504f2e447} {gopher1.png ca1f746d6f232f87fca4e4d94ef6f3ab} {gopher2.png ca1f746d6f232f87fca4e4d94ef6f3ab} {gopher3.png ca1f746d6f232f87fca4e4d94ef6f3ab} {gopher4.png ca1f746d6f232f87fca4e4d94ef6f3ab}]`

func TestStat(t *testing.T) {
	data, err := getStat(tRoot)
	if err != nil {
		log.Fatal(err)
	}
	got := fmt.Sprintf("%v", data)
	if got != tStat {
		t.Errorf("test GetStat failed - results not match\nGot:\n%v\nExpected:\n%v", got, tStat)
	}
}

func TestHash(t *testing.T) {
	got := getHash(tRoot, "empty1.txt")
	if got != tHash {
		t.Errorf("test getHash Failed - results not match\nGot:\n%v\nExpected:\n%v\n", got, tHash)
	}
}
