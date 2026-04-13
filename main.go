package main

import (
	"fmt"
	"github.com/tattoo1880/protkit/codec"
	pb "github.com/tattoo1880/protkit/proto"
)

func main() {
	user := &pb.User{
		Id:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Status:   pb.UserStatus_USER_STATUS_ACTIVE,
	}

	data, _ := codec.Marshal(user)
	fmt.Println(string(data))
	got := &pb.User{}
	err1 := codec.Unmarshal(data, got)
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println(got)
}
