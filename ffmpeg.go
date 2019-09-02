package geektime

import (
	"os/exec"
)

func m3u8ToMP4(m3u8Path, outputPath string) (string, error) {
	cmd := exec.Command("ffmpeg", "-i", m3u8Path, "-c", "copy", "-bsf:a", "aac_adtstoasc", outputPath)
	data, err := cmd.Output()
	return string(data), err
}

func m3u8ToMP3(m3u8Path, outputPath string) (string, error) {
	cmd := exec.Command("ffmpeg", "-i", m3u8Path, "-acodec", "mp3", "-ab", "257k", outputPath)
	data, err := cmd.Output()
	return string(data), err
}
