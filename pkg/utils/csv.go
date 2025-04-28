package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func WriteToCsv[T any](data []T, filePath string) error {
	if len(data) == 0 {
		return nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 获取结构体类型信息
	t := reflect.TypeOf(data[0])
	var headers []string
	var fieldIndexes []int

	// 解析结构体字段和 CSV tag
	for i := range t.NumField() {
		field := t.Field(i)
		if tag := field.Tag.Get("csv"); tag != "" {
			headers = append(headers, tag)
			fieldIndexes = append(fieldIndexes, i)
		}
	}

	// 写入表头
	if err := writer.Write(headers); err != nil {
		return err
	}

	// 写入数据
	for _, item := range data {
		row := make([]string, len(fieldIndexes))
		v := reflect.ValueOf(item)

		for i, idx := range fieldIndexes {
			field := v.Field(idx)
			switch field.Kind() {
			case reflect.String:
				row[i] = field.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				row[i] = strconv.FormatInt(field.Int(), 10)
			case reflect.Float32, reflect.Float64:
				row[i] = strconv.FormatFloat(field.Float(), 'f', -1, 64)
			default:
				row[i] = fmt.Sprint(field.Interface())
			}
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func ReadFromCsv[T any](filePath string) ([]T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// 读取表头
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// 创建结果切片
	var result []T

	// 获取类型信息
	var t T
	typ := reflect.TypeOf(t)

	// 创建字段映射
	fieldMap := make(map[string]int)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if tag := field.Tag.Get("csv"); tag != "" {
			for _, header := range headers {
				if header == tag {
					fieldMap[tag] = i
					break
				}
			}
		}
	}

	// 读取数据行
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// 创建新的结构体实例
		item := reflect.New(typ).Elem()

		// 填充数据
		for i, value := range record {
			if i >= len(headers) {
				continue
			}

			if fieldIndex, ok := fieldMap[headers[i]]; ok {
				field := item.Field(fieldIndex)
				switch field.Kind() {
				case reflect.String:
					field.SetString(value)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if v, err := strconv.ParseInt(value, 10, 64); err == nil {
						field.SetInt(v)
					}
				case reflect.Float32, reflect.Float64:
					if v, err := strconv.ParseFloat(value, 64); err == nil {
						field.SetFloat(v)
					}
				}
			}
		}

		result = append(result, item.Interface().(T))
	}

	return result, nil
}
