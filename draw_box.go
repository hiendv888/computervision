package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
)

func main() {
	// URL của ảnh mẫu
	imageURL := "https://picsum.photos/800/600"
	inputFile := "input_image.jpg"
	outputFile := "output_with_box.jpg"

	// Tải ảnh từ URL
	fmt.Println("Đang tải ảnh từ:", imageURL)
	err := downloadImage(imageURL, inputFile)
	if err != nil {
		log.Fatalf("Lỗi khi tải ảnh: %v", err)
	}
	fmt.Println("Đã tải ảnh thành công:", inputFile)

	// Mở file ảnh
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Lỗi khi mở file ảnh: %v", err)
	}
	defer file.Close()

	// Decode ảnh
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Lỗi khi decode ảnh: %v", err)
	}

	// Chuyển sang RGBA để có thể vẽ
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Định nghĩa bounding box (x, y, w, h)
	x := 200
	y := 150
	w := 300
	h := 200

	// Màu đỏ cho bounding box
	boxColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	thickness := 3

	// Vẽ bounding box
	drawBoundingBox(rgba, x, y, w, h, boxColor, thickness)

	fmt.Printf("\nĐã vẽ bounding box tại: x=%d, y=%d, w=%d, h=%d\n", x, y, w, h)

	// Lưu ảnh kết quả
	err = saveImage(rgba, outputFile)
	if err != nil {
		log.Fatalf("Lỗi khi lưu ảnh: %v", err)
	}

	fmt.Println("Đã lưu ảnh kết quả:", outputFile)
}

// drawBoundingBox vẽ hình chữ nhật lên ảnh
func drawBoundingBox(img *image.RGBA, x, y, w, h int, c color.Color, thickness int) {
	// Vẽ 4 cạnh của hình chữ nhật

	// Cạnh trên
	for t := 0; t < thickness; t++ {
		for i := x; i < x+w; i++ {
			if i >= 0 && i < img.Bounds().Dx() && y+t >= 0 && y+t < img.Bounds().Dy() {
				img.Set(i, y+t, c)
			}
		}
	}

	// Cạnh dưới
	for t := 0; t < thickness; t++ {
		for i := x; i < x+w; i++ {
			if i >= 0 && i < img.Bounds().Dx() && y+h-t-1 >= 0 && y+h-t-1 < img.Bounds().Dy() {
				img.Set(i, y+h-t-1, c)
			}
		}
	}

	// Cạnh trái
	for t := 0; t < thickness; t++ {
		for i := y; i < y+h; i++ {
			if x+t >= 0 && x+t < img.Bounds().Dx() && i >= 0 && i < img.Bounds().Dy() {
				img.Set(x+t, i, c)
			}
		}
	}

	// Cạnh phải
	for t := 0; t < thickness; t++ {
		for i := y; i < y+h; i++ {
			if x+w-t-1 >= 0 && x+w-t-1 < img.Bounds().Dx() && i >= 0 && i < img.Bounds().Dy() {
				img.Set(x+w-t-1, i, c)
			}
		}
	}
}

// downloadImage tải ảnh từ URL và lưu vào file
func downloadImage(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(resp.Body)
	return err
}

// saveImage lưu ảnh vào file
func saveImage(img *image.RGBA, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Sử dụng encoder JPEG với quality 95
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}
