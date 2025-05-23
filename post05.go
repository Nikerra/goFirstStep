package post05

import(
    "database/sql"
    "errors"
    "fmt"
    "strings"

    _ "github.com/lib/pq"
)

type Userdata struct {
    ID int
    Username string
    Name string
    Surname string
    Description string
}

var (
    Hostname = ""
    Port = 2345
    Username = ""
    Password = ""
    Database = ""
)

func openConnection() (*sql.DB, error) {
    //string connection
    conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)

    //open data base
    db, err := sql.Open("postgres", conn)
    if err != nil {
        return nil, err
    }
    return db, nil
}

// Function return ID username
// -1 if user don't find
func exists(username string) int {
    username = strings.ToLower(username)

    db, err := openConnection()
    if err != nil {
        fmt.Println(err)
        return -1
    }
    defer db.Close()

    userID := -1
    statement := fmt.Sprintf(`select "id" from "users" where username = '%s'`, username)
    rows, err := db.Query(statement)

    for rows.Next() {
        var id int
        err = rows.Scan(&id)
        if err != nil {
            fmt.Println("Scan", err)
            return -1
        }
        userID = id
    }
    defer rows.Close()
    return userID
}

// AddUser add new user on database
// return new ID user
// -1 if error
func AddUser(ud Userdata) int {
    ud.Username = strings.ToLower(ud.Username)

    db, err := openConnection()
    if err != nil {
        fmt.Println(err)
        return -1
    }
    defer db.Close()

    userID := exists(ud.Username)
    if userID != -1 {
        fmt.Println("User already exists:", Username)
        return -1
    }

    insertStatement := `insert into "users" ("username") values ($1)`
    _, err = db.Exec(insertStatement, ud.Username)
    if err != nil {
        fmt.Println(err)
        return -1
    }

    userID = exists(ud.Username)
    if userID == -1 {
        return userID
    }

    insertStatement = `insert into "userdata" ("userid", "name", "surname", "description") values ($1, $2, $3, $4)`

    _, err = db.Exec(insertStatement, userID, ud.Name, ud.Surname, ud.Description)
    if err != nil {
        fmt.Println("db.Exec()", err)
        return -1
    }

    return userID
}

//deleteUser delete existing user
func DeleteUser(id int) error {
    db, err := openConnection()
    if err != nil {
        return err
    }
    defer db.Close()

    statement := fmt.Sprintf(`select "username" from "users" where id = %d`, id)
    rows, err := db.Query(statement)

    var username string
    for rows.Next() {
        err = rows.Scan(&username)
        if err != nil {
            return err
        }
    }
    defer rows.Close()

    if exists(username) != id {
        return fmt.Errorf("User with id %d does not exist", id)
    }

    deleteStatement := `delete from "userdata" where userid=$1`
    _, err = db.Exec(deleteStatement, id)
    if err != nil {
        return err
    }

    deleteStatement = `delete from "users" where id=$1`
    _, err = db.Exec(deleteStatement, id)
    if err != nil {
        return err
    }

    return nil
}

func ListUsers() ([]Userdata, error) {
    Data := []Userdata{}
    db, err := openConnection()
    if err != nil {
        return Data, err
    }
    defer db.Close()

    rows, err := db.Query(`select
    "id", "username", "name", "surname", "description"
    from "users", "userdata"
    where users.id = userdata.userid`)
    if err != nil {
        return Data, err
    }

    for rows.Next() {
        var id int
        var username string
        var name string
        var surname string
        var description string
        err = rows.Scan(&id, &username, &name, &surname, &description)
        temp := Userdata {ID: id, Username: username, Name: name, Surname: surname, Description: description}

        Data = append(Data, temp)
        if err != nil {
            return Data, err
        }
    }
    defer rows.Close()
    return Data, nil
}

//UpdateUser update data exist user
func UpdateUser(ud Userdata) error {
    db, err := openConnection()
    if err != nil {
        return err
    }
    defer db.Close()

    userID := exists(ud.Username)
    if userID == -1 {
        return errors.New("User does not exist")
    }

    ud.ID = userID
    updateStatement := `update "userdata" set "name"=$1, "surname"=$2, "description"=$3 where "userid"=$4`
    _, err = db.Exec(updateStatement, ud.Name, ud.Surname, ud.Description, ud.ID)
    if err != nil {
        return err
    }

    return nil
}
