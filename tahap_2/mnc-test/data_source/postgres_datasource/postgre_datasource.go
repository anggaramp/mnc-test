package postgres_datasource

import (
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type PostgresDatasource struct {
	Client *gorm.DB
}

type QueryOptionOrder struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type QueryOptionPage struct {
	PageIndex uint `json:"pageIndex"`
	PageSize  uint `json:"pageSize"`
}

type QueryOption struct {
	Filter       map[string]interface{} `json:"filter"`
	Order        []QueryOptionOrder     `json:"order"`
	Page         *QueryOptionPage       `json:"page"`
	DisableCount bool                   `json:"disableCount"`
}

func NewPostgresDatasource(client *gorm.DB) *PostgresDatasource {
	return &PostgresDatasource{
		Client: client,
	}
}

func (r *PostgresDatasource) GetDB(tx *gorm.DB) (db *gorm.DB) {
	if nil != tx {
		db = tx
	} else {
		db = r.Client
	}

	return
}

func (r *PostgresDatasource) Create(tx *gorm.DB, data interface{}) (err error) {
	db := r.GetDB(tx).Model(data)

	err = db.Clauses(clause.Returning{}).Create(data).Error

	return
}

func (r *PostgresDatasource) Delete(tx *gorm.DB, data interface{}, conditions []map[string]interface{}) (err error) {
	db := r.GetDB(tx)

	for _, condition := range conditions {
		var key, operator, value string

		if _value, ok := condition["key"].(string); ok {
			key = _value
		}
		if _value, ok := condition["operator"].(string); ok {
			operator = _value
		}
		if _value, ok := condition["value"].(string); ok {
			value = _value
		}

		db = db.Where(fmt.Sprintf("%s %s '%v'", key, operator, value))
	}

	err = db.Delete(data).Error

	return
}

func (r *PostgresDatasource) UpdateData(tx *gorm.DB, model string, id uint64, key string, value interface{}) (err error) {
	db := r.GetDB(tx)

	valueString := fmt.Sprintf("\"%v\"", value)
	err = db.Exec(fmt.Sprintf("UPDATE %s SET data = jsonb_set(data, ?, ?), updated_timestamp=NOW() WHERE id=?", model), pq.Array(&[]string{key}), valueString, id).Error

	return
}

func (r *PostgresDatasource) UpdateDataNotString(tx *gorm.DB, model string, id uint64, key string, value interface{}) (err error) {
	db := r.GetDB(tx)

	err = db.Exec(fmt.Sprintf("UPDATE %s SET data = jsonb_set(data, ?, ?), updated_timestamp=NOW() WHERE id=?", model), pq.Array(&[]string{key}), value, id).Error

	return
}

func (r *PostgresDatasource) UpdateRow(tx *gorm.DB, model string, id string, updates map[string]string) (err error) {
	db := r.GetDB(tx)
	var updateString string
	for field, value := range updates {
		updateString += fmt.Sprintf("%s = %s, ", field, value)
	}

	err = db.Exec(fmt.Sprintf("UPDATE %s SET %s updated_timestamp=NOW() WHERE ulid='?'", model, updateString), id).Error

	return
}

func (r *PostgresDatasource) Update(tx *gorm.DB, columnName, id string, data interface{}, properties map[string]interface{}) (err error) {
	db := r.GetDB(tx).Model(data)

	err = db.Clauses(clause.Returning{}).Where(fmt.Sprintf("%s=?", columnName), id).Updates(properties).Error

	return
}

func (r *PostgresDatasource) Get(tx *gorm.DB, entity interface{}, conditions []map[string]interface{}) (resp interface{}, err error) {
	db := r.GetDB(tx).Model(entity)

	for _, condition := range conditions {
		var key, operator, value string

		if _value, ok := condition["key"].(string); ok {
			key = _value
		}
		if _value, ok := condition["operator"].(string); ok {
			operator = _value
		}
		if _value, ok := condition["value"].(string); ok {
			value = _value
		}

		db = db.Where(fmt.Sprintf("%s %s '%s'", key, operator, value))
	}

	err = db.First(entity).Error

	resp = entity

	return
}

func (r *PostgresDatasource) GetV2(tx *gorm.DB, entity interface{}, conditions []map[string]interface{}) (err error) {
	db := r.GetDB(tx)

	for _, condition := range conditions {
		var key, operator, value string

		if _value, ok := condition["key"].(string); ok {
			key = _value
		}
		if _value, ok := condition["operator"].(string); ok {
			operator = _value
		}
		if _value, ok := condition["value"].(string); ok {
			value = _value
		}

		db = db.Where(fmt.Sprintf("%s %s '%s'", key, operator, value))
	}

	err = db.First(entity).Error

	return
}

func (r *PostgresDatasource) GetList(tx *gorm.DB, entity interface{}, data interface{}, queryOption QueryOption) (total int64, resp interface{}, err error) {
	db := r.GetDB(tx).Model(entity)

	// filter
	err = r.processQueryOptionFilter(db, queryOption.Filter)
	if nil != err {
		return
	}

	// order
	for _, order := range queryOption.Order {
		orderQuery := fmt.Sprintf("%s %s", order.Field, order.Direction)

		db = db.Order(orderQuery)
	}

	if !queryOption.DisableCount {
		err = db.Count(&total).Error
	}

	if nil != err {
		return
	}

	if len(queryOption.Order) < 1 {
		db = db.Order("created_date desc")
	}
	// page
	if nil != queryOption.Page {

		if queryOption.Page.PageIndex > 0 {
			db = db.Offset(int(queryOption.Page.PageIndex))
		}
		if queryOption.Page.PageSize > 0 {
			db = db.Limit(int(queryOption.Page.PageSize))
		}

	}

	err = db.Find(data).Error

	resp = data

	return
}

func (r *PostgresDatasource) GetListWithRaw(tx *gorm.DB, data interface{}, query *string) (resp interface{}, err error) {
	db := r.GetDB(tx)

	if nil != query {
		db = db.Raw(*query)
	}

	err = db.Find(data).Error

	resp = data

	return
}

func (r *PostgresDatasource) Query(tx *gorm.DB, data interface{}, query *string) (resp interface{}, err error) {
	db := r.GetDB(tx)

	if nil != query {
		db = db.Raw(*query)
	}

	err = db.Scan(data).Error

	resp = data

	return
}

func (r *PostgresDatasource) processQueryOptionFilter(db *gorm.DB, filter map[string]interface{}) (err error) {
	if nil == filter {
		return
	}

	for key, value := range filter {
		if 0 == strings.Index(key, "orSet") {
			queries, args := r.getWhereOrSet(value.([]interface{}))
			*db = *db.Where(strings.Join(queries, " OR "), args...)
		} else {
			query, arg := r.getWhere(value.(map[string]interface{}))
			*db = *db.Where(query, arg)
		}
	}

	return
}

func (r *PostgresDatasource) getWhereOrSet(items []interface{}) (queries []string, args []interface{}) {
	for _, item := range items {
		query, arg := r.getWhere(item.(map[string]interface{}))
		queries = append(queries, query)
		args = append(args, arg)
	}

	return
}

func (r *PostgresDatasource) getWhere(item map[string]interface{}) (query string, arg interface{}) {

	field := item["field"].(string)
	searchType := item["searchType"].(string)
	match := item["match"].(string)
	keyword := item["keyword"]

	field = strings.ReplaceAll(field, "\"", "'")

	//log.Debug("getWhere %s %s %s %v", field, searchType, match, keyword)

	switch searchType {
	case "text":
		switch match {
		case "contain":
			query = fmt.Sprintf("LOWER(%s) LIKE ?", field)
			arg = fmt.Sprintf("%%%s%%", strings.ToLower(keyword.(string)))
		case "startWith":
			query = fmt.Sprintf("LOWER(%s) LIKE ?", field)
			arg = fmt.Sprintf("%s%%", strings.ToLower(keyword.(string)))
		case "endWith":
			query = fmt.Sprintf("LOWER(%s) LIKE ?", field)
			arg = fmt.Sprintf("%%%s", strings.ToLower(keyword.(string)))
		case "exact":
			query = fmt.Sprintf("LOWER(%s) = ?", field)
			arg = strings.ToLower(keyword.(string))
		case "notEqual":
			query = fmt.Sprintf("LOWER(%s) != ?", field)
			arg = strings.ToLower(keyword.(string))
		case "gt":
			query = fmt.Sprintf("%s > ?", field)
			arg = keyword
		case "gte":
			query = fmt.Sprintf("%s >= ?", field)
			arg = keyword
		case "lt":
			query = fmt.Sprintf("%s < ?", field)
			arg = keyword
		case "lte":
			query = fmt.Sprintf("%s <= ?", field)
			arg = keyword
		}
	case "number":
		switch match {
		case "exact":
			query = fmt.Sprintf("%s = ?", field)
			arg = keyword.(uint64)
		case "notEqual":
			query = fmt.Sprintf("%s != ?", field)
			arg = keyword.(uint64)
		}
	case "list":
		switch match {
		case "contain", "overlap":
			if s, ok := keyword.(string); ok {
				query = fmt.Sprintf("%s = ?", field)
				arg = fmt.Sprintf("%v", s)
			} else if ss, ok := keyword.([]interface{}); ok {
				var symbol string
				switch match {
				case "contain":
					symbol = "@>"
				case "overlap":
					symbol = "&&"
				}

				query = fmt.Sprintf("%s %s ?", field, symbol)
				var ssString []string
				for _, s := range ss {
					ssString = append(ssString, s.(string))
				}
				arg = pq.Array(ssString)
			}
		}
	case "date":

		dateTypeName := "day"

		dateType, ok := item["dateType"]
		if ok {
			dateTypeName = dateType.(string)
		}

		switch match {
		case "gt":
			query = fmt.Sprintf("date_trunc('%s', %s) > ?", dateTypeName, field)
			arg = keyword
		case "gte":
			query = fmt.Sprintf("date_trunc('%s', %s) >= ?", dateTypeName, field)
			arg = keyword
		case "lt":
			query = fmt.Sprintf("date_trunc('%s', %s) < ?", dateTypeName, field)
			arg = keyword
		case "lte":
			query = fmt.Sprintf("date_trunc('%s', %s) <= ?", dateTypeName, field)
			arg = keyword
		}
	case "bool":
		switch match {
		case "exact":
			query = fmt.Sprintf("%s = ?", field)
			arg = keyword.(bool)
		}
	}

	//log.Debug("getWhere %s %v", query, arg)

	return
}
