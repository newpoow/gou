package gou

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/kun/any"
	"github.com/yaoapp/kun/grpc"
	"github.com/yaoapp/kun/maps"
	"github.com/yaoapp/xun/capsule"
)

func TestProcessExec(t *testing.T) {
	defer SelectPlugin("user").Client.Kill()
	res := NewProcess("plugins.user.Login", 1).Run().(*grpc.Response).MustMap()
	res2 := NewProcess("plugins.user.Login", 2).Run().(*grpc.Response).MustMap()
	assert.Equal(t, "login", res.Get("name"))
	assert.Equal(t, "login", res2.Get("name"))
	assert.Equal(t, 1, any.Of(res.Dot().Get("args.0")).CInt())
	assert.Equal(t, 2, any.Of(res2.Dot().Get("args.0")).CInt())
}

func TestProcessFind(t *testing.T) {
	res := NewProcess("models.user.Find", 1, QueryParam{}).Run().(maps.MapStr)
	assert.Equal(t, 1, any.Of(res.Dot().Get("id")).CInt())
	assert.Equal(t, "男", res.Dot().Get("extra.sex"))
}

func TestProcessGet(t *testing.T) {
	rows := NewProcess("models.user.Get", QueryParam{Limit: 2}).Run().([]maps.MapStr)
	res := maps.Map{"data": rows}.Dot()
	assert.Equal(t, 2, len(rows))
	assert.Equal(t, 1, any.Of(res.Get("data.0.id")).CInt())
	assert.Equal(t, "男", res.Get("data.0.extra.sex"))
	assert.Equal(t, 2, any.Of(res.Get("data.1.id")).CInt())
	assert.Equal(t, "女", res.Get("data.1.extra.sex"))
}

func TestProcessPaginate(t *testing.T) {
	res := NewProcess("models.user.Paginate", QueryParam{}, 1, 2).Run().(maps.MapStr).Dot()
	assert.Equal(t, 3, res.Get("total"))
	assert.Equal(t, 1, res.Get("page"))
	assert.Equal(t, 2, res.Get("pagesize"))
	assert.Equal(t, 2, res.Get("pagecnt"))
	assert.Equal(t, 2, res.Get("next"))
	assert.Equal(t, -1, res.Get("prev"))
	assert.Equal(t, 1, any.Of(res.Get("data.0.id")).CInt())
	assert.Equal(t, "男", res.Get("data.0.extra.sex"))
	assert.Equal(t, 2, any.Of(res.Get("data.1.id")).CInt())
	assert.Equal(t, "女", res.Get("data.1.extra.sex"))
}

func TestProcessCreate(t *testing.T) {
	row := maps.MapStr{
		"name":     "用户创建",
		"manu_id":  2,
		"type":     "user",
		"idcard":   "23082619820207006X",
		"mobile":   "13900004444",
		"password": "qV@uT1DI",
		"key":      "XZ12MiPp",
		"secret":   "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN",
		"status":   "enabled",
		"extra":    maps.MapStr{"sex": "女"},
	}
	id := NewProcess("models.user.Create", row).Run().(int)
	assert.Greater(t, id, 0)

	// 清空数据
	capsule.Query().Table(Select("user").MetaData.Table.Name).Where("id", id).Delete()
}

func TestProcessUpdate(t *testing.T) {
	id := NewProcess("models.user.Update", 1, maps.MapStr{"balance": 200}).Run()
	assert.Nil(t, id)

	// 恢复数据
	capsule.Query().Table(Select("user").MetaData.Table.Name).Where("id", 1).Update(maps.MapStr{"balance": 0})
}

func TestProcessSave(t *testing.T) {
	row := maps.MapStr{
		"name":     "用户创建",
		"manu_id":  2,
		"type":     "user",
		"idcard":   "23082619820207006X",
		"mobile":   "13900004444",
		"password": "qV@uT1DI",
		"key":      "XZ12MiPp",
		"secret":   "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN",
		"status":   "enabled",
		"extra":    maps.MapStr{"sex": "女"},
	}
	id := NewProcess("models.user.Save", row).Run().(int)
	assert.Greater(t, id, 0)

	// 清空数据
	capsule.Query().Table(Select("user").MetaData.Table.Name).Where("id", id).Delete()
}

