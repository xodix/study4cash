package routes

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"study4cash/DB/models"
	"study4cash/auth"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RouteData(prefix string, e *echo.Echo, db *gorm.DB) {
	dataRouter := e.Group(prefix)
	dataRouter.Use(auth.JWTMiddleware)

	dataRouter.GET("/attending", func(c *echo.Context) error { return GetAttending(c, db) })
	dataRouter.GET("/graduating", func(c *echo.Context) error { return GetGraduating(c, db) })
	dataRouter.GET("/averageWage", func(c *echo.Context) error { return GetAverageWage(c, db) })

	dataRouter.PUT("/attending/XML", func(c *echo.Context) error { return SetAttendingXML(c, db) })
	dataRouter.PUT("/graduating/XML", func(c *echo.Context) error { return SetGraduatingXML(c, db) })
	dataRouter.PUT("/averageWage/XML", func(c *echo.Context) error { return SetAverageWageXML(c, db) })

	dataRouter.PUT("/attending/JSON", func(c *echo.Context) error { return SetAttendingJSON(c, db) })
	dataRouter.PUT("/graduating/JSON", func(c *echo.Context) error { return SetGraduatingJSON(c, db) })
	dataRouter.PUT("/averageWage/JSON", func(c *echo.Context) error { return SetAverageWageJSON(c, db) })
}

func GetAttending(c *echo.Context, db *gorm.DB) error {
	attending, err := gorm.G[models.Attending](db).Order("year asc").Find(c.Request().Context())
	if err != nil {
		return c.String(500, err.Error())
	}

	return c.JSON(200, attending)
}

func GetGraduating(c *echo.Context, db *gorm.DB) error {
	graduating, err := gorm.G[models.Graduating](db).Order("year asc").Find(c.Request().Context())
	if err != nil {
		return c.String(500, err.Error())
	}

	return c.JSON(200, graduating)
}

func GetAverageWage(c *echo.Context, db *gorm.DB) error {
	averageWages, err := gorm.G[models.AverageWage](db).Order("year asc").Find(c.Request().Context())
	if err != nil {
		return c.String(500, err.Error())
	}

	return c.JSON(200, averageWages)
}

func parseXML(file multipart.File) ([]map[string]string, error) {
	decoder := xml.NewDecoder(file)

	var results []map[string]string
	var current map[string]string
	var currentTag string

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			currentTag = t.Name.Local
			if currentTag != "root" {
				current = make(map[string]string)
			}
		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" && current != nil {
				current[currentTag] = text
			}
		case xml.EndElement:
			if t.Name.Local != "root" && current != nil {
				results = append(results, current)
				current = nil
			}
		}
	}

	return results, nil
}

func parseJSON(file multipart.File) ([]map[string]string, error) {
	var payload []map[string]string
	err := json.NewDecoder(file).Decode(&payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func SetAttendingXML(c *echo.Context, db *gorm.DB) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.String(400, err.Error())
	}
	if fileHeader == nil {
		return c.String(400, "No file uploaded")
	}
	if fileHeader.Size > 50*1024 {
		return c.String(400, "File too large")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.String(500, err.Error())
	}
	defer file.Close()

	parsed, err := parseXML(file)
	if err != nil {
		return c.String(400, err.Error())
	}

	records := make([]models.Attending, 0, len(parsed)*5)
	voivodeship := ""
	for _, elem := range parsed {
		if elem["Nazwa"] != "" {
			voivodeship = elem["Nazwa"]
		}

		for key, value := range elem {
			if yearStr, ok := strings.CutPrefix(key, "Y"); ok {
				year, err := strconv.Atoi(yearStr)
				if err != nil {
					return c.String(500, err.Error())
				}
				attending, err := strconv.Atoi(value)
				if err != nil {
					return c.String(500, err.Error())
				}

				records = append(records, models.Attending{
					Voivodeship:       voivodeship,
					Year:              year,
					StudentsAttending: attending,
				})
			}
		}
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		stmt := &gorm.Statement{DB: tx}
		if err := stmt.Parse(&models.Attending{}); err != nil {
			return err
		}
		tableName := stmt.Schema.Table
		if err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error; err != nil {
			return err
		}

		if err := tx.Create(&records).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return c.String(500, err.Error())
	}

	return c.String(200, "Added attending students to the database")
}

