#!bin/bash
csv_url="https://url-of-csv-file"
now=$(date +'%Y-%m-%d__%H-%M-%S')
file_name="csv-$now.csv"
echo $file_name
curl $csv_url > $file_name
./geo-service import -f $file_name -d mongodb://localhost:27017 -u user -p pass
