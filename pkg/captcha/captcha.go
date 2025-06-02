package captcha

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	captchaStore = make(map[string]captchaData)
	storeMutex   sync.RWMutex
	// 使用数字和大写字母
	chars = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ" // 去掉容易混淆的I和O
)

type captchaData struct {
	code         string
	createdAt    time.Time
	failAttempts int // 错误尝试次数
}

func init() {
	// Clean expired captchas every 5 minutes
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cleanExpiredCaptchas()
			}
		}
	}()
}

// GenerateCaptcha generates a new captcha and returns ID and base64 image
func GenerateCaptcha() (string, string, error) {
	// Generate random ID and code
	id := generateRandomString(16)
	code := generateRandomString(4)

	// Store captcha
	storeMutex.Lock()
	captchaStore[id] = captchaData{
		code:         code,
		createdAt:    time.Now(),
		failAttempts: 0,
	}
	storeMutex.Unlock()

	// Generate image
	img := generateCaptchaImage(code)

	// Convert to base64
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", "", err
	}

	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return id, "data:image/png;base64," + b64, nil
}

// VerifyCaptcha verifies a captcha code
func VerifyCaptcha(id, code string) bool {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	data, exists := captchaStore[id]
	if !exists {
		return false
	}

	// Check if expired (5 minutes)
	if time.Since(data.createdAt) > 5*time.Minute {
		delete(captchaStore, id)
		return false
	}

	// 检查错误尝试次数是否超过限制（最多3次）
	if data.failAttempts >= 3 {
		delete(captchaStore, id)
		return false
	}

	// 验证验证码是否正确（不区分大小写）
	isValid := strings.EqualFold(data.code, code)

	if isValid {
		// 验证成功，删除验证码（一次性使用）
		delete(captchaStore, id)
	} else {
		// 验证失败，增加失败次数
		data.failAttempts++
		captchaStore[id] = data
	}

	return isValid
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func generateCaptchaImage(code string) image.Image {
	width, height := 120, 40
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充白色背景
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// 添加干扰线
	for i := 0; i < 6; i++ {
		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)
		drawLine(img, x1, y1, x2, y2, color.RGBA{uint8(rand.Intn(128)), uint8(rand.Intn(128)), uint8(rand.Intn(128)), 255})
	}

	// 添加干扰点
	for i := 0; i < 100; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		img.Set(x, y, color.RGBA{uint8(rand.Intn(128)), uint8(rand.Intn(128)), uint8(rand.Intn(128)), 255})
	}

	// 绘制字符
	for i, ch := range code {
		if ch >= '0' && ch <= '9' {
			drawDigit(img, int(ch-'0'), 15+i*25, 8)
		} else {
			drawLetter(img, ch, 15+i*25, 8)
		}
	}

	return img
}

// 绘制数字的点阵
func drawDigit(img *image.RGBA, digit int, x, y int) {
	patterns := [][]bool{
		{ // 0
			true, true, true,
			true, false, true,
			true, false, true,
			true, false, true,
			true, true, true,
		},
		{ // 1
			false, true, false,
			true, true, false,
			false, true, false,
			false, true, false,
			true, true, true,
		},
		{ // 2
			true, true, true,
			false, false, true,
			true, true, true,
			true, false, false,
			true, true, true,
		},
		{ // 3
			true, true, true,
			false, false, true,
			true, true, true,
			false, false, true,
			true, true, true,
		},
		{ // 4
			true, false, true,
			true, false, true,
			true, true, true,
			false, false, true,
			false, false, true,
		},
		{ // 5
			true, true, true,
			true, false, false,
			true, true, true,
			false, false, true,
			true, true, true,
		},
		{ // 6
			true, true, true,
			true, false, false,
			true, true, true,
			true, false, true,
			true, true, true,
		},
		{ // 7
			true, true, true,
			false, false, true,
			false, true, false,
			true, false, false,
			true, false, false,
		},
		{ // 8
			true, true, true,
			true, false, true,
			true, true, true,
			true, false, true,
			true, true, true,
		},
		{ // 9
			true, true, true,
			true, false, true,
			true, true, true,
			false, false, true,
			true, true, true,
		},
	}

	pattern := patterns[digit]
	dotSize := 3
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			if pattern[i*3+j] {
				drawDot(img, x+j*dotSize, y+i*dotSize, dotSize, color.RGBA{0, 0, 0, 255})
			}
		}
	}
}