func SetAverageWageXML(c *echo.Context, db *gorm.DB) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.String(500, err.Error())
	}
	if fileHeader == nil {
		return c.String(400, "No file uploaded")
	}
	if fileHeader.Size > 50*1024 {
		return c.String(400, "File too large")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.String(500, err.Error())
	}
	defer file.Close()

	parsed, err := parseXML(file)
	if err != nil {
		return c.String(500, err.Error())
	}

	records := make([]models.AverageWage, 0, len(parsed)*5)
	voivodeship := ""
	for _, elem := range parsed {
		if elem["Nazwa"] != "" {
			voivodeship = elem["Nazwa"]
		}
		for key, value := range elem {
			if yearStr, ok := strings.CutPrefix(key, "Y"); ok {
				year, err := strconv.Atoi(yearStr)
				if err != nil {
					return c.String(500, err.Error())
				}
				averageWage, err := strconv.ParseFloat(strings.Replace(value, ",", ".", 1), 64)
				if err != nil {
					return c.String(500, err.Error())
				}

				records = append(records, models.AverageWage{
					Voivodeship: voivodeship,
					Year:        year,
					AverageWage: averageWage,
				})
			}
		}
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		stmt := &gorm.Statement{DB: tx}
		if err := stmt.Parse(&models.AverageWage{}); err != nil {
			return err
		}
		tableName := stmt.Schema.Table
		if err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error; err != nil {
			return err
		}

		if err := tx.Create(&records).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	return c.String(200, "Added average wages students to the database")
}

func SetGraduatingXML(c *echo.Context, db *gorm.DB) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.String(500, err.Error())
	}
	if fileHeader == nil {
		return c.String(400, "No file uploaded")
	}
	if fileHeader.Size > 50*1024 {
		return c.String(400, "File too large")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.String(500, err.Error())
	}
	defer file.Close()

	parsed, err := parseXML(file)
	if err != nil {
		return c.String(500, err.Error())
	}

	records := make([]models.Graduating, 0, len(parsed)*5)
	voivodeship := ""
	for _, elem := range parsed {
		if elem["Nazwa"] != "" {
			voivodeship = elem["Nazwa"]
		}
		for key, value := range elem {
			if strings.HasPrefix(key, "Y") {
				year, err := strconv.Atoi(strings.TrimPrefix(key, "Y"))
				if err != nil {
					return c.String(500, err.Error())
				}
				graduating, err := strconv.Atoi(value)
				if err != nil {
					return c.String(500, err.Error())
				}

				records = append(records, models.Graduating{
					Voivodeship:        voivodeship,
					Year:               year,
					StudentsGraduating: graduating,
				})
			}
		}
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		stmt := &gorm.Statement{DB: tx}
		if err := stmt.Parse(&models.Graduating{}); err != nil {
			return err
		}
		tableName := stmt.Schema.Table
		if err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error; err != nil {
			return err
		}

		if err := tx.Create(&records).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	return c.String(200, "Added graduating students to the database")
}

