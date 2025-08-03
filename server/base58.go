package server

import (
	"errors"
	"math/big"
)

// Base58 字母表 (Bitcoin风格)
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var base58Map [256]int

func init() {
	// 初始化字符映射表
	for i := range base58Map {
		base58Map[i] = -1
	}
	for i, char := range base58Alphabet {
		base58Map[char] = i
	}
}

// Base58Encode 编码字节数组为Base58字符串
func Base58Encode(input []byte) string {
	if len(input) == 0 {
		return ""
	}

	// 计算前导零字节数
	zeroCount := 0
	for zeroCount < len(input) && input[zeroCount] == 0 {
		zeroCount++
	}

	// 转换为大整数
	num := new(big.Int).SetBytes(input)
	
	// 转换为base58
	var result []byte
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := new(big.Int)
	
	for num.Cmp(zero) > 0 {
		num.DivMod(num, base, mod)
		result = append(result, base58Alphabet[mod.Int64()])
	}

	// 添加前导1字符（对应前导零字节）
	for i := 0; i < zeroCount; i++ {
		result = append(result, base58Alphabet[0])
	}

	// 反转结果
	for i := 0; i < len(result)/2; i++ {
		result[i], result[len(result)-1-i] = result[len(result)-1-i], result[i]
	}

	return string(result)
}

// Base58Decode 解码Base58字符串为字节数组
func Base58Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil
	}

	// 计算前导1字符数
	zeroCount := 0
	for zeroCount < len(input) && input[zeroCount] == base58Alphabet[0] {
		zeroCount++
	}

	// 转换为大整数
	num := big.NewInt(0)
	base := big.NewInt(58)
	
	for _, char := range input {
		if char > 255 || base58Map[char] == -1 {
			return nil, errors.New("invalid base58 character")
		}
		num.Mul(num, base)
		num.Add(num, big.NewInt(int64(base58Map[char])))
	}

	// 转换为字节数组
	decoded := num.Bytes()
	
	// 添加前导零字节
	result := make([]byte, zeroCount+len(decoded))
	copy(result[zeroCount:], decoded)
	
	return result, nil
}