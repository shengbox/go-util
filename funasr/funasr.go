package funasr

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	SendChunkSize = 960                        // size of data to send in bytes
	BatchSize     = 120                        // number of messages to send at once
	Hotwords      = ""                         // hotwords (optional)
	Host          = "audio.shuzigongsheng.com" // server IP
)

type TextSeg struct {
	StampSents []struct {
		Start   int64  `json:"start"`
		End     int64  `json:"end"`
		Punc    string `json:"punc"`
		TextSeg string `json:"text_seg"`
	} `json:"stamp_sents"`
	Text string `json:"text"`
}

type Message struct {
	Mode       string `json:"mode"`
	ChunkSize  []int  `json:"chunk_size"`
	WavName    string `json:"wav_name"`
	Hotwords   string `json:"hotwords"`
	WavFormat  string `json:"wav_format"`
	IsSpeaking bool   `json:"is_speaking"`
}

type ResultHandler func(res TextSeg)

func SpeechToText(path string, f ResultHandler) (string, error) {
	header := http.Header{}
	header.Set("Origin", fmt.Sprintf("https://%s", Host)) // 设置来源
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://%s/", Host), header)
	if err != nil {
		log.Println("连接失败:", err)
		return "", err
	}
	log.Println("连接成功!")
	defer conn.Close()

	result := ""

	wait := &sync.WaitGroup{}

	// 接收消息
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var res TextSeg
			json.Unmarshal(message, &res)
			result = result + res.Text
			if f != nil {
				f(res)
			}
			wait.Done()
		}
	}()

	err = SendAudioData(conn, wait, path)
	if err != nil {
		return "", err
	}

	log.Println("发送完数据,请等候,正在识别...")
	wait.Wait()

	return result, nil
}

func SendAudioData(conn *websocket.Conn, wait *sync.WaitGroup, path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Println("Error opening WAV file:", err)
		return err
	}
	defer file.Close()

	// Check if the file is a WAV file
	wavHeader := make([]byte, 44)
	_, err = file.Read(wavHeader)
	if err != nil {
		log.Println(err)
	}

	// Skip the WAV header (44 bytes)
	_, err = file.Seek(44, io.SeekStart)
	if err != nil {
		log.Println("Error seeking to audio data:", err)
		return err
	}

	start := `{"mode":"offline","chunk_size":[5,10,5],"wav_name":"h5","hotwords":"","wav_format":"","chunk_interval":10, "is_speaking":true}`
	end := `{"mode":"offline","chunk_size":[5,10,5],"wav_name":"h5","hotwords":"","wav_format":"","chunk_interval":10, "is_speaking":false}`
	conn.WriteMessage(websocket.TextMessage, []byte(start))

	i := 0
	// Send audio data in chunks
	buffer := make([]byte, SendChunkSize)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading audio data:", err)
			break
		}
		conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
		i = i + 1
		if i%BatchSize == 0 {
			conn.WriteMessage(websocket.TextMessage, []byte(end))
			wait.Add(1)
			wait.Wait()

			conn.WriteMessage(websocket.TextMessage, []byte(start))
		}
	}

	wait.Add(1)
	conn.WriteMessage(websocket.TextMessage, []byte(end))
	return nil
}
