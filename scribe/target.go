package scribe

type target struct {
	isFile   bool
	rootPath string
	fileSize int64
}

func createTarget(root string, fileSize int64) (*target, error) {
	err := CheckPath(root)
	if err != nil {
		return nil, err
	}

	target := &target{
		isFile: true,
		rootPath: root,
		fileSize: fileSize,
	}

	return target, nil
}
