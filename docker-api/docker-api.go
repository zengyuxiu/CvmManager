package docker_api

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"regexp"
)

const NameLength = 5

func DockerStatus(sqliteDatabase *sql.DB) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error(err)
		return err
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	for _, container := range containers {
		stats, err := cli.ContainerStats(ctx, container.ID, false)
		if err != nil {
			log.Error(err)
			return err
		}
		totalStat := types.StatsJSON{}
		decode := json.NewDecoder(stats.Body)
		err = decode.Decode(&totalStat)
		if err != nil {
			log.Error(err)
			return err
		} else {
			InsertStatDocker(sqliteDatabase, &totalStat)
		}
	}
	return nil
}
func InsertStatDocker(db *sql.DB, totalStat *types.StatsJSON) {
	/*	id := stat.ID
		mmr_usg := stat.MemoryStats.Usage
		mmr_musg := stat.MemoryStats.MaxUsage
		mmr_lmt := stat.MemoryStats.Limit
		cpu_usg := stat.CPUStats.CPUUsage.TotalUsage
		cpu_nln := stat.CPUStats.OnlineCPUs
		cpu_per := calculateCPUPercent(stat)
		str_r := stat.StorageStats.ReadSizeBytes
		str_w := stat.StorageStats.WriteSizeBytes
		tmstp:= stat.Read*/
	InsertSQL := `INSERT INTO docker_stats(
                    ID,
                    memory_usage,
                    memory_max_usage,
                    memory_limit,
                    cpu_usage,
                    cpu_online,
                    storage_write,
                    storage_read,
                    time_stamp 
                    ) VALUES (?, ?, ?, ?, ?, ?,  ?, ?, ?)`
	statement, err := db.Prepare(InsertSQL)
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(
		totalStat.ID, totalStat.MemoryStats.Usage, totalStat.MemoryStats.MaxUsage,
		totalStat.MemoryStats.Limit, totalStat.CPUStats.CPUUsage.TotalUsage,
		uint64(totalStat.CPUStats.OnlineCPUs), totalStat.StorageStats.ReadSizeBytes,
		totalStat.StorageStats.WriteSizeBytes, totalStat.Read,
	)
	if err != nil {
		log.Error(err)
	}
}
func calculateCPUPercent(v *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}

func DockerCreate(InstanceNum int, ImageID string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error(err)
		return err
	}
	var TestLaber = map[string]string{
		"ContainerType": "TestContainer",
	}
	Config := container.Config{
		AttachStderr: false,
		AttachStdin:  false,
		AttachStdout: false,
		Tty:          true,
		Labels:       TestLaber,
		Cmd: strslice.StrSlice{
			"/bin/bash",
		},
		Image: ImageID,
	}
	HostConfig := container.HostConfig{
		NetworkMode: "default",
	}
	for i := 0; i < InstanceNum; i++ {
		name := "test" + GetRandomString(NameLength)
		container_created, err := cli.ContainerCreate(ctx, &Config, &HostConfig,
			nil, nil, name)
		if err != nil {
			log.Error(err)
		}
		_, _ = cli.ContainerWait(ctx, container_created.ID, container.WaitConditionNotRunning)
		startoption := types.ContainerStartOptions{}
		err = cli.ContainerStart(ctx, container_created.ID, startoption)
		if err != nil {
			log.Error(err)
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

func DockerDelete() error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error(err)
		return err
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	for _, ctn := range containers {
		if match, _ := regexp.Match("/test(.*)", []byte(ctn.Names[0])); match == true {
			if err := cli.ContainerStop(ctx, ctn.ID, nil); err != nil {
				log.Error(err)
			}
			_, _ = cli.ContainerWait(ctx, ctn.ID, container.WaitConditionNotRunning)
			args := filters.NewArgs()
			args.Add("label", "ContainerType=TestContainer")
			_, err := cli.ContainersPrune(ctx, args)
			if err != nil {
				log.Error(err)
			}
		}
	}
	return nil
}
