package gen

import (
	"fmt"
	"strings"

	"github.com/zhufuyi/goctl/pkg/replacer"
	"github.com/zhufuyi/pkg/gofile"
)

// 微服务子模块名称
const (
	// ModuleSponge sponge 模块，http和grpc共有的服务
	ModuleSponge = "sponge"

	// ModuleModel model 模块, http和grpc共有
	ModuleModel = "model"

	// ModuleDao dao 模块，http和grpc共有
	ModuleDao = "dao"

	// ModuleHandler handler 模块，只属于http子模块
	ModuleHandler = "handler"

	// ModuleHTTP http 服务模块，只属于http子模块
	ModuleHTTP = "http"

	// ModuleProto proto 模块，只属于grpc子模块
	ModuleProto = "proto"

	// ModuleService service 模块，只属于grpc子模块
	ModuleService = "service"

	// ModuleGRPC grpc服务 模块，只属于grpc子模块
	ModuleGRPC = "grpc"
)

// MicroServiceGroupModules 微服务子模块群主
var MicroServiceGroupModules = []string{ModuleModel, ModuleDao, ModuleHandler, ModuleHTTP, ModuleSponge, ModuleProto, ModuleService, ModuleGRPC}

var allModules = []string{ModuleModel, ModuleDao, ModuleHandler, ModuleHTTP, ModuleSponge, ModuleProto, ModuleService, ModuleGRPC}

var (
	// 指定文件替换标记
	modelFile     = "model/userExample.go"
	modelFileMark = "// todo generate model codes to here"

	daoFile     = "dao/userExample.go"
	daoFileMark = "// todo generate the update fields code to here"

	handlerFile     = "handler/userExample.go"
	handlerFileMark = "// todo generate the request and response struct to here"

	mainFile     = "sponge/main.go"
	mainFileMark = "// todo generate the code to register http and grpc services here"

	protoFile     = "v1/userExample.proto"
	protoFileMark = "// todo generate the protobuf code here"

	serviceFile     = "service/userExample_test.go"
	serviceFileMark = "// todo generate the service struct code here"

	// 清除标记的模板代码片段标记
	startMark     = []byte("// delete the templates code start")
	endMark       = []byte("// delete the templates code end")
	grpcStartMark = []byte("// grpc import start")
	grpcEndMark   = []byte("// grpc import end")

	selfPackageName = "github.com/zhufuyi/goctl"
)

func adjustmentOfIDType(handlerCodes string) string {
	return idTypeToStr(idTypeFixToUint64(handlerCodes))
}

func idTypeFixToUint64(handlerCodes string) string {
	subStart := "ByIDRequest struct {"
	subEnd := "`" + `json:"id" binding:""` + "`"
	if subBytes := gofile.FindSubBytesNotIn([]byte(handlerCodes), []byte(subStart), []byte(subEnd)); len(subBytes) > 0 {
		old := subStart + string(subBytes) + subEnd
		newStr := subStart + "\n\tID uint64 " + subEnd + " // uint64 id\n"
		handlerCodes = strings.ReplaceAll(handlerCodes, old, newStr)
	}

	return handlerCodes
}

func idTypeToStr(handlerCodes string) string {
	subStart := "ByIDRespond struct {"
	subEnd := "`" + `json:"id"` + "`"
	if subBytes := gofile.FindSubBytesNotIn([]byte(handlerCodes), []byte(subStart), []byte(subEnd)); len(subBytes) > 0 {
		old := subStart + string(subBytes) + subEnd
		newStr := subStart + "\n\tID string " + subEnd + " // covert to string id\n"
		handlerCodes = strings.ReplaceAll(handlerCodes, old, newStr)
	}

	return handlerCodes
}

func addTheDeleteFields(r replacer.Replacer, filename string) []replacer.Field {
	var fields []replacer.Field

	data, err := r.ReadFile(filename)
	if err != nil {
		fmt.Printf("read the file '%s' error: %v\n", filename, err)
		return fields
	}
	if subBytes := gofile.FindSubBytes(data, startMark, endMark); len(subBytes) > 0 {
		fields = append(fields,
			replacer.Field{ // 清除标记的模板代码
				Old: string(subBytes),
				New: "",
			},
		)
	}

	return fields
}

func addTheDeleteGrpcFields(r replacer.Replacer, filename string) []replacer.Field {
	var fields []replacer.Field

	data, err := r.ReadFile(filename)
	if err != nil {
		fmt.Printf("read the file '%s' error: %v\n", filename, err)
		return fields
	}
	if subBytes := gofile.FindSubBytes(data, grpcStartMark, grpcEndMark); len(subBytes) > 0 {
		fields = append(fields,
			replacer.Field{ // 清除标记的模板代码
				Old: string(subBytes),
				New: "",
			},
		)
	}

	return fields
}

const httpServerRegisterCode = `func registerServers() []app.IServer {
	var servers []app.IServer

	// 创建http服务
	httpAddr := ":" + strconv.Itoa(config.Get().HTTP.Port)
	httpServer := server.NewHTTPServer(httpAddr,
		server.WithHTTPReadTimeout(time.Second*time.Duration(config.Get().HTTP.ReadTimeout)),
		server.WithHTTPWriteTimeout(time.Second*time.Duration(config.Get().HTTP.WriteTimeout)),
		server.WithHTTPIsProd(config.Get().App.Env == "prod"),
	)
	servers = append(servers, httpServer)

	return servers
}`

const grpcServerRegisterCode = `func registerServers() []app.IServer {
	var servers []app.IServer

	// 创建grpc服务
	grpcAddr := ":" + strconv.Itoa(config.Get().Grpc.Port)
	grpcServer := server.NewGRPCServer(grpcAddr, grpcOptions()...)
	servers = append(servers, grpcServer)

	return servers
}

func grpcOptions() []server.GRPCOption {
	var opts []server.GRPCOption

	if config.Get().App.EnableRegistryDiscovery {
		iRegistry, instance := getETCDRegistry(
			config.Get().Etcd.Addrs,
			config.Get().App.Name,
			[]string{fmt.Sprintf("grpc://%s:%d", config.Get().App.Host, config.Get().Grpc.Port)},
		)
		opts = append(opts, server.WithRegistry(iRegistry, instance))
	}

	return opts
}

func getETCDRegistry(etcdEndpoints []string, instanceName string, instanceEndpoints []string) (registry.Registry, *registry.ServiceInstance) {
	serviceInstance := registry.NewServiceInstance(instanceName, instanceEndpoints)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	})
	if err != nil {
		panic(err)
	}
	iRegistry := etcd.New(cli)

	return iRegistry, serviceInstance
}`
