When you start with main folder, you do 

    go mod init sms_gateway

After that you execute

    go mod tidy

Thereafter, when you are create a new package, let us say logger (sub folder/package)

    sms_gateway/
	    |__	logger/
			    |__ logger.go

Inside logger package / folder, if you create a file called 

    logger.go 

and add external libraries, then you need to run below at sms_gateway folder level 

    go mod tidy
