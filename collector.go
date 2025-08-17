package collector

/*
import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)


import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

var logger *slog.Logger

type Core map[string]any // набор произвольных значений = map[string]interface{}

type CPUUsage struct {
	Cores []Core
}

type CPU struct {
	CoreID      string
	UserTime    uint64
	NiceTime    uint64
	SystemTime  uint64
	IdleTime    uint64
	IOWaitTime  uint64
	IRQTime     uint64
	SoftIRQTime uint64
}

func CPUUse() []Core {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		logger.Error(err.Error())
	}
	defer file.Close()

	var cores []Core
	core := make(Core) // map Core

	scanner := bufio.NewScanner(file) // построчное чтение
	for scanner.Scan() {
		line := scanner.Text() // читаем каждую строку
		if line == "" {
			// fmt.Println(core)
			if len(core) > 0 {
				cores = append(cores, core)
				core = make(Core)
			}
			continue // переход к следующей строке файла
		}
		parts := strings.SplitN(line, ":", 2) // возвращает срезы подстрок между ":" (разделённая на 2 части)
		// fmt.Println(parts)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])

			nval, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 32)
			if err != nil {
				core[key] = strings.TrimSpace(parts[1])
			} else {
				core[key] = nval
			}
		}
	}
	// fmt.Println(cores, len(cores))
	if len(cores) > 0 {
		cores = append(cores, core)
	}
	return cores
}

func CPUTime() ([]uint64, []uint64, []uint64, []uint64, []uint64, []uint64, []uint64) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		logger.Error(err.Error())
	}
	defer file.Close()

	//	var totals []uint64 // сумма всего времени, проведённого в разных состояниях
	var usert []uint64
	var nice []uint64
	var system []uint64
	var idles []uint64 // время простоя при чтении /proc/stat
	var iowait []uint64
	var irq []uint64
	var softirq []uint64

	scanner := bufio.NewScanner(file)
	var cores []Core
	core := make(Core) // map Core

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.HasPrefix(line, "cpu") {
			parts := strings.Fields(line) // разделение строки на срез подстрок по пробелам
*/

/*
	column 0: user – time spent in user mode
	column 1: nice – time spent processing nice processes in user mode
	column 2: system – time spent executing kernel code
	column 3: idle – time spent idle
	column 4: iowait – time spent waiting for I/O
	column 5: irq – time spent servicing interrupts
	column 6: softirq – time spent servicing software interrupts
	column 7: steal – time stolen from a virtual machine
	column 8: guest – time spent running a virtual CPU for a guest operating system
	column 9: guest_nice – time spent running a virtual CPU for a “niced” guest operating system
*/
/*
			for i, part := range parts {
				if
				fmt.Println(i, part)
				nval, err := strconv.ParseUint(part, 10, 64)

				if err != nil {
					logger.Error(err.Error())
				}

				if i == 0 {
					usert = append(usert, uint64(nval))
				} else if i == 1 {
					nice = append(nice, uint64(nval))
				} else if i == 2 {
					system = append(system, uint64(nval))
				} else if i == 3 {
					idles = append(idles, uint64(nval))
				} else if i == 4 {
					iowait = append(iowait, uint64(nval))
				} else if i == 5 {
					irq = append(irq, uint64(nval))
				} else if i == 6 {
					softirq = append(softirq, uint64(nval))
				}
			}
		}
	}
	// cpu := CPU{}
	return usert, nice, system, idles, iowait, irq, softirq // arrays with 3 values (cpu, cpu0, cpu1)
}

/*
func CPUUsePerCent(usert, nice, system, idles, iowait, irq, softirq []uint64) []uint64 {
	for i := 0; i < len(usert); i++ {
		prevNI := usert[i] + nice[i] + system[i] + irq[i] + softirq[i]
	}
}

func NewCPU() *CPU {
	cpu := CPU{}
	cpu.Cores = CPUUse()
	cpu.startIdles, cpu.startTotals = CPUTime()

	// fmt.Println(cpu.Cores[0]["model name"])

	if model, ok := cpu.Cores[0]["model name"].(string); ok {
		cpu.ModelName = model
	}
	for _, core := range cpu.Cores {
		if freq, ok := core["cpu MHz"].(float64); ok {
			cpu.Freq = append(cpu.Freq, freq)
		}
	}
	// fmt.Println("Iter")
	// fmt.Println(cpu.ModelName, cpu.Usage, cpu.Freq)
	return &cpu
}

func (c *CPU) Refresh() {
	for i, core := range CPUUse() {
		if freq, ok := core["cpu MHz"].(float64); ok {
			c.Freq[i] = freq
		}
	}
}

func (c *CPU) UseRefresh() {
	idles, totals := CPUTime()
	// fmt.Println(idles, totals)
	if len(c.Usage) != len(idles) {
		// fmt.Println(len(idles), len(c.Usage))
		c.Usage = make([]float64, len(idles))
	}
	for i := 0; i < len(idles); i++ {
		fmt.Println("Idles", idles[i], c.startIdles[i])
		idleDelta := idles[i] - c.startIdles[i]
		// fmt.Println(c.startIdles[i])
		fmt.Println("Totals", totals[i], c.startTotals[i])
		totalDelta := totals[i] - c.startTotals[i]
		// fmt.Println(totalDelta)
		if totalDelta == 0 {
			c.Usage[i] = 0
			continue
		}
		// fmt.Println(totalDelta, idleDelta)
		// c.Usage[i] = 100.0*1.0 - float64(idleDelta)/float64(totalDelta)
		// CPU Usage (%) = 100 * (total time - idle time) / total time
		c.Usage[i] = 100.0 * (float64(totalDelta) - float64(idleDelta)) / float64(totalDelta)
	}
	// fmt.Println(c.Usage)
	c.startIdles = idles
	c.startTotals = totals
}

func MonitoringCPU() {
	for {
		cpu := NewCPU()
		cpu.Refresh()
		//slog.Info("New CPU:", cpu)
		cpu.UseRefresh()
		//slog.Info()
	}
}
*/
// Print system resource usage every 10 seconds.
//func System() {
//for {
//cpuUsage := runtime.ReadCPUUsage()
//fmt.Printf("CPU Usage: %.2f%%\n", cpuUsage*100)
//time.Sleep(time.Second)

//}

/*
	mem := &runtime.MemStats{}

	for {
		cpu := runtime.NumCPU()
		log.Println("CPU:", cpu)

		rot := runtime.NumGoroutine()
		log.Println("Goroutine:", rot)

		// Byte
		runtime.ReadMemStats(mem)
		log.Println("Memory:", mem.Alloc)

		time.Sleep(2 * time.Second)
		log.Println("-------")
	}
}


type Collector struct {
	logger *slog.Logger
}

func NewCollector(logger *slog.Logger) *Collector {
	return &Collector{
		logger: logger,
	}
}

func (c *Collector) Collect(ctx context.Context) (*MetricSet, error) {
	c.logger.Debug("Start collecting metrics")
	metrics := &MetricSet{
		Timestamp: time.Now(),
	}

} */
