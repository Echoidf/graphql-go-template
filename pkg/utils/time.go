package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type TimeRange struct {
	Start int
	End   int
}

func ProcessTimeRange(s string) []TimeRange {
	// 先按逗号分割多个时间段
	ranges := strings.Split(s, ",")
	result := make([]TimeRange, 0, len(ranges))

	for _, r := range ranges {
		// 分割开始和结束时间
		times := strings.Split(strings.TrimSpace(r), "-")
		if len(times) != 2 {
			continue
		}

		// 解析开始时间
		startParts := strings.Split(times[0], ":")
		if len(startParts) != 3 {
			continue
		}
		startHour, _ := strconv.Atoi(startParts[0])
		startMin, _ := strconv.Atoi(startParts[1])
		startSec, _ := strconv.Atoi(startParts[2])

		// 解析结束时间
		endParts := strings.Split(times[1], ":")
		if len(endParts) != 3 {
			continue
		}
		endHour, _ := strconv.Atoi(endParts[0])
		endMin, _ := strconv.Atoi(endParts[1])
		endSec, _ := strconv.Atoi(endParts[2])

		// 转换为秒数
		startTime := startHour*3600 + startMin*60 + startSec
		endTime := endHour*3600 + endMin*60 + endSec

		result = append(result, TimeRange{
			Start: startTime,
			End:   endTime,
		})
	}

	return result
}

func ParseIntervalString(s string) (int, int, int, int, error) {
	// 使用正则表达式匹配格式
	re := regexp.MustCompile(`^\[(\d+)-(\d+)\):(\d+):(\d+)$`)
	matches := re.FindStringSubmatch(s)

	if len(matches) != 5 {
		return 0, 0, 0, 0, fmt.Errorf("invalid format")
	}

	// 将字符串转换为整数
	var nums [4]int
	var err error
	for i := 0; i < 4; i++ {
		nums[i], err = strconv.Atoi(matches[i+1])
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("invalid number: %v", err)
		}
	}

	return nums[0], nums[1], nums[2], nums[3], nil
}
