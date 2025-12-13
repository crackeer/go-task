package downloadwork

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var httpClient *resty.Client

func init() {
	httpClient = resty.NewWithClient(&http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		}})
	httpClient.SetTimeout(60 * time.Second)
	httpClient.SetRetryCount(3)
	httpClient.SetRetryWaitTime(20 * time.Second)
}

func getWorkCubeSize(work []byte) []string {
	if gr := gjson.GetBytes(work, "panorama.list.0.size_list"); gr.Exists() {
		var retData []string
		for _, item := range gr.Array() {
			retData = append(retData, item.String())
		}
		return retData
	}
	return []string{"2048"}
}

func parseJSONFiles(urlList []string, baseURL string, saveDir string, sendFunc func(string)) ([]string, error) {
	retData := []string{}
	for index, item := range urlList {
		sendFunc(fmt.Sprintf("--> parse_json [%d/%d] %s", index+1, len(urlList), item))
		bytes, err := downloadFileToContent(baseURL+item, filepath.Join(saveDir, item))
		if err != nil {
			return nil, fmt.Errorf("download %s error:%s", baseURL+item, err.Error())
		}
		data := map[string]interface{}{}
		if err := json.Unmarshal(bytes, &data); err != nil {
			return nil, fmt.Errorf("json unmarshal content %s error: %s", item, err.Error())
		}
		parts := parseTilesetNames(data)
		for _, part := range parts {
			retData = append(retData, comparePath(item, part))
		}
	}
	return retData, nil
}

func splitJson(urlList []string) ([]string, []string) {
	var (
		result  []string
		result1 []string
	)
	for _, v := range urlList {
		if strings.HasSuffix(v, ".json") {
			result = append(result, v)
		} else {
			result1 = append(result1, v)
		}
	}
	return result1, result
}

func getPanoramaURLS(bytes []byte, cubeSize []string) []string {
	panoramas := gjson.GetBytes(bytes, "panorama.list").Array()
	fmt.Println("Panorama Count:", len(panoramas))
	urlList := []string{}
	for _, panorama := range panoramas {
		for _, items := range []string{
			"back", "front", "left", "right", "up", "down",
		} {
			if gr := panorama.Get(items); gr.Exists() {
				urlList = append(urlList, getURLs(gr.String(), cubeSize)...)
			}
		}
	}
	return urlList
}

func getURLs(path string, cubeSize []string) []string {
	urlList := []string{}
	for _, size := range cubeSize {
		realPath := strings.Replace(path, "cube_2048", "cube_"+size, 1)
		urlList = append(urlList, realPath)
	}
	return urlList
}

func getSetModelURLS(bytes []byte) ([]string, []byte) {
	urlList := []string{}
	if gr := gjson.GetBytes(bytes, "model.file_url"); gr.Exists() {
		urlList = append(urlList, gr.String())
	}

	materialBaseURL := gjson.GetBytes(bytes, "model.material_base_url").String()
	material_textures := gjson.GetBytes(bytes, "model.material_textures").Array()
	for _, item := range material_textures {
		urlList = append(urlList, materialBaseURL+item.String())
	}

	layers := gjson.GetBytes(bytes, "model.layers").Array()
	for index, item := range layers {
		if gr := item.Get("tileset_url"); gr.Exists() {
			relativePath := getRelativePath(gr.String())
			if len(relativePath) > 0 {
				urlList = append(urlList, relativePath)
				bytes, _ = sjson.SetBytes(bytes, fmt.Sprintf("model.layers.%d.tileset_url", index), relativePath)
			}
		}
	}

	b3mdMappingsUrl := gjson.GetBytes(bytes, "model.tiles.b3md_mappings_url").String()
	if len(b3mdMappingsUrl) > 0 {
		if tmp := getRelativePath(b3mdMappingsUrl); len(tmp) > 0 {
			urlList = append(urlList, tmp)
		}
	}
	tilesetURL := gjson.GetBytes(bytes, "model.tiles.tileset_url").String()
	if len(tilesetURL) > 0 {
		if tmp := getRelativePath(tilesetURL); len(tmp) > 0 {
			urlList = append(urlList, tmp)
			bytes, _ = sjson.SetBytes(bytes, "model.tiles.tileset_url", tmp)
		}
	}
	return urlList, bytes
}

func getRelativePath(path string) string {

	seperators := []string{"mesh", "point_cloud", "model", "lod", "areophoto"}
	for _, seperator := range seperators {
		if strings.HasPrefix(path, seperator+"/") {
			return path
		}
		if parts := strings.Split(path, "/"+seperator+"/"); len(parts) == 2 {
			return seperator + "/" + parts[1]
		}
	}
	return ""
}

func parseTilesetNames(value interface{}) []string {
	if val, ok := value.(string); ok {
		if strings.HasSuffix(val, ".glb") || strings.HasSuffix(val, ".json") || strings.HasSuffix(val, ".pnts") || strings.HasSuffix(val, ".b3dm") {
			return []string{val}
		}
	}
	retData := make([]string, 0)
	if val, ok := value.(map[string]interface{}); ok {
		for _, v := range val {
			retData = append(retData, parseTilesetNames(v)...)
		}
	}
	if val, ok := value.([]interface{}); ok {
		for _, v := range val {
			retData = append(retData, parseTilesetNames(v)...)
		}
	}
	return retData
}

func comparePath(basePath, path2 string) string {
	if len(basePath) < 1 {
		return path2
	}
	parts := strings.Split(basePath, "/")
	parts[len(parts)-1] = path2
	return strings.Join(parts, "/")
}

func downloadFile(urlString string, dest string) error {
	if info, err := os.Stat(dest); err == nil && info.Size() > 0 {
		fmt.Println("File already exists:", dest, ", skipping download")
		return nil
	}
	fmt.Println("downloadFile:", urlString, "-->", dest)
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}
	var (
		downloadError error
	)
	for i := 0; i < 3; i++ {
		_, downloadError = httpClient.R().SetOutput(dest).Get(urlString)
		if downloadError == nil {
			return nil
		}
		fmt.Println("downloadError", downloadError.Error(), ", sleep 5s and retry")
		time.Sleep(5 * time.Second)
	}

	return downloadError
}

func downloadFileToContent(urlString string, dest string) ([]byte, error) {
	if err := downloadFile(urlString, dest); err != nil {
		return nil, err
	}
	return os.ReadFile(dest)
}
