package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"userProfile/pkg/userProfile"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"

	"log"
)

func main() {
	backend := flag.String("b", "localhost:10000", "address of userProfile backend")

	conn, err := grpc.Dial(*backend, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to %s: %v", *backend, err)
	}
	defer conn.Close()

	client := userProfile.NewUserProfileClient(conn)
	ctx := context.Background()
	//callRegisterInfo(ctx, client)
	//time.Sleep(time.Second * 1)
	callGetClientInfo(ctx, client)
}

func callRegisterInfo(ctx context.Context, client userProfile.UserProfileClient) {
	params := &userProfile.RegisterRequest{
		RequestTimestamp: &timestamp.Timestamp{},
		UserProfile: &userProfile.UserProfile{
			ClientID: "123",
			ClientInfo: &userProfile.ClientInfo{
				ClientVersion:     "1.0.0",
				ClientVersionCode: 1,
				Model:             "test",
				Language:          userProfile.Language_ENGLISH,
			},
			AndroidInfo: &userProfile.AndroidInfo{
				SdkVersion: 21,
				Dpi:        480,
				Cpu:        "x86",
				Gpu:        "Mali-G72",
			},
			Features: []*userProfile.Feature{
				{
					Name: "com.test",
				},
				{
					Name: "com.test2",
				},
			},
			Libraries: []*userProfile.Library{
				{
					Name: "com.test",
				},
				{
					Name: "com.test2",
				},
			},
		},
	}

	res, err := client.RegisterClientInfo(ctx, params)

	if err != nil {
		log.Fatalf("could not call callRegisterInfo %v", err)
	}
	fmt.Println("input:")
	spew.Dump(params)
	fmt.Println("output:")
	spew.Dump(res)
	//resJson, err := json.Marshal(res)
	//if err != nil {
	//	log.Fatalf("could marshal callRegisterInfo response: %v", err)
	//}
	//err = ioutil.WriteFile("response.json", resJson, 0644)
}

func callGetClientInfo(ctx context.Context, client userProfile.UserProfileClient) {
	params := &userProfile.ClientInfoRequest{
		RequestTimestamp: &timestamp.Timestamp{},
		ClientID:         "123",
		Contains: &userProfile.Contains{
			ClientInfo:  true,
			AndroidInfo: true,
			Library:     true,
			Feature:     true,
		},
	}

	res, err := client.GetClientInfo(ctx, params)

	if err != nil {
		log.Fatalf("could not call callRegisterInfo %v", err)
	}
	fmt.Println("input:")
	spew.Dump(params)
	fmt.Println("output:")
	spew.Dump(res)
	resJson, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("could marshal callRegisterInfo response: %v", err)
	}
	err = ioutil.WriteFile("response.json", resJson, 0644)
}
