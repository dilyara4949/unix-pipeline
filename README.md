# Unix-pipeline

## Description

A pipeline operator in Golang that mimics Unix pipelines, taking a string from stdout containing a list of commands and arguments split by the | operator. The commands include cat for reading, creating, and concatenating files, grep for searching text patterns within files, and sort for sorting the contents of a text file. 

## Usage

file.txt with the following content:
 ```
 user1 data1
user3 data3
user2 data2
```

example how to run program:

``` 
% go run .
cat file.txt | grep user | sort
user1 data1
user2 data2
user3 data3
```
