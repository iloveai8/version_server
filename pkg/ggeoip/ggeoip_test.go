package ggeoip

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"game_slots_vsn/pkg/utils"
	"gitlab.ftsview.com/fotoable-go/ggeoip"
	"os"
	"sort"
	"testing"
)

const path = "E:\\go\\src\\game_slots_vsn\\pkg\\ggeoip\\data\\"
const name = "ggeoip.mmdb"

func init() {
	ctx := context.TODO()
	config := &ggeoip.GeoIpConfig{
		Ctx:                ctx,
		FileName:           "ggeoip.mmdb",
		Path:               "E:\\go\\src\\game_slots_vsn\\pkg\\ggeoip\\data\\",
		Scope:              0,
		UpdateIntervalHour: 1,
		Fun:                isSuccess,
	}
	_, err := os.Stat(path + name)
	if err == nil {
		fmt.Println("load file:", path+name)
		ggeoip.LoadLocalFile(ctx, config, false)
	} else {
		util := ggeoip.NewGeoIpUtil(
			ctx,
			config,
		)
		fmt.Println("download file:", path+name)
		util.GeoIPInit()
	}
}

func TestGetIP(t *testing.T) {
	fmt.Println(ggeoip.GetCountryAndCityByIP("1.202.246.19"))
}

func TestGip_GetCountryByIP(t *testing.T) {
	type Person struct {
		ArchiveID string `json:"archive_id"`
		Ip        string `json:"Ip"`
		Country   string `json:"Country"`
	}

	// Open the CSV rFile
	rFile, err := os.Open("E:\\go\\src\\game_slots_vsn\\pkg\\ggeoip\\data\\Ip.csv")
	if err != nil {
		fmt.Println("Error opening rFile:", err)
		return
	}
	defer rFile.Close()

	reader := csv.NewReader(rFile)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}
	//cAndOList := make([]int64, 10)
	cnCounter := 0
	otherCounter := 0
	ips := make(map[string]int32)
	ipIDList := make(map[string][]string, 10)
	otherIDList := make([]string, 10)
	people := make(map[string]Person)
	person := Person{}
	for _, row := range records {
		archiveID := row[0]
		ip := row[1]
		//fmt.Printf("archiveID:%s ip:%s", archiveID, ip)
		//fmt.Println()
		country := ggeoip.GetCountryByIP(ip)
		if country == "CN" {
			person.ArchiveID = archiveID
			person.Ip = ip
			person.Country = country

			people[archiveID] = person

			cnCounter++

			ips[ip]++
			if utils.MatchIp(ip) {
				ipIDList[ip] = append(ipIDList[ip], archiveID)
			} else {
				otherIDList = append(otherIDList, archiveID)
			}
			fmt.Printf("archiveID:%s ip:%s\t\n", archiveID, ip)
		} else {
			otherCounter++
		}
	}

	fmt.Printf("cnCounter:%d oCount:%d\t\n", cnCounter, otherCounter)

	for ipS, v := range ips {
		fmt.Printf("ip:%s count:%d \t\n", ipS, v)
	}

	for ipS, v := range ipIDList {
		sort.Strings(v) // sort the slice in ascending order
		v1 := deduplicate(v)
		fmt.Printf("ip:%s size:%d archive_id_list:%v \t\n", ipS, len(v1), v1)
	}

	sort.Strings(otherIDList) // sort the slice in ascending order
	otherIDList1 := deduplicate(otherIDList)
	fmt.Printf("size:%d archive_id_list:%v \t\n", len(otherIDList1), otherIDList1)

	// Create a new JSON rFile
	wFile, err1 := os.Create("E:\\go\\src\\game_slots_vsn\\pkg\\ggeoip\\data\\cn.json")
	if err1 != nil {
		fmt.Println("Error creating rFile:", err)
		return
	}
	defer wFile.Close()
	encoder := json.NewEncoder(wFile)
	err = encoder.Encode(people)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}

}
func deduplicate(slice []string) []string {
	uniqueMap := make(map[string]bool)
	dedupedSlice := make([]string, 0)
	for _, item := range slice {
		if !uniqueMap[item] {
			uniqueMap[item] = true
			dedupedSlice = append(dedupedSlice, item)
		}
	}
	return dedupedSlice
}
