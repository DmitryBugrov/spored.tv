package main

import (
	//	"bufio"
	"fmt"
	"strings"
	//	"math/rand"
	"net/http"
	//	"os"
	"regexp"
	"time"
	//	"regexp/syntax"
	"encoding/xml"
	"os"

	"spored.tv/siteparser"
)

var (
	Client *http.Client
	Domain string
	err    error
	//	channelList  []Channel
	programmList []Program
	tv           TV

//	blacklist []string
)

type TV struct {
	XMLName     xml.Name  `xml:"tv"`
	Generator   string    `xml:"generator-info-name,attr"`
	ChannelList []Channel `xml:"channel"`
	Created     string    `xml:"created-by"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Id      string   `xml:"id,attr"`

	DisplayName string `xml:"display-name"`
	Lang        string `xml:"display-name>lang,attr"`

	Url string `xml:"url"`
}

//type DisplayNameStruct struct {
//	DisplayName string `xml:"display-name"`
//	Lang        string `xml:"lang,attr"`
//}

type Program struct {
	channel string
	start   string
	stop    string
	title   string
	desc    string
}

func main() {
	fmt.Println(time.Now(), "Starting...")

	firstPage := "http://www.spored.tv/"
	Domain = "http://www.spored.tv"
	//create http client with cooks
	jar := SiteParser.NewJar()
	Client = &http.Client{Jar: jar}

	//get first page
	page := SiteParser.GetPage(Client, firstPage)

	//get list of chanel from left block of site
	tv.ChannelList = GetChannelList(page)
	tv.Generator = "Alternet"
	tv.Created = firstPage
	for currentChanell := 0; currentChanell < 1; currentChanell++ { //len(channelList); i++ {
		//get channel page
		page := SiteParser.GetPage(Client, tv.ChannelList[currentChanell].Url)
		channel := string(GetStationHeader(page))
		//get urls for days
		daysURL := GetDaysURL(page)
		for currentDay := 0; currentDay < len(daysURL); currentDay++ {
			var pr Program

			page := SiteParser.GetPage(Client, daysURL[currentDay])
			day := GetDaySelected(page)

			//get container with list of programs
			block := SiteParser.FindTegBlockByParam(page, []byte("id"), []byte("ScheduleItemsContainer"))

			//get list of programs
			items := SiteParser.FindTegBlocksByParam(block, []byte("class"), []byte("ScheduleItem"))
			for currentItem := 0; currentItem < len(items); currentItem++ {
				pr.desc = ""
				pr.title = ""
				pr.channel = channel

				//parsing program title
				pr.title = GetTitle(items[currentItem])
				if pr.title == "" {
					fmt.Println("Error parsing title:", string(items[currentItem]))
				}
				//parsing program description
				pr.desc = GetDescription(items[currentItem])
				if pr.desc == "" {
					fmt.Println("Error parsing description:", string(items[currentItem]))
				}

				//parsing start time
				time := GetStartTime(items[currentItem])
				if time == "" {
					fmt.Println("Error parsing time:", string(items[currentItem]))
				} else {
					pr.start = day + time + "00 +0200"
				}

				//parsing stop time
				if currentItem < len(items)-1 {
					time := GetStartTime(items[currentItem+1])
					if time == "" {
						fmt.Println("Error parsing stop time :", string(items[currentItem+1]))
					} else {
						pr.stop = day + time + "00 +0200"
					}
				}
				programmList = append(programmList, pr)
				//	fmt.Println(pr.title, "\t\t", pr.start, "\t", pr.stop, "\t\t", pr.desc)
			}
		}
	}
	fmt.Println(tv)
	output, err := xml.MarshalIndent(tv, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	os.Stdout.Write(output)

}

func GetStationHeader(page []byte) []byte {
	block := SiteParser.FindTegBlockByParam(page, []byte("id"), []byte("StationHeader"))
	stationHeader := SiteParser.GetBlocks(block, []byte("<h1>"), []byte("</h1>"))
	return stationHeader[0]
}

func GetChannelList(page []byte) []Channel {
	block := SiteParser.FindTegBlockByParam(page, []byte("id"), []byte("MenuContainer"))
	items := SiteParser.GetBlocks(block, []byte("<a"), []byte("/a>"))
	var channelList []Channel
	for i := 0; i < len(items); i++ {
		var new_channel Channel
		new_channel.DisplayName = string(SiteParser.GetBlocks(items[i], []byte("title=\""), []byte("\""))[0])
		new_channel.Id = new_channel.DisplayName
		new_channel.Url = SiteParser.GetURL(items[i], Domain)
		new_channel.Lang = "sr"
		channelList = append(channelList, new_channel)
		//		fmt.Println(string(cl[i].url))
	}
	return channelList
}

func GetDaysURL(page []byte) []string {
	block := SiteParser.FindTegBlockByParam(page, []byte("id"), []byte("StationDays"))
	items := SiteParser.GetBlocks(block, []byte("<a"), []byte("/a>"))
	var urls []string
	for i := 0; i < len(items); i++ {
		urls = append(urls, SiteParser.GetURL(items[i], Domain))
	}
	return urls
}

func GetDaySelected(page []byte) string {
	block := SiteParser.FindTegBlockByParam(page, []byte("id"), []byte("DaySelected"))
	urlWithDate := SiteParser.GetURL(block, Domain)
	//	fmt.Println(urlWithDate)
	re := regexp.MustCompile(".*([0-9][0-9])-([0-9][0-9])-([0-9][0-9][0-9][0-9])")
	dateArray := re.FindStringSubmatch(urlWithDate)
	//	fmt.Println(dateArray)
	date := ""
	if len(dateArray) == 4 {
		date = dateArray[3] + dateArray[2] + dateArray[1]
	}
	return date

}

func GetTitle(page []byte) string {
	var result = ""
	blockForTitle := SiteParser.FindTegBlockByParam(page, []byte("id"), []byte("ProgramDescriptionLink"))
	title := SiteParser.GetBlocks(blockForTitle, []byte(">"), []byte("<"))
	//if title don't exist
	if len(title) > 0 {
		result = string(title[0])
	} else {

		title = SiteParser.GetBlocks(page, []byte("/span>"), []byte("<"))
		if len(title) > 0 {
			result = strings.TrimSpace(string(title[0]))
		} else {
			fmt.Println("error parsing title:", string(page), "=")
		}
	}
	return result

}

func GetDescription(page []byte) string {
	var result = ""
	blockForDescription := SiteParser.FindTegBlockByParam(page, []byte("id"), []byte("DivProgramDescription_"))
	description := SiteParser.GetBlocks(blockForDescription, []byte(">"), []byte("</div>"))
	if len(description) > 0 {
		result = string(description[0])
	}
	return result
}

func GetStartTime(page []byte) string {
	var result = ""
	block := SiteParser.GetBlocks(page, []byte("<span"), []byte("</span>"))
	if len(block) == 1 {
		re := regexp.MustCompile(".*([0-9][0-9]):([0-9][0-9])")
		timeArray := re.FindStringSubmatch(string(block[0]))
		if len(timeArray) == 3 {
			result = timeArray[1] + timeArray[2]
		}
	}
	return result
}
