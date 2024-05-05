package helper

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ZipFiles(path, filename string, files []string, rootPath string) error {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	newZipFile, err := os.Create(filepath.Join(path, filename+".zip"))
	if err != nil {
		return err
	}

	defer newZipFile.Close()
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 把files添加到zip中
	for _, file := range files {
		zipFile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipFile.Close()
		info, err := zipFile.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.Replace(file, rootPath, "", -1)
		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipFile); err != nil {
			return err
		}
	}
	return nil
}

func UnzipFiles(filename, path string) error {
	// 打开zip文件
	r, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	// 遍历zip中的所有文件
	for _, f := range r.Reader.File {
		zipped, err := f.Open()
		if err != nil {
			return err
		}

		// 创建文件路径
		path := filepath.Join(path, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			// 创建文件
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(file, zipped)
			file.Close()
			zipped.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
