package geektimedl

import (
	"os/exec"
)

func m3u8ToMP4(m3u8Path, outputPath string) (string, error) {
	cmd := exec.Command("ffmpeg", "-i", m3u8Path, "-c", "copy", "-bsf:a", "aac_adtstoasc", outputPath)
	data, err := cmd.Output()
	return string(data), err
}
