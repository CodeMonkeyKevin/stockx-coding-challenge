# StockX Coding Challenge

# Setup and Starting the API
The project contains Dockerfile will install Postgresql, create the necessary database, table, function, trigger, build and exec the API.

###### Build Image:
```sh
$ docker build -t stockx .
```
###### Run Container:
```sh
$ docker run --rm -v stockxcc:/var/lib/postgresql -p 3000:3000 -p 5432:5432 --name stockxcc stockx
```

```-p 5432:5432``` can be removed if you do not wish to connect Postgres inside the contrainer to inspect the database/table schema.

# API Calls

##### Adding Shoe Data: POST /shoes 

`POST /shoes` requires JSON payload:
```JSON
{
    "shoe": "AJ 1 Mid Cool Blue",
    "trueToSizeVal": 4
}
```

###### Request:
```curl
$ curl -v -X POST -H "Content-Type: application/json" localhost:3000/shoes -d '{"shoe": "AJ 1 Mid Cool Blue","trueToSizeVal": 4}'
```

###### Response:
```curl
HTTP/1.1 201 Created
Content-Type: application/json

{
    "id":34,
    "shoe":"AJ 1 Mid Cool Blue","trueToSizeData":[4],
    "trueToSizeCalculation":4
}
```

##### Getting Shoes Data: GET /shoes
`GET /shoes` returns a JSON array of all shoes data:

###### Request:
```curl
$ curl -v localhost:3000/shoes/1
```

###### Response:
```curl
HTTP/1.1 200 OK
Content-Type: application/json

[
    {
        "id":1,
        "shoe":"AJ MID",
        "trueToSizeData":[4,4,4],
        "trueToSizeCalculation":4
    },{
        "id":34,
        "shoe":"AJ 1 Mid Cool Blue",
        "trueToSizeData":[4],
        "trueToSizeCalculation":4
    }
]
```

##### Getting Shoe Data: GET /shoes/:ID
`GET /shoes/:ID` returns a JSON object of shoe by ID:

###### Request:
```curl
$ curl -v localhost:3000/shoes/1
```

###### Response:
```curl
HTTP/1.1 200 OK
Content-Type: application/json

{
    "id":34,
    "shoe":"AJ 1 Mid Cool Blue",
    "trueToSizeData":[4],
    "trueToSizeCalculation":4
}
```

##### Remove Shoe Data: DELETE /shoes/:ID
`DELETE /shoes/:ID` will delete shoe data:

###### Request:
```curl
$ curl -v -X DELETE localhost:3000/shoes/4
```

###### Response:
```curl
HTTP/1.1 200 OK
Content-Type: application/json

{"result":"success"}
```
