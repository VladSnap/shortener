del shortener.log
go build -o cmd\shortener\shortener.exe cmd\shortener\main.go

shortenertest -test.v -test.run=^TestIteration1$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -source-path=C:\Develop\go\projects\shortener
shortenertest -test.v -test.run=^TestIteration2$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -source-path=C:\Develop\go\projects\shortener
shortenertest -test.v -test.run=^TestIteration3$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -source-path=C:\Develop\go\projects\shortener
shortenertest -test.v -test.run=^TestIteration4$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -source-path=C:\Develop\go\projects\shortener -server-port=8080
shortenertest -test.v -test.run=^TestIteration5$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -source-path=C:\Develop\go\projects\shortener -server-port=8080
shortenertest -test.v -test.run=^TestIteration6$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -source-path=C:\Develop\go\projects\shortener
shortenertest -test.v -test.run=^TestIteration7$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -source-path=C:\Develop\go\projects\shortener
shortenertest -test.v -test.run=^TestIteration8$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -source-path=C:\Develop\go\projects\shortener
shortenertest -test.v -test.run=^TestIteration9$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -source-path=C:\Develop\go\projects\shortener -file-storage-path=D://shortenDataBase.json
shortenertest -test.v -test.run=^TestIteration10$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -source-path=C:\Develop\go\projects\shortener -database-dsn="host=localhost user=postgres password=qwerty dbname=shortener sslmode=disable"
::shortenertest -test.v -test.run=^TestIteration11$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -database-dsn="host=localhost user=postgres password=qwerty dbname=shortener sslmode=disable"
shortenertest -test.v -test.run=^TestIteration12$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -database-dsn="host=localhost user=postgres password=qwerty dbname=shortener sslmode=disable"
shortenertest -test.v -test.run=^TestIteration13$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -database-dsn="host=localhost user=postgres password=qwerty dbname=shortener sslmode=disable"
shortenertest -test.v -test.run=^TestIteration14$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -database-dsn="host=localhost user=postgres password=qwerty dbname=shortener sslmode=disable"
shortenertest -test.v -test.run=^TestIteration15$ -binary-path=C:\Develop\go\projects\shortener\cmd\shortener\shortener.exe -database-dsn="host=localhost user=postgres password=qwerty dbname=shortener sslmode=disable"

del cmd\shortener\shortener.exe

pause