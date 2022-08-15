package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/unix"
)

// CountLine returns line count of given file.
func CountLine(filePath string) (count int, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return count, err
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	br := bufio.NewReader(file)
	for {
		if _, _, err = br.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
		count++
	}
	return count, nil
}

// Copy file from srcFilePath to dstFilePath.
func Copy(srcFilePath string, dstFilePath string) (err error) {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	r := bufio.NewReader(srcFile)

	dstFile, err := os.OpenFile(dstFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(dstFile)

	_, err = io.Copy(w, r)

	defer func() {
		closeSrcFileErr := srcFile.Close()
		closeDstFileErr := dstFile.Close()

		if closeSrcFileErr != nil || closeDstFileErr != nil {
			err = fmt.Errorf("error closeSrcFileErr: %v, closeDstFileErr: %v", closeSrcFileErr, closeDstFileErr)
		}
	}()
	return err
}

// Read whole content string of a file.
func Read(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)

	return string(bytes), err
}

// Write string data into file.
// It creates file if not exists, and overwrite whole content in case file already exists.
func Write(filePath string, data string) error {
	return ioutil.WriteFile(filePath, []byte(data), 0644)
}

// Exists checks if a file or directory exists.
func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsReadable checks if a file or directory can be read.
func IsReadable(filePath string) bool {
	return syscall.Access(filePath, unix.R_OK) == nil
}

// Rename a file or directory.
func Rename(oldFilePath string, newFilePath string) error {
	return os.Rename(oldFilePath, newFilePath)
}

// Remove removes given filePath and any children it contains.
func Remove(filePath string) error {
	return os.RemoveAll(filePath)
}

// MakeDir creates a directory recursively.
func MakeDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// ClearDir removes all files in a directory.
func ClearDir(dirPath string) (err error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := dir.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		if err = os.RemoveAll(filepath.Join(dirPath, name)); err != nil {
			return err
		}
	}
	return nil
}

// GetAllFiles returns all files in a directory.
// If suffix is not empty, it returns only files of specified suffix.
func GetAllFiles(dirPath string, suffix string) (filePaths []string, err error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := dir.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	filesInDir, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, file := range filesInDir {
		fileName := file.Name()
		fileName = filepath.Join(dirPath, file.Name())
		if suffix != "" {
			if path.Ext(fileName) != suffix {
				continue
			}
		}
		filePaths = append(filePaths, fileName)
	}
	return filePaths, nil
}

// AppendString appends string data to a file.
// It creates distFile in case not exists, and truncates distFile in case already exists.
func AppendString(filePath string, data string) (err error) {
	dstFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dstFile.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	writer := bufio.NewWriter(dstFile)
	if _, err = writer.WriteString(data); err != nil {
		return err
	}
	return writer.Flush()
}
