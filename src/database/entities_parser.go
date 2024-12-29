package database

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func entityToMap(entity interface{}) (string, map[string]interface{}, error) {
	logger.Infow("Converting entity to map")

	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		logger.Errorw("Entity must be a struct or a pointer to a struct", "actualKind", v.Kind())
		return "", nil, errors.New("entity must be a struct or a pointer to a struct")
	}

	var dbName string
	fieldsMap := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)

		if field.Name == "DB_NAME" {
			dbName = field.Tag.Get("db")
			continue
		}

		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}

		if value.IsValid() && value.Interface() != nil {
			fieldsMap[dbTag] = value.Interface()
		}
	}

	if dbName == "" {
		logger.Errorw("Struct must have an unexported field 'DB_NAME'", "entity", entity)
		return "", nil, errors.New("struct must have an unexported field 'DB_NAME'")
	}

	logger.Infow("Entity mapped successfully", "dbName", dbName)
	return dbName, fieldsMap, nil
}

func generateSQLUpdateFromEntityById(entity interface{}) (string, []interface{}, error) {
	logger.Infow("Generating SQL UPDATE query from entity")

	dbName, fieldsMap, err := entityToMap(entity)
	if err != nil {
		logger.Errorw("Error in entityToMap", "error", err)
		return "", nil, err
	}

	idValue, hasID := fieldsMap["id"]
	if !hasID || idValue == nil {
		logger.Errorw("The struct must have an 'id' field to use as a condition", "entity", entity)
		return "", nil, errors.New("the struct must have an 'id' field to use as a condition")
	}

	var setClauses []string
	var values []interface{}

	for field, value := range fieldsMap {
		if field == "id" || field == "created_at" {
			continue
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = ?", field))
		values = append(values, value)
	}

	if len(setClauses) == 0 {
		logger.Errorw("No fields to update", "entity", entity)
		return "", nil, errors.New("no fields to update")
	}

	setClause := strings.Join(setClauses, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", dbName, setClause)
	values = append(values, idValue)

	logger.Infow("Generated SQL UPDATE query", "query", query)
	return query, values, nil
}

func generateSQLInsertFromEntity(entity interface{}) (string, []interface{}, error) {
	logger.Infow("Generating SQL INSERT query from entity")

	dbName, fieldsMap, err := entityToMap(entity)
	if err != nil {
		logger.Errorw("Error in entityToMap", "error", err)
		return "", nil, err
	}

	var columns []string
	var placeholders []string
	var values []interface{}

	for field, value := range fieldsMap {
		if field == "id" || field == "created_at" || value == "" {
			continue
		}

		columns = append(columns, field)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	if len(columns) == 0 {
		logger.Errorw("No fields to insert", "entity", entity)
		return "", nil, errors.New("no fields to insert")
	}

	columnClause := strings.Join(columns, ", ")
	placeholderClause := strings.Join(placeholders, ", ")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", dbName, columnClause, placeholderClause)

	logger.Infow("Generated SQL INSERT query", "query", query)
	return query, values, nil
}

func executeGeneratedSQL(query string, values []interface{}) (sql.Result, error) {
	logger.Infow("Executing query", "query", query)
	result, err := db.Exec(query, values...)
	if err != nil {
		logger.Errorw("Error executing query", "error", err)
		return nil, err
	}

	logger.Infow("Query executed successfully")
	return result, nil
}

func InsertEntityById(entity interface{}) (sql.Result, error) {
	logger.Infow("Inserting entity by ID")

	query, args, err := generateSQLInsertFromEntity(entity)
	if err != nil {
		logger.Errorw("Error generating SQL insert", "error", err)
		return nil, err
	}

	result, err := executeGeneratedSQL(query, args)
	if err != nil {
		logger.Errorw("Error executing SQL insert", "error", err)
		return nil, err
	}

	logger.Infow("Entity inserted successfully")
	return result, nil
}

func UpdateEntityById(entity interface{}) (sql.Result, error) {
	logger.Infow("Updating entity by ID")

	query, args, err := generateSQLUpdateFromEntityById(entity)
	if err != nil {
		logger.Errorw("Error generating SQL update", "error", err)
		return nil, err
	}

	result, err := executeGeneratedSQL(query, args)
	if err != nil {
		logger.Errorw("Error executing SQL update", "error", err)
		return nil, err
	}

	logger.Infow("Entity updated successfully")
	return result, nil
}

func GetAllEntities[T any]() ([]T, error) {
	logger.Infow("Fetching all entities from table")

	dbName, _, err := entityToMap(new(T))
	if err != nil {
		logger.Errorw("Error in entityToMap", "error", err)
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM " + dbName)
	if err != nil {
		logger.Errorw("Error executing SELECT query", "error", err, "query", dbName)
		return nil, err
	}
	defer rows.Close()

	var entities []T

	for rows.Next() {
		newEntity := new(T)

		var dest []interface{}

		val := reflect.ValueOf(newEntity).Elem()
		numFields := val.NumField()

		for i := 0; i < numFields; i++ {
			field := val.Type().Field(i)

			if field.Name == "DB_NAME" {
				continue
			}

			dest = append(dest, val.Field(i).Addr().Interface())
		}

		err := rows.Scan(dest...)
		if err != nil {
			logger.Errorw("Error scanning row", "error", err)
			return nil, err
		}

		entities = append(entities, *newEntity)
	}

	if err := rows.Err(); err != nil {
		logger.Errorw("Error during rows iteration", "error", err)
		return nil, err
	}

	logger.Infow("Successfully fetched entities", "count", len(entities), "table", dbName)
	return entities, nil
}