// 绘制字母的点阵
func drawLetter(img *image.RGBA, letter rune, x, y int) {
	patterns := map[rune][]bool{
		'A': {
			true, true, true,
			true, false, true,
			true, true, true,
			true, false, true,
			true, false, true,
		},
		'B': {
			true, true, false,
			true, false, true,
			true, true, false,
			true, false, true,
			true, true, false,
		},
		'C': {
			true, true, true,
			true, false, false,
			true, false, false,
			true, false, false,
			true, true, true,
		},
		'D': {
			true, true, false,
			true, false, true,
			true, false, true,
			true, false, true,
			true, true, false,
		},
		'E': {
			true, true, true,
			true, false, false,
			true, true, false,
			true, false, false,
			true, true, true,
		},
		'F': {
			true, true, true,
			true, false, false,
			true, true, false,
			true, false, false,
			true, false, false,
		},
		'G': {
			true, true, true,
			true, false, false,
			true, false, true,
			true, false, true,
			true, true, true,
		},
		'H': {
			true, false, true,
			true, false, true,
			true, true, true,
			true, false, true,
			true, false, true,
		},
		'J': {
			false, false, true,
			false, false, true,
			false, false, true,
			true, false, true,
			true, true, true,
		},
		'K': {
			true, false, true,
			true, false, true,
			true, true, false,
			true, false, true,
			true, false, true,
		},
		'L': {
			true, false, false,
			true, false, false,
			true, false, false,
			true, false, false,
			true, true, true,
		},
		'M': {
			true, false, true,
			true, true, true,
			true, false, true,
			true, false, true,
			true, false, true,
		},
		'N': {
			true, false, true,
			true, true, true,
			true, true, true,
			true, false, true,
			true, false, true,
		},
		'P': {
			true, true, true,
			true, false, true,
			true, true, true,
			true, false, false,
			true, false, false,
		},
		'Q': {
			true, true, true,
			true, false, true,
			true, false, true,
			true, true, false,
			false, true, true,
		},
		'R': {
			true, true, true,
			true, false, true,
			true, true, false,
			true, false, true,
			true, false, true,
		},
		'S': {
			true, true, true,
			true, false, false,
			true, true, true,
			false, false, true,
			true, true, true,
		},
		'T': {
			true, true, true,
			false, true, false,
			false, true, false,
			false, true, false,
			false, true, false,
		},
		'U': {
			true, false, true,
			true, false, true,
			true, false, true,
			true, false, true,
			true, true, true,
		},
		'V': {
			true, false, true,
			true, false, true,
			true, false, true,
			true, false, true,
			false, true, false,
		},
		'W': {
			true, false, true,
			true, false, true,
			true, false, true,
			true, true, true,
			true, false, true,
		},
		'X': {
			true, false, true,
			true, false, true,
			false, true, false,
			true, false, true,
			true, false, true,
		},
		'Y': {
			true, false, true,
			true, false, true,
			false, true, false,
			false, true, false,
			false, true, false,
		},
		'Z': {
			true, true, true,
			false, false, true,
			false, true, false,
			true, false, false,
			true, true, true,
		},
	}

	pattern, exists := patterns[letter]
	if !exists {
		return
	}

	dotSize := 3
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			if pattern[i*3+j] {
				drawDot(img, x+j*dotSize, y+i*dotSize, dotSize, color.RGBA{0, 0, 0, 255})
			}
		}
	}
}

// 绘制点
func drawDot(img *image.RGBA, x, y, size int, c color.Color) {
	for dy := 0; dy < size; dy++ {
		for dx := 0; dx < size; dx++ {
			img.Set(x+dx, y+dy, c)
		}
	}
}

// 绘制直线
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	steep := dy > dx

	if steep {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
	}
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	dx = x2 - x1
	dy = abs(y2 - y1)
	err := dx / 2
	ystep := 1
	if y1 >= y2 {
		ystep = -1
	}

	y := y1
	for x := x1; x <= x2; x++ {
		if steep {
			img.Set(y, x, c)
		} else {
			img.Set(x, y, c)
		}
		err -= dy
		if err < 0 {
			y += ystep
			err += dx
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func cleanExpiredCaptchas() {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	now := time.Now()
	for id, data := range captchaStore {
		if now.Sub(data.createdAt) > 5*time.Minute {
			delete(captchaStore, id)
		}
	}
}
