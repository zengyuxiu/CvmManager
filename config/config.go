package config

import (
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type NetworkConfig struct {
	XMLName xml.Name `xml:"xml"`
	Body    struct {
		Route  string `xml:"route"`
		Bridge []struct {
			Name     string `xml:"name"`
			Type     string `xml:"type"`
			Address  string `xml:"address"`
			Netmask  string `xml:"netmask"`
			OspfArea string `xml:"ospf_area"`
		} `xml:"bridge"`
		Instance []struct {
			Name        string `xml:"name"`
			RouteConfig string
		} `xml:"instance"`
		Link []struct {
			Instance string `xml:"instance"`
			Bridge   string `xml:"bridge"`
		} `xml:"link"`
	} `xml:"body"`
}

func GetConfig() *NetworkConfig {
	config_xml := NetworkConfig{}
	file, err := os.Open("../example/TestEnv.xml")
	if err != nil {
		log.Error(err)
	}
	bytes, _ := ioutil.ReadAll(file)
	err = xml.Unmarshal(bytes, &config_xml)
	if err != nil {
		log.Error(err)
	}
	//fmt.Printf(config_xml.Body.Route)
	return &config_xml
}

func (config *NetworkConfig) RouteConfigGen() {
	if config.Body.Route == "ospfv2" {
	}
}

func (config *NetworkConfig) OspfdConfig() {
	for instance_id, instance := range config.Body.Instance {
		config.Body.Instance[instance_id].RouteConfig = fmt.Sprintf("router ospf\n"+
			"ospf router-id 0.0.0.%d\n", instance_id)
		for _, link := range config.Body.Link {
			if instance.Name == link.Instance {
				for _, bridge := range config.Body.Bridge {
					if bridge.Name == link.Bridge {
						config.Body.Instance[instance_id].RouteConfig += fmt.Sprintf("network "+
							"%s/%s area %s\n", bridge.Address, bridge.Netmask, bridge.OspfArea)
					}
				}
			}
		}
	}
	/*	for _,instance := range config.Body.Instance{
		fmt.Print(instance.RouteConfig)
	}*/
}