func TestProcessDelete(t *testing.T) {
	row := maps.MapStr{
		"name":     "用户创建",
		"manu_id":  2,
		"type":     "user",
		"idcard":   "23082619820207006X",
		"mobile":   "13900004444",
		"password": "qV@uT1DI",
		"key":      "XZ12MiPp",
		"secret":   "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN",
		"status":   "enabled",
		"extra":    maps.MapStr{"sex": "女"},
	}

	user := Select("user")
	id := user.MustSave(row)
	NewProcess("models.user.Delete", id).Run()

	// 清空数据
	capsule.Query().Table(Select("user").MetaData.Table.Name).Where("id", id).Delete()
}

func TestProcessDestroy(t *testing.T) {
	row := maps.MapStr{
		"name":     "用户创建",
		"manu_id":  2,
		"type":     "user",
		"idcard":   "23082619820207006X",
		"mobile":   "13900004444",
		"password": "qV@uT1DI",
		"key":      "XZ12MiPp",
		"secret":   "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN",
		"status":   "enabled",
		"extra":    maps.MapStr{"sex": "女"},
	}

	user := Select("user")
	id := user.MustSave(row)
	NewProcess("models.user.Destroy", id).Run()

	// 清空数据
	capsule.Query().Table(Select("user").MetaData.Table.Name).Where("id", id).Delete()
}

func TestProcessInsert(t *testing.T) {

	content := `{
		"columns": ["user_id", "province", "city", "location"],
		"rows": [
			[4, "北京市", "丰台区", "银海星月9号楼9单元9层1024室"],
			[4, "天津市", "塘沽区", "益海星云7号楼3单元1003室"]
		]
	}`

	payload := map[string]interface{}{}
	err := jsoniter.UnmarshalFromString(content, &payload)
	if err != nil {
		assert.Nil(t, err)
		return
	}

	NewProcess("models.address.Insert", payload["columns"], payload["rows"]).Run()

	// 清理数据
	address := Select("address")
	capsule.Query().Table(address.MetaData.Table.Name).Where("user_id", 4).Delete()
}

func TestProcessUpdateWhere(t *testing.T) {
	effect := NewProcess("models.user.UpdateWhere",
		QueryParam{
			Limit: 1,
			Wheres: []QueryWhere{
				{
					Column: "id",
					Value:  1,
				},
			},
		},
		maps.MapStr{
			"balance": 200,
		},
	).Run().(int)

	user := Select("user")
	row := user.MustFind(1, QueryParam{})

	// 恢复数据
	capsule.Query().Table(user.MetaData.Table.Name).Where("id", 1).Update(maps.MapStr{"balance": 0})
	assert.Equal(t, any.Of(row.Get("balance")).CInt(), 200)
	assert.Equal(t, 1, effect)
}
func TestProcessDeleteWhere(t *testing.T) {

	columns := []string{"name", "manu_id", "type", "idcard", "mobile", "password", "key", "secret", "status"}
	rows := [][]interface{}{
		{"用户创建1", 5, "user", "23082619820207006X", "13900004444", "qV@uT1DI", "XZ12MiP1", "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN", "enabled"},
		{"用户创建2", 5, "user", "33082619820207006X", "13900005555", "qV@uT1DI", "XZ12MiP2", "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN", "enabled"},
		{"用户创建3", 5, "user", "43082619820207006X", "13900006666", "qV@uT1DI", "XZ12MiP3", "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN", "enabled"},
	}

	user := Select("user")
	user.Insert(columns, rows)
	param := QueryParam{Wheres: []QueryWhere{
		{
			Column: "manu_id",
			Value:  5,
		},
	}}
	effect := NewProcess("models.user.DeleteWhere", param).Run().(int)

	// 清理数据
	capsule.Query().Table(user.MetaData.Table.Name).Where("name", "like", "用户创建%").Delete()
	assert.Equal(t, effect, 3)
}

func TestProcessDestroyWhere(t *testing.T) {

	columns := []string{"name", "manu_id", "type", "idcard", "mobile", "password", "key", "secret", "status"}
	rows := [][]interface{}{
		{"用户创建1", 5, "user", "23082619820207006X", "13900004444", "qV@uT1DI", "XZ12MiP1", "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN", "enabled"},
		{"用户创建2", 5, "user", "33082619820207006X", "13900005555", "qV@uT1DI", "XZ12MiP2", "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN", "enabled"},
		{"用户创建3", 5, "user", "43082619820207006X", "13900006666", "qV@uT1DI", "XZ12MiP3", "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN", "enabled"},
	}

	user := Select("user")
	user.Insert(columns, rows)
	param := QueryParam{Wheres: []QueryWhere{
		{
			Column: "manu_id",
			Value:  5,
		},
	}}
	effect := NewProcess("models.user.DestroyWhere", param).Run().(int)

	// 清理数据
	assert.Equal(t, effect, 3)
}