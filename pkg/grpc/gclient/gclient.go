package gclient

import (
	"fmt"
	"goshop/service-order/pkg/grpc/etcd3"
	"goshop/service-order/pkg/utils"
	"log"
	"strings"

	"github.com/shinmigo/pb/memberpb"
	"github.com/shinmigo/pb/productpb"
	"github.com/shinmigo/pb/shoppb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var (
	ProductTag            productpb.TagServiceClient
	ProductParam          productpb.ParamServiceClient
	ProductKind           productpb.KindServiceClient
	Member                memberpb.MemberServiceClient
	ProductCategoryClient productpb.CategoryServiceClient
	ProductSpecClient     productpb.SpecServiceClient
	ProductClient         productpb.ProductServiceClient
	ShopUser              shoppb.UserServiceClient
	ShopCarrier           shoppb.CarrierServiceClient
)

func DialGrpcService() {
	shop()
	pms()
	crm()
}

func shop() {
	r := etcd3.NewResolver(utils.C.Etcd.Host)
	resolver.Register(r)
	fmt.Println(utils.C.GrpcClient.Name["shop"])
	conn, err := grpc.Dial(r.Scheme()+"://author/"+utils.C.GrpcClient.Name["shop"], grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		log.Panicf("grpc没有连接上%s, err: %v \n", utils.C.GrpcClient.Name["shop"], err)
	}
	fmt.Printf("连接成功：%s, host分别为: %s \n", utils.C.GrpcClient.Name["shop"], strings.Join(utils.C.Etcd.Host, ","))
	ShopUser = shoppb.NewUserServiceClient(conn)
	ShopCarrier = shoppb.NewCarrierServiceClient(conn)
}

func crm() {
	r := etcd3.NewResolver(utils.C.Etcd.Host)
	resolver.Register(r)
	conn, err := grpc.Dial(r.Scheme()+"://author/"+utils.C.GrpcClient.Name["crm"], grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		log.Panicf("grpc没有连接上%s, err: %v \n", utils.C.GrpcClient.Name["crm"], err)
	}
	fmt.Printf("连接成功：%s, host分别为: %s \n", utils.C.GrpcClient.Name["crm"], strings.Join(utils.C.Etcd.Host, ","))
	Member = memberpb.NewMemberServiceClient(conn)
}

func pms() {
	r := etcd3.NewResolver(utils.C.Etcd.Host)
	resolver.Register(r)

	//这里后面会有多个grpc服务，
	conn, err := grpc.Dial(r.Scheme()+"://author/"+utils.C.GrpcClient.Name["pms"], grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		log.Panicf("grpc没有连接上%s, err: %v \n", utils.C.GrpcClient.Name["pms"], err)
	}
	fmt.Printf("连接成功：%s, host分别为: %s \n", utils.C.GrpcClient.Name["pms"], strings.Join(utils.C.Etcd.Host, ","))
	ProductTag = productpb.NewTagServiceClient(conn)
	ProductParam = productpb.NewParamServiceClient(conn)
	ProductCategoryClient = productpb.NewCategoryServiceClient(conn)
	ProductKind = productpb.NewKindServiceClient(conn)
	ProductSpecClient = productpb.NewSpecServiceClient(conn)
	ProductClient = productpb.NewProductServiceClient(conn)
}
