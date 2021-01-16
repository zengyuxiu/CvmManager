package lxc_api

import (
	"CvmManager/config"
	"database/sql"
	"fmt"
	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var NameLength = 5

func init() {
	rand.Seed(time.Now().UnixNano())
}

func LxcStatus(sqliteDatabase *sql.DB) error {
	c, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		log.Error(err)
	}
	containers, err := c.GetContainers()
	for _, container := range containers {
		state, _, err := c.GetInstanceState(container.Name)
		if err != nil {
			log.Error(err)
		} else {
			InsertStatLxc(sqliteDatabase, state, container.Name, time.Now().UTC())
		}
	}
	return nil
}

//goland:noinspection SqlResolve
func InsertStatLxc(db *sql.DB, stat *api.InstanceState, name string, time time.Time) {
	var InsertSQL = `INSERT INTO lxc_stats(
                    name,
                    cpu_usage,
                    memory_usage,
                    memory_max_usage,
                    time_stamp
                    ) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(InsertSQL)
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(name, stat.CPU.Usage,
		stat.Memory.Usage, stat.Memory.UsagePeak,
		time,
	)
	if err != nil {
		log.Error(err)
	}
}

func LxcCreate(InstanceNumber int, ImageFingerprint string) error {
	c, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		log.Error(err)
	}
	for i := 1; i <= InstanceNumber; i++ {
		name := "test" + GetRandomString(NameLength)
		request := api.ContainersPost{
			Name: name,
			Source: api.ContainerSource{
				Type:        "image",
				Fingerprint: ImageFingerprint,
			},
		}
		op, err := c.CreateContainer(request)
		if err != nil {
			log.Error(err)
		}
		err = op.Wait()
		if err != nil {
			log.Error(err)
		}

		reqState := api.InstanceStatePut{
			Action:  "start",
			Timeout: -1,
			Force:   true,
		}

		op, err = c.UpdateInstanceState(name, reqState, "")
		if err != nil {
			log.Error(err)
		}
	}

	// Wait for the operation to complete
	return nil
}

func LxcDelete() error {
	c, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		log.Error(err)
	}
	containers, err := c.GetContainers()
	for _, container := range containers {
		if err != nil {
			log.Error(err)
			return err
		} else if match, _ := regexp.Match("test(.*)", []byte(container.Name)); match == true {
			reqState := api.InstanceStatePut{
				Action:  "stop",
				Timeout: -1,
				Force:   true,
			}
			op, err := c.UpdateInstanceState(container.Name, reqState, "")
			if err != nil {
				log.Errorf("%s:%s", container.Name, err)
			}
			err = op.Wait()
			if err != nil {
				log.Errorf("%s:%s", container.Name, err)
			}
			op, err = c.DeleteInstance(container.Name)
			if err != nil {
				log.Errorf("%s:%s", container.Name, err)
			}
			// Wait for the operation to complete
		}
	}
	return nil
}

func GetRandomString(n int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < n; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func BridgeCreate(config *config.NetworkConfig) {
	c, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		log.Error(err)
	}
	for _, bridge := range config.Body.Bridge {
		if bridge.Type == "lxdbrd" {
			BrConfig := make(map[string]string)
			BrConfig["ipv4.nat"] = "false"
			BrConfig["ipv4.address"] = fmt.Sprintf("%s/%s", bridge.Address, bridge.Netmask)
			err := c.CreateNetwork(api.NetworksPost{
				Name: bridge.Name,
				Type: "bridge",
			})
			if err != nil {
				log.Error(err)
			}
			err = c.UpdateNetwork(bridge.Name, api.NetworkPut{Config: BrConfig}, "")
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func BridgeAttach(config *config.NetworkConfig) {
	for _, link := range config.Body.Link {
		args := fmt.Sprintf("network attach %s %s --force-local", link.Bridge, link.Instance)
		//args := fmt.Sprintf("network attach lxdbr0 route0 --force-local")
		cmd := exec.Command("lxc", strings.Split(args, " ")...)
		err := cmd.Run()
		if err != nil {
			log.Error(err)
			fmt.Print(link.Instance)
		}
	}
}

func InstanceOspfConfig(config *config.NetworkConfig) {
	for _, container := range config.Body.Instance {
		config_file := fmt.Sprintf("service integrated-vtysh-config\n%s", container.RouteConfig)
		file, err := os.Create("/tmp/frr.conf")
		if err != nil {
			log.Error(err)
		} else {
			_, err = file.Write([]byte(config_file))
		}
		args := fmt.Sprintf("file push /tmp/frr.conf %s/etc/frr/frr.conf --force-local", container.Name)
		cmd := exec.Command("lxc", strings.Split(args, " ")...)
		err = cmd.Run()
		if err != nil {
			log.Error(err)
		}
		args = fmt.Sprintf("exec %s -- vtysh -b  --force-local", container.Name)
		cmd = exec.Command("lxc", strings.Split(args, " ")...)
		err = cmd.Run()
		if err != nil {
			log.Error(err)
		}
	}
}

func CreateRoute(networkConfig *config.NetworkConfig) {
	c, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		log.Error(err)
	}
	for _, route := range networkConfig.Body.Instance {
		request := api.ContainersPost{
			Name: route.Name,
			Source: api.ContainerSource{
				Type:        "image",
				Fingerprint: "d95b3e925df1",
			},
		}
		op, err := c.CreateContainer(request)
		if err != nil {
			log.Error(err)
		}
		err = op.Wait()
		if err != nil {
			log.Error(err)
		}

		reqState := api.InstanceStatePut{
			Action:  "start",
			Timeout: -1,
			Force:   true,
		}

		op, err = c.UpdateInstanceState(route.Name, reqState, "")
		if err != nil {
			log.Error(err)
		}
	}
}

func RestartRoute(config *config.NetworkConfig) {
	for _, container := range config.Body.Instance {
		args := fmt.Sprintf("restart %s --force-local", container.Name)
		cmd := exec.Command("lxc", strings.Split(args, " ")...)
		err := cmd.Run()
		if err != nil {
			log.Error(err)
		}
	}
}
