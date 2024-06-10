/*
RoWhoIs item to package converter

Made by github.com/aut-mn
*/
package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var item = flag.Int("i", 0, "Specify the item")
var verboseMode = flag.Bool("v", false, "Enable verbose mode")
var filepath = flag.String("o", ".", "Specify the file path")
var legacy = flag.Bool("l", false, "Convert items to use legacy rbxm types")

var ItemMap = map[int]string{
	2:  "TShirt",
	8:  "Hat",
	11: "Shirt",
	12: "Pants",
	18: "Face",
	41: "Hat",
	42: "Hat",
	43: "Hat",
	44: "Hat",
	45: "Hat",
	46: "Hat",
	47: "Hat",
}

type ThumbnailResponse struct {
	Data []Thumbnail `json:"data"`
}

type Thumbnail struct {
	TargetID int    `json:"targetId"`
	State    string `json:"state"`
	ImageURL string `json:"imageUrl"`
	Version  string `json:"version"`
}

type Creator struct {
	Id               int64  `json:"Id"`
	Name             string `json:"Name"`
	CreatorType      string `json:"CreatorType"`
	CreatorTargetId  int64  `json:"CreatorTargetId"`
	HasVerifiedBadge bool   `json:"HasVerifiedBadge"`
}

type ItemDetails struct {
	TargetId                  int64    `json:"TargetId"`
	ProductType               string   `json:"ProductType"`
	AssetId                   int64    `json:"AssetId"`
	ProductId                 int64    `json:"ProductId"`
	Name                      string   `json:"Name"`
	Description               string   `json:"Description"`
	AssetTypeId               int      `json:"AssetTypeId"`
	Creator                   Creator  `json:"Creator"`
	IconImageAssetId          int      `json:"IconImageAssetId"`
	Created                   string   `json:"Created"`
	Updated                   string   `json:"Updated"`
	PriceInRobux              int      `json:"PriceInRobux"`
	PriceInTickets            *int     `json:"PriceInTickets"`
	Sales                     int      `json:"Sales"`
	IsNew                     bool     `json:"IsNew"`
	IsForSale                 bool     `json:"IsForSale"`
	IsPublicDomain            bool     `json:"IsPublicDomain"`
	IsLimited                 bool     `json:"IsLimited"`
	IsLimitedUnique           bool     `json:"IsLimitedUnique"`
	Remaining                 *int     `json:"Remaining"`
	MinimumMembershipLevel    int      `json:"MinimumMembershipLevel"`
	ContentRatingTypeId       int      `json:"ContentRatingTypeId"`
	SaleAvailabilityLocations []string `json:"SaleAvailabilityLocations"`
	SaleLocation              *string  `json:"SaleLocation"`
	CollectibleItemId         *int     `json:"CollectibleItemId"`
	CollectibleProductId      *int     `json:"CollectibleProductId"`
	CollectiblesItemDetails   *string  `json:"CollectiblesItemDetails"`
}

func addToZip(zipWriter *zip.Writer, filename string, data []byte) {
	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		slog.Log(context.TODO(), 8, err.Error())
		os.Exit(1)
	}
	_, err = fileWriter.Write(data)
	if err != nil {
		slog.Log(context.TODO(), 8, err.Error())
		os.Exit(1)
	}
}

