# Hotel reservation backend 

# Project environment variables
```
HTTP_LISTEN_ADDRESS=:3000
JWT_SECRET=somethingsupersecretthatNOBODYKNOWS
MONGO_DB_NAME=hotel-reservation
MONGO_DB_URL=mongodb://localhost:27017
```

## Project outlines
- users -> book room from an hotel
- admins -> going too check reservation/bookings
- authen and autoriz -> JWT tokens
- Hotels -> CRUD API -> JSON
- Rooms -> CRUD API -> JSON
- Scripts -> database management -> seeding, migration





Installing MongoDB
```
go get go.mongodb.org/mongo-driver/mongo
```

### gofiber
Documentation
```
https://gofiber.io
```

Installing go fiber
```
go get github.com/gofiber/fiber/v2
```

## Docker
### Installing mongodb as a Docker container
```
docker run --name mongodb -d mongo:latest -p 27017:27017    
```
