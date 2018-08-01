package scribe

type target struct {
	isFile   bool
	rootPath string
	fileSize int64
}

func createTarget(root string, fileSize int64) (*target, error) {
	target := &target{
		isFile:   true,
		rootPath: root,
		fileSize: fileSize,
	}

	return target, nil
}
