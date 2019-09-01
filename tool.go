package geektimedl

import "path/filepath"

func doesDownloadAllTS(m3u8Path string) (bool, error) {
	ts, err := parseM3u8(m3u8Path)
	if err != nil {
		return false, err
	}

	dir := filepath.Dir(m3u8Path)
	for _, t := range ts {
		if ok, err := fileExists(filepath.Join(dir, t)); err == nil {
			if !ok {
				return false, nil
			}
		} else {
			return false, err
		}
	}

	return true, nil
}
