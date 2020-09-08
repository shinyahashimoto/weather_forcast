package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Forcast struct {
	Hour []Hour `json:"hourly"`
	Day  []Day  `json:"daily"`

	TimeZone string `json:"timezone"`
}

type Hour struct {
	Time    int64     `json:"dt"`
	Temp    float32   `json:"temp"`
	Weather []Weather `json:"weather"`
}

type Day struct {
	Time int64 `json:"dt"`
	Temp struct {
		MinTemp float32 `json:"min"`
		MaxTemp float32 `json:"max"`
	} `json:"temp"`
	Weather []Weather `json:"weather"`
}

type Weather struct {
	Info   string `json:"main"`
	Detail string `json:"description"`
}

func main() {
	lambda.Start(postLineMessage)
}

func postLineMessage() {
	cityNameInfo, weatherForcastInfo, message := createWeatherForcast()

	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
	if err != nil {
		fmt.Println(err)
	}
	if _, err := bot.PushMessage(os.Getenv("USER_ID"), linebot.NewTextMessage(cityNameInfo), linebot.NewTextMessage(weatherForcastInfo), linebot.NewTextMessage(message)).Do(); err != nil {
		fmt.Println(err)
	}
}

func getWeatherForcast() *Forcast {
	//緯度
	const LATITUDE = "35.465786"
	//経度
	const LONGITUDE = "139.622313"
	base_url := "https://api.openweathermap.org/data/2.5/onecall"
	url := fmt.Sprintf("%s?lat=%s&lon=%s&exclude=current&units=metric&lang=ja&appid=%s", base_url, LATITUDE, LONGITUDE, os.Getenv("API_FORCAST_KEY"))

	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var forcast Forcast
	json.Unmarshal([]byte(body), &forcast)

	return (&forcast)
}

func createWeatherForcast() (string, string, string) {
	forcast := getWeatherForcast()
	//プログラム実行が前日なので、翌日を取得
	day := time.Now().AddDate(0, 0, 1).Day()

	hasUmbrella := 0

	var buffer bytes.Buffer
	firstInfo := true

	for _, list := range forcast.Hour {
		time := time.Unix(list.Time, 0)
		if day == time.Day() && time.Hour() > 6 {

			hasUmbrella += changeWeatherName(&list.Weather[0])
			forcastInfo := ""
			if firstInfo == true {
				// int の時は%d
				// 初めは改行したくないため
				forcastInfo = fmt.Sprintf("%d：%s,気温は%s˚C", time.Hour(), list.Weather[0].Info, fmt.Sprintf("%.1f", list.Temp))
				firstInfo = false
			} else {
				forcastInfo = fmt.Sprintf("\n%d：%s,気温は%s˚C", time.Hour(), list.Weather[0].Info, fmt.Sprintf("%.1f", list.Temp))
			}
			firstString := strings.HasPrefix(forcastInfo, "0")
			// 06時の場合０を取り除く
			if firstString {
				forcastInfo = fmt.Sprintf(forcastInfo[1:])
			}
			buffer.WriteString(forcastInfo)
		}
	}
	var message string = ""

	if hasUmbrella == 0 {
		message = "良い1日をお過ごしください！"
	} else {
		message = "雨予報なので傘を持っていきましょう！"
	}

	cityNameInfo := "横浜の天気予報です！"

	return cityNameInfo, buffer.String(), message
}

func changeWeatherName(Weather *Weather) int {
	hasUmbrella := 0

	var weatherNameList map[string]string = map[string]string{"Clear": "晴れ", "Clouds": "曇り", "Rain": "雨", "Drizzle": "霧雨", "Thunderstorm": "雷雨", "Snow": "雪"}

	for key, value := range weatherNameList {
		if Weather.Info == key {
			Weather.Info = value

			if key == "Rain" || key == "Drizzle" || key == "Thunderstorm" || key == "Snow" {
				hasUmbrella++
			}
		}
	}
	return hasUmbrella
}
