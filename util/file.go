package util

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v4"
)

// QuickZip
//
//	@param srcDir
//	@param dest
//	@return error
func QuickZip(srcDir, dest string) error {

	file, err := os.Open(srcDir)
	if err != nil {
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("file stat error: %v", err)
	}

	fileMap := map[string]string{}
	if fileStat.IsDir() {
		fileMap = GetDirFilesAsMap(srcDir)
	} else {
		// srcDir是一个文件
		_, name := filepath.Split(srcDir)
		fileMap[name] = srcDir
	}

	fileMapRevert := map[string]string{}

	for key, value := range fileMap {
		fileMapRevert[value] = key
	}

	// map files on disk to their paths in the archive
	files, err := archiver.FilesFromDisk(nil, fileMapRevert)
	if err != nil {
		return err
	}

	// create the output file we'll write to
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// we can use the CompressedArchive type to gzip a tarball
	// (compression is not required; you could use Tar directly)
	format := archiver.Archive{
		Compression: nil,
		Archival:    archiver.Zip{},
	}

	// create the archive
	err = format.Archive(context.Background(), out, files)
	if err != nil {
		return err
	}
	return nil
}

func GetDirFilesAsMap(dirPath string) map[string]string {
	fileMap := map[string]string{}
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(dirPath, path)
			fileMap[relPath] = path
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking the path:", err)
		return nil
	}
	return fileMap
}

func Json2CsvFile(jsonData []byte, csvFilePath string) error {
	// 解析JSON数据
	var data []map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	// 创建CSV文件
	file, err := os.Create(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	// 写入CSV头
	headers := []string{}
	for key := range data[0] {
		headers = append(headers, key)
	}
	file.WriteString(strings.Join(headers, ",") + "\n")

	// 写入CSV数据
	for _, item := range data {
		values := []string{}
		for _, key := range headers {
			values = append(values, fmt.Sprintf("%v", item[key]))
		}
		file.WriteString(strings.Join(values, ",") + "\n")
	}

	return nil
}

// ReadJsonFile
//
//	@param filePath
//	@param data
//	@return error
func ReadJsonFile(filePath string, data interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to decode JSON data: %v", err)
	}

	return nil
}
