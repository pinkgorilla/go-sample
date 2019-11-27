# Table Description

## User 

menyimpan data profile `user`

| Field | Type | Nullable | Description | note
| --- | --- | :-:| --- | --- |
| userId | `int` | false | unique identifier for user
| name | `varchar` | false | nama lengkap user, dipakai sebagai display name
| age | `int` | false | usia user 17 ~ 100 
| option | `varchar(256)` | false | user option | deprecated
| metadata | `json` | false | data dinamis user | please refer [`User.metadata`](##User.metadata)

## User.metadata

data dinamis user

| Field | Type | Nullable | Description | note
| --- | --- | :-:| --- | --- |
| field1 | `int` | false | field1 of User.metadata
| field2 | `varchar(256)` | false | field2 of User.metadata | some note

 