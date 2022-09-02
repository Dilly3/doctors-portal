package database

import (
	_ "database/sql"
	"fmt"
	"os"

	"github.com/dilly3/doctors-portal/models"
	"github.com/dilly3/doctors-portal/utils"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresDb struct {
	DB *gorm.DB
}

func SetupDB() *gorm.DB {
	utils.LoadEnv()
	password := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB")
	root := os.Getenv("DB_ROOTS")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", root, password, dbDatabase, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database: ", err)
		os.Exit(1)

	}
	err = db.AutoMigrate(&models.Doctor{}, &models.Patient{}, &models.Appointment{})
	if err != nil {
		fmt.Printf("failed to migrate: %v", err)
		os.Exit(1)
	}
	return db
}
func NewDB() DataStore {
	return &PostgresDb{
		DB: SetupDB(),
	}
}

func (p PostgresDb) FindUserByEmailandUserName(email string, username string) (*models.Patient, error) {
	user := &models.Patient{}
	// SELECT * from Patient table where email = ?
	err := p.DB.Where("email = ?", email).First(user).Error
	if err != nil {
		err = p.DB.Where("username = ?", username).First(user).Error
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (p PostgresDb) CreatePatientInTable(user models.Patient) {
	//MYSQL: INSERT IN patient TABLE...
	if err := p.DB.Create(&user).Error; err != nil {
		fmt.Println(err)
		return
	}
}

func (p PostgresDb) FindDocByEmailandUserName(email string, username string) (*models.Doctor, error) {
	user := &models.Doctor{}
	// SELECT * from Doctor table where email = ?
	err := p.DB.Where("email = ?", "username = ?", email, username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p PostgresDb) FindDoctorByID(id string) *models.Doctor {
	user := &models.Doctor{}
	err := p.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil
	}
	return user
}

func (p PostgresDb) CreateDocInTable(user models.Doctor) {
	if err := p.DB.Create(&user).Error; err != nil {
		fmt.Println(err)
		return
	}
}

func (p PostgresDb) Authenticate(username, password string) (*models.Doctor, error) {
	//Retrieve the username and hashed password associated with the given username.
	//If matching username exists, return the ErrMismatchedHashAndPassword error.
	user := &models.Doctor{}
	err := p.DB.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	// Check whether the hashed password and plain-text password provided match
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, bcrypt.ErrMismatchedHashAndPassword
	}
	return user, nil
}
func (p PostgresDb) UpdateDoctorWorkingHours(docID string, starttime string, closetime string) (*models.Doctor, error) {
	st := "UPDATE doctors SET string_start = ?, string_close = ? WHERE id = ?"
	p.DB.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&models.Doctor{}, docID)
	pgdb := p.DB.Exec(st, starttime, closetime, docID).Error
	if pgdb != nil {
		return nil, pgdb
	}
	return p.FindDoctorByID(docID), nil
}
func (p PostgresDb) UpdateDoctorWorkingHoursInt(docID string, starttime int, closetime int) (*models.Doctor, error) {
	st := "UPDATE doctors SET start_time = ?, close_time = ? WHERE id = ?"
	p.DB.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&models.Doctor{}, docID)
	pgdb := p.DB.Exec(st, starttime, closetime, docID).Error
	if pgdb != nil {
		return nil, pgdb
	}
	return p.FindDoctorByID(docID), nil
}

func (p PostgresDb) FindDoctorByUsername(username string) (*models.Doctor, error) {
	user := &models.Doctor{}
	err := p.DB.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p PostgresDb) AuthenticatePatient(username, password string) (*models.Patient, error) {
	//Retrieve the username and hashed password associated with the given username.
	//If matching username exists, return the ErrMismatchedHashAndPassword error.
	user := &models.Patient{}
	err := p.DB.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	// Check whether the hashed password and plain-text password provided match
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, bcrypt.ErrMismatchedHashAndPassword
	}
	return user, nil
}
func (p PostgresDb) FindPatientByUsername(username string) (*models.Patient, error) {
	user := &models.Patient{}

	err := p.DB.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p PostgresDb) GetAllDoctors() []models.Doctor {
	var users []models.Doctor
	PgDB, _ := p.DB.DB()

	st := "SELECT * FROM doctors"
	rows, err := PgDB.Query(st)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for rows.Next() {
		var r models.Doctor
		err := rows.Scan(&r.ID, &r.Username, &r.Name, &r.Age, &r.Email, &r.PasswordHash, &r.Specialty, &r.StartTime, &r.CloseTime, &r.StringStart, &r.StringClose)
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, r)
	}

	return users
}

func (p PostgresDb) CreateAppointmentInTable(user models.Appointment) {
	//MYSQL: INSERT IN patient TABLE...
	if err := p.DB.Create(&user).Error; err != nil {
		fmt.Println(err)
		return
	}
}

func (p PostgresDb) FindPatientAppointmentsByPatientID(id string) []models.Appointment {
	appointments := []models.Appointment{}

	_ = p.DB.Model(models.Appointment{}).Where("patient_id = ?", id).Find(&appointments)

	return appointments
}

func (p PostgresDb) DeleteAppointmentbyID(id string) {
	p.DB.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&models.Appointment{}, id)
	p.DB.Where("id = ?", id).Delete(&models.Appointment{})

}
func (p PostgresDb) CheckAppointmentIsValidWithDoctorID(id string, time string) bool {
	appointments := models.Appointment{}
	_ = p.DB.Model(models.Appointment{}).Where("doctor_id = ? AND appointment_hour = ?", id, time).Find(&appointments)
	return len(appointments.ID) < 1

}
func (p PostgresDb) FindDoctorAppointmentsByDoctorID(id string) []models.Appointment {
	appointments := []models.Appointment{}

	_ = p.DB.Model(models.Appointment{}).Where("doctor_id = ?", id).Find(&appointments)

	return appointments
}
