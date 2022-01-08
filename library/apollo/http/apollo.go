//vim ~/.bash_profile
//export RUNTIME_CLUSTER=development     #填写对应的集群名称
//export CONFIG_CENTER_URL=http://development.apollo.com   #填写开发环境配置中心域名
//export CONFIG_CENTER_TOKEN=abc123   #填写正确的token信息
//保存退出
//source ~/.bash_profile  #使刚才的环境变量生效
//
//验证是否配置成功
//echo $RUNTIME_CLUSTER
//echo CONFIG_CENTER_URL
//echo CONFIG_CENTER_TOKEN
//echo CONFIG_CENTER_APPID

package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

type AppoloConf struct {
	AppId          string            `json:"appId"`
	Cluster        string            `json:"cluster"`
	NamespaceName  string            `json:"namespaceName"`
	Configurations map[string]string `json:"configurations"`
	ReleaseKey     string            `json:"releaseKey"`
}

var apolloOnce sync.Once
var apolloConfigs map[string]string

func DoLoadApolloConf(host, service, cluster, token string, space []string) (map[string]string, error) {
	client := &http.Client{
		Timeout: time.Second * 60, // Notifications由于服务端会hold住请求60秒，所以请确保客户端访问服务端的超时时间要大于60秒。
	}

	cfgMap := make(map[string]string)

	for _, v := range space {
		ServiceUriFmt := fmt.Sprintf("%s/configs/%s/%s/%s", host, service, cluster, v)
		req, err := http.NewRequest("GET", ServiceUriFmt, nil)
		query := req.URL.Query()
		query.Add("token", token)
		req.URL.RawQuery = query.Encode()
		response, err := client.Do(req)

		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("http status is not ok.")
		}

		responseData, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return nil, err
		}
		var ac AppoloConf
		err = json.Unmarshal(responseData, &ac)
		if err != nil {
			return nil, err
		}
		cfgMap = ac.Configurations
	}

	return cfgMap, nil
}

// LoadApolloConf
// example: LoadApolloConf("777", []string{"application"})
func LoadApolloConf(service string, space []string) (map[string]string, error) {
	var err error
	apolloOnce.Do(func() {
		apolloConfigs = make(map[string]string)
		host := os.Getenv("CONFIG_CENTER_URL")
		cluster := os.Getenv("RUNTIME_CLUSTER")
		token := os.Getenv("CONFIG_CENTER_TOKEN")
		apolloConfigs, err = DoLoadApolloConf(host, service, cluster, token, space)
		fmt.Println("get apollo")
	})

	return apolloConfigs, err
}
