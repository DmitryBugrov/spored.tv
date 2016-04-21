package main

import (
	//	"bufio"
	"fmt"
	"strings"
	//	"math/rand"
	"net/http"
	//	"os"
	"time"

	"spored.tv/siteparser"
)

var (
	Client       *http.Client
	Domain       string
	err          error
	channelList  []Channel
	programmList []Program

//	blacklist []string
)

type Channel struct {
	id          string
	displayName string
	url         string
}

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
	channelList := GetChannelList(page)
	for currentChanell := 0; currentChanell < 1; currentChanell++ { //len(channelList); i++ {
		//get channel page
		page := SiteParser.GetPage(Client, channelList[currentChanell].url)
		//get urls for days
		daysURL := GetDaysURL(page)
		for currentDay := 0; currentDay < len(daysURL); currentDay++ {
			var pr Program
			page := SiteParser.GetPage(Client, daysURL[currentDay])
			pr.channel = string(GetStationHeader(page))
			block := SiteParser.FindTegBlockByParam(page, []byte("id"), []byte("ScheduleItemsContainer"))
			//items := SiteParser.GetBlocks(block, []byte("class=\"ScheduleItem"), []byte("/div>"))
			items := SiteParser.FindTegBlocksByParam(items[currentItem], []byte("class"), []byte("ScheduleItem"))
			for currentItem := 0; currentItem < len(items[]); currentItem++ {
				blockForTitle := SiteParser.FindTegBlockByParam(items[currentItem], []byte("class"), []byte("ProgramDescriptionLink"))

				//parsing program title
				title := SiteParser.GetBlocks(blockForTitle, []byte(">"), []byte("<"))
				//if title don't exist
				if len(title) > 0 {
					pr.title = string(title[0])
				} else {

					title = SiteParser.GetBlocks(items[currentItem], []byte("/span>"), []byte("<"))
					if len(title) > 0 {
						pr.desc = strings.TrimSpace(string(title[0]))
					} else {
						fmt.Println("error parsing title:", string(blockForTitle))
					}
				}

				//parsing program description
				blockForDescription := SiteParser.FindTegBlockByParam(items[currentItem], []byte("id"), []byte("DivProgramDescription_"))
				description := SiteParser.GetBlocks(blockForDescription, []byte("\""), []byte("\""))
				if len(description) > 0 {
					pr.desc = string(description[0])
				} else {
					pr.desc = ""
				}

				fmt.Println(pr.title, "\t\t\t", pr.desc)
			}
		}
	}
	//	fmt.Println(string(GetStationHeader(page)))
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
		new_channel.displayName = string(SiteParser.GetBlocks(items[i], []byte("title=\""), []byte("\""))[0])
		new_channel.id = new_channel.displayName
		new_channel.url = SiteParser.GetURL(items[i], Domain)
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
