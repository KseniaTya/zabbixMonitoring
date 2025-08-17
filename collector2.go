package collector

import (
	"bufio"
	"context"
	"fmt"
	"time"
	//"encoding/json"
	//"time"

	//	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const maxPercent = 100.0

var logger *slog.Logger

type Core map[string]any // набор произвольных значений = map[string]interface{}

// ------------------------------------

// Collector отвечает за сбор системных метрик
type Collector struct {
	logger *slog.Logger
}

// New создает новый экземпляр сборщика метрик
func New(logger *slog.Logger) *Collector {
	return &Collector{
		logger: logger,
	}
}

// Collect собирает все системные метрики
func (c *Collector) Collect(ctx context.Context, prev, cur []CPUCoreUsage) (MetricSet, error) {
	c.logger.Debug("Starting metrics collection")

	metrics := &MetricSet{
		Timestamp: time.Now(),
	}

	// Используем каналы для параллельного сбора метрик
	type result struct {
		name string
		err  error
	}

	results := make(chan result, 2)

	// Собираем CPU метрики
	go func() {
		cpuMetrics := GetCPU(prev, cur)
		metricsCPU := CPUMetrics{}
		metricsCPU.UsagePercent = cpuMetrics
		metrics.CPU = metricsCPU
		results <- result{name: "CPU", err: nil}
	}()

	// Собираем метрики памяти
	go func() {

		total, used, free, avail, per := GetMemUsage()
		metricsMem := MemoryMetrics{}
		metricsMem.TotalBytes = total
		metricsMem.UsedBytes = used
		metricsMem.FreeBytes = free
		metricsMem.AvailableBytes = avail
		metricsMem.UsagePercent = per

		metrics.Memory = metricsMem

		results <- result{name: "Memory", err: nil}
	}()

	// Ждем завершения всех горутин
	var errors []string
	for i := 0; i < 2; i++ {
		select {
		case res := <-results:
			if res.err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", res.name, res.err))
				c.logger.Warn("Failed to collect metrics",
					slog.String("component", res.name))
				//slog.Error("Error:", res.err))
			}
		case <-ctx.Done():
			return *metrics, ctx.Err() // nil, ctx.Err()
		}
	}

	if len(errors) == 2 {
		return *metrics, fmt.Errorf("failed to collect all metrics: %v", errors) //nil, fmt.Errorf("failed to collect all metrics: %v", errors)
	}

	c.logger.Debug("Metrics collection completed",
		slog.Int("errors", len(errors)),
		slog.Time("timestamp", metrics.Timestamp))

	return *metrics, nil
}

// ---------------------------

type CPU struct {
	NumCores int
	//CoresFreq  []float64
	CoresUsage []CoreCPU
}

type CoreCPU struct {
	CoreID    string
	Freq      float64
	ModelName string
	Usage     float64
}

type CPUCoreUsage struct {
	CoreID      string
	UserTime    uint64
	NiceTime    uint64
	SystemTime  uint64
	IdleTime    uint64
	IOWaitTime  uint64
	IRQTime     uint64
	SoftIRQTime uint64
}

func GetCPU(prev, cur []CPUCoreUsage) float64 { //*CPU {
	c, num := ModelFreq()
	c1 := CPUUse(prev, cur)
	fmt.Println("Number of cores:", num)
	a := c1[0].Usage

	cpu := CPU{}
	for i, core := range c1 {
		//fmt.Println(core)
		//fmt.Println(c[i])
		if core.CoreID == "cpu" {
			fmt.Println("Total CPU usage", "(", core.CoreID, "):", core.Usage)
			a = core.Usage
		} else {
			fmt.Println("Frequency for", core.CoreID, c[i-1].Freq)
			fmt.Println("CPU usage for", core.CoreID, ":", core.Usage)
		}
	}
	cpu.NumCores = num
	cpu.CoresUsage = c1

	//fmt.Println(cpu.NumCores)
	//return &cpu
	return a
}

func ModelFreq() ([]CoreCPU, int) {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		logger.Error(err.Error())
	}
	defer file.Close()

	var cores []CoreCPU
	var numCores int
	core := CoreCPU{} // map Core

	scanner := bufio.NewScanner(file) // построчное чтение

	for scanner.Scan() {
		line := scanner.Text() // читаем каждую строку
		if line == "" {
			// fmt.Println(core)
			//if core.ModelName != "" {
			cores = append(cores, core)
			core = CoreCPU{}
			//}
			continue // переход к следующей строке файла
		}
		parts := strings.SplitN(line, ":", 2) // возвращает срезы подстрок между ":" (разделённая на 2 части)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			//fmt.Println(key, parts[1])
			if key == "core id" {
				core.CoreID = strings.TrimSpace(parts[1])
			} else if key == "model name" {
				core.ModelName = strings.TrimSpace(parts[1])
			} else if key == "cpu MHz" {
				f, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err != nil {
					logger.Error(err.Error())
				}
				core.Freq = f
			} else if key == "cpu cores" {
				numCores, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
			}
		}

	}
	cores = append(cores, core)
	return cores, numCores
}

