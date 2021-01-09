package lxc_api

import (
	"database/sql"
	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"regexp"
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
	for i := 0; i <= InstanceNumber; i++ {
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
