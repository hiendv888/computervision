package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
)

func main() {
	// URL của ảnh mẫu (sử dụng ảnh từ picsum.photos)
	imageURL := "https://picsum.photos/800/600"
	imageFile := "sample_image.jpg"

	// Tải ảnh từ URL
	fmt.Println("Đang tải ảnh từ:", imageURL)
	err := downloadImage(imageURL, imageFile)
	if err != nil {
		log.Fatalf("Lỗi khi tải ảnh: %v", err)
	}
	fmt.Println("Đã tải ảnh thành công:", imageFile)

	// Mở file ảnh
	file, err := os.Open(imageFile)
	if err != nil {
		log.Fatalf("Lỗi khi mở file ảnh: %v", err)
	}
	defer file.Close()

	// Decode ảnh để lấy thông tin
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Lỗi khi decode ảnh: %v", err)
	}

	// Lấy kích thước ảnh
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// In ra kết quả
	fmt.Println("\n=== THÔNG TIN ẢNH ===")
	fmt.Printf("Height: %d pixels\n", height)
	fmt.Printf("Width: %d pixels\n", width)

	// Lấy giá trị pixel tại tọa độ (100, 200)
	x, y := 100, 200

	// Kiểm tra xem tọa độ có nằm trong ảnh không
	if x >= 0 && x < width && y >= 0 && y < height {
		// Lấy màu tại pixel
		color := img.At(x, y)

		// Chuyển đổi sang RGBA để lấy giá trị từng kênh màu
		r, g, b, _ := color.RGBA()

		// RGBA() trả về giá trị 16-bit (0-65535), chuyển về 8-bit (0-255)
		r8 := uint8(r >> 8)
		g8 := uint8(g >> 8)
		b8 := uint8(b >> 8)

		// In ra theo định dạng RGB
		fmt.Printf("[%d,%d,%d]\n", r8, g8, b8)
	} else {
		fmt.Printf("Tọa độ (%d, %d) nằm ngoài phạm vi ảnh!\n", x, y)
	}

	// In ra toàn bộ ảnh theo định dạng mảng 2D
	fmt.Println("\n=== TOÀN BỘ PIXEL CỦA ẢNH (Định dạng: [[[R,G,B], ...], ...]) ===")
	fmt.Println("[")

	for row := 0; row < height; row++ {
		fmt.Print("  [")
		for col := 0; col < width; col++ {
			color := img.At(col, row)
			r, g, b, _ := color.RGBA()

			// Chuyển về 8-bit
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			fmt.Printf("[%d,%d,%d]", r8, g8, b8)

			if col < width-1 {
				fmt.Print(", ")
			}
		}
		fmt.Print("]")

		if row < height-1 {
			fmt.Println(",")
		} else {
			fmt.Println()
		}
	}

	fmt.Println("]")

}

// downloadImage tải ảnh từ URL và lưu vào file
func downloadImage(url, filepath string) error {
	// Tạo HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Kiểm tra status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Tạo file để lưu ảnh
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy dữ liệu từ response vào file
	_, err = out.ReadFrom(resp.Body)
	return err
}
