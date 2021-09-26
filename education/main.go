package main

import (
	"education/sdkInit"
	"education/service"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	cc_name    = "simplecc"
	cc_version = "1.0.0"
)

var serviceSetup *service.ServiceSetup

func goSdkInit() {
	// init orgs information
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: "fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
	}

	// init sdk env info
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    "fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    "chaincode/",
		ChaincodeVersion: cc_version,
	}

	// sdk setup
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Println(">> create chaincode lifecycle error:", err)
		os.Exit(-1)
	}

	// invoke chaincode set status
	fmt.Println(">> 通过链码外部服务设置链码状态......")

	serviceSetup, err = service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk)
	if err != nil {
		fmt.Println()
		os.Exit(-1)
	}
}

func ginInit() {
	r := gin.Default()
	r.GET("/add", func(c *gin.Context) {
		name := c.DefaultQuery("Name", "no-name")
		gender := c.DefaultQuery("Gender", "unknown")
		age := c.DefaultQuery("Age", "0")
		td := service.Test{
			Name:   name,
			Gender: gender,
			Age:    age,
		}
		res, err := serviceSetup.ServiceTestAdd(td)
		if err != nil {
			c.String(http.StatusOK, err.Error())
		} else {
			c.String(http.StatusOK, "add user success\n"+res)
		}
	})
	r.GET("/query", func(c *gin.Context) {
		name := c.DefaultQuery("Name", "no-name")
		res, err := serviceSetup.ServiceQueryTest(name)
		if err != nil {
			c.String(http.StatusOK, err.Error())
		} else {
			c.String(http.StatusOK, res)
		}
	})
	r.Run(":9000")
}

func main() {
	goSdkInit()
	ginInit()
}
