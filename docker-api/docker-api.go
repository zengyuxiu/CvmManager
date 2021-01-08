package docker_api

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

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
                    cpu_usage_percent,
                    storage_write,
                    storage_read,
                    time_stamp 
                    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(InsertSQL)
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(
		totalStat.ID, totalStat.MemoryStats.Usage, totalStat.MemoryStats.MaxUsage,
		totalStat.MemoryStats.Limit, totalStat.CPUStats.CPUUsage.TotalUsage,
		uint64(totalStat.CPUStats.OnlineCPUs),calculateCPUPercent(totalStat), totalStat.StorageStats.ReadSizeBytes,
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
