package models

import (
	"bufio"
	"github.com/fabric-app/pkg/util/hash"
	"github.com/fabric-app/pkg/util/rand"
	"github.com/jinzhu/gorm"
	"os"
	"path"
	"time"
)

const HEADER_IMAGE_PATH = "./test/header/images/"

// user table structure
type User struct {
	Model
	Username string `json:"username"`
	Identity string `json:"identity"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
	CaSecure string `json:"ca_secure"`
	Secret   string `json:"secret"`
	Address  string `json:"address"`
	Header   string `json:"header"`
}

// create user info
func NewUser(user *User) (int, error) {
	err := db.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// del user
func DelUser(user *User) (int, error) {
	err := db.Delete(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// find user by id
func FindUserById(id int) (User, error) {
	var user User
	err := db.First(&user, id).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, err
}

// by name
func FindUserByName(name string) (User, error) {
	var user User
	err := db.Where("username = ?", name).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, err
}

// find user by email
func FindUserByEmail(e string) (User, error) {
	var user User
	err := db.Where("email = ?", e).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, err
}

// update user info
func UpdateUserInfo(newUser *User) (int, error) {
	var oldUser User
	err := db.First(&oldUser, newUser.ID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	oldUser.Username = newUser.Username
	oldUser.Email = newUser.Email
	oldUser.Address = newUser.Address
	oldUser.Phone = newUser.Phone
	oldUser.ModifiedOn = int(time.Now().Unix())
	err = db.Save(oldUser).Error
	if err != nil {
		return 0, nil
	}
	return oldUser.ID, nil
}

// update user secret
func UpdateUserSecret(user *User) (int, error) {
	var secretString string
	for {
		secretString = rand.RandStringBytesMaskImprSrcUnsafe(5)
		if user.Secret != secretString {
			break
		}
	}
	db.First(user)
	user.Secret = secretString
	err := db.Save(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// update user header
func UpdateUserheader(name, header string) (int, error) {
	var user User
	db.Where("username = ?", name).First(&user)
	user.Header = header
	err := db.Save(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// update user password
func UpdateUserNewPassword(user *User, newPassword string) (int, error) {
	var secretString string
	for {
		secretString = rand.RandStringBytesMaskImprSrcUnsafe(5)
		if user.Secret != secretString {
			break
		}
	}
	db.First(user)
	user.Secret = secretString
	user.Password = hash.EncodeMD5(newPassword)
	err := db.Save(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// store user header
func SaveUserHeader(username string, data []byte) (int, error) {
	path := path.Join(HEADER_IMAGE_PATH, username, ".jpg")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		return 0, err
	}
	writer := bufio.NewWriter(file)
	count, err := writer.Write(data)
	if err != nil {
		return 0, err
	}
	writer.Flush()
	return count, err
}
