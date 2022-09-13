#!/bin/bash

case "$1" in
    "research")
        db="research"
        ;;
    "production")
        echo "This will erase all production data. This action CANNOT be undone."
        echo "Do you wish to continue?"
        select yn in Yes No
        do
            case "$yn" in
                "No")
                    exit
                    ;;
                "Yes")
                    break
                    ;;
            esac
        done
        db="production"
        ;;
    *)
        echo "Error: invalid target database. Expected either: 'production' or 'research'"
        exit -1
        ;;
esac

echo $db

echo "Running Global Init..."
echo "\i globalInit.sql" | psql $db

echo "Copying exercise types..."
echo "\\copy ExerciseTypes(_type,description) FROM 'ExerciseTypesInit.csv' CSV HEADER;" | psql $db

echo "Copying exercise focus..."
echo "\\copy ExerciseFocus(focus) FROM 'ExerciseFocusInit.csv' CSV HEADER;" | psql $db

accId=`psql $db -c "SELECT id FROM ExerciseTypes WHERE _type='Accessory'"`
accId=`echo $accId | grep -o --regexp=" [0-9] "`
compAccId=`psql $db -c "SELECT id FROM ExerciseTypes WHERE _type='Compound Accessory'"`
compAccId=`echo $compAccId | grep -o --regexp=" [0-9] "`
mainCompAccId=`psql $db -c "SELECT id FROM ExerciseTypes WHERE _type='Main Compound Accessory'"`
mainCompAccId=`echo $mainCompAccId | grep -o --regexp=" [0-9] "`
mainCompId=`psql $db -c "SELECT id FROM ExerciseTypes WHERE _type='Main Compound'"`
mainCompId=`echo $mainCompId | grep -o --regexp=" [0-9] "`

squatFocId=`psql $db -c "SELECT id FROM ExerciseFocus WHERE focus='Squat'"`
squatFocId=`echo $squatFocId | grep -o --regexp=" [0-9] "`
benchFocId=`psql $db -c "SELECT id FROM ExerciseFocus WHERE focus='Bench'"`
benchFocId=`echo $benchFocId | grep -o --regexp=" [0-9] "`
deadliftFocId=`psql $db -c "SELECT id FROM ExerciseFocus WHERE focus='Deadlift'"`
deadliftFocId=`echo $deadliftFocId | grep -o --regexp=" [0-9] "`

cp ExercisesInit.csv ExercisesInit.tmp.csv

sed -i "s/Main Compound Accessory/$mainCompAccId/g" ExercisesInit.tmp.csv
sed -i "s/Main Compound/$mainCompId/g" ExercisesInit.tmp.csv
sed -i "s/Compound Accessory/$compAccId/g" ExercisesInit.tmp.csv
sed -i "s/Accessory/$accId/g" ExercisesInit.tmp.csv

sed -i "s/,Squat/,$squatFocId/g" ExercisesInit.tmp.csv
sed -i "s/,Bench/,$benchFocId"/g ExercisesInit.tmp.csv
sed -i "s/,Deadlift/,$deadliftFocId"/g ExercisesInit.tmp.csv


echo "Copying exercises..."
echo "\\copy Exercises(name,typeID,focusID) FROM 'ExercisesInit.tmp.csv' CSV HEADER;" | psql $db

rm ExercisesInit.tmp.csv
