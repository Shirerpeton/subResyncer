package syncer

import (
	"time"
	"strings"
	"fmt"
	"strconv"
	"os"
	"errors"
)

func getDurationFromTimestamp(timestamp string) (time.Duration, error) {
	parts := strings.Split(timestamp, ":")
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		err = fmt.Errorf("can't convert timestamp to duration, error in hours: %v", err)
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		err = fmt.Errorf("can't convert timestamp to duration, error in minutes: %v", err)
		return 0, err
	}
	var sep string
	if strings.Contains(parts[2], ".") {
		sep = "."
	} else if strings.Contains(parts[2], ",") {
		sep = ","
	}
	parts = strings.Split(parts[2], sep)
	seconds, err := strconv.Atoi(parts[0])
	if err != nil {
		err = fmt.Errorf("can't convert timestamp to duration, error in seconds: %v", err)
		return 0, err
	}
	if len(parts[1]) == 2 {
		parts[1] += "0"
	}
	miliseconds, err := strconv.Atoi(parts[1])
	if err != nil {
		err = fmt.Errorf("can't convert timestamp to duration, error in miliseconds: %v", err)
		return 0, err
	}
	duration := hours*int(time.Hour) + minutes*int(time.Minute) + seconds*int(time.Second) + miliseconds*int(time.Millisecond);
	return time.Duration(duration), nil
}

func syncSub(path string, shift time.Duration) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	contentStr := string(content)
	if strings.HasSuffix(path, ".ass") {
		return syncAssSub(contentStr, shift)
	} else {
		return syncSrtSub(contentStr, shift)
	}
}

func getAssTimestampFromDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	centiseconds := int(d.Milliseconds() / 10) % 100
	timestamp := fmt.Sprintf("%d:%d02:%02d.%02d", hours, minutes, seconds, centiseconds)
	return timestamp
}

func getSrtTimestampFromDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	miliseconds := int(d.Milliseconds()) % 1000
	timestamp := fmt.Sprintf("%02d:%02d:%02d,%d", hours, minutes, seconds, miliseconds)
	return timestamp
}

func syncAssSub(content string, shift time.Duration) (string, error) {
	lines := make([]string, 0)
	for line := range strings.SplitSeq(content, "\n") {
		if !strings.HasPrefix(line, "Dialogue: ") {
			lines = append(lines, line)
		}
		parts := strings.Split(line, ",")
		if len(parts) < 10 {
			return "", errors.New("malformed subtitle file")
		}
		start, err := getDurationFromTimestamp(parts[1])
		if err != nil {
			return "", err
		}
		end, err := getDurationFromTimestamp(parts[2])
		if err != nil {
			return "", err
		}
		start += shift
		end += shift
		parts[1] = getAssTimestampFromDuration(start)
		parts[2] = getAssTimestampFromDuration(end)
		newLine := strings.Join(parts, ",")
		lines = append(lines, newLine)
	}
	result := strings.Join(lines, "\n")
	return result, nil
}

func syncSrtSub(content string, shift time.Duration) (string, error) {
	lines := make([]string, 0)
	for line := range strings.SplitSeq(content, "\n") {
		if !strings.Contains(line, "-->") {
			lines = append(lines, line)
			continue
		}
		parts := strings.Split(line, " ")
		start, err := getDurationFromTimestamp(parts[0])
		if err != nil {
			return "", err
		}
		end, err := getDurationFromTimestamp(parts[2])
		if err != nil {
			return "", err
		}
		start += shift
		end += shift
		parts[0] = getSrtTimestampFromDuration(start)
		parts[2] = getSrtTimestampFromDuration(end)
		newLine := strings.Join(parts, " ")
		lines = append(lines, newLine)
	}
	result := strings.Join(lines, "\n")
	return result, nil
}

func Sync(file string, shift time.Duration) (string, error) {
	result, err := syncSub(file, shift)
	if err != nil {
		return "", err
	}
	return result, nil
}
