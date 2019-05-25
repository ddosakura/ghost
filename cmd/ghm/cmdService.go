package main

import (
	"fmt"

	"net/url"

	"github.com/ddosakura/ghost/cmd"
	"github.com/ddosakura/ghost/cmd/proto/sign"
	proto "github.com/golang/protobuf/proto"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2"
)

// Cmd
var (
	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "manage status of service",
		Long:  `Service start/stop/status (recommend starting as system service).`,
		RunE: func(c *cobra.Command, args []string) error {
			if u, e := url.Parse("http://" + upService); e == nil {
				// pretty.Println(u)
				if u.User == nil {
					upServiceData.user = url.UserPassword("ghost", "123456")
				} else {
					upServiceData.user = u.User
				}
				upServiceData.host = u.Host
			} else {
				return e
			}
			// pretty.Println(upServiceData)
			switch serviceSign {
			case "start":
				session, e := mgo.Dial(mongoURL)
				// TODO: learn mongodb's mode
				session.SetMode(mgo.Monotonic, true)
				//pretty.Println(session)
				if e != nil {
					return e
				}
				defer session.Close()

				m := newRepo(session)
				if e = m.init(dbName); e != nil {
					return e
				}

				fmt.Println("finish")
				//s, e := net.Listen("tcp", serviceAddr)
				//defer s.Close()
			case "stop":
				bs, e := proto.Marshal(&sign.Request{
					Type: sign.Type_STOP,
				})
				pretty.Println(bs)
				if e != nil {
					return e
				}
			case "status":
				bs, e := proto.Marshal(&sign.Request{
					Type: sign.Type_STATUS,
				})
				pretty.Println(bs)
				if e != nil {
					return e
				}
			default:
				return cmd.ErrUnknowServiceSign
			}
			return nil
		},
	}

	serviceSign   string // 服务操作
	serviceAddr   string // 服务监听
	upService     string // 上游地址 (include Auth)
	upServiceData struct {
		user *url.Userinfo
		host string
	}
	useKCP bool // TODO: kcp

	mongoURL string // db URL (include Auth)
	dbName   string // db Name
)

func init() {
	serviceCmd.PersistentFlags().StringVarP(
		&serviceSign,
		"sign", "s",
		"start",
		"sign sended to service",
	)
	serviceCmd.PersistentFlags().StringVarP(
		&serviceAddr,
		"addr", "a",
		cmd.AddrOfMaster,
		"address of this service",
	)
	serviceCmd.PersistentFlags().StringVarP(
		&upService,
		"up", "u",
		"ghost:123456@g",
		"address of up service - [user:pass@]host[:port] or init top domain - [user:pass@]domain",
	)
	serviceCmd.PersistentFlags().BoolVarP(
		&useKCP,
		"kcp", "",
		false,
		"replacing TCP with KCP",
	)

	serviceCmd.PersistentFlags().StringVarP(
		&mongoURL,
		"mongo", "m",
		"mongodb://root:123456@127.0.0.1:27017",
		"the url of MongoDB - [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]",
	)
	serviceCmd.PersistentFlags().StringVarP(
		&dbName,
		"db", "",
		"ghost-master",
		"the name of db",
	)
}
