const appUser = process.env.APP_DB_USER;
const appPass = process.env.APP_DB_PASS;
const dbName = process.env.MONGO_INITDB_DATABASE;

db = db.getSiblingDB(dbName);

db.createUser({
  user: appUser,
  pwd: appPass,
  roles: [
    {
      role: "readWrite",
      db: dbName
    }
  ]
});