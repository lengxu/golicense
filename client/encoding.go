package client

import (
	"compress/gzip"
	"encoding/json"
	"bytes"
	"fmt"
)

// EncodeToString 将结构体编码为压缩的base64字符串
func EncodeToString(data interface{}) (string, error) {
	// 1. JSON序列化
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %v", err)
	}

	// 2. Gzip压缩
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write(jsonData); err != nil {
		return "", fmt.Errorf("failed to compress data: %v", err)
	}
	if err := gzWriter.Close(); err != nil {
		return "", fmt.Errorf("failed to close gzip writer: %v", err)
	}

	// 3. Base58编码
	encoded := Base58Encode(buf.Bytes())
	
	// 4. 添加标识前缀
	return "REQ:" + encoded, nil
}

// DecodeFromString 从压缩的base64字符串解码为结构体
func DecodeFromString(encoded string, result interface{}) error {
	// 1. 检查并移除前缀
	if len(encoded) < 4 {
		return fmt.Errorf("invalid encoded string: too short")
	}
	
	prefix := encoded[:4]
	data := encoded[4:]
	
	if prefix != "REQ:" && prefix != "LIC:" {
		return fmt.Errorf("invalid encoded string: unknown prefix %s", prefix)
	}

	// 2. Base58解码
	compressed, err := Base58Decode(data)
	if err != nil {
		return fmt.Errorf("failed to decode base58: %v", err)
	}

	// 3. Gzip解压
	gzReader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(gzReader); err != nil {
		return fmt.Errorf("failed to decompress data: %v", err)
	}

	// 4. JSON反序列化
	if err := json.Unmarshal(buf.Bytes(), result); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	return nil
}

// EncodeLicenseToString 将license编码为字符串（带LIC前缀）
func EncodeLicenseToString(data interface{}) (string, error) {
	// 1. JSON序列化
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %v", err)
	}

	// 2. Gzip压缩
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write(jsonData); err != nil {
		return "", fmt.Errorf("failed to compress data: %v", err)
	}
	if err := gzWriter.Close(); err != nil {
		return "", fmt.Errorf("failed to close gzip writer: %v", err)
	}

	// 3. Base58编码
	encoded := Base58Encode(buf.Bytes())
	
	// 4. 添加license标识前缀
	return "LIC:" + encoded, nil
}