package git

import (
	"fmt"
	"os"
)

func Init(path string) (err error) {
	dirs := []string{"refs/heads", "refs/tags", "object/info", "object/pack"}
	for _, v := range dirs {
		err = os.MkdirAll(fmt.Sprintf("%s/%s", path, v), 0770)
		if err != nil {
			return err
		}
	}
	return nil
}