func CPUTime() []CPUCoreUsage {
	file, err := os.Open("/proc/stat")
	if err != nil {
		logger.Error(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var cores []CPUCoreUsage

	for scanner.Scan() {
		core := CPUCoreUsage{}
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") {
			parts := strings.Fields(line) // разделение строки на срез подстрок по пробелам
			for i, part := range parts {
				switch i {
				case 0:
					core.CoreID = strings.TrimSpace(part)
				case 1:
					val, err := strconv.ParseUint(part, 10, 64)
					if err != nil {
						logger.Error(err.Error())
					} else {
						core.UserTime = val
					}
				case 2:
					val, err := strconv.ParseUint(part, 10, 64)
					if err != nil {
						logger.Error(err.Error())
					} else {
						core.NiceTime = val
					}
				case 3:
					val, err := strconv.ParseUint(part, 10, 64)
					if err != nil {
						logger.Error(err.Error())
					} else {
						core.SystemTime = val
					}
				case 4:
					val, err := strconv.ParseUint(part, 10, 64)
					if err != nil {
						logger.Error(err.Error())
					} else {
						core.IdleTime = val
					}
				case 5:
					val, err := strconv.ParseUint(part, 10, 64)
					if err != nil {
						logger.Error(err.Error())
					} else {
						core.IOWaitTime = val
					}
				case 6:
					val, err := strconv.ParseUint(part, 10, 64)
					if err != nil {
						logger.Error(err.Error())
					} else {
						core.IRQTime = val
					}
				case 7:
					val, err := strconv.ParseUint(part, 10, 64)
					if err != nil {
						logger.Error(err.Error())
					} else {
						core.SoftIRQTime = val
					}

				}
			}
			cores = append(cores, core)
		}
	}
	return cores
}

func CPUUse(prev, cur []CPUCoreUsage) []CoreCPU {
	arrCoCPU := make([]CoreCPU, len(prev))
	for i, core := range prev {
		c := CoreCPU{}
		prevIdle := core.IdleTime + core.IOWaitTime
		curIdle := cur[i].IdleTime + cur[i].IOWaitTime

		prevNonIdle := core.UserTime + core.NiceTime + core.SystemTime + core.IdleTime + core.SoftIRQTime
		curNonIdle := cur[i].UserTime + cur[i].NiceTime + cur[i].SystemTime + cur[i].IdleTime + cur[i].SoftIRQTime

		prevTotal := prevNonIdle + prevIdle
		curTotal := curNonIdle + curIdle

		deltaTotal := curTotal - prevTotal
		deltaIdle := curIdle - prevIdle

		//fmt.Println("Total", deltaTotal)
		//fmt.Println("Idle", deltaIdle)

		//usage := (float64(deltaTotal-deltaIdle) / float64(deltaTotal)) * maxPercent
		c.CoreID = core.CoreID
		c.Usage = (float64(deltaTotal-deltaIdle) / float64(deltaTotal)) * maxPercent

		arrCoCPU[i] = c
	}
	return arrCoCPU
}

/*
	column 0: cpu id
	column 1: user – time spent in user mode
	column 2: nice – time spent processing nice processes in user mode
	column 3: system – time spent executing kernel code
	column 4: idle – time spent idle
	column 5: iowait – time spent waiting for I/O
	column 6: irq – time spent servicing interrupts
	column 7: softirq – time spent servicing software interrupts
	column 8: steal – time stolen from a virtual machine
	column 9: guest – time spent running a virtual CPU for a guest operating system
	column 10: guest_nice – time spent running a virtual CPU for a “niced” guest operating system
*/
// usert, nice, system, idles, iowait, irq, softirq
// 1 2 3 4 5 6 7

func MemUsage() {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		logger.Error(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) // построчное чтение
	mem := map[string]uint64{}
	for scanner.Scan() {
		line := scanner.Text()                // читаем каждую строку
		parts := strings.SplitN(line, ":", 2) // возвращает срезы подстрок между ":" (разделённая на 2 части)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			a := strings.TrimSpace(parts[1])
			a = strings.Replace(a, "kB", "", -1)
			a = strings.TrimSpace(a)
			b, _ := strconv.ParseUint(a, 10, 64)
			//if err != nil {
			//	logger.Error(err.Error())
			//}
			mem[key] = b
		}
	}
	fmt.Println("Total memory:", mem["MemTotal"], "kByte")
	fmt.Println("Free memory:", mem["MemFree"], "kByte                  ", float64(mem["MemFree"])/float64(mem["MemTotal"]), "%")
	fmt.Println("Available memory:", mem["MemAvailable"], "kByte")
	fmt.Println("Used memory:", mem["MemTotal"]-mem["MemFree"], "kByte")

	fmt.Println(" ")

	fmt.Println("Swap total memory:", mem["SwapTotal"], "kByte")
	fmt.Println("Swap free memory:", mem["SwapFree"], "kByte")
}

func GetMemUsage() (total, used, free, available uint64, percent float64) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		logger.Error(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) // построчное чтение
	mem := map[string]uint64{}
	for scanner.Scan() {
		line := scanner.Text()                // читаем каждую строку
		parts := strings.SplitN(line, ":", 2) // возвращает срезы подстрок между ":" (разделённая на 2 части)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			a := strings.TrimSpace(parts[1])
			a = strings.Replace(a, "kB", "", -1)
			a = strings.TrimSpace(a)
			b, _ := strconv.ParseUint(a, 10, 64)
			mem[key] = b
		}
	}
	total = mem["MemTotal"]
	free = mem["MemFree"]
	percent = float64(mem["MemFree"]) / float64(mem["MemTotal"])
	available = mem["MemAvailable"]
	used = mem["MemTotal"] - mem["MemFree"]

	return total, used, free, available, percent
}
