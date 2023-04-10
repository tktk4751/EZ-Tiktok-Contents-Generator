package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	inputFile := "src/input/input.mp4"
	outputFileName := "output_basic_size" + time.Now().Format("20060102150405") + ".mp4"
	outputFile := filepath.Join("src/output", outputFileName)
	backgroundColor := "2f2f2f" //Feel free to change the color code
	text := "今日のハイライト"          //Feel free to change text

	err := resizeAndAddText(inputFile, outputFile, text, backgroundColor)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func resizeAndAddText(inputFile, outputFile, text, backgroundColor string) error {
	fontfile := "src/font/JNRfont_n.ttf"
	fontsize := 120
	textColor := "white"
	borderColor := "Red"
	borderWidth := 2
	shadowColor := "black@0.5"
	shadowXOffset := 3
	shadowYOffset := 3

	textFilter := fmt.Sprintf("drawtext=text='%s':x=(w-text_w)/2:y=(h-text_h)/8:fontsize=%d:fontfile=%s:fontcolor=%s:bordercolor=%s:borderw=%d:shadowcolor=%s:shadowx=%d:shadowy=%d", text, fontsize, fontfile, textColor, borderColor, borderWidth, shadowColor, shadowXOffset, shadowYOffset)

	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-vf", fmt.Sprintf("scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:%s,%s", backgroundColor, textFilter),
		"-strict", "-2", outputFile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg command execution failed: %w", err)
	}
	return nil
}
