package helpers

import (
	"charts/domain/diff"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

func GenerateCacheKey(jsonStr []byte) string {
	hash := sha256.Sum256(jsonStr)
	return hex.EncodeToString(hash[:])
}

func FindStatus(diff *diff.CommentsDiff, comment string) (string, error) {
	var status string
	var jsonData  []byte
	var resultMap map[string]interface{}

	jsonData = diff.Result
	if err := json.Unmarshal(jsonData, &resultMap); err != nil {
		return "", err
	}

	result, ok := resultMap["status"].(map[string]interface{})
	if !ok {
		return "", errors.New("error: the 'status' field missing or has an incorrect format")
	}

	switch comment {
	case "new":
		status, ok = result["new"].(string)
		if !ok {
			return "", errors.New("error: the 'new' field missing or has an incorrect format")
		}
	case "old":
		status, ok = result["old"].(string)
		if !ok {
			return "", errors.New("error: the 'new' field missing or has an incorrect format")
		}
	}

	return status, nil
}
