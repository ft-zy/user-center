package model

import "time"

// CustomTime 是一个包装了time.Time的结构体，用于自定义JSON的编码和解码格式
type CustomTime struct {
	time.Time
}

// MarshalJSON 实现了CustomTime的自定义JSON编码
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	// 使用自定义的格式来格式化时间
	return []byte(ct.Time.Format("\"2006-01-02 15:04:05\"")), nil
}

// UnmarshalJSON 实现了CustomTime的自定义JSON解码
// 注意：这里假设输入的JSON字符串已经按照自定义格式格式化
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// 去掉JSON字符串的引号
	data = data[1 : len(data)-1]
	parsedTime, err := time.Parse("2006-01-02 15:04:05", string(data))
	if err != nil {
		return err
	}
	ct.Time = parsedTime
	return nil
}
