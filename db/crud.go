/**
 * @Time: 2022/3/5 20:16
 * @Author: yt.yin
 */

package db

import (
	"encoding/json"
	"errors"
	"github.com/goworkeryyt/go-core/global"
    "github.com/goworkeryyt/go-toolbox/page"
	"reflect"
	"strconv"
	"time"
)

// FindPage 分页查询 v-空对象指针
func FindPage(v interface{},rows interface{} ,pageInfo *page.PageInfo) (err error,pageBean *page.PageBean) {
	if pageInfo == nil {
		return errors.New("入参pageInfo不能为空指针 "),nil
	}
	pageBean = &page.PageBean{Page: pageInfo.Current,PageSize: pageInfo.RowCount }
	var total int64
	db := global.DB.Model(v)
	typeOf := reflect.TypeOf(v)
	if typeOf.Kind() == reflect.String {
		db = global.DB.Table(asString(v))
	}
	andCons := pageInfo.AndParams
	orCons := pageInfo.OrParams
	orderStr := pageInfo.OrderStr
	if andCons != nil && len(andCons) > 0 {
		for k, v := range andCons {
			db = db.Where(k, v)
		}
	}
	if orCons != nil && len(orCons) > 0 {
		for k, v := range orCons {
			db = db.Or(k, v)
		}
	}
	db.Count(&total)
	if len(orderStr) > 0 {
		err = db.Limit(pageBean.PageSize).Offset((pageBean.Page - 1) * pageBean.PageSize).Order(orderStr).Find(rows).Error
	}else {
		err = db.Limit(pageBean.PageSize).Offset((pageBean.Page - 1) * pageBean.PageSize).Find(rows).Error
	}
	pageBean.Rows = rows
	pageBean.Total = total
	return
}

//  其他类型转String
func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		return time.Time.Format(v, "2006-01-02 15:04:05")
	case bool:
		return strconv.FormatBool(v)
	default:
		{
			b, _ := json.Marshal(v)
			return string(b)
		}
	}
}

