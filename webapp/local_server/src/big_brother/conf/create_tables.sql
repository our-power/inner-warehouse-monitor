CREATE TABLE user (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT UNIQUE, passwd TEXT, email TEXT, register_time TEXT, last_login TEXT, role_id INTEGER);
CREATE TABLE role (id INTEGER PRIMARY KEY AUTOINCREMENT, role_type TEXT, permission TEXT);
CREATE TABLE trace (id INTEGER PRIMARY KEY AUTOINCREMENT, user TEXT, do_what TEXT, that_time TEXT);
INSERT INTO role(role_type, permission) VALUES ("user_admin", "view|add|modify|del|admin");
INSERT INTO user(name, passwd, role_id) VALUES ("root", "03941b924d12454219648d61a7b025e1", 0);