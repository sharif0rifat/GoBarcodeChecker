# GoBarcodeChecker

1. run the program by 'go run .\Server.go'  
2. Then it  will start listening to port 8081
3. First you need to go to  'http://localhost:8081/GenerateBarcodes'  this link to generate barcode
it generate 100 files in 'Database' folder with some texts (refered as barcode)
4. then you need to go to this link 'http://localhost:8081/SearchBarcode?barcode=some_value'  to search bar code 
5. You can cehck by copying one of the line from those database files and and check