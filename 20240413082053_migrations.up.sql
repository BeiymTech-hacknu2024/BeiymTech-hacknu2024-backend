CREATE TABLE IF NOT EXISTS users (
    ID SERIAL PRIMARY KEY,
    Role VARCHAR(50) NOT NULL,
    Name VARCHAR(50) NOT NULL,
    Surname VARCHAR(50) NOT NULL,
    Email VARCHAR(100) NOT NULL,
    Password VARCHAR(100) NOT NULL,
    FOREIGN KEY (Role) REFERENCES roles(Rolename),

)


CREATE TABLE IF NOT EXISTS roles (
    ID SERIAL PRIMARY KEY,
    Rolename VARCHAR(50) NOT NULL,
    UNIQUE(Rolename)
)



CREATE TABLE IF NOT EXISTS reports (
    UserID INTEGER NOT NULL
    ID SERIAL PRIMARY KEY,
    USER VARCHAR(50) NOT NULL,
    FOREIGN KEY (UserID) REFERENCES users(ID)

); 

CREATE TABLE IF NOT EXISTS assignments {
    ID SERIAL PRIMARY KEY,
    Question VARCHAR(100) NOT NULL,
    Answer VARCHAR(50) NOT NULL,

}