func client(method string, url string, jsonInput interface{}) ([]byte, error) {
	if *verboseMode {
		slog.Log(context.TODO(), 0, fmt.Sprintf("%s -> %s", method, url))
	}
	httpClient := &http.Client{}
	var reqBody io.Reader
	if jsonInput != nil {
		jsonBytes, err := json.Marshal(jsonInput)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	if jsonInput != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func extractID(s string, regex string) string {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(s)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}

func main() {
	flag.Parse()
	if *item == 0 {
		slog.Log(context.TODO(), 8, "you must specify an item using the -i flag")
		os.Exit(1)
	}
	_, err := os.Stat(*filepath)
	if err != nil {
		slog.Log(context.TODO(), 8, "the file path specified does not exist")
		os.Exit(1)
	}
	itemDetails, err := client("GET", fmt.Sprintf("https://economy.roblox.com/v2/assets/%d/details", *item), nil)
	if err != nil {
		slog.Log(context.TODO(), 8, err.Error())
		os.Exit(1)
	}
	var details ItemDetails
	err = json.Unmarshal(itemDetails, &details)
	if err != nil {
		slog.Log(context.TODO(), 8, err.Error())
	}
	if _, ok := ItemMap[details.AssetTypeId]; !ok {
		slog.Log(context.TODO(), 8, "the item is not supported")
		os.Exit(1)
	}
	if *verboseMode {
		slog.Log(context.TODO(), 0, "reevaluating item name for file use...")
	}
	itemType, _ := ItemMap[details.AssetTypeId]
	itemName := strings.Replace(details.Name, " ", "_", -1)
	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	if err != nil {
		slog.Log(context.TODO(), 8, err.Error())
		os.Exit(1)
	}
	itemName = reg.ReplaceAllString(itemName, "")
	if *verboseMode {
		slog.Log(context.TODO(), 0, "fetching base item details...")
	}
	rbxmReturn, err := client("GET", fmt.Sprintf("https://assetdelivery.roblox.com/v1/asset?id=%d", details.AssetId), nil)
	if err != nil {
		slog.Log(context.TODO(), 8, err.Error())
		os.Exit(1)
	}
	rbxmReturnStr := string(rbxmReturn)
	if !strings.HasPrefix(rbxmReturnStr, "<roblox") {
		slog.Log(context.TODO(), 8, "this item is not a valid rbxm file- is it too new?")
		os.Exit(1)
	}
	var fetchedMesh, fetchedTexture, formattedRbxm, fetchedThumbnail []byte
	switch itemType {
	case "Hat":
		if *verboseMode {
			slog.Log(context.TODO(), 0, "parsing retrieved xml...")
		}
		meshID := extractID(rbxmReturnStr, `<Content name="MeshId"><url>http://www.roblox.com/asset/\?id=(\d+)</url></Content>`)
		textureID := extractID(rbxmReturnStr, `<Content name="TextureId"><url>http://www.roblox.com/asset/\?id=(\d+)</url></Content>`)
		if *legacy {
			if *verboseMode {
				slog.Log(context.TODO(), 0, "converting to legacy rbxm...")
			}
			rbxmReturnStr = strings.Replace(rbxmReturnStr, "Accessory", "Hat", 1)
			re := regexp.MustCompile(`<Item class="Hat" referent="[^"]*">`)
			rbxmReturnStr = re.ReplaceAllString(rbxmReturnStr, `<Item class="Hat">`)
		}
		formattedRbxm = []byte(rbxmReturnStr)
		// TODO: Route through channels to avoid blocking
		if *verboseMode {
			slog.Log(context.TODO(), 0, "fetching mesh, texture, and thumbnail...")
		}
		fetchedMesh, _ = client("GET", fmt.Sprintf("https://assetdelivery.roblox.com/v1/asset?id=%s", meshID), nil)
		fetchedTexture, _ = client("GET", fmt.Sprintf("https://assetdelivery.roblox.com/v1/asset?id=%s", textureID), nil)
		fetchedThumbnail, _ = client("GET", fmt.Sprintf("https://thumbnails.roblox.com/v1/assets?assetIds=%d&returnPolicy=PlaceHolder&size=250x250&format=Png&isCircular=false", details.AssetId), nil)
		if *verboseMode {
			slog.Log(context.TODO(), 0, "done! mapping to variables")
		}
		var thumbnailResponse ThumbnailResponse
		if json.Unmarshal(fetchedThumbnail, &thumbnailResponse) != nil {
			slog.Log(context.TODO(), 8, err.Error())
		}
		if len(thumbnailResponse.Data) > 0 {
			thumbnail := thumbnailResponse.Data[0]
			fetchedThumbnail, err = client("GET", thumbnail.ImageURL, nil)
		} else {
			slog.Log(context.TODO(), 6, "no thumbnail found")
		}
	case "TShirt":
		if *verboseMode {
			slog.Log(context.TODO(), 0, "fetching texture...")
		}
		fetchedThumbnail, err = client("GET", fmt.Sprintf("https://assetdelivery.roblox.com/v1/asset?id=%d", details.AssetId), nil)
		if err != nil {
			slog.Log(context.TODO(), 8, err.Error())
			os.Exit(1)
		}
		formattedRbxm = []byte(rbxmReturnStr)
	case "Shirt", "Pants":
		if *verboseMode {
			slog.Log(context.TODO(), 0, "fetching texture and thumbnail...")
		}
		fetchedTexture, err = client("GET", fmt.Sprintf("https://assetdelivery.roblox.com/v1/asset?id=%d", details.AssetId), nil)
		if err != nil {
			slog.Log(context.TODO(), 8, err.Error())
			os.Exit(1)
		}
		fetchedThumbnail, _ = client("GET", fmt.Sprintf("https://thumbnails.roblox.com/v1/assets?assetIds=%d&returnPolicy=PlaceHolder&size=250x250&format=Png&isCircular=false", details.AssetId), nil)
		var thumbnailResponse ThumbnailResponse
		if json.Unmarshal(fetchedThumbnail, &thumbnailResponse) != nil {
			slog.Log(context.TODO(), 8, err.Error())
		}
		if len(thumbnailResponse.Data) > 0 {
			thumbnail := thumbnailResponse.Data[0]
			fetchedThumbnail, _ = client("GET", thumbnail.ImageURL, nil)
		} else {
			slog.Log(context.TODO(), 6, "no thumbnail found")
		}
	}
	if *verboseMode {
		slog.Log(context.TODO(), 0, fmt.Sprintf("creating writer and writing to %s/%s.zip", *filepath, itemName))
	}
	zipFile, err := os.Create(fmt.Sprintf("%s/%s.zip", *filepath, itemName))
	if err != nil {
		slog.Log(context.TODO(), 8, err.Error())
		os.Exit(1)
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	if itemType != "TShirt" {
		addToZip(zipWriter, fmt.Sprintf("%s_desc.txt", itemName), []byte(details.Description))
		addToZip(zipWriter, fmt.Sprintf("textures/%s.png", itemName), fetchedTexture)
	}
	if itemType == "Hat" {
		addToZip(zipWriter, fmt.Sprintf("fonts/%s.mesh", itemName), fetchedMesh)
	}
	addToZip(zipWriter, fmt.Sprintf("%s.png", itemName), fetchedThumbnail)
	addToZip(zipWriter, fmt.Sprintf("%s.rbxm", itemName), []byte(formattedRbxm))
	slog.Log(context.TODO(), 0, "done!")
}
