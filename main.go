package main

import (
	// "context"
	// "fmt"
	// "io"
	// "fmt"
	"log"
	// "os"
	"net/http"
	// "encoding/json"
	"strings"

	// "github.com/docker/docker/api/types"
	// "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	// "github.com/gorilla/mux"
)

type RestDockerHandler struct{
	client *client.Client
}

type Container struct{
	image string
	localport string
	hostport string
}

func main() {
	//dockerclientを起動し、ハンドラーを作成
	RestDockerHandler := InitClient()
	//デフォルトハンドラーに登録
	http.Handle("/", RestDockerHandler)
	// ポート番号とインターフェイス（親ハンドラー）を結び付ける
	http.ListenAndServe(":8000", nil)
}

//ListenAndServeのServe関数がhandleを呼び起こしてスレッドを作成する
//handlerの引数の形が固定されるのであれば、メソッドにしてかませればいいじゃん！！
//コンテナIDはstring
func (docker *RestDockerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET"{
		//URIは/allもしくは/{id}を想定
		path := r.URL.Path
		//分割されたものはリストとして格納
		parts := strings.Split(path, "/")
		//partsのリストの0番目は空
		// /all⇒　"","all"
		if len(parts) == 2{
			if parts[1] == "all"{
				GetAllContainer(docker.client)
			}else{
				id := parts[1]
				GetContainer(docker.client,id)
			}
		}
	}else if r.Method == "POST"{
		//URIは/{image}/{localport}/{hostport}を想定
		path := r.URL.Path
		parts := strings.Split(path, "/")
		// fmt.Printf("%d", len(parts))	
		if len(parts) == 3{
			localport := parts[2]
			image := parts[1]
			config := NewContainer(image,localport,"80")
			CreateContainer(docker.client,config)
			// fmt.Fprintln(w,"hello")
		}else if len(parts) == 4{
			hostport := parts[3]
			localport := parts[2]
			image := parts[1]
			config := NewContainer(image,localport,hostport)
			CreateContainer(docker.client,config)
			// fmt.Fprintln(w,"hello")
		}
	}else if r.Method == "PUT"{
		//URIは/allもしくは/{id}を想定
		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) == 2{
			if parts[1] == "all"{
				UpdateAllContainer(docker.client)
				// fmt.Fprintln(w,"hello")
			}else{
				id := parts[1]
				UpdateContainer(docker.client,id)	
				// fmt.Fprintln(w,"hello")	
			}
		}
	}else if r.Method == "DELETE"{
		path := r.URL.Path
		//URIは/allもしくは/{id}を想定
		parts := strings.Split(path, "/")
		if len(parts) == 2{
			if parts[1] == "all"{
				DeleteAllContainer(docker.client)
				// fmt.Fprintln(w,"hello")
			}else{
				id := parts[1]
				DeleteContainer(docker.client,id)
				// fmt.Fprintln(w,"hello")	
			}
		}
	}

}

func Dockeractivate() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

func InitClient() *RestDockerHandler{
	return &RestDockerHandler{
		client: Dockeractivate(),
	}
}

func NewContainer(dockerimage string, localport string, hostport string) *Container{
	return &Container{
		image: dockerimage,
		localport: localport,
		hostport: hostport,
	}
}


	// エンドポイントとハンドラ関数（子ハンドラー）を関連付けというよりかは、アクセスによってインターフェイスが仕訳けているイメージカナ
	// http.HandleFunc("/", Defaulthandler).Methods("GET")
	// http.HandleFunc("/all", GetAllContainer{}).Methods("GET")
	// http.HandleFunc("/{id}", GetContainer{}).Methods("GET")
	// http.HandleFunc("/{image}/{hostport}", CreateContainer{}).Methods("POST")
	// http.HandleFunc("/{id}", UpdateContainer{}).Methods("PUT")
	// http.HandleFunc("/all", UpdateAllContainer{}).Methods("PUT")
	// http.HandleFunc("/{id}", DeleteContainer{}).Methods("DELETE")
	// http.HandleFunc("/all", DeleteAllContainer{}).Methods("DELETE")
	//おそらく、二つ目の引数はServeHTTPメソッドを実装している構造体（のインターフェイス型）