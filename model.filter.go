package gou

import (
	"github.com/yaoapp/kun/exception"
	"github.com/yaoapp/kun/maps"
	"github.com/yaoapp/xun/dbal"
)

// FliterIn 输入前过滤解码
func (mod *Model) FliterIn(row maps.MapStrAny) {
	for name, value := range row {
		column, has := mod.Columns[name]

		// 删除无效字段
		if !has {
			row.Del(name)
			continue
		}

		// 过滤输入信息
		column.FliterIn(value, row)
	}
}

// FliterSelect 选择字段
func (mod *Model) FliterSelect(alias string, columns []interface{}) []interface{} {
	res := []interface{}{}
	for _, col := range columns {

		if _, ok := col.(dbal.Expression); ok {
			res = append(res, col)
			continue
		}

		name, ok := col.(string)
		if !ok {
			continue
		}

		column, has := mod.Columns[name]
		if !has {
			continue
		}

		// alias.field
		field := name
		if alias != "" {
			field = alias + "." + name
		}

		// 加密字段
		if column.Crypt == "AES" {
			icrypt := SelectCrypt(column.Crypt)
			raw, err := icrypt.Decode(field)
			if err != nil {
				exception.Err(err, 500).Throw()
			}
			raw = raw + " as " + name
			res = append(res, dbal.Raw(raw))
		} else {
			res = append(res, field)
		}
	}
	return res
}

// FliterWhere 选项
func (mod *Model) FliterWhere(alias string, col interface{}) interface{} {
	if _, ok := col.(dbal.Expression); ok {
		return col
	}

	name, ok := col.(string)
	if !ok {
		return col
	}

	column, has := mod.Columns[name]
	if !has {
		return col
	}

	// alias.field
	if alias != "" {
		name = alias + "." + name
	}

	// 加密字段
	if column.Crypt == "AES" {
		icrypt := SelectCrypt(column.Crypt)
		raw, err := icrypt.Decode(name)
		if err != nil {
			exception.Err(err, 500).Throw()
		}
		return dbal.Raw(raw)
	}

	return name
}

// FliterOut 输出前过滤解码
func (mod *Model) FliterOut(row maps.MapStrAny) {
	for name, value := range row {
		column, has := mod.Columns[name]

		// 删除无效字段
		if !has {
			continue
		}

		// 过滤输入信息
		column.FliterOut(value, row)
	}
}
