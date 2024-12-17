package database

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func entityToMap(entity interface{}) (string, map[string]interface{}, error) {
	logger.Println("Converting entity to map...")

	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		logger.Printf("Error: Entity must be a struct or a pointer to a struct. Got: %s\n", v.Kind())
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
		logger.Println("Error: Struct must have an unexported field 'DB_NAME'")
		return "", nil, errors.New("struct must have an unexported field 'DB_NAME'")
	}

	logger.Printf("Entity mapped successfully. DB name: %s\n", dbName)
	return dbName, fieldsMap, nil
}

func generateSQLUpdateFromEntityById(entity interface{}) (string, []interface{}, error) {
	logger.Println("Generating SQL UPDATE query from entity...")

	dbName, fieldsMap, err := entityToMap(entity)
	if err != nil {
		logger.Printf("Error in entityToMap: %v\n", err)
		return "", nil, err
	}

	idValue, hasID := fieldsMap["id"]
	if !hasID || idValue == nil {
		logger.Println("Error: The struct must have an 'id' field to use as a condition")
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
		logger.Println("Error: No fields to update")
		return "", nil, errors.New("no fields to update")
	}

	setClause := strings.Join(setClauses, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", dbName, setClause)
	values = append(values, idValue)

	logger.Printf("Generated SQL UPDATE query: %s\n", query)
	return query, values, nil
}

func generateSQLInsertFromEntity(entity interface{}) (string, []interface{}, error) {
	logger.Println("Generating SQL INSERT query from entity...")

	dbName, fieldsMap, err := entityToMap(entity)
	if err != nil {
		logger.Printf("Error in entityToMap: %v\n", err)
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
		logger.Println("Error: No fields to insert")
		return "", nil, errors.New("no fields to insert")
	}

	columnClause := strings.Join(columns, ", ")
	placeholderClause := strings.Join(placeholders, ", ")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", dbName, columnClause, placeholderClause)

	logger.Printf("Generated SQL INSERT query: %s\n", query)
	return query, values, nil
}

func executeGeneratedSQL(query string, values []interface{}) (sql.Result, error) {
	logger.Printf("Executing query: %s\n", query)
	result, err := db.Exec(query, values...)
	if err != nil {
		logger.Printf("Error executing query: %v\n", err)
		return nil, err
	}

	logger.Println("Query executed successfully.")
	return result, nil
}

func InsertEntityById(entity interface{}) (sql.Result, error) {
	logger.Println("Inserting entity by ID...")

	query, args, err := generateSQLInsertFromEntity(entity)
	if err != nil {
		logger.Printf("Error generating SQL insert: %v\n", err)
		return nil, err
	}

	result, err := executeGeneratedSQL(query, args)
	if err != nil {
		logger.Printf("Error executing SQL insert: %v\n", err)
		return nil, err
	}

	logger.Println("Entity inserted successfully.")
	return result, nil
}

func UpdateEntityById(entity interface{}) (sql.Result, error) {
	logger.Println("Updating entity by ID...")

	query, args, err := generateSQLUpdateFromEntityById(entity)
	if err != nil {
		logger.Printf("Error generating SQL update: %v\n", err)
		return nil, err
	}

	result, err := executeGeneratedSQL(query, args)
	if err != nil {
		logger.Printf("Error executing SQL update: %v\n", err)
		return nil, err
	}

	logger.Println("Entity updated successfully.")
	return result, nil
}

func GetAllEntities[T any]() ([]T, error) {
	logger.Println("Fetching all entities from table...")

	dbName, _, err := entityToMap(new(T))
	if err != nil {
		logger.Printf("Error in entityToMap: %v\n", err)
		return nil, err
	}

	query := fmt.Sprintf("SELECT * FROM %s", dbName)

	rows, err := db.Query(query)
	if err != nil {
		logger.Printf("Error executing SELECT query: %v\n", err)
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
			logger.Printf("Error scanning row: %v\n", err)
			return nil, err
		}

		entities = append(entities, *newEntity)
	}

	if err := rows.Err(); err != nil {
		logger.Printf("Error during rows iteration: %v\n", err)
		return nil, err
	}

	logger.Printf("Successfully fetched %d entities from the table %s.\n", len(entities), dbName)
	return entities, nil
}
