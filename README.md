CSV to JSON line

This code will convert CSV to JSON line file. In order to run this code, you need to clone the code first, then in the terminal exectute:
"go build -o csvtojl main.go"

This will give you an executable "csvtojl"

Now, copy your csv file to a same directory as your executable. Then run:

"./csvtojl inputFile.csv outputFile.jl"

This will create a new JSON Line file based on your CSV file.

There is a separate tests that we can run or modify.

To run the tests, in terminal, run:

"go test -v"
