Coded by: Tahseen Jamal

This sms gateway has been written in Go programming language. It has 3 primary modules

    1. Logger
    2. Message Broker
    3. SMPP Connector
    4. HTTP handler

The main application initializes them and uses them, except for Logger, which is initialized as Singleton and then any module calling gets the single instance of the Logger.

Pull the code from master using

    git clone

Post that you need to run

    go mod init sms_gateway
    go mod tidy

After that open the configuration file

    main.properties

It has configuration related to all the modules

    consumer.properties.filename=logger/logger.properties
    
    #ActiveMQ properties
    activemq.broker.url=localhost:61613
    activemq.broker.username=admin
    activemq.broker.password=admin
    
    
    #SMPP properties
    smpp.host=localhost
    smpp.port=2775
    smpp.systemId=
    smpp.password=
    smpp.systemType=
    smpp.sourceTon=5
    smpp.sourceNpi=1
    smpp.destinationTon=1
    smpp.destinationNpi=1
    #smpp.windowSize=
    smpp.prefixPlus=true
    
    #SMPP Client TPS properties
    smpp.tps=5
    
    
    #SMPP Black Hour properties
    #Format: HH:mm (24 hour format)
    #Don't use quotes
    smpp.morning=09:00
    smpp.evening=18:00


