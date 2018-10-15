# StockX Coding Challenge

### Setup and Starting the API
The project contains Dockerfile will install Postgresql, create the necessary database, table, function, trigger, build and exec the API. Git clone or download this repo and cd into root dir of this repo.

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
    "shoe":"AJ 1 Mid Cool Blue",
    "trueToSizeData":[4],
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

<br>


#### DATABASE

To connect to Postgres instance in the docker container:

```psql -h localhost -U docker -d stockxcc_kp```

Password is ```docker```

<br>

Table ```shoes``` schema:

```
CREATE TABLE shoes (
    id                            SERIAL,
    name                          TEXT NOT NULL,
    "trueToSizeData"              int[] NOT NULL DEFAULT '{}',
    "trueToSizeCalculation"       numeric(14,13) NOT NULL DEFAULT 0.00,
    CONSTRAINT shoes_pkey PRIMARY KEY (id)
);
```

Note: ```trueToSizeData``` field is an postgres Array data type. When updating this array postgres append_arr function is used. This is avoid having to rewriting the whole array updated with new data.


Function ```CalculationTrueToSize``` and Trigger ```update_trueToSizeCalculation```:

```
CREATE OR REPLACE FUNCTION CalculationTrueToSize()
RETURNS TRIGGER AS $$
DECLARE
    arrLen integer;
    arrSum integer;
BEGIN
    SELECT cardinality("trueToSizeData") INTO arrLen
    FROM shoes WHERE id=new.id;
    
    SELECT SUM(UNNEST(t)) INTO arrSum
    FROM (SELECT UNNEST("trueToSizeData") FROM shoes WHERE id=new.id) t;
    
    UPDATE shoes SET "trueToSizeCalculation" = (arrSum::NUMERIC/arrLen::NUMERIC) WHERE id=new.id;
    
    RETURN new;
END;
$$ LANGUAGE plpgsql;
```

```
CREATE TRIGGER update_trueToSizeCalculation 
AFTER UPDATE OF "trueToSizeData" ON shoes
FOR EACH ROW EXECUTE PROCEDURE CalculationTrueToSize();
```

```CalculationTrueToSize()``` function is triggered when new data is added to the ```trueToSizeData``` column. This function calculates and updates the ```trueToSizeCalculation```.
