package storage

import (
	"io"
	"net/url"
	"os"
)

const (
	objectStorageEndpointConfigKey  = "object_storage.endpoint"
	objectStorageAccessKeyConfigKey = "object_storage.access_key"
	objectStorageSecretKeyConfigKey = "object_storage.secret_key"
	objectStorageBucketConfigKey    = "object_storage.bucket"
	objectStorageBasePathConfigKey  = "object_storage.base_path"
)

// Object represents the object on the storage
type Object interface {
	io.ReadCloser
	io.Seeker
	Stat() (os.FileInfo, error)
}

// ObjectStorage represents an object storage to handle a bucket and files
type ObjectStorage interface {
	Open(path string) (Object, error)
	// Save store a object, if size is unknown set -1
	Save(path string, r io.Reader, size int64) (int64, error)
	Stat(path string) (os.FileInfo, error)
	Delete(path string) error
	URL(path, name string, reqParams url.Values) (*url.URL, error)
	IterateObjects(path string, iterator func(path string, obj Object) error) error
}

// Copy copies a file from source ObjectStorage to dest ObjectStorage
func Copy(dstStorage ObjectStorage, dstPath string, srcStorage ObjectStorage, srcPath string) (int64, error) {
	f, err := srcStorage.Open(srcPath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	size := int64(-1)
	fsinfo, err := f.Stat()
	if err == nil {
		size = fsinfo.Size()
	}

	return dstStorage.Save(dstPath, f, size)
}

// Clean delete all the objects in this storage
func Clean(storage ObjectStorage) error {
	return storage.IterateObjects("", func(path string, obj Object) error {
		_ = obj.Close()
		return storage.Delete(path)
	})
}

// SaveFrom saves data to the ObjectStorage with path p from the callback
func SaveFrom(objStorage ObjectStorage, path string, callback func(w io.Writer) error) error {
	pr, pw := io.Pipe()
	defer pr.Close()
	go func() {
		defer pw.Close()
		if err := callback(pw); err != nil {
			_ = pw.CloseWithError(err)
		}
	}()

	_, err := objStorage.Save(path, pr, -1)
	return err
}
