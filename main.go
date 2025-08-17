package main

import (
	"context"
	//	"fmt"
	"log/slog"
	//	"time"
	"withoutZab/internal/collector"
	"withoutZab/internal/zabbix1"
)

func main() {
	var prev []collector.CPUCoreUsage
	var cur []collector.CPUCoreUsage

	// Создание сессии
	cl := collector.Collector{}
	logger := slog.Logger{}
	zabbix1.NewClient("http://zb.iridium-soft.com", "root", "9r3xi7Pp2zkJ", 10, &logger)
	c := zabbix1.Client{}
	ctx := context.Context(context.Background())
	err := c.Initialize(ctx, "Zabbix server")
	if err != nil {
		slog.Error(err.Error())
	}
	for {
		prev = collector.CPUTime()
		cur = collector.CPUTime()
		metrics := collector.MetricSet{}
		metrics, err := cl.Collect(ctx, prev, cur)
		if err != nil {
			slog.Error(err.Error())
		}
		err = c.SendMetrics(ctx, &metrics)
		if err != nil {
			slog.Error(err.Error())
		}
		cur = prev
		prev = collector.CPUTime()
	}

}

/*
import (
	//"flag"
	"fmt"
	"log/slog"
	"time"
	"withoutZab/internal/collector"
	"withoutZab/internal/zabbix1"
)

var logger *slog.Logger

var (
	//ZabServer = flag.String("zabbix", "", "https://zb.iridium-soft.com/")
	//HostName  = flag.String("host", "", "Zabbix server")
	//ZabServer = "zb.iridium-soft.com"
	ZabServer = "127.0.0.1"
	HostName  = "Zabbix server"
	PORT      = 10050
	ZabHeader = "ZBXD\x01"
)

func main() {
	//	var prev []collector.CPUCoreUsage
	//	var cur []collector.CPUCoreUsage
	//go func() {
	//	prev = collector.CPUTime()
	//	time.Sleep(1 * time.Second)
	//}()

	//	go collector.MemUsage()
	//	cur = collector.CPUTime()
	//	fmt.Println("")
	//	go collector.GetCPU(prev, cur)
	//	cpu := collector.GetCPU(prev, cur)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		var prev []collector.CPUCoreUsage
		var cur []collector.CPUCoreUsage
		prev = collector.CPUTime()
		cur = collector.CPUTime()
		fmt.Println("")
		go collector.GetCPU(prev, cur)
		//cpu := collector.GetCPU(prev, cur)
		times := collector.CPUTime()
		userT := times[0].UserTime

		select {
		case <-ticker.C:
			//usage := cpu.CoresUsage[1].Usage
			metrics := []zabbix1.Metrics{
				{
					Host:  HostName,
					Key:   "system.cpu.util[,user]",
					Value: userT,
				},
			}
			err := zabbix1.SendMetrics(ZabServer, metrics)
			if err != nil {
				slog.Error(err.Error())
			}

		}
		prev = cur
	}
}

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
	var logger slog.Logger
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		logger.Error(err.Error())
	}
	defer file.Close()

	type Core map[string]any
	var cores []Core
	core := make(Core) // map Core

	scanner := bufio.NewScanner(file) // построчное чтение
	i := 0
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
			//fmt.Println(parts[1])
			a := parts[1]
			//fmt.Println(i, key, a)
			if key == "cpu cores" {
				fmt.Println(i, key, a)
			} else if key == "model name" {
				fmt.Println(i, key, a)
			} else if key == "core id" {
				fmt.Println(i, key, a)
			} else if key == "cpu MHz" {
				fmt.Println(i, key, a)
			}

			nval, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 32)
			if err != nil {
				core[key] = strings.TrimSpace(parts[1])
			} else {
				core[key] = nval
			}
		}
		i++
	}
	// fmt.Println(cores, len(cores))
	if len(cores) > 0 {
		cores = append(cores, core)
	}
*/

/*
	go collector.CPUUse()

	go collector.CPUUse()
	for {
		go func() {
			time.Sleep(10 * time.Second)
		}()
	}

	   	stat := clocks{
	   		maxClock: 0,
	   		avgClock: 0,
	   	}
	   	for {
	   		stat.lsClock = []int{}
	   		cpu(&stat)
	   		cpuPrint(&stat)
	   		time.Sleep(time.Millisecond * 500)
	   	}
	   }

	   // Given a file and a string to look for, find all lines
	   // in the file containing the search string as a substring
	   // Returns a pointer to the array of strings and an error field
	   func badGrep(file string, exp string) (*[]string, error) {
	   	var x []string
	   	f, err := os.Open(file)
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	defer f.Close()
	   	scanner := bufio.NewScanner(f)
	   	for scanner.Scan() {
	   		line := scanner.Text()
	   		if strings.Contains(line, exp) {
	   			x = append(x, line)
	   		}
	   	}
	   	return &x, nil
	   }

	   // Process data obtained through badGrep
	   func calcMHz(arr *[]string, stat *clocks) {
	   	s := 0
	   	for _, v := range *arr {
	   		curr := strings.Split(v, ": ")[1]
	   		curr = strings.Split(curr, ".")[0]
	   		i, err := strconv.Atoi(curr)
	   		if err != nil {
	   			log.Fatal(err)
	   		}
	   		if i >= stat.maxClock {
	   			stat.maxClock = i
	   		}
	   		s += i
	   		stat.lsClock = append(stat.lsClock, i)
	   	}
	   	stat.avgClock = s / len(stat.lsClock)
	   }

	   // Read Clock information from /proc/cpuinfo
	   // and have calcMHZ() parse it.
	   func cpu(stat *clocks) {
	   	var (
	   		file string = "/proc/cpuinfo"
	   		exp  string = "cpu MHz"
	   	)
	   	arr, err := badGrep(file, exp)
	   	if err != nil {
	   		panic("something went wrong when reading /proc/cpuinfo")
	   	}
	   	calcMHz(arr, stat)
	   }

	   // clear terminal (linux only)
	   func callClear() {
	   	cmd := exec.Command("clear")
	   	cmd.Stdout = os.Stdout
	   	cmd.Run()
	   }

	   // Print info to stdout
	   func cpuPrint(stat *clocks) {
	   	callClear()
	   	println("Max Clock: ", stat.maxClock)
	   	println("Avg Clock: ", stat.avgClock)
	   	println("\nClocks:")
	   	for _, v := range stat.lsClock {
	   		println(v)
	   	}
*/
