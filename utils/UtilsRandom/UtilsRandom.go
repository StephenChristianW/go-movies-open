package UtilsRandom

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// RandomAlphaNumString 随机生成指定长度的字母数字字符串（包含大写、小写字母和数字）
// length: 要生成的字符串长度
// 返回值: 随机生成的字符串，可能的错误（如 rand 出错）
func RandomAlphaNumString(length int) (string, error) {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[int(num.Int64())]
	}
	return string(result), nil
}

// RandomInt 生成 [0, max) 范围内的随机整数
// max: 随机数上限（不包含 max）
// 返回值: 随机整数，可能的错误
func RandomInt(max int) (int, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(num.Int64()), nil
}

// RandomDigitString 随机生成指定长度的数字字符串（仅包含 0-9）
// length: 要生成的字符串长度
// 返回值: 随机生成的数字字符串，可能的错误
func RandomDigitString(length int) (string, error) {
	const letters = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[int(num.Int64())]
	}
	return string(result), nil
}

// GetFiveMultipleCrypto 在 [min, max] 范围内生成 5 的倍数随机整数
// min: 范围最小值（包含）
// max: 范围最大值（包含）
// 返回值: 随机 5 的倍数整数，可能的错误（如范围内没有 5 的倍数或 rand 出错）
func GetFiveMultipleCrypto(min, max int) (int, error) {
	if min > max {
		return 0, fmt.Errorf("min 不能大于 max")
	}

	// 找到范围内第一个 >= min 的 5 的倍数
	start := min
	if start%5 != 0 {
		start += 5 - start%5
	}

	// 找到范围内最后一个 <= max 的 5 的倍数
	end := max
	if end%5 != 0 {
		end -= end % 5
	}

	if start > end {
		return 0, fmt.Errorf("该范围内没有 5 的倍数")
	}

	// 计算总共有多少个 5 的倍数
	count := (end-start)/5 + 1

	// 使用 crypto/rand 生成 [0, count) 的随机索引
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(count)))
	if err != nil {
		return 0, err
	}

	return start + int(nBig.Int64())*5, nil
}
