package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"image/png"
	"os"

	"github.com/go-errors/errors"
	openai "github.com/sashabaranov/go-openai"
)

func createImageUrl(prompt string) (string, error) {
	ctx := context.Background()
	reqUrl := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	}

	respUrl, err := c.CreateImage(ctx, reqUrl)
	if err != nil {
		return "", err
	}

	return respUrl.Data[0].URL, nil
}

func createImageBytes(prompt string) ([]byte, error) {
	ctx := context.Background()
	reqBase64 := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}

	respBase64, err := c.CreateImage(ctx, reqBase64)
	if err != nil {
		return nil, errors.Errorf("Image creation error: %v\n", err)
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		return nil, errors.Errorf("Base64 decode error: %v\n", err)
	}

	return imgBytes, nil
}

func saveImage(imgBytes []byte, filename string) error {
	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		return err
	}

	return nil
}
