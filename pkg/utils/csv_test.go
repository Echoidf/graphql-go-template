package utils

import (
	"testing"
)

type testStruct struct {
	Name  string  `csv:"name"`
	Age   int     `csv:"age"`
	Price float64 `csv:"price"`
}

func TestWriteToCsv(t *testing.T) {
	var data []testStruct
	for i := range 10 {
		data = append(data, testStruct{
			Name:  "test",
			Age:   i+1,
			Price: float64(i),
		})
	}
	
	if err := WriteToCsv(data, "test.csv"); err != nil {
		t.Errorf("Failed to write CSV: %v", err)
	}
	
	// 测试读取
	readData, err := ReadFromCsv[testStruct]("test.csv")
	if err != nil {
		t.Errorf("Failed to read CSV: %v", err)
	}
	
	// 验证数据
	if len(readData) != len(data) {
		t.Errorf("Expected %d records, got %d", len(data), len(readData))
	}
	
	for i, item := range readData {
		if item.Name != data[i].Name || item.Age != data[i].Age || item.Price != data[i].Price {
			t.Errorf("Record %d mismatch: expected %v, got %v", i, data[i], item)
		}
	}
}