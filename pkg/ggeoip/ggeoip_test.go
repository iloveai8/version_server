package ggeoip

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gitlab.ftsview.com/fotoable-go/ggeoip"
	"os"
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
	type person struct {
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
	people := make(map[string]person)
	for _, row := range records {
		country := ggeoip.GetCountryByIP(row[1])
		if country == "CN" {
			person := person{ArchiveID: row[0], Ip: row[1], Country: country}
			people[row[0]] = person
		}
	}

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
