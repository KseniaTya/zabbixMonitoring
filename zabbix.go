package zabbix1

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	//"io"
	"net"
	"strconv"
	//"time"
)

var (
	ZabServer = flag.String("zabbix", "", "https://zb.iridium-soft.com/")
	HostName  = flag.String("host", "", "Zabbix server")
	PORT      = 10050
	ZabHeader = "ZBXD\x01"
)

type Metrics struct {
	Host  string
	Key   string
	Value uint64 // float64
}

// Функция для отправки метрик
func SendMetrics(server string, metrics []Metrics) error {
	conn, err := net.Dial("tcp", server+":"+strconv.Itoa(PORT))
	if err != nil {
		return err
	}
	defer conn.Close()

	//buffer := &bytes.Buffer{}

	//var seeker io.Seeker = buffer
	var buffer bytes.Buffer
	buffer.WriteString(ZabHeader)

	// Версия протокола
	binary.Write(&buffer, binary.LittleEndian, uint64(1))

	// Заглушка для длины данных
	binary.Write(&buffer, binary.LittleEndian, uint64(0))

	// Формируем JSON
	jsonData := fmt.Sprintf(`{"request":"sender data","data":[`)
	for i, metric := range metrics {
		jsonData += fmt.Sprintf(`{"host":"%s","key":"%s","value":"%f"}`,
			metric.Host, metric.Key, metric.Value)
		if i < len(metrics)-1 {
			jsonData += ","
		}
	}
	jsonData += "]}"

	buffer.WriteString(jsonData)

	// Возвращаемся и записываем правильную длину данных
	dataLength := uint64(buffer.Len() - 13)
	//buffer.(9, 0)
	binary.Write(&buffer, binary.LittleEndian, dataLength)

	// Отправляем данные
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	// Читаем ответ
	reader := bufio.NewReader(conn)
	header := make([]byte, 13)
	reader.Read(header)

	// Проверяем заголовок
	if string(header[:5]) != ZabHeader {
		return fmt.Errorf("неверный ответ от сервера")
	}

	// Читаем версию и длину ответа
	//version := binary.LittleEndian.Uint64(header[5:13])
	responseLength := binary.LittleEndian.Uint64(header[5:13])

	// Читаем JSON ответ
	response := make([]byte, responseLength)
	reader.Read(response)
	fmt.Printf("Ответ сервера: %s\n", response)

	return nil
}

/*

func connectZab() {

		zabServer := flag.String("zabbix", "", "https://zb.iridium-soft.com/")
		hostName := flag.String("host", "", "Zabbix server")

		client :=


}

func sendMetrics(server string, metrics []Metrics) error {
	// Создаем соединение с Zabbix сервером
	conn, err := net.Dial("dns", server+":"+strconv.Itoa(PORT))
	if err != nil {
		return err
	}
	defer conn.Close()

	// Формируем пакет
	var buffer bytes.Buffer
	buffer.Write([]byte(zabHeader))

	// Записываем версию протокола (версия 1)
	binary.Write(&buffer, binary.LittleEndian, uint64(1))

	// Записываем длину данных
	dataLength := uint64(0)
	binary.Write(&buffer, binary.LittleEndian, dataLength)

	// Формируем JSON с метриками
	jsonData := fmt.Sprintf(`{"request":"sender data","data":[`)
	for i, metric := range metrics {
		jsonData += fmt.Sprintf(`{"host":"%s","key":"%s","value":"%s"}`,
			metric.Host, metric.Key, metric.Value)
		if i < len(metrics)-1 {
			jsonData += ","
		}
	}
	jsonData += "]}"

	// Записываем JSON в буфер
	buffer.WriteString(jsonData)

	// Возвращаемся и записываем правильную длину данных
	dataLength = uint64(buffer.Len() - 13) // 13 = len(ZBXD_HEADER) + 2*uint64
	buffer.Seek(9, 0)
	binary.Write(&buffer, binary.LittleEndian, dataLength)

	// Отправляем данные
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	// Читаем ответ
	reader := bufio.NewReader(conn)
	header := make([]byte, 13)
	reader.Read(header)

	// Проверяем заголовок
	if string(header[:5]) != zabHeader {
		return fmt.Errorf("неверный ответ от сервера")
	}

	// Читаем версию и длину ответа
	version := binary.LittleEndian.Uint64(header[5:13])
	responseLength := binary.LittleEndian.Uint64(header[5:13])

	// Читаем JSON ответ
	response := make([]byte, responseLength)
	reader.Read(response)
	fmt.Printf("Ответ сервера: %s\n", response)

	return nil
}
*/