func SetAttendingJSON(c *echo.Context, db *gorm.DB) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.String(400, err.Error())
	}
	if fileHeader == nil {
		return c.String(400, "No file uploaded")
	}
	if fileHeader.Size > 50*1024 {
		return c.String(400, "File too large")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.String(500, err.Error())
	}
	defer file.Close()

	parsed, err := parseJSON(file)
	if err != nil {
		return c.String(400, err.Error())
	}

	records := make([]models.Attending, 0, len(parsed)*5)
	voivodeship := ""
	for _, elem := range parsed {
		if elem["Nazwa"] != "" {
			voivodeship = elem["Nazwa"]
		}

		for key, value := range elem {
			if yearStr, ok := strings.CutPrefix(key, "Y"); ok {
				year, err := strconv.Atoi(yearStr)
				if err != nil {
					return c.String(500, err.Error())
				}
				attending, err := strconv.Atoi(value)
				if err != nil {
					return c.String(500, err.Error())
				}

				records = append(records, models.Attending{
					Voivodeship:       voivodeship,
					Year:              year,
					StudentsAttending: attending,
				})
			}
		}
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		stmt := &gorm.Statement{DB: tx}
		if err := stmt.Parse(&models.Attending{}); err != nil {
			return err
		}
		tableName := stmt.Schema.Table
		if err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error; err != nil {
			return err
		}

		if err := tx.Create(&records).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return c.String(500, err.Error())
	}

	return c.String(200, "Added attending students to the database")
}

func SetAverageWageJSON(c *echo.Context, db *gorm.DB) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.String(500, err.Error())
	}
	if fileHeader == nil {
		return c.String(400, "No file uploaded")
	}
	if fileHeader.Size > 50*1024 {
		return c.String(400, "File too large")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.String(500, err.Error())
	}
	defer file.Close()

	parsed, err := parseJSON(file)
	if err != nil {
		return c.String(500, err.Error())
	}

	records := make([]models.AverageWage, 0, len(parsed)*5)
	voivodeship := ""
	for _, elem := range parsed {
		if elem["Nazwa"] != "" {
			voivodeship = elem["Nazwa"]
		}
		for key, value := range elem {
			if yearStr, ok := strings.CutPrefix(key, "Y"); ok {
				year, err := strconv.Atoi(yearStr)
				if err != nil {
					return c.String(500, err.Error())
				}
				averageWage, err := strconv.ParseFloat(strings.Replace(value, ",", ".", 1), 64)
				if err != nil {
					return c.String(500, err.Error())
				}

				records = append(records, models.AverageWage{
					Voivodeship: voivodeship,
					Year:        year,
					AverageWage: averageWage,
				})
			}
		}
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		stmt := &gorm.Statement{DB: tx}
		if err := stmt.Parse(&models.AverageWage{}); err != nil {
			return err
		}
		tableName := stmt.Schema.Table
		if err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error; err != nil {
			return err
		}

		if err := tx.Create(&records).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	return c.String(200, "Added average wages students to the database")
}

func SetGraduatingJSON(c *echo.Context, db *gorm.DB) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.String(500, err.Error())
	}
	if fileHeader == nil {
		return c.String(400, "No file uploaded")
	}
	if fileHeader.Size > 50*1024 {
		return c.String(400, "File too large")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.String(500, err.Error())
	}
	defer file.Close()

	parsed, err := parseJSON(file)
	if err != nil {
		return c.String(500, err.Error())
	}

	records := make([]models.Graduating, 0, len(parsed)*5)
	voivodeship := ""
	for _, elem := range parsed {
		if elem["Nazwa"] != "" {
			voivodeship = elem["Nazwa"]
		}
		for key, value := range elem {
			if strings.HasPrefix(key, "Y") {
				year, err := strconv.Atoi(strings.TrimPrefix(key, "Y"))
				if err != nil {
					return c.String(500, err.Error())
				}
				graduating, err := strconv.Atoi(value)
				if err != nil {
					return c.String(500, err.Error())
				}

				records = append(records, models.Graduating{
					Voivodeship:        voivodeship,
					Year:               year,
					StudentsGraduating: graduating,
				})
			}
		}
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		stmt := &gorm.Statement{DB: tx}
		if err := stmt.Parse(&models.Graduating{}); err != nil {
			return err
		}
		tableName := stmt.Schema.Table
		if err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error; err != nil {
			return err
		}

		if err := tx.Create(&records).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	return c.String(200, "Added graduating students to the database")
}
