package db

type User struct {
    Id       int   `json:"id" gorm:"praimaly_key"`
    Name     string `json:"name"`
    Email	 string `json:"email"`
    Password string `json:"password"`
}

func CreateUser(user *User) {
    db.Create(user)
}

func FindUser(u *User) User {
    var user User
    db.Where(u).First(&user)
    return user
}