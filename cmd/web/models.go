package main

import (
	"errors"

	"gorm.io/gorm"
)

var db *gorm.DB

// Owner - схема таблицы owners в БД
type Owner struct {
	ID         int    `json:"-"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

// Car - схема таблицы cars в БД
type Car struct {
	RegNum  string `json:"reg_num" gorm:"primaryKey"`
	Mark    string `json:"mark"`
	Model   string `json:"model"`
	Year    int    `json:"year,omitempty"`
	Owner   Owner  `json:"owner"`
	OwnerID int    `json:"-"`
}

// InitConn инициализирует соединение с БД
func InitConn(dbPool *gorm.DB) {
	db = dbPool
}

// MigrateModels в начале программы один раз создаёт таблицы в БД, если их еще нет
func MigrateModels() (bool, error) {
	if !db.Migrator().HasTable(&Car{}) {
		app.InfoLog.Println("Table Car does not exist, creating...")
		if err := db.AutoMigrate(&Car{}); err != nil {
			return false, err
		}
	} else {
		app.InfoLog.Println("Table Car already exists skipping...")
	}
	if !db.Migrator().HasTable(&Owner{}) {
		app.InfoLog.Println("Table Owner does not exist, creating...")
		if err := db.AutoMigrate(&Owner{}); err != nil {
			return false, err
		}
	} else {
		app.InfoLog.Println("Table Owner already exists skipping...")
	}
	return true, nil
}

// DeleteCarByRegNum удаляет записи с первичным ключом равным regNum(string)
func DeleteCarByRegNum(regNum string) (int64, error) {
	res := db.Delete(&Car{}, "reg_num = ?", regNum)
	if res.Error != nil {
		return 0, res.Error
	}

	return res.RowsAffected, nil

}

// InsertCarInfo добавляет объект Car в БД, создавая также Owner, если его нет в БД
func InsertCarInfo(car Car) (string, error) {
	owner, err := FindOwner(car.Owner)
	if err != nil {
		return "", err
	}
	if owner.ID != 0 {
		car.Owner = owner
		car.OwnerID = owner.ID
	}
	app.InfoLog.Println(car.Owner, owner)
	res := db.Create(&car)
	if res.Error != nil {
		return "", res.Error
	}

	return car.RegNum, nil
}

// GetCarsByFilters возвращает все записи Car, поля которых совпадают c полями переданного аргумента
// ВНИМАНИЕ: поля с нулевыми значениями (0, "", false и т.д) не учитываются в запросе
func GetCarsByFilters(car Car) ([]Car, error) {
	var cars []Car
	var res *gorm.DB
	//res := db.Where(car).Preload("Owner").Find(&cars)
	fullname := car.Owner.Name + car.Owner.Surname + car.Owner.Patronymic
	if fullname != "" {
		name_search := ""
		name_search_arr := make([]interface{}, 0, 3)
		if len(car.Owner.Name) > 0 {
			name_search += " AND owners.name = ?"
			name_search_arr = append(name_search_arr, car.Owner.Name)
		}
		if len(car.Owner.Surname) > 0 {
			name_search += " AND owners.surname = ?"
			name_search_arr = append(name_search_arr, car.Owner.Surname)
		}
		if len(car.Owner.Patronymic) > 0 {
			name_search += " AND owners.patronymic = ?"
			name_search_arr = append(name_search_arr, car.Owner.Patronymic)
		}
		res = db.Preload("Owner").Joins(
			"JOIN owners ON owners.id = cars.owner_id"+name_search,
			name_search_arr...).Where(car).Find(&cars)
	} else {
		res = db.Where(car).Preload("Owner").Find(&cars)
	}

	if res.Error != nil {
		return cars, res.Error
	}

	return cars, nil
}

// GetCarsByFilters возвращает объект Car с первичным ключом, равным regNum
// или ошибку: внутренняя ошибка БД или ошибка, если не найдено ни одного совпадения
func GetCarByRegNum(regNum string) (*Car, error) {
	var car Car
	res := db.Preload("Owner").First(&car, "reg_num = ?", regNum) // чтобы избежать SQL-инъекций
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {

		return nil, res.Error
	} else if res.Error != nil {
		return nil, res.Error
	}

	return &car, nil
}

// UpdateCar обновляет объект Car в БД по первичному ключу, равному car.RegNum
func UpdateCar(car Car) error {
	res := db.Model(&car).Updates(car)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

// FindOwner находит владельца по его первичному ключу, равному owner.ID
func FindOwner(owner Owner) (Owner, error) {
	var foundOwner Owner

	res := db.Where(owner).First(&foundOwner)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) { // владельца еще  нет, ну и ладно: это не ошибка
		return foundOwner, nil
	} else if res.Error != nil {
		return foundOwner, res.Error
	} else {
		return foundOwner, nil
	}

}
